// This file is safe to edit. Once it exists it will not be overwritten

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

package restapi

import (
	"crypto/tls"
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/rs/cors"

	"github.com/ossf/scorecard-webapp/app/generated/restapi/operations"
	"github.com/ossf/scorecard-webapp/app/generated/restapi/operations/badge"
	"github.com/ossf/scorecard-webapp/app/generated/restapi/operations/results"
	"github.com/ossf/scorecard-webapp/app/server"
)

//nolint:lll // generated code
//go:generate swagger generate server --target ../../generated --name Scorecard --spec ../../../openapi.yaml --principal interface{}

//go:embed static
var staticDir embed.FS

func configureFlags(api *operations.ScorecardAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ScorecardAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	// api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	api.ResultsGetResultHandler = results.GetResultHandlerFunc(server.GetResultHandler)
	api.ResultsPostResultHandler = results.PostResultHandlerFunc(server.PostResultsHandler)
	api.BadgeGetBadgeHandler = badge.GetBadgeHandlerFunc(server.GetBadgeHandler)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
//
//nolint:lll // generated code.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return cors.Default().Handler(serveStatic(handler))
}

func serveStatic(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/docs":
			http.Redirect(w, r, "/", http.StatusFound)
			return
		// TODO: find a more generic solution.
		case "/",
			"/favicon.ico",
			"/swagger-ui.css",
			"/index.css",
			"/favicon-32x32.png",
			"/favicon-16x16.png",
			"/swagger-ui-bundle.js",
			"/swagger-ui-standalone-preset.js",
			"/swagger-initializer.js":
			static, err := fs.Sub(staticDir, "static")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			http.FileServer(http.FS(static)).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
