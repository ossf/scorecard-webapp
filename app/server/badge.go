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

package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	_ "gocloud.dev/blob/gcsblob" // Needed to link in GCP drivers.

	"github.com/ossf/scorecard-webapp/app/generated/models"
	"github.com/ossf/scorecard-webapp/app/generated/restapi/operations/badge"
)

const (
	shieldsURL = "https://img.shields.io/ossf-scorecard"
	badgeLabel = "openssf scorecard"
)

func GetBadgeHandler(params badge.GetBadgeParams) middleware.Responder {
	host := params.Platform
	orgName := params.Org
	repoName := params.Repo
	parsedURL, err := url.Parse(fmt.Sprintf("%s/%s/%s/%s?label=%s", shieldsURL, host, orgName, repoName, badgeLabel))
	if err != nil {
		return badge.NewGetBadgeDefault(http.StatusInternalServerError).WithPayload(&models.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, producer runtime.Producer) {
		http.Redirect(rw, params.HTTPRequest, parsedURL.String(), http.StatusFound)
	})
}
