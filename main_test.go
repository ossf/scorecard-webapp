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
	sarifpayload, _ := ioutil.ReadFile("testdata/validSig-invalidWkflw.sarif")
	jsonpayload, _ := ioutil.ReadFile("testdata/scorecard-results.json")
	payload := ScorecardOutput{SarifOutput: string(sarifpayload), JsonOutput: string(jsonpayload)}
	payloadbytes, err := json.Marshal(payload)
	assert.Equal(t, err, nil)

	r, _ := http.NewRequest("POST", "/verify", bytes.NewBuffer(payloadbytes))
	r.Header = http.Header{"X-Repository": []string{"rohankh532/scorecard-OIDC-test"}, "X-Branch": []string{"refs/heads/main"}}
	w := httptest.NewRecorder()

	verifySignature(w, r)

	// Only the invalid workflow file error code is allowed
	assert.Equal(t, http.StatusNotAcceptable, w.Code)
}

func TestVerifySignatureInvalidRepo(t *testing.T) {
	sarifpayload, _ := ioutil.ReadFile("testdata/validSig-invalidWkflw.sarif")
	jsonpayload, _ := ioutil.ReadFile("testdata/scorecard-results.json")
	payload := ScorecardOutput{SarifOutput: string(sarifpayload), JsonOutput: string(jsonpayload)}
	payloadbytes, err := json.Marshal(payload)
	assert.Equal(t, err, nil)

	r, _ := http.NewRequest("POST", "/verify", bytes.NewBuffer(payloadbytes))
	r.Header = http.Header{"X-Repository": []string{"rohankh532/invalid-repo"}, "X-Branch": []string{"refs/heads/main"}}
	w := httptest.NewRecorder()

	verifySignature(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestVerifyValidWorkflow(t *testing.T) {
	workflowContent, _ := ioutil.ReadFile("testdata/workflow-valid.yml")
	err := verifyScorecardWorkflow(string(workflowContent))
	assert.Equal(t, err, nil)
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
		err := verifyScorecardWorkflow(string(workflowContent))
		assert.NotEqual(t, err, nil)
	}
}

func TestGetScore(t *testing.T) {
	r, _ := http.NewRequest("GET", "/projects/github.com/org/repo", nil)
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
