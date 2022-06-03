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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	endpts := struct {
		GetRepoResults string `json:"get_repo_results"`
		GetRepoBadge   string `json:"get_repo_badge"`
	}{
		GetRepoResults: "/projects/{host}/{owner}/{repository}",
		GetRepoBadge:   "/projects/{host}/{owner}/{repository}/badge",
	}
	endptsBytes, err := json.MarshalIndent(endpts, "", " ")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := fmt.Fprint(w, string(endptsBytes)); err != nil {
		log.Fatal(err)
	}
}
