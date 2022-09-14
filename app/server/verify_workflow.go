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
	"fmt"
	"regexp"

	"github.com/rhysd/actionlint"
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
	errUnallowedStepName            = errors.New("job has unallowed step")
	errScorecardJobEnvVars          = errors.New("scorecard job contains env vars")
	errScorecardJobDefaults         = errors.New("scorecard job must not have defaults set")
	errEmptyStepUses                = errors.New("scorecard job must only have steps with `uses`")
)

func verifyScorecardWorkflow(workflowContent string) error {
	// Verify workflow contents using actionlint.
	workflow, lintErrs := actionlint.Parse([]byte(workflowContent))
	if lintErrs != nil || workflow == nil {
		return fmt.Errorf("%w: %v", errActionlintParse, lintErrs)
	}

	// Verify that there are no global env vars or defaults.
	if workflow.Env != nil || workflow.Defaults != nil {
		return fmt.Errorf("%w", errGlobalVarsOrDefaults)
	}

	if workflow.Permissions != nil {
		globalPerms := workflow.Permissions
		// Verify that the all scope, if set, isn't write-all.
		if globalPerms.All != nil && globalPerms.All.Value == "write-all" {
			return fmt.Errorf("%w", errGlobalWriteAll)
		}

		// Verify that there are no global permissions (including id-token) set to write.
		for globalPerm, val := range globalPerms.Scopes {
			if val.Value.Value == "write" {
				return fmt.Errorf("%w: permission for %v is set to write",
					errGlobalWrite, globalPerm)
			}
		}
	}

	// Find the (first) job with a step that calls scorecard-action.
	scorecardJob := findScorecardJob(workflow.Jobs)
	if scorecardJob == nil {
		return fmt.Errorf("%w", errScorecardJobNotFound)
	}

	// Make sure other jobs don't have id-token permissions.
	for _, job := range workflow.Jobs {
		if job != scorecardJob && job.Permissions != nil {
			idToken := job.Permissions.Scopes["id-token"]
			if idToken != nil && idToken.Value.Value == "write" {
				return fmt.Errorf("%w", errNonScorecardJobHasTokenWrite)
			}
		}
	}

	// Verify that there is no job container or services.
	if scorecardJob.Container != nil || len(scorecardJob.Services) > 0 {
		return fmt.Errorf("%w", errJobHasContainerOrServices)
	}

	labels := scorecardJob.RunsOn.Labels
	if len(labels) != 1 {
		return fmt.Errorf("%w", errScorecardJobRunsOn)
	}
	jobEnv := labels[0].Value
	if jobEnv != "ubuntu-latest" && jobEnv != "ubuntu-20.04" && jobEnv != "ubuntu-18.04" {
		return fmt.Errorf("%w", errScorecardJobRunsOn)
	}

	// Verify that there are no job env vars set.
	if scorecardJob.Env != nil {
		return fmt.Errorf("%w", errScorecardJobEnvVars)
	}

	// Verify that there are no job defaults set.
	if scorecardJob.Defaults != nil {
		return fmt.Errorf("%w", errScorecardJobDefaults)
	}

	// Get steps in job.
	steps := scorecardJob.Steps

	// Verify that steps are valid (checkout, scorecard-action, upload-artifact, upload-sarif).
	for _, step := range steps {
		stepUses := getStepUses(step)
		if stepUses == nil {
			return fmt.Errorf("%w", errEmptyStepUses)
		}
		stepName := getStepName(stepUses.Value)

		switch stepName {
		case
			"actions/checkout",
			"ossf/scorecard-action",
			"actions/upload-artifact",
			"github/codeql-action/upload-sarif",
			"step-security/harden-runner",
			// Needed for e2e tests
			"gcr.io/openssf/scorecard-action",
			"actions-ecosystem/action-create-issue":
		default:
			return fmt.Errorf("%w: %s", errUnallowedStepName, stepName)
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
			stepName := getStepName(stepUses.Value)
			if stepName == "ossf/scorecard-action" ||
				stepName == "gcr.io/openssf/scorecard-action" {
				return job
			}
		}
	}
	return nil
}

func getStepName(step string) string {
	// Check for `uses: ossf/scorecard-action@ref`.
	reRef := regexp.MustCompile(`^([^@]*)@.*$`)
	refMatches := reRef.FindStringSubmatch(step)
	if len(refMatches) > 1 {
		return refMatches[1]
	}

	// Check for `uses: docker://gcr.io/openssf/scorecard-action:tag`.
	reDocker := regexp.MustCompile(`^docker://([^:]*):.*$`)
	dockerMatches := reDocker.FindStringSubmatch(step)
	if len(dockerMatches) > 1 {
		return dockerMatches[1]
	}
	return ""
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
