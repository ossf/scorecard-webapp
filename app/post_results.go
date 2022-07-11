// Copyright 2022 Security Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/gorilla/mux"
	"gocloud.dev/blob"
)

const (
	// OID source: https://github.com/sigstore/fulcio/blob/96ef49cc7662912ba37d46f738757e8d8d5b5355/docs/oid-info.md#L33
	// TODO: retrieve these by name.
	fulcioRepoRefKey  = "1.3.6.1.4.1.57264.1.6"
	fulcioRepoPathKey = "1.3.6.1.4.1.57264.1.5"
	fulcioRepoSHAKey  = "1.3.6.1.4.1.57264.1.3"
	resultsBucket     = "gs://ossf-scorecard-results"
	resultsFile       = "results.json"
)

var (
	errWritingBucket            = errors.New("error writing to GCS bucket")
	errMultipleCerts            = errors.New("multiple certificates found for the entry")
	errEmptyCertRef             = errors.New("cert has empty repository ref")
	errEmptyCertPath            = errors.New("cert has empty repository path")
	errCertMissingURI           = errors.New("certificate has no URIs")
	errCertWorkflowPathEmpty    = errors.New("cert workflow path is empty")
	errMismatchedCertAndRequest = errors.New("repository and branch of cert doesn't match that of request")
)

type ScorecardResult struct {
	Result      string `json:"result"`
	Branch      string `json:"branch"`
	AccessToken string `json:"accessToken"`
}

type certInfo struct {
	repoFullName  string
	repoBranchRef string
	repoSHA       string
	workflowPath  string
}

type tlogEntry struct {
	Body           string `json:"body"`
	IntegratedTime int64  `json:"integratedTime"`
}

//go:embed fulcio_v1.crt.pem
var fulcioRoot []byte

//go:embed fulcio__intermediate.crt.pem
var fulcioIntermediate []byte

