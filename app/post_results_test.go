// Copyright 2022 Security Scorecard Authors
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

package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostResultsHandler(t *testing.T) {
	t.Parallel()
	// Should pass entry, cert, and workflow verification but fail GCS upload.
	jsonpayload, _ := ioutil.ReadFile("testdata/results/results.json")
	payload := ScorecardOutput{JSONOutput: string(jsonpayload)}
	payloadbytes, err := json.Marshal(payload)
	assert.Equal(t, err, nil)

	r, _ := http.NewRequestWithContext(context.Background(), "POST", "/verify", bytes.NewBuffer(payloadbytes))
	r.Header = http.Header{
		"X-Repository": []string{"rohankh532/scorecard-OIDC-test"},
		"X-Branch":     []string{"refs/heads/main"},
	}
	w := httptest.NewRecorder()

	PostResultsHandler(w, r)

	// Only the GCS upload error code is allowed
	errMsg := strings.TrimSuffix(w.Body.String(), "\n")
	assert.True(t, strings.HasPrefix(errMsg, errorWritingBucket.Error()))
}

func TestPostResultsInvalidRepo(t *testing.T) {
	t.Parallel()
	jsonpayload, _ := ioutil.ReadFile("testdata/results/results.json")
	payload := ScorecardOutput{JSONOutput: string(jsonpayload)}
	payloadbytes, err := json.Marshal(payload)
	assert.Equal(t, err, nil)

	r, _ := http.NewRequestWithContext(context.Background(), "POST", "/verify", bytes.NewBuffer(payloadbytes))
	r.Header = http.Header{"X-Repository": []string{"rohankh532/invalid-repo"}, "X-Branch": []string{"refs/heads/main"}}
	w := httptest.NewRecorder()

	PostResultsHandler(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
