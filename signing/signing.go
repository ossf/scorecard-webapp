package signing

import (
	"context"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/ossf/scorecard/v2/cron/data"
)

type ScorecardOutput struct {
	SarifOutput string
	JsonOutput  string
}

func GetResults(w http.ResponseWriter, r *http.Request) {
	// Get params to build GCS filepath.
	ctx := context.Background()
	bucketURL := "gs://ossf-scorecard-results"
	host := mux.Vars(r)["host"]
	orgName := mux.Vars(r)["orgName"]
	repoName := mux.Vars(r)["repoName"]
	resultsFile := filepath.Join(host, orgName, repoName, "results.json")

	// Sanitize input and log query.
	resultsFileCleaned := filepath.Clean(resultsFile)
	matched, err := filepath.Match("*/*/*/results.json", resultsFileCleaned)
	if err != nil || !matched {
		http.Error(w, "error verifying filepath format", http.StatusInternalServerError)
		log.Println(matched, err)
		return
	}
	log.Printf("Querying GCS bucket for: %s", resultsFileCleaned)

	// Query GCS bucket.
	resultsBytes, err := data.GetBlobContent(ctx, bucketURL, resultsFileCleaned)
	if err != nil {
		http.Error(w, "error pulling from GCS bucket", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Write results json to response.
	w.Header().Set("Content-Type", "application/json")
	w.Write(resultsBytes)
}
