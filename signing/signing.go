package signing

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
	"github.com/ossf/scorecard/v2/cron/data"
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

type ScorecardOutput struct {
	SarifOutput string
	JsonOutput  string
}

// verifySignature receives the scorecard analysis payload, looks up its associated tlog entry and
// certificate, and extracts the repository's workflow file to ensure its legitimacy.
func VerifySignature(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading http request body", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Unmarshal body
	var scorecardOutput ScorecardOutput
	err = json.Unmarshal(reqBody, &scorecardOutput)
	if err != nil {
		http.Error(w, "error unmarshalling request body", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Get most recent Rekor entry uuid.
	rekorClient, err := rekor.NewClient(options.DefaultRekorURL)
	if err != nil {
		http.Error(w, "error initializing Rekor client", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// TODO: also process the jsonoutput
	uuids, err := cosign.FindTLogEntriesByPayload(ctx, rekorClient, []byte(scorecardOutput.SarifOutput))
	if err != nil || len(uuids) == 0 {
		http.Error(w, "error fetching tlog entries", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	uuid := uuids[len(uuids)-1] // ignore past entries.

	// Verify tlog entry to make sure it is actually in the log.
	entry, err := verifyTLogEntry(ctx, rekorClient, uuid)
	if err != nil {
		http.Error(w, "error verifying tlog entry", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Extract certificate and get repo reference & path.
	certs, err := extractCerts(entry)
	if err != nil || len(certs) == 0 {
		http.Error(w, "error extracting certificate from entry", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if len(certs) > 1 {
		http.Error(w, "multiple certificates found for the entry", http.StatusInternalServerError)
		log.Println("error: multiple certificates found for the entry")
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
		log.Println("error: repoRef or repoPath are empty")
		return
	}

	// Split repo path into owner and repo name.
	ownerName := repoPath[:strings.Index(repoPath, "/")]
	repoName := repoPath[strings.Index(repoPath, "/")+1:]

	// Verify that the repository and branch of the cert and request are equal.
	reqRepo := r.Header["X-Repository"]
	reqBranch := r.Header["X-Branch"]
	if len(reqRepo) == 0 || len(reqBranch) == 0 || reqRepo[0] != repoPath || reqBranch[0] != repoRef {
		http.Error(w, "repository and branch of cert doesn't match that of request", http.StatusInternalServerError)
		return
	}

	// Make github client.
	client := github.NewClient(nil)

	// Get all workflows in the repository.
	workflows, _, err := client.Actions.ListWorkflows(ctx, ownerName, repoName, nil)
	if err != nil {
		http.Error(w, "error listing workflows from repo", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	wkflws := workflows.Workflows
	wkflwPath := ""
	for _, wkflw := range wkflws {
		if *wkflw.Name == "Scorecards supply-chain security" {
			wkflwPath = *wkflw.Path
			break
		}
	}
	if wkflwPath == "" {
		http.Error(w, "error finding scorecard workflow in repository", http.StatusInternalServerError)
		log.Println("error finding scorecard workflow in repository")
		return
	}

	// Get workflow file from repo reference.
	// TODO: use GITHUB_TOKEN from workflow to make the api call.
	opts := &github.RepositoryContentGetOptions{Ref: repoRef}
	contents, _, _, err := client.Repositories.GetContents(ctx, ownerName, repoName, wkflwPath, opts)
	if err != nil {
		http.Error(w, "error downloading workflow contents from repo", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	workflowContent, err := contents.GetContent()
	if err != nil {
		http.Error(w, "error decoding workflow contents", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Verify scorecard workflow.
	err = verifyScorecardWorkflow(workflowContent)
	if err != nil {
		http.Error(w, "workflow could not be verified", http.StatusNotAcceptable)
		log.Println(err)
		return
	}

	// Save scorecard results (results.sarif, results.json, score.txt) to GCS
	bucketURL := "gs://ossf-scorecard-results"
	folderPath := fmt.Sprintf("%s/%s", "github.com", repoPath)
	sarifPath := fmt.Sprintf("%s/results.sarif", folderPath)
	jsonPath := fmt.Sprintf("%s/results.json", folderPath)

	err1 := data.WriteToBlobStore(ctx, bucketURL, sarifPath, []byte(scorecardOutput.SarifOutput))
	err2 := data.WriteToBlobStore(ctx, bucketURL, jsonPath, []byte(scorecardOutput.JsonOutput))
	if err1 != nil || err2 != nil {
		http.Error(w, "error writing to GCS bucket", http.StatusNotAcceptable)
		log.Println(err1, err2)
		return
	}

	// Write response.
	w.Write([]byte(fmt.Sprintf("Successfully verified and uploaded scorecard results for repo %s on branch %s", []byte(repoName), []byte(repoRef))))
	w.WriteHeader(http.StatusOK)
}

func verifyScorecardWorkflow(workflowContent string) error {
	// Verify workflow contents using actionlint.
	workflow, lintErrs := actionlint.Parse([]byte(workflowContent))
	if lintErrs != nil {
		return fmt.Errorf("actionlint errors parsing workflow: %v", lintErrs)
	}

	// Extract main job
	jobs := workflow.Jobs
	if len(jobs) != 1 {
		return errors.New("number of jobs isn't 1")
	}
	analysisJob := jobs["analysis"]
	if analysisJob == nil {
		return errors.New("workflow doens't have analysis job")
	}

	// Verify that there is no container or services.
	if analysisJob.Container != nil || len(analysisJob.Services) > 0 {
		return errors.New("workflow contains container or service")
	}

	// Verify that the workflow runs on ubuntu-latest and nothing else.
	if analysisJob.RunsOn != nil {
		labels := analysisJob.RunsOn.Labels
		if len(labels) == 0 || len(labels) > 1 || labels[0].Value != "ubuntu-latest" {
			return errors.New("workflow doesn't run solely on ubuntu-latest")
		}
	} else {
		return errors.New("no RunsOn found in workflow")
	}

	// Verify that there are no env vars set.
	if analysisJob.Env != nil {
		return errors.New("workflow contains env vars")
	}

	// Verify that there are no defaults set.
	if analysisJob.Defaults != nil {
		return errors.New("workflow has defaults set")
	}

	// Get steps in job.
	steps := analysisJob.Steps

	if steps == nil || len(steps) > 4 {
		return fmt.Errorf("workflow has an invalid number of steps: %d", len(steps))
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
			"github/codeql-action/upload-sarif",
			"rohankh532/scorecard-action": // TODO: remove later, for debugging
		default:
			return fmt.Errorf("workflow has invalid step name: %s", stepName)
		}
	}

	return nil
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
