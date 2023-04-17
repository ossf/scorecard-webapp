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
	"os"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestVerifyValidWorkflows(t *testing.T) {
	t.Parallel()
	workflowFiles := []string{
		"testdata/workflow-valid.yml",
		"testdata/workflow-valid-noglobalperm.yml",
		"testdata/workflow-valid-e2e.yml",
	}

	for _, workflowFile := range workflowFiles {
		workflowContent, _ := os.ReadFile(workflowFile)
		err := verifyScorecardWorkflow(string(workflowContent), nil)
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
	}

	for _, workflowFile := range workflowFiles {
		workflowContent, _ := os.ReadFile(workflowFile)
		err := verifyScorecardWorkflow(string(workflowContent), nil)
		assert.NotEqual(t, err, nil, workflowFile)
	}
}

//nolint
func FuzzVerifyWorkflow(f *testing.F) {
	f.Fuzz(func(t *testing.T, data string) {
		if !utf8.ValidString(data) {
			t.Skip()
		}
		verifyScorecardWorkflow(data, nil)
	})
}
