// Copyright 2021 Security Scorecard Authors
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

// Package main implements the scorecard.dev webapp.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ossf/scorecard-webapp/signing"
)

func main() {
	fmt.Printf("Starting HTTP server on port 8080 ...\n")

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homepage)
	r.HandleFunc("/projects/{host}/{orgName}/{repoName}", signing.GetResults).Methods("GET")
	http.Handle("/", r)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func homepage(w http.ResponseWriter, r *http.Request) {
	endpts := struct {
		GetRepoResults string `json:"get_repo_results"`
	}{
		// TODO: make this domain specific.
		GetRepoResults: "/projects{/host}{/owner}{/repository}",
	}
	endptsBytes, err := json.MarshalIndent(endpts, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := fmt.Fprint(w, string(endptsBytes)); err != nil {
		log.Fatal(err)
	}
}
