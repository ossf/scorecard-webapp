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
	"context"
	"io"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ossf/scorecard-webapp/app/generated/models"
)

var _ = Describe("E2E Test: extractAndVerifyCertForPayload", func() {
	AssertValidPayload := func(filename string) {
		It("Should successfully extract cert for payload", func() {
			testFile, err := os.Open(filename)
			Expect(err).Should(BeNil())

			payload, err := io.ReadAll(testFile)
			Expect(err).Should(BeNil())

			_, errCertExtract := extractAndVerifyCertForPayload(context.Background(), payload, noTlogIndex)
			Expect(errCertExtract).Should(BeNil())
		})
	}
	Context("E2E Test: Validate functionality", func() {
		AssertValidPayload("testdata/results/valid-payload.json")
	})
	Context("E2E Test: Validate functionality", func() {
		AssertValidPayload("testdata/results/valid-payload-2.json")
	})
})

// readGitHubTokens is used to authenticate the github client in the getAndVerifyWorkflowContent tests
// The CI/CD will have GITHUB_TOKEN available.
func readGitHubTokens() (string, bool) {
	githubAuthTokens := []string{"GITHUB_AUTH_TOKEN", "GITHUB_TOKEN", "GH_TOKEN", "GH_AUTH_TOKEN"}
	for _, name := range githubAuthTokens {
		if token, exists := os.LookupEnv(name); exists && token != "" {
			return token, exists
		}
	}
	return "", false
}

var _ = Describe("E2E Test: getAndVerifyWorkflowContent", func() {
	AssertValidWorkflowContent := func(filename string) {
		It("Should successfully extract cert and verify workflow for payload", func() {
			testFile, err := os.Open(filename)
			Expect(err).Should(BeNil())

			payload, err := io.ReadAll(testFile)
			Expect(err).Should(BeNil())

			ctx := context.Background()
			cert, errCertExtract := extractAndVerifyCertForPayload(ctx, payload, noTlogIndex)
			Expect(errCertExtract).Should(BeNil())

			info, errCertExtractInfo := extractCertInfo(cert)
			Expect(errCertExtractInfo).Should(BeNil())

			token, _ := readGitHubTokens()
			scorecardResult := &models.VerifiedScorecardResult{
				AccessToken: token,
				Branch:      "main",
				Result:      string(payload),
			}
			Expect(getAndVerifyWorkflowContent(ctx, scorecardResult, info)).Should(BeNil())
		})
	}
	Context("E2E Test: Validate functionality intra-repo", func() {
		AssertValidWorkflowContent("testdata/results/reusable-workflow-intra-repo-results.json")
	})
	Context("E2E Test: Validate functionality inter-repo", func() {
		AssertValidWorkflowContent("testdata/results/reusable-workflow-inter-repo-results.json")
	})
})
