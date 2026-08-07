// Copyright 2022 OpenSSF Scorecard Authors
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

package server

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

type allowListVerifier struct {
	allowed map[string]bool
}

func (a *allowListVerifier) contains(c commit) (bool, error) {
	return a.allowed[c.hash], nil
}

var allowCommitVerifier = &allowListVerifier{
	allowed: map[string]bool{
		"dd2c410b088af7c0dc8046f3ac9a8f4148492a95": true,
		"ec3a7ce113134d7a93b817d10a8272cb61118579": true,
		"c8416b0b2bf627c349ca92fc8e3de51a64b005cf": true,
		"82c141cc518b40d92cc801eee768e7aafc9c2fa2": true,
		"5f532563584d71fdef14ee64d17bafb34f751ce5": true,
	},
}

func TestVerifyValidWorkflows(t *testing.T) {
	t.Parallel()
	workflowFiles := []string{
		"testdata/workflow-valid.yml",
		"testdata/workflow-valid-noglobalperm.yml",
		"testdata/workflow-valid-e2e.yml",
		"testdata/workflow-valid-e2e-ghcr.yml",
		"testdata/workflow-valid-tagged-action.yml",
	}

	for _, workflowFile := range workflowFiles {
		workflowContent, _ := os.ReadFile(workflowFile)
		err := verifyScorecardWorkflow(string(workflowContent), allowCommitVerifier)
		if err != nil {
			t.Errorf("expected - %v, got - %v", nil, err)
		}
	}
}

func TestVerifyInvalidWorkflows(t *testing.T) {
	t.Parallel()
	workflowFiles := []string{
		"testdata/workflow-invalid-formatting.yml",
		"testdata/workflow-invalid-container.yml",
		"testdata/workflow-invalid-services.yml",
		"testdata/workflow-invalid-runson.yml",
		"testdata/workflow-invalid-envvars.yml",
		"testdata/workflow-invalid-diffsteps.yml",
		"testdata/workflow-invalid-defaults.yml",
		"testdata/workflow-invalid-global-perm.yml",
		"testdata/workflow-invalid-global-env.yml",
		"testdata/workflow-invalid-global-defaults.yml",
		"testdata/workflow-invalid-otherjob.yml",
		"testdata/workflow-invalid-global-idtoken.yml",
		"testdata/workflow-invalid-empty.yml",
		"testdata/workflow-invalid-missing-scorecard.yml",
		"testdata/workflow-invalid-missing-runson.yml",
		"testdata/workflow-invalid-multiple-labels.yml",
		"testdata/workflow-invalid-nil-steps.yml",
		"testdata/workflow-invalid-execaction.yml",
		"testdata/workflow-invalid-imposter-commit.yml",
	}

	for _, workflowFile := range workflowFiles {
		workflowContent, _ := os.ReadFile(workflowFile)
		err := verifyScorecardWorkflow(string(workflowContent), allowCommitVerifier)
		assert.NotEqual(t, err, nil, workflowFile)
	}
}

func TestVerifyRunnerLabels(t *testing.T) {
	t.Parallel()
	// Start from a known-valid workflow and swap only the runs-on label so each
	// case exercises the runner-label check in isolation.
	base, err := os.ReadFile("testdata/workflow-valid.yml")
	if err != nil {
		t.Fatal(err)
	}
	const origLabel = "runs-on: ubuntu-latest"
	if !strings.Contains(string(base), origLabel) {
		t.Fatalf("fixture no longer contains %q; update this test", origLabel)
	}

	tests := []struct {
		label   string
		wantErr bool
	}{
		{label: "ubuntu-latest", wantErr: false},
		{label: "ubuntu-latest-arm", wantErr: false},
		{label: "ubuntu-22.04", wantErr: false},
		{label: "ubuntu-24.04", wantErr: false},
		{label: "ubuntu-26.04", wantErr: false}, // ossf/scorecard-action#1684
		{label: "ubuntu-22.04-arm", wantErr: false},
		{label: "ubuntu-24.04-arm", wantErr: false},
		{label: "ubuntu-20.04", wantErr: true}, // end-of-life, below minimum
		{label: "ubuntu-18.04", wantErr: true}, // end-of-life, below minimum
		{label: "macos-13", wantErr: true},
		{label: "self-hosted", wantErr: true},
		{label: "ubuntu-", wantErr: true},
		{label: "Ubuntu-latest", wantErr: true}, // labels are case-sensitive
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			t.Parallel()
			content := strings.Replace(string(base), origLabel, "runs-on: "+tt.label, 1)
			err := verifyScorecardWorkflow(content, allowCommitVerifier)
			if !tt.wantErr {
				assert.NoError(t, err)
				return
			}
			// An unsupported runner must surface as a verificationError wrapping
			// errInvalidRunnerLabel so the API returns 400, not 500.
			assert.ErrorIs(t, err, errInvalidRunnerLabel)
			var vErr verificationError
			assert.ErrorAs(t, err, &vErr)
		})
	}
}

// suffix may not be the best term, but maps the final part of a path to a response file.
// this is helpful when multiple API calls need to be made.
// e.g. a call to /foo/bar/some/endpoint would have "endpoint" as a suffix.
type suffixStubTripper struct {
	// key is suffix, value is response file.
	responsePaths map[string]string
}

func (s suffixStubTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	pathParts := strings.Split(r.URL.Path, "/")
	suffix := pathParts[len(pathParts)-1]
	f, err := os.Open(s.responsePaths[suffix])
	if err != nil {
		return nil, err
	}
	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       f,
	}, nil
}

func FuzzVerifyWorkflow(f *testing.F) {
	testfiles := []string{
		"testdata/workflow-valid.yml",
		"testdata/workflow-valid-noglobalperm.yml",
		"testdata/workflow-valid-e2e.yml",
		"testdata/workflow-valid-e2e-ghcr.yml",
		"testdata/workflow-valid-tagged-action.yml",
	}
	for _, file := range testfiles {
		content, err := os.ReadFile(file)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(string(content))
	}
	f.Fuzz(func(t *testing.T, data string) {
		if !utf8.ValidString(data) {
			t.Skip()
		}
		verifyScorecardWorkflow(data, allowCommitVerifier)
	})
}
