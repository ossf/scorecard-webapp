package signing

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestVerifySignature(t *testing.T) {
	// Should pass entry, cert, and workflow verification but fail GCS upload.
	sarifpayload, _ := ioutil.ReadFile("../testdata/results/results.sarif")
	jsonpayload, _ := ioutil.ReadFile("../testdata/results/results.json")
	payload := ScorecardOutput{SarifOutput: string(sarifpayload), JsonOutput: string(jsonpayload)}
	payloadbytes, err := json.Marshal(payload)
	assert.Equal(t, err, nil)

	r, _ := http.NewRequest("POST", "/verify", bytes.NewBuffer(payloadbytes))
	r.Header = http.Header{"X-Repository": []string{"rohankh532/scorecard-OIDC-test"}, "X-Branch": []string{"refs/heads/main"}}
	w := httptest.NewRecorder()

	VerifySignatureHandler(w, r)

	// Only the GCS upload error code is allowed
	err_msg := strings.TrimSuffix(w.Body.String(), "\n")
	assert.True(t, strings.HasPrefix(err_msg, errorWritingBucket.Error()))
}

func TestVerifySignatureInvalidRepo(t *testing.T) {
	sarifpayload, _ := ioutil.ReadFile("../testdata/results/results.sarif")
	jsonpayload, _ := ioutil.ReadFile("../testdata/results/results.json")
	payload := ScorecardOutput{SarifOutput: string(sarifpayload), JsonOutput: string(jsonpayload)}
	payloadbytes, err := json.Marshal(payload)
	assert.Equal(t, err, nil)

	r, _ := http.NewRequest("POST", "/verify", bytes.NewBuffer(payloadbytes))
	r.Header = http.Header{"X-Repository": []string{"rohankh532/invalid-repo"}, "X-Branch": []string{"refs/heads/main"}}
	w := httptest.NewRecorder()

	VerifySignatureHandler(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVerifyValidWorkflow(t *testing.T) {
	workflowContent, _ := ioutil.ReadFile("../testdata/workflow-valid.yml")
	err := verifyScorecardWorkflow(string(workflowContent))
	assert.Equal(t, err, nil)
}

func TestVerifyInvalidWorkflows(t *testing.T) {
	workflowFiles := []string{
		"../testdata/workflow-invalid-formatting.yml",
		"../testdata/workflow-invalid-container.yml",
		"../testdata/workflow-invalid-services.yml",
		"../testdata/workflow-invalid-runson.yml",
		"../testdata/workflow-invalid-envvars.yml",
		"../testdata/workflow-invalid-manysteps.yml",
		"../testdata/workflow-invalid-diffsteps.yml",
		"../testdata/workflow-invalid-defaults.yml",
		"../testdata/workflow-invalid-global-perm.yml",
		"../testdata/workflow-invalid-global-env.yml",
		"../testdata/workflow-invalid-global-defaults.yml",
	}

	for _, workflowFile := range workflowFiles {
		workflowContent, _ := ioutil.ReadFile(workflowFile)
		err := verifyScorecardWorkflow(string(workflowContent))
		assert.NotEqual(t, err, nil, workflowFile)
	}
}

func TestValidGetResults(t *testing.T) {
	r, _ := http.NewRequest("GET", "/projects", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"host":     "github.com",
		"orgName":  "rohankh532",
		"repoName": "scorecard-OIDC-test",
	}

	r = mux.SetURLVars(r, vars)

	// Should fail connecting to GCS but pass filepath screening.
	GetResults(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)

}

func TestInvalidGetResults(t *testing.T) {
	r, _ := http.NewRequest("GET", "/projects", nil)
	w := httptest.NewRecorder()

	varMapsInvalidURL := [...]map[string]string{
		{
			"host":     "../github.com",
			"orgName":  "rohankh532",
			"repoName": "scorecard-OIDC-test",
		},
		{
			"host":     "malicious/github.com",
			"orgName":  "rohankh532",
			"repoName": "scorecard-OIDC-test",
		},
		{
			"host":     "malicious\\/github.com",
			"orgName":  "rohankh532",
			"repoName": "scorecard-OIDC-test",
		},
	}

	for _, varMap := range varMapsInvalidURL {
		// Verify that the url sanitization fails.
		r = mux.SetURLVars(r, varMap)
		GetResults(w, r)

		err_msg := strings.TrimSuffix(w.Body.String(), "\n")

		assert.Equal(t, errorVerifyingFilepath.Error(), err_msg)
		w = httptest.NewRecorder() // Reset the recorder.
	}

	varMapsInvalidRepo := [...]map[string]string{
		{
			"host":     "github.com",
			"orgName":  "invalid",
			"repoName": "scorecard-OIDC-test",
		},
		{
			"host":     "github.com",
			"orgName":  "rohankh532",
			"repoName": "invalid",
		},
		{
			"host":     "invalid",
			"orgName":  "rohankh532",
			"repoName": "scorecard-OIDC-test",
		},
	}

	for _, varMap := range varMapsInvalidRepo {
		// Verify that the file doesn't exist in the bucket.
		r = mux.SetURLVars(r, varMap)
		GetResults(w, r)

		err_msg := strings.TrimSuffix(w.Body.String(), "\n")

		assert.Equal(t, errorPullingBucket.Error(), err_msg)
		w = httptest.NewRecorder() // Reset the recorder.
	}

}
