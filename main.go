// Copyright 2021 Security Scorecard Authors
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

// Package main implements the scorecard.dev webapp.
package main

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/google/go-github/v42/github"
	"github.com/google/trillian/merkle/logverifier"
	"github.com/google/trillian/merkle/rfc6962"
	"github.com/rhysd/actionlint"
	"github.com/sigstore/cosign/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/cmd/cosign/cli/rekor"
	"github.com/sigstore/cosign/pkg/cosign"
	"github.com/sigstore/cosign/pkg/cosign/bundle"
	"github.com/sigstore/cosign/pkg/cosign/tuf"
	"github.com/sigstore/rekor/pkg/generated/client"
	"github.com/sigstore/rekor/pkg/generated/client/entries"
	"github.com/sigstore/rekor/pkg/generated/models"
	"github.com/sigstore/rekor/pkg/types"
	hashedrekord "github.com/sigstore/rekor/pkg/types/hashedrekord/v0.0.1"
	rekord "github.com/sigstore/rekor/pkg/types/rekord/v0.0.1"
	"github.com/sigstore/sigstore/pkg/cryptoutils"
)

// verifySignature receives the scorecard analysis payload, looks up its associated tlog entry and
// certificate, and extracts the repository's workflow file to ensure its legitimacy.
func verifySignature(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading http request body", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Get most recent Rekor entry uuid.
	rekorClient, err := rekor.NewClient(options.DefaultRekorURL)
	if err != nil {
		http.Error(w, "error initializing Rekor client", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	uuids, err := cosign.FindTLogEntriesByPayload(ctx, rekorClient, reqBody)
	if err != nil || len(uuids) == 0 {
		http.Error(w, "error fetching tlog entries", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	uuid := uuids[len(uuids)-1] // ignore past entries.

	// Verify tlog entry to make sure it is actually in the log.
	entry, err := verifyTLogEntry(ctx, rekorClient, uuid)
	if err != nil {
		http.Error(w, "error verifying tlog entry", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Extract certificate and get repo reference & path.
	certs, err := extractCerts(entry)
	if err != nil || len(certs) == 0 {
		http.Error(w, "error extracting certificate from entry", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	if len(certs) > 1 {
		http.Error(w, "multiple certificates found for the entry", http.StatusInternalServerError)
		return
	}
	cert := certs[0]
	var repoRef string
	var repoPath string
	for _, ext := range cert.Extensions {
		// OID source: https://github.com/sigstore/fulcio/blob/96ef49cc7662912ba37d46f738757e8d8d5b5355/docs/oid-info.md#L33
		// TODO: retrieve these by name.
		if ext.Id.String() == "1.3.6.1.4.1.57264.1.6" {
			repoRef = string(ext.Value)
		}
		if ext.Id.String() == "1.3.6.1.4.1.57264.1.5" {
			repoPath = string(ext.Value)
		}
	}

	if len(repoRef) == 0 || len(repoPath) == 0 {
		http.Error(w, "error extracting repo ref or path from certificate", http.StatusInternalServerError)
		return
	}

	// Split repo path into owner and repo name.
	ownerName := repoPath[:strings.Index(repoPath, "/")]
	repoName := repoPath[strings.Index(repoPath, "/")+1:]

	// Verify that the repository and branch of the cert and request are equal.
	reqRepo := r.Header["Repository"]
	reqBranch := r.Header["Branch"]
	if len(reqRepo) == 0 || len(reqBranch) == 0 || reqRepo[0] != repoPath || reqBranch[0] != repoRef {
		http.Error(w, "repository and branch of cert doesn't match that of request", http.StatusInternalServerError)
		return
	}

	// Get workflow file from repo reference.
	// TODO: use GITHUB_TOKEN from workflow to make the api call.
	client := github.NewClient(nil)
	opts := &github.RepositoryContentGetOptions{Ref: repoRef}
	contents, _, _, err := client.Repositories.GetContents(ctx, ownerName, repoName, ".github/workflows/scorecards.yml", opts)
	if err != nil {
		http.Error(w, "error downloading workflow contents from repo", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	workflowContent, err := contents.GetContent()
	if err != nil {
		http.Error(w, "error decoding workflow contents", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Verify scorecard workflow.
	verified := verifyScorecardWorkflow(workflowContent)
	fmt.Println(verified)

	// Save blob to GCS

	// Next: badging...
}

func getScore(w http.ResponseWriter, r *http.Request) {
	scorecardData := struct{ Score int }{Score: 1}
	jData, err := json.Marshal(scorecardData)
	if err != nil {
		http.Error(w, "error marshalling struct", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func main() {
	http.HandleFunc("/", httpHandler)
	fmt.Printf("Starting HTTP server on port 8080 ...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	http.HandleFunc("/score", getScore) // TODO: organize in handler.

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Hello world!" +
			" This site is still under construction." +
			" Please check back again later.")); err != nil {
			log.Printf("error during Write: %v", err)
		}
	case http.MethodPost:
		verifySignature(w, r)
	default:
		http.Error(w, "only GET method is allowed", http.StatusMethodNotAllowed)
	}
}

func verifyScorecardWorkflow(workflowContent string) bool {
	// Verify workflow contents using actionlint.
	workflow, lintErrs := actionlint.Parse([]byte(workflowContent))
	if lintErrs != nil {
		return false
	}

	// Extract main job
	jobs := workflow.Jobs
	if len(jobs) != 1 {
		return false
	}
	analysisJob := jobs["analysis"]
	if analysisJob == nil {
		return false
	}

	// Verify that there is no container or services.
	if analysisJob.Container != nil || len(analysisJob.Services) > 0 {
		return false
	}

	// Verify that the workflow runs on ubuntu-latest and nothing else.
	if analysisJob.RunsOn != nil {
		labels := analysisJob.RunsOn.Labels
		if len(labels) == 0 || len(labels) > 1 || labels[0].Value != "ubuntu-latest" {
			return false
		}
	} else {
		return false
	}

	// Verify that there are no env vars set.
	if analysisJob.Env != nil {
		return false
	}

	// Verify that there are no defaults set.
	if analysisJob.Defaults != nil {
		return false
	}

	// Get steps in job.
	steps := analysisJob.Steps

	if steps == nil || len(steps) > 4 {
		return false
	}

	// Verify that steps are valid (checkout, scorecard-action, upload-artifact, upload-sarif).
	for _, step := range steps {
		stepName := step.Exec.(*actionlint.ExecAction).Uses.Value
		stepName = stepName[:strings.Index(stepName, "@")] // get rid of commit sha.

		switch stepName {
		case
			"actions/checkout",
			"ossf/scorecard-action",
			"actions/upload-artifact",
			"github/codeql-action/upload-sarif":
		default:
			return false
		}
	}

	return true
}

// Source: https://github.com/sigstore/cosign/blob/18d2ce0b458018951f7356db911467a427a8dffe/pkg/cosign/tlog.go#L247
func verifyTLogEntry(ctx context.Context, rekorClient *client.Rekor, uuid string) (*models.LogEntryAnon, error) {
	params := entries.NewGetLogEntryByUUIDParamsWithContext(ctx)
	params.EntryUUID = uuid

	lep, err := rekorClient.Entries.GetLogEntryByUUID(params)
	if err != nil {
		return nil, err
	}

	if len(lep.Payload) != 1 {
		return nil, errors.New("UUID value can not be extracted")
	}
	e := lep.Payload[params.EntryUUID]
	if e.Verification == nil || e.Verification.InclusionProof == nil {
		return nil, errors.New("inclusion proof not provided")
	}

	hashes := [][]byte{}
	for _, h := range e.Verification.InclusionProof.Hashes {
		hb, _ := hex.DecodeString(h)
		hashes = append(hashes, hb)
	}

	rootHash, err := hex.DecodeString(*e.Verification.InclusionProof.RootHash)
	if err != nil {
		return nil, errors.New("error decoding root hash string")
	}
	leafHash, err := hex.DecodeString(params.EntryUUID)
	if err != nil {
		return nil, errors.New("error decoding leaf hash string")
	}

	v := logverifier.New(rfc6962.DefaultHasher)
	if err := v.VerifyInclusionProof(*e.Verification.InclusionProof.LogIndex, *e.Verification.InclusionProof.TreeSize, hashes, rootHash, leafHash); err != nil {
		return nil, err
	}

	// Verify rekor's signature over the SET.
	payload := bundle.RekorPayload{
		Body:           e.Body,
		IntegratedTime: *e.IntegratedTime,
		LogIndex:       *e.LogIndex,
		LogID:          *e.LogID,
	}

	rekorPubKeys, err := cosign.GetRekorPubs(ctx)
	if err != nil {
		return nil, err
	}
	var entryVerError error
	for _, pubKey := range rekorPubKeys {
		entryVerError = cosign.VerifySET(payload, []byte(e.Verification.SignedEntryTimestamp), pubKey.PubKey)
		// Return once the SET is verified successfully.
		if entryVerError == nil {
			if pubKey.Status != tuf.Active {
				fmt.Fprintf(os.Stderr, "**Info** Successfully verified Rekor entry using an expired verification key\n")
			}
			return &e, nil
		}
	}
	return nil, err
}

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
