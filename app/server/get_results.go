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
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/gcsblob" // Needed to link in GCP drivers.

	"github.com/ossf/scorecard-webapp/app/generated/models"
	"github.com/ossf/scorecard-webapp/app/generated/restapi/operations/results"
)

const (
	scorecardResultBucketURL     = "gs://ossf-scorecard-results"
	scorecardCronResultBucketURL = "gs://ossf-scorecard-cron-results"
)

var errInvalidInputs = errors.New("invalid inputs provided")

func GetResultHandler(params results.GetResultParams) middleware.Responder {
	res, err := getResults(params.Platform, params.Org, params.Repo, params.Commit)

	if errors.Is(err, errNotFound) {
		return results.NewGetResultNotFound()
	}
	if errors.Is(err, errInvalidInputs) {
		return results.NewGetResultBadRequest()
	}
	if err == nil {
		var ret models.ScorecardResult
		if err = ret.UnmarshalBinary(res); err == nil {
			return results.NewGetResultOK().WithPayload(&ret)
		}
	}

	return results.NewGetResultDefault(http.StatusInternalServerError).WithPayload(&models.Error{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	})
}

func getResults(host, orgName, repoName string, commit *string) ([]byte, error) {
	// Sanitize input and log query.
	cleanResultsFile, err := sanitizeInputs(host, orgName, repoName, commit)
	if err != nil {
		return nil, err
	}
	log.Printf("Querying GCS bucket for: %s", cleanResultsFile)

	// Query GCS bucket.
	ctx := context.Background()
	if bucket, err := blob.OpenBucket(ctx, scorecardResultBucketURL); err == nil {
		if results, err := bucket.ReadAll(ctx, cleanResultsFile); err == nil {
			return results, nil
		}
	}

	cleanResultsFile2, err := sanitizeInputs(host, orgName, repoName, commit)
	if err != nil {
		return nil, err
	}

	// Try the backup cron bucket.
	if bucket, err := blob.OpenBucket(ctx, scorecardCronResultBucketURL); err == nil {
		if results, err := bucket.ReadAll(ctx, cleanResultsFile2); err == nil {
			return results, nil
		}
	}

	return nil, errNotFound
}

func sanitizeInputs(host, orgName, repoName string, commit *string) (string, error) {
	resultsFile := filepath.Join(host, orgName, repoName, "results.json")
	if commit != nil {
		resultsFile = filepath.Join(host, orgName, repoName, *commit, "results.json")
	}
	cleanResultsFile := filepath.Clean(resultsFile)
	cleanResultsFile = strings.Replace(cleanResultsFile, "\n", "", -1)
	cleanResultsFile = strings.Replace(cleanResultsFile, "\r", "", -1)
	var matched bool
	var err error
	if commit == nil {
		matched, err = filepath.Match("*/*/*/results.json", cleanResultsFile)
	} else {
		matched, err = filepath.Match("*/*/*/*/results.json", cleanResultsFile)
	}
	if err != nil || !matched {
		return "", errInvalidInputs
	}
	return cleanResultsFile, nil
}
