package signing

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, http.StatusInternalServerError, w.Code)

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
