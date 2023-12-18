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
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/rhysd/actionlint"
)

const (
	workflowRestrictionLink = "https://github.com/ossf/scorecard-action#workflow-restrictions"
)

var (
	errActionlintParse              = errors.New("errors during actionlint.Parse")
	errGlobalVarsOrDefaults         = errors.New("workflow contains global env vars or defaults")
	errGlobalWriteAll               = errors.New("global perm is set to write-all")
	errGlobalWrite                  = errors.New("global perm is set to write")
	errScorecardJobNotFound         = errors.New("workflow has no job that calls ossf/scorecard-action")
	errNonScorecardJobHasTokenWrite = errors.New("workflow has a non-scorecard job with id-token permissions")
	errJobHasContainerOrServices    = errors.New("job contains container or service")
	errScorecardJobRunsOn           = errors.New("scorecard job should have exactly 1 'Ubuntu' virtual environment")
	errInvalidRunnerLabel           = errors.New("scorecard job has invalid runner label")
	errUnallowedStepName            = errors.New("job has unallowed step")
	errScorecardJobEnvVars          = errors.New("scorecard job contains env vars")
	errScorecardJobDefaults         = errors.New("scorecard job must not have defaults set")
	errEmptyStepUses                = errors.New("scorecard job must only have steps with `uses`")
	errNoDefaultBranch              = errors.New("no default branch")

	reCommitSHA = regexp.MustCompile(`^[0-9a-fA-F]{40}$`)
)

// TODO(#290): retrieve the runners dynamically.
// List below is from https://docs.github.com/en/actions/using-github-hosted-runners/about-github-hosted-runners.
var ubuntuRunners = map[string]bool{
	"ubuntu-latest": true,
	"ubuntu-22.04":  true,
	"ubuntu-20.04":  true,
	"ubuntu-18.04":  true,
}

type commitVerifier interface {
	contains(owner, repo, hash string) (bool, error)
}

type verificationError struct {
	e error
}

func (ve verificationError) Error() string {
	return fmt.Sprintf("workflow verification failed: %v, see %s for details.", ve.e, workflowRestrictionLink)
}

func (ve verificationError) Unwrap() error {
	return ve.e
}

type imposterCommitError struct {
	action, ref string
}

func (i imposterCommitError) Error() string {
	return fmt.Sprintf("imposter commit: %s does not belong to %s", i.ref, i.action)
}

