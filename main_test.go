package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequest(t *testing.T) {
	// Using a scorecard results file from a previous run that was successfully signed.
	// This isn't expected to pass verifyScorecardWorkflow because the workflow it was
	// generated from contains extra steps to call cosign.
	payload, _ := ioutil.ReadFile("testdata/results.sarif")
	r, _ := http.NewRequest("POST", "/projects", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()

	verifySignature(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVerifyValidWorkflow(t *testing.T) {
	workflowContent, _ := ioutil.ReadFile("testdata/workflow-valid.yml")
	res := verifyScorecardWorkflow(string(workflowContent))
	assert.Equal(t, res, true)
}
