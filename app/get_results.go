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

package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/gcsblob" // Needed to link in GCP drivers.
)

const (
	scorecardResultBucketURL     = "gs://ossf-scorecard-results"
	scorecardCronResultBucketURL = "gs://ossf-scorecard-cron-results"
)

var errInvalidInputs = errors.New("invalid inputs provided")

func GetResultsHandler(w http.ResponseWriter, r *http.Request) {
	host := mux.Vars(r)["host"]
	orgName := mux.Vars(r)["orgName"]
	repoName := mux.Vars(r)["repoName"]
	commit := r.URL.Query().Get("commit")
	results, err := getResults(host, orgName, repoName, commit)
	if err == nil {
		_, err := w.Write(results)
		if err != nil {
			log.Printf("err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		return
	}
	// TODO: Remove if here once all errors have been migrated to
	// map to a http status code.
	if errors.Is(err, errNotFound) {
		errHandler(w, err)
		return
	}

	if errors.Is(err, errInvalidInputs) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("err: %v", err)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
	log.Printf("err: %v", err)
}

func getResults(host, orgName, repoName, commit string) ([]byte, error) {
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

func sanitizeInputs(host, orgName, repoName, commit string) (string, error) {
	resultsFile := filepath.Join(host, orgName, repoName, commit, "results.json")
	if commit == "" {
		resultsFile = filepath.Join(host, orgName, repoName, "results.json")
	}
	cleanResultsFile := filepath.Clean(resultsFile)
	cleanResultsFile = strings.Replace(cleanResultsFile, "\n", "", -1)
	cleanResultsFile = strings.Replace(cleanResultsFile, "\r", "", -1)
	var matched bool
	var err error
	if commit == "" {
		matched, err = filepath.Match("*/*/*/results.json", cleanResultsFile)
	} else {
		matched, err = filepath.Match("*/*/*/*/results.json", cleanResultsFile)
	}
	if err != nil || !matched {
		return "", errInvalidInputs
	}
	return cleanResultsFile, nil
}