func verifyScorecardWorkflow(workflowContent string, verifier commitVerifier) error {
	// Verify workflow contents using actionlint.
	workflow, lintErrs := actionlint.Parse([]byte(workflowContent))
	if lintErrs != nil || workflow == nil {
		return fmt.Errorf("%w: %v", errActionlintParse, lintErrs)
	}

	// Verify that there are no global env vars or defaults.
	if workflow.Env != nil || workflow.Defaults != nil {
		return verificationError{e: errGlobalVarsOrDefaults}
	}

	if workflow.Permissions != nil {
		globalPerms := workflow.Permissions
		// Verify that the all scope, if set, isn't write-all.
		if globalPerms.All != nil && globalPerms.All.Value == "write-all" {
			return verificationError{e: errGlobalWriteAll}
		}

		// Verify that there are no global permissions (including id-token) set to write.
		for globalPerm, val := range globalPerms.Scopes {
			if val.Value.Value == "write" {
				return verificationError{e: fmt.Errorf("%w: permission for %v is set to write",
					errGlobalWrite, globalPerm)}
			}
		}
	}

	// Find the (first) job with a step that calls scorecard-action.
	scorecardJob := findScorecardJob(workflow.Jobs)
	if scorecardJob == nil {
		return verificationError{e: errScorecardJobNotFound}
	}

	// Make sure other jobs don't have id-token permissions.
	for _, job := range workflow.Jobs {
		if job != scorecardJob && job.Permissions != nil {
			idToken := job.Permissions.Scopes["id-token"]
			if idToken != nil && idToken.Value.Value == "write" {
				return verificationError{e: errNonScorecardJobHasTokenWrite}
			}
		}
	}

	// Verify that there is no job container or services.
	if scorecardJob.Container != nil || len(scorecardJob.Services) > 0 {
		return verificationError{e: errJobHasContainerOrServices}
	}

	labels := scorecardJob.RunsOn.Labels
	if len(labels) != 1 {
		return verificationError{e: errScorecardJobRunsOn}
	}
	label := labels[0].Value
	if _, ok := ubuntuRunners[label]; !ok {
		return fmt.Errorf("%w: '%s'", errInvalidRunnerLabel, label)
	}

	// Verify that there are no job env vars set.
	if scorecardJob.Env != nil {
		return verificationError{e: errScorecardJobEnvVars}
	}

	// Verify that there are no job defaults set.
	if scorecardJob.Defaults != nil {
		return verificationError{e: errScorecardJobDefaults}
	}

	// Get steps in job.
	steps := scorecardJob.Steps

	// Verify that steps are valid (checkout, scorecard-action, upload-artifact, upload-sarif).
	for _, step := range steps {
		stepUses := getStepUses(step)
		if stepUses == nil {
			return verificationError{e: errEmptyStepUses}
		}
		stepName, ref := parseStep(stepUses.Value)

		switch stepName {
		case
			"actions/checkout",
			"ossf/scorecard-action",
			"actions/upload-artifact",
			"github/codeql-action/upload-sarif",
			"step-security/harden-runner":
			if isCommitHash(ref) {
				s := strings.Split(stepName, "/")
				// no need to length check as the step name is one of the ones above
				owner, repo := s[0], s[1]
				contains, err := verifier.contains(owner, repo, ref)
				if err != nil {
					return err
				}
				if !contains {
					return verificationError{e: imposterCommitError{ref: ref, action: stepName}}
				}
			}
		// Needed for e2e tests
		case "gcr.io/openssf/scorecard-action":
		default:
			return verificationError{e: fmt.Errorf("%w: %s", errUnallowedStepName, stepName)}
		}
	}

	return nil
}

// Finds the job with a step that calls ossf/scorecard-action.
func findScorecardJob(jobs map[string]*actionlint.Job) *actionlint.Job {
	for _, job := range jobs {
		if job == nil {
			continue
		}
		for _, step := range job.Steps {
			stepUses := getStepUses(step)
			if stepUses == nil {
				continue
			}
			stepName, _ := parseStep(stepUses.Value)
			if stepName == "ossf/scorecard-action" ||
				stepName == "gcr.io/openssf/scorecard-action" {
				return job
			}
		}
	}
	return nil
}

func parseStep(step string) (name, ref string) {
	// Check for `uses: ossf/scorecard-action@ref`.
	reRef := regexp.MustCompile(`^([^@]*)@(.*)$`)
	refMatches := reRef.FindStringSubmatch(step)
	if len(refMatches) > 2 {
		return refMatches[1], refMatches[2]
	}

	// Check for `uses: docker://gcr.io/openssf/scorecard-action:tag`.
	reDocker := regexp.MustCompile(`^docker://([^:]*):.*$`)
	dockerMatches := reDocker.FindStringSubmatch(step)
	if len(dockerMatches) > 1 {
		// TODO don't currently need ref for the docker images
		return dockerMatches[1], ""
	}
	return "", ""
}

func getStepUses(step *actionlint.Step) *actionlint.String {
	if step.Exec == nil {
		return nil
	}
	execAction, exists := step.Exec.(*actionlint.ExecAction)
	if !exists || execAction == nil {
		return nil
	}
	return execAction.Uses
}

func isCommitHash(s string) bool {
	return reCommitSHA.MatchString(s)
}

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
	// github/codeql-action has commits from their v1 and v2 release branch that don't show up in the default branch
	// this isn't the best approach for now, but theres no universal "does this commit belong to this repo" call
	if owner == "github" && repo == "codeql-action" {
		contains, err = g.branchContains("releases/v2", owner, repo, hash)
		if err != nil {
			return false, err
		}
		if !contains {
			contains, err = g.branchContains("releases/v1", owner, repo, hash)
		}
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
