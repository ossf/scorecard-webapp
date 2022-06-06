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
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/google/go-github/v42/github"
	"github.com/gorilla/mux"
	"github.com/sigstore/cosign/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/cmd/cosign/cli/rekor"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/rekor/pkg/generated/client"
	"github.com/sigstore/rekor/pkg/generated/models"
	"github.com/sigstore/rekor/pkg/types"
	hashedrekord "github.com/sigstore/rekor/pkg/types/hashedrekord/v0.0.1"
	rekord "github.com/sigstore/rekor/pkg/types/rekord/v0.0.1"
	"github.com/sigstore/sigstore/pkg/cryptoutils"
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
	if errors.Is(err, errMismatchedCertAndRequest) {
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
	rekorClient, err := rekor.NewClient(options.DefaultRekorURL)
	if err != nil {
		return fmt.Errorf("error initializing Rekor client: %w", err)
	}

	ctx := context.Background()
	cert, err := extractAndVerifyCertForPayload(ctx, rekorClient, []byte(scorecardResult.Result))
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

func extractAndVerifyCertForPayload(ctx context.Context,
	rekorClient *client.Rekor, payload []byte,
) (*x509.Certificate, error) {
	// Get most recent Rekor entry uuid.
	uuids, err := cosign.FindTLogEntriesByPayload(ctx, rekorClient, payload)
	if err != nil || len(uuids) == 0 {
		return nil, fmt.Errorf("error finding tlog entries corresponding to payload: %w", err)
	}
	uuid := uuids[len(uuids)-1] // ignore past entries.

	// Get tlog entry from the UUID.
	entry, err := cosign.GetTlogEntry(ctx, rekorClient, uuid)
	if err != nil {
		return nil, fmt.Errorf("error fetching tlog entry: %w", err)
	}

	// Verify tlog entry to make sure it is actually in the log.
	err = cosign.VerifyTLogEntry(ctx, rekorClient, entry)
	if err != nil {
		return nil, fmt.Errorf("error verifying tlog entry: %w", err)
	}

	// Extract certificate.
	certs, err := extractCerts(entry)
	if err != nil || len(certs) == 0 {
		return nil, fmt.Errorf("error extracting certificate from entry: %w", err)
	}
	if len(certs) > 1 {
		return nil, fmt.Errorf("%w", errMultipleCerts)
	}

	cert := certs[0]
	// Verify that cert isn't expired.
	if err = cosign.CheckExpiry(cert, time.Unix(*entry.IntegratedTime, 0)); err != nil {
		return nil, fmt.Errorf("error during cosign.CheckExpiry: %w", err)
	}
	return cert, nil
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

// nolint:lll
// Source: https://github.com/sigstore/cosign/blob/18d2ce0b458018951f7356db911467a427a8dffe/cmd/cosign/cli/verify/verify_blob.go#L321
func extractCerts(e *models.LogEntryAnon) ([]*x509.Certificate, error) {
	b, err := base64.StdEncoding.DecodeString(e.Body.(string))
	if err != nil {
		return nil, err
	}

	pe, err := models.UnmarshalProposedEntry(bytes.NewReader(b), runtime.JSONConsumer())
	if err != nil {
		return nil, err
	}

	eimpl, err := types.NewEntry(pe)
	if err != nil {
		return nil, err
	}

	var publicKeyB64 []byte
	switch e := eimpl.(type) {
	case *rekord.V001Entry:
		publicKeyB64, err = e.RekordObj.Signature.PublicKey.Content.MarshalText()
		if err != nil {
			return nil, err
		}
	case *hashedrekord.V001Entry:
		publicKeyB64, err = e.HashedRekordObj.Signature.PublicKey.Content.MarshalText()
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unexpected tlog entry type")
	}

	publicKey, err := base64.StdEncoding.DecodeString(string(publicKeyB64))
	if err != nil {
		return nil, err
	}

	certs, err := cryptoutils.UnmarshalCertificatesFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	if len(certs) == 0 {
		return nil, errors.New("no certs found in pem tlog")
	}

	return certs, err
}
