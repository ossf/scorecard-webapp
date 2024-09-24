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
	"net/http"
	"os"

	"github.com/google/go-github/v65/github"
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
	ctx := context.Background()
	getPayload := func(filename string) []byte {
		testFile, err := os.Open(filename)
		Expect(err).Should(BeNil())
		payload, err := io.ReadAll(testFile)
		Expect(err).Should(BeNil())
		return payload
	}
	extractCertInfo := func(payload []byte) certInfo {
		cert, errCertExtract := extractAndVerifyCertForPayload(ctx, payload, noTlogIndex)
		Expect(errCertExtract).Should(BeNil())
		info, errCertExtractInfo := extractCertInfo(cert)
		Expect(errCertExtractInfo).Should(BeNil())
		return info
	}
	AssertValidWorkflowContent := func(filename string) {
		It("Should pass workflow verifcation", func() {
			payload := getPayload(filename)
			info := extractCertInfo(payload)
			token, _ := readGitHubTokens()
			scorecardResult := &models.VerifiedScorecardResult{
				AccessToken: token,
				Branch:      "main",
				Result:      string(payload),
				TlogIndex:   noTlogIndex,
			}
			Expect(getAndVerifyWorkflowContent(ctx, scorecardResult, info)).Should(BeNil())
		})
	}
	AssertInvalidWorkflowContent := func(filename string, errSubstr string) {
		It("Should fail workflow verifcation", func() {
			payload := getPayload(filename)
			info := extractCertInfo(payload)
			token, _ := readGitHubTokens()
			scorecardResult := &models.VerifiedScorecardResult{
				AccessToken: token,
				Branch:      "main",
				Result:      string(payload),
			}
			err := getAndVerifyWorkflowContent(ctx, scorecardResult, info)
			Expect(err).Should(MatchError(ContainSubstring(errSubstr)))
		})
	}
	Context("E2E Test: Validate functionality intra-repo", func() {
		AssertValidWorkflowContent("testdata/results/reusable-workflow-intra-repo-results.json")
	})
	Context("E2E Test: Validate functionality inter-repo", func() {
		AssertValidWorkflowContent("testdata/results/reusable-workflow-inter-repo-results.json")
	})
	Context("E2E Test: Fail on imposter commit", func() {
		AssertInvalidWorkflowContent("testdata/results/imposter-commit-results.json", "imposter commit")
	})
})

// helper function to setup a github verifier with an appropriately set token.
func getGithubVerifier() githubVerifier {
	httpClient := http.DefaultClient
	token, _ := readGitHubTokens()
	if token != "" {
		httpClient.Transport = githubTransport{
			token: token,
		}
	}
	return githubVerifier{
		ctx:    context.Background(),
		client: github.NewClient(httpClient),
	}
}

var _ = Describe("E2E Test: githubVerifier_contains", func() {
	Context("E2E Test: Validate known good commits", func() {
		It("can detect actions/upload-artifact v3-node20 commits", func() {
			gv := getGithubVerifier()
			c, err := gv.contains(commit{"actions", "upload-artifact", "97a0fba1372883ab732affbe8f94b823f91727db"})
			Expect(err).Should(BeNil())
			Expect(c).To(BeTrue())
		})

		It("can detect github/codeql-action backport commits", func() {
			gv := getGithubVerifier()
			c, err := gv.contains(commit{"github", "codeql-action", "a82bad71823183e5b120ab52d521460ecb0585fe"})
			Expect(err).Should(BeNil())
			Expect(c).To(BeTrue())
		})
	})
})
