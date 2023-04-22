// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/ytka/gcs-proxy-cloud-run/config"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	// initialize
	log.Print("starting server...")
	handler := http.HandlerFunc(ProxyHTTPGCS)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Warn().Msgf("defaulting to port %s", port)
	}

	// Initialize
	if err := config.Setup(); err != nil {
		log.Fatal().Msgf("main setup: %v", err)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	http2server := &http2.Server{}
	if err := http.ListenAndServe(":"+port, h2c.NewHandler(handler, http2server)); err != nil {
		log.Fatal().Msgf("main: %v", err)
	}
}

// ProxyHTTPGCS is the entry point for the cloud function, providing a proxy that
// permits HTTP protocol usage of a GCS bucket's contents.
func ProxyHTTPGCS(output http.ResponseWriter, input *http.Request) {
	ctx := context.Background()
	// route HTTP methods to appropriate handlers.
	switch input.Method {
	case http.MethodGet:
		config.GET(ctx, output, input)
	case http.MethodHead:
		config.HEAD(ctx, output, input)
	case http.MethodOptions:
		config.OPTIONS(ctx, output, input)
	default:
		http.Error(output, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
