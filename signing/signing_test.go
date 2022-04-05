package signing

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

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

	VerifySignature(w, r)

	// Only the GCS upload error code is allowed
	assert.Equal(t, http.StatusUnauthorized, w.Code)
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

	VerifySignature(w, r)

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
		"../testdata/workflow-invalid-jobs.yml",
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
		assert.NotEqual(t, err, nil)
	}
}
