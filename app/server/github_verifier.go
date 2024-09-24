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
	"fmt"
	"net/http"

	"github.com/google/go-github/v65/github"
)

type githubVerifier struct {
	ctx    context.Context
	client *github.Client
}

// contains makes two "core" API calls: one for the default branch, and one to compare the target hash to a branch
// if the repo is "github/codeql-action", also check releases/v1 before failing.
func (g *githubVerifier) contains(owner, repo, hash string) (bool, error) {
	defaultBranch, err := g.defaultBranch(owner, repo)
	if err != nil {
		return false, err
	}
	contains, err := g.branchContains(defaultBranch, owner, repo, hash)
	if err != nil {
		return false, err
	}
	if contains {
		return true, nil
	}

	switch {
	// github/codeql-action has commits from their release branches that don't show up in the default branch
	// this isn't the best approach for now, but theres no universal "does this commit belong to this repo" call
	case owner == "github" && repo == "codeql-action":
		releaseBranches := []string{"releases/v3", "releases/v2", "releases/v1"}
		for _, branch := range releaseBranches {
			contains, err = g.branchContains(branch, owner, repo, hash)
			if err != nil {
				return false, err
			}
			if contains {
				return true, nil
			}
		}

	// add fallback lookup for actions/upload-artifact v3/node20 branch
	// https://github.com/actions/starter-workflows/pull/2348#discussion_r1536228344
	case owner == "actions" && repo == "upload-artifact":
		contains, err = g.branchContains("v3/node20", owner, repo, hash)
	}
	return contains, err
}

func (g *githubVerifier) defaultBranch(owner, repo string) (string, error) {
	githubRepository, _, err := g.client.Repositories.Get(g.ctx, owner, repo)
	if err != nil {
		return "", fmt.Errorf("fetching repository info: %w", err)
	}
	if githubRepository == nil || githubRepository.DefaultBranch == nil {
		return "", errNoDefaultBranch
	}
	return *githubRepository.DefaultBranch, nil
}

func (g *githubVerifier) branchContains(branch, owner, repo, hash string) (bool, error) {
	opts := &github.ListOptions{PerPage: 1}
	diff, resp, err := g.client.Repositories.CompareCommits(g.ctx, owner, repo, branch, hash, opts)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			// NotFound can be returned for some divergent cases: "404 No common ancestor between ..."
			return false, nil
		}
		return false, fmt.Errorf("error comparing revisions: %w", err)
	}

	// Target should be behind or at the base ref if it is considered contained.
	return diff.GetStatus() == "behind" || diff.GetStatus() == "identical", nil
}
