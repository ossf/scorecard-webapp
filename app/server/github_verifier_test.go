// Copyright 2024 OpenSSF Scorecard Authors
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
	"net/http"
	"testing"

	"github.com/google/go-github/v65/github"
)

func Test_githubVerifier_contains_codeql_v1(t *testing.T) {
	t.Parallel()
	httpClient := http.Client{
		Transport: suffixStubTripper{
			responsePaths: map[string]string{
				"codeql-action":   "./testdata/api/github/repository.json",     // api call which finds the default branch
				"main...somehash": "./testdata/api/github/divergent.json",      // doesnt belong to default branch
				"v3...somehash":   "./testdata/api/github/divergent.json",      // doesnt belong to releases/v3 branch
				"v2...somehash":   "./testdata/api/github/divergent.json",      // doesnt belong to releases/v2 branch
				"v1...somehash":   "./testdata/api/github/containsCommit.json", // belongs to releases/v1 branch
			},
		},
	}
	client := github.NewClient(&httpClient)
	gv := githubVerifier{
		ctx:    context.Background(),
		client: client,
	}
	got, err := gv.contains("github", "codeql-action", "somehash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != true {
		t.Errorf("expected to contain hash, but it didnt")
	}
}

func Test_githubVerifier_contains_codeql_v2(t *testing.T) {
	t.Parallel()
	httpClient := http.Client{
		Transport: suffixStubTripper{
			responsePaths: map[string]string{
				"codeql-action":   "./testdata/api/github/repository.json",     // api call which finds the default branch
				"main...somehash": "./testdata/api/github/divergent.json",      // doesnt belong to default branch
				"v3...somehash":   "./testdata/api/github/divergent.json",      // doesnt belong to releases/v3 branch either
				"v2...somehash":   "./testdata/api/github/containsCommit.json", // belongs to releases/v2 branch
			},
		},
	}
	client := github.NewClient(&httpClient)
	gv := githubVerifier{
		ctx:    context.Background(),
		client: client,
	}
	got, err := gv.contains("github", "codeql-action", "somehash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != true {
		t.Errorf("expected to contain hash, but it didnt")
	}
}
