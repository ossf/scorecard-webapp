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

package server

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSanitizePath(t *testing.T) {
	t.Parallel()
	testSHA := "sha1"
	testcases := []struct {
		name     string
		host     string
		orgName  string
		repoName string
		commit   *string
		wantPath string
		wantErr  error
	}{
		{
			name:     "error on non-absolute paths",
			host:     "../github.com",
			orgName:  "org",
			repoName: "repo",
			wantErr:  errInvalidInputs,
		},
		{
			name:     "handles tabs and newlines",
			host:     "github.com\n",
			orgName:  "\rorg",
			repoName: "repo\r",
			wantPath: "github.com/org/repo/results.json",
		},
		{
			name:     "error on separator characters",
			host:     "github.com/g",
			orgName:  "org",
			repoName: "repo",
			wantErr:  errInvalidInputs,
		},
		{
			name:     "error on escaped separators",
			host:     "github.com\\/g",
			orgName:  "org",
			repoName: "repo",
			wantErr:  errInvalidInputs,
		},
		{
			name:     "handles commit",
			host:     "github.com",
			orgName:  "org",
			repoName: "repo",
			commit:   &testSHA,
			wantPath: "github.com/org/repo/sha1/results.json",
		},
	}
	for _, tt := range testcases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotPath, gotErr := sanitizeInputs(tt.host, tt.orgName, tt.repoName, tt.commit)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, gotErr)
			}
			if !cmp.Equal(gotPath, tt.wantPath) {
				t.Errorf("expected %s, got %s", tt.wantPath, gotPath)
			}
		})
	}
}
