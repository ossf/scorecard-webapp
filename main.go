// Copyright 2021 OpenSSF Scorecard Authors
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

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-openapi/loads"
	flag "github.com/spf13/pflag"

	"github.com/ossf/scorecard-webapp/app/generated/restapi"
	"github.com/ossf/scorecard-webapp/app/generated/restapi/operations"
)

func main() {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	var server *restapi.Server // make sure init is called

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage:\n")
		fmt.Fprint(os.Stderr, "  scorecard-server [OPTIONS]\n\n")

		title := "OpenSSF Scorecard API"
		fmt.Fprint(os.Stderr, title+"\n\n")
		desc := "API to interact with a project's published Scorecard result"
		if desc != "" {
			fmt.Fprintf(os.Stderr, desc+"\n\n")
		}
		fmt.Fprintln(os.Stderr, flag.CommandLine.FlagUsages())
	}
	// parse the CLI flags
	flag.Parse()

	api := operations.NewScorecardAPI(swaggerSpec)
	// get server with flag values filled out
	server = restapi.NewServer(api)
	defer func() {
		if err := server.Shutdown(); err != nil {
			log.Println(err)
		}
	}()

	server.ConfigureAPI()
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
