package signing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetResults(t *testing.T) {
	r, _ := http.NewRequest("GET", "/projects", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"host":     "github.com",
		"orgName":  "rohankh532",
		"repoName": "scorecard-OIDC-test",
	}

	r = mux.SetURLVars(r, vars)

	// Verify that corresponding results to this repo are found.
	GetResults(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

}
