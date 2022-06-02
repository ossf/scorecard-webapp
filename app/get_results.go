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

const scorecardResultBucketURL = "gs://ossf-scorecard-results"

var (
	errRetrievingData = errors.New("error pulling from GCS bucket")
	errInvalidInputs  = errors.New("invalid inputs provided")
	errBucketNotFound = errors.New("bucket not found")
)

func GetResultsHandler(w http.ResponseWriter, r *http.Request) {
	host := mux.Vars(r)["host"]
	orgName := mux.Vars(r)["orgName"]
	repoName := mux.Vars(r)["repoName"]
	results, err := getResults(host, orgName, repoName)
	if err == nil {
		_, err := w.Write(results)
		if err != nil {
			log.Printf("err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	if errors.Is(err, errInvalidInputs) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("err: %v", err)
	} else if errors.Is(err, errBucketNotFound) ||
		errors.Is(err, errRetrievingData) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("err: %v", err)
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func getResults(host, orgName, repoName string) ([]byte, error) {
	// Sanitize input and log query.
	cleanResultsFile, err := sanitizeInputs(host, orgName, repoName)
	if err != nil {
		return nil, err
	}
	log.Printf("Querying GCS bucket for: %s", cleanResultsFile)

	// Query GCS bucket.
	ctx := context.Background()
	bucket, err := blob.OpenBucket(ctx, scorecardResultBucketURL)
	if err != nil {
		return nil, errBucketNotFound
	}
	results, err := bucket.ReadAll(ctx, cleanResultsFile)
	if err != nil {
		return nil, errRetrievingData
	}
	return results, nil
}

func sanitizeInputs(host, orgName, repoName string) (string, error) {
	resultsFile := filepath.Join(host, orgName, repoName, "results.json")
	cleanResultsFile := filepath.Clean(resultsFile)
	cleanResultsFile = strings.Replace(cleanResultsFile, "\n", "", -1)
	cleanResultsFile = strings.Replace(cleanResultsFile, "\r", "", -1)
	matched, err := filepath.Match("*/*/*/results.json", cleanResultsFile)
	if err != nil || !matched {
		return "", errInvalidInputs
	}
	return cleanResultsFile, nil
}