func PostResultsHandler(w http.ResponseWriter, r *http.Request) {
	// Sanity check
	host := mux.Vars(r)["host"]
	orgName := mux.Vars(r)["orgName"]
	repoName := mux.Vars(r)["repoName"]

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprintf(w, "error reading request body")
		if err != nil {
			log.Printf("error during Write: %v", err)
		}
		return
	}
	var scorecardResult ScorecardResult
	if err := json.Unmarshal(reqBody, &scorecardResult); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprint(w, "error unmarshaling input JSON")
		if err != nil {
			log.Printf("error during Write: %v", err)
		}
		return
	}

	// Process
	err = processRequest(host, orgName, repoName, scorecardResult)
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		return
	}
	if errors.Is(err, errMismatchedCertAndRequest) ||
		errors.Is(err, errGlobalVarsOrDefaults) ||
		errors.Is(err, errGlobalWriteAll) ||
		errors.Is(err, errGlobalWrite) ||
		errors.Is(err, errScorecardJobNotFound) ||
		errors.Is(err, errNonScorecardJobHasTokenWrite) ||
		errors.Is(err, errJobHasContainerOrServices) ||
		errors.Is(err, errScorecardJobRunsOn) ||
		errors.Is(err, errUnallowedStepName) ||
		errors.Is(err, errScorecardJobEnvVars) ||
		errors.Is(err, errScorecardJobDefaults) ||
		errors.Is(err, errEmptyStepUses) {
		w.WriteHeader(http.StatusBadRequest)
		if _, errWrite := fmt.Fprintf(w, "%v", err); errWrite != nil {
			log.Printf("%v", errWrite)
		}
		return
	}
	log.Printf("%v", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func processRequest(host, org, repo string, scorecardResult ScorecardResult) error {
	ctx := context.Background()
	cert, err := extractAndVerifyCertForPayload(ctx, []byte(scorecardResult.Result))
	if err != nil {
		return fmt.Errorf("error extracting cert: %w", err)
	}

	info, err := extractCertInfo(cert)
	if err != nil {
		return fmt.Errorf("error extracting cert info: %w", err)
	}
	if info.repoFullName != fmt.Sprintf("%s/%s", org, repo) ||
		(info.repoBranchRef != scorecardResult.Branch &&
			info.repoBranchRef != fmt.Sprintf("refs/heads/%s", scorecardResult.Branch)) {
		return fmt.Errorf("%w", errMismatchedCertAndRequest)
	}

	if err := getAndVerifyWorkflowContent(ctx, org, repo, scorecardResult, info); err != nil {
		return fmt.Errorf("error verifying workflow: %w", err)
	}

	// Save scorecard results (results.json, score.txt) to GCS
	bucketURL := resultsBucket
	objectPath := fmt.Sprintf("%s/%s/%s/%s", host, org, repo, resultsFile)

	if err := writeToBlobStore(ctx, bucketURL, objectPath, []byte(scorecardResult.Result)); err != nil {
		return fmt.Errorf(errWritingBucket.Error()+": %v, %v", err)
	}
	return nil
}

func getAndVerifyWorkflowContent(ctx context.Context,
	org, repo string, scorecardResult ScorecardResult, info certInfo,
) error {
	// Get the corresponding GitHub repository.
	httpClient := http.DefaultClient
	if scorecardResult.AccessToken != "" {
		httpClient.Transport = githubTransport{
			token: scorecardResult.AccessToken,
		}
	}
	client := github.NewClient(httpClient)
	repoClient, _, err := client.Repositories.Get(ctx, org, repo)
	if err != nil {
		return fmt.Errorf("error getting repository: %w", err)
	}

	// Verify that the branch from the results files is the repo's default branch.
	defaultBranch := repoClient.GetDefaultBranch()
	if scorecardResult.Branch != defaultBranch &&
		scorecardResult.Branch != fmt.Sprintf("refs/heads/%s", defaultBranch) {
		return fmt.Errorf("branch of cert isn't the repo's default branch")
	}

	// Get workflow file from cert commit SHA.
	opts := &github.RepositoryContentGetOptions{Ref: info.repoSHA}
	contents, _, _, err := client.Repositories.GetContents(ctx, org, repo, info.workflowPath, opts)
	if err != nil {
		return fmt.Errorf("error downloading workflow contents from repo: %v", err)
	}

	workflowContent, err := contents.GetContent()
	if err != nil {
		return fmt.Errorf("error decoding workflow contents: %v", err)
	}

	// Verify scorecard workflow.
	if err := verifyScorecardWorkflow(workflowContent); err != nil {
		return fmt.Errorf("workflow could not be verified: %v", err)
	}
	return nil
}

func writeToBlobStore(ctx context.Context, bucketURL, filename string, data []byte) error {
	bucket, err := blob.OpenBucket(ctx, bucketURL)
	if err != nil {
		return fmt.Errorf("error from blob.OpenBucket: %w", err)
	}
	defer bucket.Close()

	blobWriter, err := bucket.NewWriter(ctx, filename, nil)
	if err != nil {
		return fmt.Errorf("error from bucket.NewWriter: %w", err)
	}
	if _, err = blobWriter.Write(data); err != nil {
		return fmt.Errorf("error from blobWriter.Write: %w", err)
	}
	if err := blobWriter.Close(); err != nil {
		return fmt.Errorf("error from blobWriter.Close: %w", err)
	}
	return nil
}

func extractAndVerifyCertForPayload(ctx context.Context, payload []byte) (*x509.Certificate, error) {
	// Get most recent Rekor entry uuid.
	uuids, err := getUUIDsByPayload(ctx, payload)
	if err != nil || len(uuids) == 0 {
		return nil, fmt.Errorf("error finding tlog entries corresponding to payload: %w", err)
	}
	// TODO(#135): We can't simply take the latest UUID. Either:
	// (a) iterate through all returned UUIDs to find the right one.
	// (b) send tlog index in the POST payload to identify the corresponding UUID.
	uuid := uuids[len(uuids)-1] // ignore past entries.

	// Get tlog entry from the UUID.
	entry, err := getTLogEntry(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("error fetching tlog entry: %w", err)
	}

	// Verify tlog entry to make sure it is actually in the log.
	// TODO(#113): What exactly are we verifying here and how can we do it
	// based on the info returned from the API?
	//err = cosign.VerifyTLogEntry(ctx, rekorClient, entry)
	//if err != nil {
	//	return nil, fmt.Errorf("error verifying tlog entry: %w", err)
	//}

	// Extract certificate.
	certs, err := extractCerts(entry)
	if err != nil || len(certs) == 0 {
		return nil, fmt.Errorf("error extracting certificate from entry: %w", err)
	}
	if len(certs) > 1 {
		return nil, fmt.Errorf("%w", errMultipleCerts)
	}

	cert := certs[0]

	// Verify the certificate against Fulcio Root CA
	roots, intermediates, err := getFulcioRoots()
	if err != nil {
		return nil, fmt.Errorf("retrieving Fulcio root: %w", err)
	}
	if _, err := cert.Verify(x509.VerifyOptions{
		CurrentTime:   cert.NotBefore,
		Roots:         roots,
		Intermediates: intermediates,
		KeyUsages: []x509.ExtKeyUsage{
			x509.ExtKeyUsageCodeSigning,
		},
	}); err != nil {
		return nil, fmt.Errorf("verifying Fulcio issued certificate: %w", err)
	}

	// Verify that cert isn't expired.
	integratedTime := time.Unix(entry.IntegratedTime, 0)
	if cert.NotAfter.Before(integratedTime) {
		return nil, fmt.Errorf("certificate expired before signatures were entered in log: %s is before %s",
			cert.NotAfter, integratedTime)
	}
	if cert.NotBefore.After(integratedTime) {
		return nil, fmt.Errorf("certificate was issued after signatures were entered in log: %s is after %s",
			cert.NotAfter, integratedTime)
	}
	return cert, nil
}

func getUUIDsByPayload(ctx context.Context, payload []byte) ([]string, error) {
	payloadSHA := sha256.Sum256(payload)
	rekorPayload := struct {
		Hash string `json:"hash"`
	}{
		Hash: fmt.Sprintf("sha256:%s", hex.EncodeToString(payloadSHA[:])),
	}
	jsonPayload, err := json.Marshal(rekorPayload)
	if err != nil {
		return nil, fmt.Errorf("marshaling json payload: %w", err)
	}

	rekorReq, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		"https://rekor.sigstore.dev/api/v1/index/retrieve",
		bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("creating new HTTP request: %w", err)
	}
	rekorReq.Header.Add("Content-Type", "application/json")
	rekorReq.Header.Add("accept", "application/json")
	resp, err := http.DefaultClient.Do(rekorReq)
	if err != nil {
		return nil, fmt.Errorf("looking up Rekor index: %w", err)
	}
	defer resp.Body.Close()

	var rekorResult []string
	if err := json.NewDecoder(resp.Body).Decode(&rekorResult); err != nil {
		return nil, fmt.Errorf("decoding Rekor response: %w", err)
	}

	return rekorResult, nil
}

