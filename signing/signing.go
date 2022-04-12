package signing

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
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/google/go-github/v42/github"
	"github.com/gorilla/mux"
	"github.com/ossf/scorecard/v2/cron/data"
	"github.com/rhysd/actionlint"
	"github.com/sigstore/cosign/cmd/cosign/cli/options"
	"github.com/sigstore/cosign/cmd/cosign/cli/rekor"
	"github.com/sigstore/cosign/pkg/cosign"
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

	// Unmarshal request body.
	var scorecardOutput ScorecardOutput
	err = json.Unmarshal(reqBody, &scorecardOutput)
	if err != nil {
		http.Error(w, "error unmarshalling request body", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Lookup results payload to get the repo info from the corresponding entry & cert.
	sarifRepoPath, sarifRepoRef, sarifRepoSHA, workflowPath, err1 := lookupPayload(ctx, []byte(scorecardOutput.SarifOutput))
	jsonRepoPath, jsonRepoRef, jsonRepoSHA, _, err2 := lookupPayload(ctx, []byte(scorecardOutput.JsonOutput))
	if err1 != nil || err2 != nil {
		http.Error(w, "error looking up payloads and processing/verifying entries & certs", http.StatusInternalServerError)
		log.Printf("sarif: %v, json: %v", err1, err2)
		return
	}

	// Verify that the sarif results and json results are from the same repo, ref, and SHA.
	if sarifRepoPath != jsonRepoPath || sarifRepoRef != jsonRepoRef || sarifRepoSHA != jsonRepoSHA {
		http.Error(w, "sarif and json results correspond to different repos/refs/SHAs.", http.StatusInternalServerError)
		return
	}

	repoPath := sarifRepoPath
	repoRef := sarifRepoRef
	repoSHA := sarifRepoSHA

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

	// Get the corresponding GitHub repository.
	// TODO: use GITHUB_TOKEN from workflow to make the api call.
	client := github.NewClient(nil)
	repo, resp, err := client.Repositories.Get(ctx, ownerName, repoName)
	if err != nil {
		http.Error(w, "error getting repository", http.StatusInternalServerError)
		log.Println(err, resp)
		return
	}

	// Verify that the branch from the results files is the repo's default branch.
	defaultBranch := "refs/heads/" + repo.GetDefaultBranch()
	if defaultBranch != repoRef {
		http.Error(w, "branch of cert isn't the repo's default branch", http.StatusInternalServerError)
		log.Println(defaultBranch, repoRef)
		return
	}

	// Get workflow file from cert commit SHA.
	opts := &github.RepositoryContentGetOptions{Ref: repoSHA}
	contents, _, _, err := client.Repositories.GetContents(ctx, ownerName, repoName, workflowPath, opts)
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
		http.Error(w, "workflow could not be verified", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Save scorecard results (results.sarif, results.json, score.txt) to GCS
	bucketURL := "gs://ossf-scorecard-results"
	folderPath := fmt.Sprintf("%s/%s", "github.com", repoPath)
	sarifPath := fmt.Sprintf("%s/results.sarif", folderPath)
	jsonPath := fmt.Sprintf("%s/results.json", folderPath)

	err1 = data.WriteToBlobStore(ctx, bucketURL, sarifPath, []byte(scorecardOutput.SarifOutput))
	err2 = data.WriteToBlobStore(ctx, bucketURL, jsonPath, []byte(scorecardOutput.JsonOutput))
	if err1 != nil || err2 != nil {
		http.Error(w, "error writing to GCS bucket", http.StatusUnauthorized)
		log.Println(err1, err2)
		return
	}

	// Write response.
	w.Write([]byte(fmt.Sprintf("Successfully verified and uploaded scorecard results for repo %s on branch %s", []byte(repoName), []byte(repoRef))))
	w.WriteHeader(http.StatusOK)
}

func lookupPayload(ctx context.Context, payload []byte) (repoPath, repoRef, repoSHA, workflowPath string, err error) {
	// Get most recent Rekor entry uuid.
	rekorClient, err := rekor.NewClient(options.DefaultRekorURL)
	if err != nil {
		return "", "", "", "", fmt.Errorf("error initializing Rekor client: %v", err)
	}

	uuids, err := cosign.FindTLogEntriesByPayload(ctx, rekorClient, []byte(payload))
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

func verifyScorecardWorkflow(workflowContent string) error {
	// Verify workflow contents using actionlint.
	workflow, lintErrs := actionlint.Parse([]byte(workflowContent))
	if lintErrs != nil {
		return fmt.Errorf("actionlint errors parsing workflow: %v", lintErrs)
	}

	// Verify that there are no global env vars or defaults.
	if workflow.Env != nil || workflow.Defaults != nil {
		return errors.New("workflow contains global env vars or defaults")
	}

	// Verify that the only global permission set is read-all.
	if workflow.Permissions.All.Value != "read-all" {
		return errors.New("workflow permission isn't read-all")
	}

	// Extract main (analysis) job
	jobsMap := workflow.Jobs
	jobs := reflect.ValueOf(jobsMap).MapKeys()

	if len(jobs) != 1 {
		return errors.New("number of jobs isn't 1")
	}
	job := jobsMap[jobs[0].String()]
	if job == nil {
		return errors.New("workflow doens't have a main job")
	}

	// Verify that there is no job container or services.
	if job.Container != nil || len(job.Services) > 0 {
		return errors.New("workflow contains container or service")
	}

	// Verify that the workflow runs on ubuntu-latest and nothing else.
	if job.RunsOn == nil {
		return errors.New("no RunsOn found in workflow")
	} else {
		labels := job.RunsOn.Labels
		if len(labels) == 0 || len(labels) > 1 || labels[0].Value != "ubuntu-latest" {
			return errors.New("workflow doesn't run solely on ubuntu-latest")
		}
	}

	// Verify that there are no job env vars set.
	if job.Env != nil {
		return errors.New("workflow contains env vars")
	}

	// Verify that there are no job defaults set.
	if job.Defaults != nil {
		return errors.New("workflow has defaults set")
	}

	// Get steps in job.
	steps := job.Steps

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

var errorPullingBucket = errors.New("error pulling from GCS bucket")
var errorVerifyingFilepath = errors.New("error verifying filepath format")

func GetResults(w http.ResponseWriter, r *http.Request) {
	host := mux.Vars(r)["host"]
	orgName := mux.Vars(r)["orgName"]
	repoName := mux.Vars(r)["repoName"]
	results, err := getResults(host, orgName, repoName)

	if err == errorVerifyingFilepath {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("err: %v", err)
		return
	}
	if err == errorPullingBucket {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Printf("err: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(results)
}

func getResults(host, orgName, repoName string) (results []byte, err error) {
	// Get params to build GCS filepath.
	ctx := context.Background()
	bucketURL := "gs://ossf-scorecard-results"
	resultsFile := filepath.Join(host, orgName, repoName, "results.json")

	// Sanitize input and log query.
	resultsFile = filepath.Clean(resultsFile)
	matched, err := filepath.Match("*/*/*/results.json", resultsFile)
	if err != nil || !matched {
		return nil, errorVerifyingFilepath
	}

	if len(resultsFile) >= 256 {
		return nil, fmt.Errorf("filepath (%v) is greater than the Linux maximum of 256", resultsFile[:256])
	}

	resultsFileEscaped := strings.Replace(resultsFile, "\n", "", -1)
	resultsFileEscaped = strings.Replace(resultsFileEscaped, "\r", "", -1)
	log.Printf("Querying GCS bucket for: %s", resultsFileEscaped)

	// Query GCS bucket.
	results, err = data.GetBlobContent(ctx, bucketURL, resultsFileEscaped)
	if err != nil {
		return nil, errorPullingBucket
	}
	return results, nil
}
