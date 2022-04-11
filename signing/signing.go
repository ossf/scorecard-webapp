package signing

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/ossf/scorecard/v2/cron/data"
)

type ScorecardOutput struct {
	JsonOutput string
}

var errorPullingBucket = errors.New("error pulling from GCS bucket")
var errorVerifyingFilepath = errors.New("error verifying filepath format")

func GetResults(w http.ResponseWriter, r *http.Request) {
	host := mux.Vars(r)["host"]
	orgName := mux.Vars(r)["orgName"]
	repoName := mux.Vars(r)["repoName"]
	results, err := getResults(host, orgName, repoName)

	if err == errorVerifyingFilepath {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("err: %v", err)
		return
	}
	if err == errorPullingBucket {
		http.Error(w, err.Error(), http.StatusNotFound)
		log.Printf("err: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(results)
}

func getResults(host, orgName, repoName string) (results []byte, err error) {
	// Get params to build GCS filepath.
	ctx := context.Background()
	bucketURL := "gs://ossf-scorecard-results"
	resultsFile := filepath.Join(host, orgName, repoName, "results.json")

	// Sanitize input and log query.
	resultsFile = filepath.Clean(resultsFile)
	matched, err := filepath.Match("*/*/*/results.json", resultsFile)
	if err != nil || !matched {
		return nil, errorVerifyingFilepath
	}

	if len(resultsFile) >= 256 {
		return nil, fmt.Errorf("filepath (%v) is greater than the Linux maximum of 256", resultsFile[:256])
	}

	log.Printf("Querying GCS bucket for: %s", resultsFile)

	// Query GCS bucket.
	results, err = data.GetBlobContent(ctx, bucketURL, resultsFile)
	if err != nil {
		return nil, errorPullingBucket
	}
	return results, nil
}
