package main

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
	// Using a scorecard results file from a previous run that was successfully signed.
	// This isn't expected to pass verifyScorecardWorkflow because the workflow it was
	// generated from contains extra steps to call cosign.
	payload, _ := ioutil.ReadFile("testdata/validSig-invalidWkflw.sarif")

	r, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(payload))
	r.Header = http.Header{"Repository": []string{"rohankh532/scorecard-OIDC-test"}, "Branch": []string{"refs/heads/main"}}
	w := httptest.NewRecorder()

	verifySignature(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVerifySignatureInvalidRepo(t *testing.T) {
	payload, _ := ioutil.ReadFile("testdata/validSig-invalidWkflw.sarif")

	r, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(payload))
	r.Header = http.Header{"Repository": []string{"rohankh532/invalid-repo"}, "Branch": []string{"refs/heads/main"}}
	w := httptest.NewRecorder()

	verifySignature(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVerifyValidWorkflow(t *testing.T) {
	workflowContent, _ := ioutil.ReadFile("testdata/workflow-valid.yml")
	res := verifyScorecardWorkflow(string(workflowContent))
	assert.Equal(t, res, true)
}

func TestVerifyInvalidWorkflows(t *testing.T) {
	workflowFiles := []string{
		"testdata/workflow-invalid-formatting.yml",
		"testdata/workflow-invalid-jobs.yml",
		"testdata/workflow-invalid-analysisjob.yml",
		"testdata/workflow-invalid-container.yml",
		"testdata/workflow-invalid-services.yml",
		"testdata/workflow-invalid-runson.yml",
		"testdata/workflow-invalid-envvars.yml",
		"testdata/workflow-invalid-manysteps.yml",
		"testdata/workflow-invalid-diffsteps.yml",
		"testdata/workflow-invalid-defaults.yml",
	}

	for _, workflowFile := range workflowFiles {
		workflowContent, _ := ioutil.ReadFile(workflowFile)
		res := verifyScorecardWorkflow(string(workflowContent))
		assert.Equal(t, res, false)
	}
}

func TestGetScore(t *testing.T) {
	r, _ := http.NewRequest("GET", "/score", nil)
	w := httptest.NewRecorder()

	getScore(w, r)

	scorecardData := struct{ Score int }{}

	res := w.Body.Bytes()
	err := json.Unmarshal(res, &scorecardData)
	assert.Equal(t, err, nil)

	// Endpoint currently always returns 1.
	expectedData := struct{ Score int }{Score: 1}
	assert.Equal(t, scorecardData, expectedData)
}