func getTLogEntry(ctx context.Context, uuid string) (*tlogEntry, error) {
	rekorReq, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf("https://rekor.sigstore.dev/api/v1/log/entries/%s", uuid),
		nil)
	if err != nil {
		return nil, fmt.Errorf("creating new HTTP request: %w", err)
	}
	rekorReq.Header.Add("accept", "application/json")
	resp, err := http.DefaultClient.Do(rekorReq)
	if err != nil {
		return nil, fmt.Errorf("looking up Rekor index: %w", err)
	}
	defer resp.Body.Close()

	var rekorResult map[string]tlogEntry
	if err := json.NewDecoder(resp.Body).Decode(&rekorResult); err != nil {
		return nil, fmt.Errorf("decoding Rekor response: %w", err)
	}

	for _, res := range rekorResult {
		return &tlogEntry{
			Body:           res.Body,
			IntegratedTime: res.IntegratedTime,
		}, nil
	}

	return nil, fmt.Errorf("unexpected error")
}

func extractCertInfo(cert *x509.Certificate) (certInfo, error) {
	ret := certInfo{}
	// Get repo reference & path from cert.
	for _, ext := range cert.Extensions {
		if ext.Id.String() == fulcioRepoRefKey {
			if len(ext.Value) == 0 {
				return ret, fmt.Errorf("%w", errEmptyCertRef)
			}
			ret.repoBranchRef = string(ext.Value)
		}
		if ext.Id.String() == fulcioRepoPathKey {
			if len(ext.Value) == 0 {
				return ret, fmt.Errorf("%w", errEmptyCertPath)
			}
			ret.repoFullName = string(ext.Value)
		}
		if ext.Id.String() == fulcioRepoSHAKey {
			ret.repoSHA = string(ext.Value)
		}
	}

	// Get workflow job ref from the certificate.
	if len(cert.URIs) == 0 {
		return ret, errCertMissingURI
	}
	workflowRef := cert.URIs[0].Path
	if len(workflowRef) == 0 {
		return ret, errCertWorkflowPathEmpty
	}

	// Remove repo path from workflow filepath.
	ind := strings.Index(workflowRef, ret.repoFullName) + len(ret.repoFullName) + 1
	ret.workflowPath = workflowRef[ind:]

	// Remove repo ref tag from workflow filepath.
	ret.workflowPath = ret.workflowPath[:strings.Index(ret.workflowPath, "@")]
	return ret, nil
}

func extractCerts(entry *tlogEntry) ([]*x509.Certificate, error) {
	b, err := base64.StdEncoding.DecodeString(entry.Body)
	if err != nil {
		return nil, err
	}

	var entryBody struct {
		Spec struct {
			Signature struct {
				PublicKey struct {
					Content string `json:"content"`
				} `json:"publicKey"`
			} `json:"signature"`
		} `json:"spec"`
	}
	if err := json.Unmarshal(b, &entryBody); err != nil {
		return nil, err
	}

	publicKeyB64 := entryBody.Spec.Signature.PublicKey.Content
	publicKey, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return nil, err
	}

	remaining := publicKey
	var result []*x509.Certificate
	for len(remaining) > 0 {
		var certDer *pem.Block
		certDer, remaining = pem.Decode(remaining)
		if certDer == nil {
			return nil, fmt.Errorf("error during PEM decoding: %w", err)
		}

		cert, err := x509.ParseCertificate(certDer.Bytes)
		if err != nil {
			return nil, fmt.Errorf("error during certificate parsing: %w", err)
		}
		result = append(result, cert)
	}
	return result, nil
}

func getFulcioRoots() (*x509.CertPool, *x509.CertPool, error) {
	rootPool := x509.NewCertPool()
	intermediatePool := x509.NewCertPool()

	if ok := rootPool.AppendCertsFromPEM(fulcioRoot); !ok {
		return nil, nil, fmt.Errorf("unmarshalling fulcio root")
	}

	if ok := intermediatePool.AppendCertsFromPEM(fulcioIntermediate); !ok {
		return nil, nil, fmt.Errorf("unmarshalling fulcio intermediate")
	}

	return rootPool, intermediatePool, nil
}
