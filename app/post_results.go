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
	"github.com/sigstore/cosign/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/cmd/cosign/cli/rekor"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/rekor/pkg/generated/models"
	"github.com/sigstore/rekor/pkg/types"
	hashedrekord "github.com/sigstore/rekor/pkg/types/hashedrekord/v0.0.1"
	rekord "github.com/sigstore/rekor/pkg/types/rekord/v0.0.1"
	"github.com/sigstore/sigstore/pkg/cryptoutils"
	"gocloud.dev/blob"
)

type ScorecardOutput struct {
	JSONOutput string
}

func PostResultsHandler(w http.ResponseWriter, r *http.Request) {
	err := handler(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading http request body: %v", err)
	}

	// Unmarshal request body.
	var scorecardOutput ScorecardOutput
	err = json.Unmarshal(reqBody, &scorecardOutput)
	if err != nil {
		return fmt.Errorf("error unmarshalling request body: %v", err)
	}

	// Grab parameters.
	reqRepo := r.Header["X-Repository"]
	reqBranch := r.Header["X-Branch"]
	if len(reqRepo) == 0 || len(reqBranch) == 0 {
		return fmt.Errorf("error: empty parameters")
	}

	err = verifySignature(ctx, scorecardOutput, reqRepo[0], reqBranch[0])
	if err != nil {
		return err
	}

	// Write response.
	w.Write([]byte(
		fmt.Sprintf(
			"Successfully verified and uploaded scorecard results for repo %s on branch %s",
			reqRepo[0], reqBranch[0])))

	return nil
}

var errorWritingBucket = errors.New("error writing to GCS bucket")

