package signing

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	filePath := fmt.Sprintf("%s/%s/%s/results.json", host, orgName, repoName)

	// Sanitize input and log query.
	escapedFilePath := strings.Replace(filePath, "\n", "", -1)
	escapedFilePath = strings.Replace(escapedFilePath, "\r", "", -1)
	log.Printf("Querying GCS bucket for: %s", escapedFilePath)

	// Query GCS bucket.
	resultsBytes, err := data.GetBlobContent(ctx, bucketURL, filePath)
	if err != nil {
		http.Error(w, "error pulling from GCS bucket", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Write results json to response.
	w.Header().Set("Content-Type", "application/json")
	w.Write(resultsBytes)
}
