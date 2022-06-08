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
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	_ "gocloud.dev/blob/gcsblob" // Needed to link in GCP drivers.
)

const (
	shieldsURL = "https://img.shields.io/ossf-scorecard"
	badgeLabel = "openssf scorecard"
)

func GetBadgeHandler(w http.ResponseWriter, r *http.Request) {
	host := mux.Vars(r)["host"]
	orgName := mux.Vars(r)["orgName"]
	repoName := mux.Vars(r)["repoName"]
	parsedURL, err := url.Parse(fmt.Sprintf("%s/%s/%s/%s?label=%s", shieldsURL, host, orgName, repoName, badgeLabel))
	if err != nil {
		log.Printf("parsing badge URL: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, parsedURL.String(), http.StatusFound)
}