// verifySignature receives the scorecard analysis payload, looks up its associated tlog entry and
// certificate, and extracts the repository's workflow file to ensure its legitimacy.
func verifySignature(ctx context.Context, scorecardOutput ScorecardOutput, reqRepo, reqBranch string) error {
	// Lookup results payload to get the repo info from the corresponding entry & cert.
	repoPath, repoRef, repoSHA, workflowPath, err := lookupPayload(ctx, []byte(scorecardOutput.JSONOutput))
	if err != nil {
		return fmt.Errorf("error looking up json results: %v", err)
	}

	// Split repo path into owner and repo name.
	ownerName := repoPath[:strings.Index(repoPath, "/")]
	repoName := repoPath[strings.Index(repoPath, "/")+1:]

	// Verify that the repository and branch of the cert and request are equal.
	if len(reqRepo) == 0 || len(reqBranch) == 0 || reqRepo != repoPath || reqBranch != repoRef {
		return fmt.Errorf("repository and branch of cert doesn't match that of request")
	}

	// Get the corresponding GitHub repository.
	// TODO: use GITHUB_TOKEN from workflow to make the api call.
	client := github.NewClient(nil)
	repo, _, err := client.Repositories.Get(ctx, ownerName, repoName)
	if err != nil {
		return fmt.Errorf("error getting repository: %v", err)
	}

	// Verify that the branch from the results files is the repo's default branch.
	defaultBranch := "refs/heads/" + repo.GetDefaultBranch()
	if defaultBranch != repoRef {
		return fmt.Errorf("branch of cert isn't the repo's default branch")
	}

	// Get workflow file from cert commit SHA.
	opts := &github.RepositoryContentGetOptions{Ref: repoSHA}
	contents, _, _, err := client.Repositories.GetContents(ctx, ownerName, repoName, workflowPath, opts)
	if err != nil {
		return fmt.Errorf("error downloading workflow contents from repo: %v", err)
	}

	workflowContent, err := contents.GetContent()
	if err != nil {
		return fmt.Errorf("error decoding workflow contents: %v", err)
	}

	// Verify scorecard workflow.
	err = verifyScorecardWorkflow(workflowContent)
	if err != nil {
		return fmt.Errorf("workflow could not be verified: %v", err)
	}

	// Save scorecard results (results.json, score.txt) to GCS
	bucketURL := "gs://ossf-scorecard-results"
	folderPath := fmt.Sprintf("%s/%s", "github.com", repoPath)
	jsonPath := fmt.Sprintf("%s/results.json", folderPath)

	err = writeToBlobStore(ctx, bucketURL, jsonPath, []byte(scorecardOutput.JSONOutput))

	if err != nil {
		return fmt.Errorf(errorWritingBucket.Error()+": %v, %v", err)
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

func lookupPayload(ctx context.Context, payload []byte) (repoPath, repoRef, repoSHA, workflowPath string, err error) {
	// Get most recent Rekor entry uuid.
	rekorClient, err := rekor.NewClient(options.DefaultRekorURL)
	if err != nil {
		return "", "", "", "", fmt.Errorf("error initializing Rekor client: %v", err)
	}

	uuids, err := cosign.FindTLogEntriesByPayload(ctx, rekorClient, payload)
	if err != nil || len(uuids) == 0 {
		return "", "", "", "", fmt.Errorf("error finding tlog entries corresponding to payload: %v", err)
	}
	uuid := uuids[len(uuids)-1] // ignore past entries.

	// Get tlog entry from the UUID.
	entry, err := cosign.GetTlogEntry(ctx, rekorClient, uuid)
	if err != nil {
		return "", "", "", "", fmt.Errorf("error fetching tlog entry: %v", err)
	}

	// Verify tlog entry to make sure it is actually in the log.
	err = cosign.VerifyTLogEntry(ctx, rekorClient, entry)
	if err != nil {
		return "", "", "", "", fmt.Errorf("error verifying tlog entry: %v", err)
	}

	// Extract certificate.
	certs, err := extractCerts(entry)
	if err != nil || len(certs) == 0 {
		return "", "", "", "", fmt.Errorf("error extracting certificate from entry: %v", err)
	}
	if len(certs) > 1 {
		return "", "", "", "", fmt.Errorf("multiple certificates found for the entry: %v", err)
	}

	cert := certs[0]

	// Verify that cert isn't expired.
	if err = cosign.CheckExpiry(cert, time.Unix(*entry.IntegratedTime, 0)); err != nil {
		return "", "", "", "", fmt.Errorf("cosign certificate is expired: %v", err)
	}

	// Get repo reference & path from cert.
	for _, ext := range cert.Extensions {
		// OID source: https://github.com/sigstore/fulcio/blob/96ef49cc7662912ba37d46f738757e8d8d5b5355/docs/oid-info.md#L33
		// TODO: retrieve these by name.
		if ext.Id.String() == "1.3.6.1.4.1.57264.1.6" {
			repoRef = string(ext.Value)
		}
		if ext.Id.String() == "1.3.6.1.4.1.57264.1.5" {
			repoPath = string(ext.Value)
		}
		if ext.Id.String() == "1.3.6.1.4.1.57264.1.3" {
			repoSHA = string(ext.Value)
		}
	}

	if len(repoRef) == 0 || len(repoPath) == 0 {
		return "", "", "", "", fmt.Errorf("error extracting repo ref or path from certificate %v", err)
	}

	// Get workflow job ref from the certificate.
	certURIs := cert.URIs
	if len(certURIs) == 0 {
		return "", "", "", "", errors.New("error: certificate has no URIs")
	}
	workflowRef := certURIs[0].Path
	if len(workflowRef) == 0 {
		return "", "", "", "", errors.New("error: workflow path is empty")
	}

	// Remove repo path from workflow filepath.
	ind := strings.Index(workflowRef, repoPath) + len(repoPath) + 1
	workflowPath = workflowRef[ind:]

	// Remove repo ref tag from workflow filepath.
	workflowPath = workflowPath[:strings.Index(workflowPath, "@")]

	return repoPath, repoRef, repoSHA, workflowPath, nil
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
