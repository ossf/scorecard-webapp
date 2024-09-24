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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/go-github/v65/github"
	"golang.org/x/mod/semver"
)

var errInvalidCodeQLVersion = errors.New("codeql version invalid")

type githubVerifier struct {
	ctx               context.Context
	client            *github.Client
	cachedCommits     map[commit]bool
	codeqlActionMajor string
}

// returns a new githubVerifier, with an instantiated map.
// most uses should use this constructor.
func newGitHubVerifier(ctx context.Context, client *github.Client) *githubVerifier {
	verifier := githubVerifier{
		ctx:    ctx,
		client: client,
	}
	verifier.cachedCommits = map[commit]bool{}
	return &verifier
}

// contains may make several "core" API calls:
//   - one to get repository tags (we expect most cases to only need this call)
//   - two to get the default branch name and check it
//   - up to 10 requests when checking previous release branches
func (g *githubVerifier) contains(c commit) (bool, error) {
	// return the cached answer if we've seen this commit before
	if contains, ok := g.cachedCommits[c]; ok {
		return contains, nil
	}

	// fetch 100 most recent tags first, as this should handle the most common scenario
	if err := g.getTags(c.owner, c.repo); err != nil {
		return false, err
	}

	// check cache again now that it's populated with tags
	if contains, ok := g.cachedCommits[c]; ok {
		return contains, nil
	}

	// check default branch
	defaultBranch, err := g.defaultBranch(c.owner, c.repo)
	if err != nil {
		return false, err
	}
	contains, err := g.branchContains(defaultBranch, c.owner, c.repo, c.hash)
	if err != nil {
		return false, err
	}
	if contains {
		return true, nil
	}

	// finally, check the most recent 10 release branches. This limit is arbitrary and can be adjusted in the future.
	const lookback = 10
	return g.checkReleaseBranches(c.owner, c.repo, c.hash, defaultBranch, lookback)
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
	contains := diff.GetStatus() == "behind" || diff.GetStatus() == "identical"
	if contains {
		g.markContains(owner, repo, hash)
	}
	return contains, nil
}

func (g *githubVerifier) getTags(owner, repo string) error {
	// 100 releases is the maximum supported by /repos/{owner}/{repo}/tags endpoint
	opts := github.ListOptions{PerPage: 100}
	tags, _, err := g.client.Repositories.ListTags(g.ctx, owner, repo, &opts)
	if err != nil {
		return fmt.Errorf("fetch tags: %w", err)
	}
	isCodeQL := owner == "github" && repo == "codeql-action"
	for _, t := range tags {
		if t.Commit != nil && t.Commit.SHA != nil {
			g.markContains(owner, repo, *t.Commit.SHA)
		}
		// store the highest major version for github/codeql-action
		// this helps check release branches later, due to their release process.
		if isCodeQL {
			version := t.GetName()
			if semver.IsValid(version) && semver.Compare(version, g.codeqlActionMajor) == 1 {
				g.codeqlActionMajor = semver.Major(version)
			}
		}
	}
	return nil
}

func (g *githubVerifier) markContains(owner, repo, sha string) {
	commit := commit{
		owner: owner,
		repo:  repo,
		hash:  strings.ToLower(sha),
	}
	g.cachedCommits[commit] = true
}

// check the most recent release branches, ignoring the default branch which was already checked.
func (g *githubVerifier) checkReleaseBranches(owner, repo, hash, defaultBranch string, limit int) (bool, error) {
	var (
		analyzedBranches int
		branches         []string
		err              error
	)

	switch {
	// special case: github/codeql-action releases all come from "main", even though the tags are on different branches
	case owner == "github" && repo == "codeql-action":
		branches, err = g.getCodeQLBranches()
		if err != nil {
			return false, err
		}
	default:
		branches, err = g.getBranches(owner, repo)
		if err != nil {
			return false, err
		}
		// we may have discovered more commit SHAs from the release process
		c := commit{
			owner: owner,
			repo:  repo,
			hash:  hash,
		}
		if contains, ok := g.cachedCommits[c]; ok {
			return contains, nil
		}
	}

	for _, branch := range branches {
		if analyzedBranches >= limit {
			break
		}
		if branch == defaultBranch {
			continue
		}
		analyzedBranches++
		contains, err := g.branchContains(branch, owner, repo, hash)
		if err != nil {
			return false, err
		}
		if contains {
			return true, nil
		}
	}

	return false, nil
}

// returns the integer version from the expected format: "v1", "v2", "v3" ..
func parseCodeQLVersion(version string) (int, error) {
	if !strings.HasPrefix(version, "v") {
		return 0, fmt.Errorf("%w: %s", errInvalidCodeQLVersion, version)
	}
	major, err := strconv.Atoi(version[1:])
	if major < 1 || err != nil {
		return 0, fmt.Errorf("%w: %s", errInvalidCodeQLVersion, version)
	}
	return major, nil
}

// these branches follow the releases/v3 pattern, so we can make assumptions about what they're called.
// this should be called after g.getTags(), because it requires g.codeqlActionMajor to be set.
func (g *githubVerifier) getCodeQLBranches() ([]string, error) {
	if g.codeqlActionMajor == "" {
		return nil, nil
	}
	version, err := parseCodeQLVersion(g.codeqlActionMajor)
	if err != nil {
		return nil, err
	}
	branches := make([]string, version)
	// descending order (e..g releases/v5, releases/v4, ... releases/v1)
	for i := 0; i < version; i++ {
		branches[i] = "releases/v" + strconv.Itoa(version-i)
	}
	return branches, nil
}

func (g *githubVerifier) getBranches(owner, repo string) ([]string, error) {
	var branches []string
	seen := map[string]struct{}{}

	// 100 releases is the maximum supported by /repos/{owner}/{repo}/releases endpoint
	opts := github.ListOptions{PerPage: 100}
	releases, _, err := g.client.Repositories.ListReleases(g.ctx, owner, repo, &opts)
	if err != nil {
		return nil, fmt.Errorf("fetch releases: %w", err)
	}

	for _, r := range releases {
		if r.TargetCommitish != nil {
			if isCommitHash(*r.TargetCommitish) {
				// if a commit, we know it's in the repo
				g.markContains(owner, repo, *r.TargetCommitish)
			} else {
				// otherwise we have a release branch to check
				if _, ok := seen[*r.TargetCommitish]; !ok {
					seen[*r.TargetCommitish] = struct{}{}
					branches = append(branches, *r.TargetCommitish)
				}
			}
		}
	}

	return branches, nil
}
