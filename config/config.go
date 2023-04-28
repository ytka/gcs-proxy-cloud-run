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
package config

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ytka/gcs-proxy-cloud-run/backends/gcs"
	"github.com/ytka/gcs-proxy-cloud-run/backends/proxy"
	"github.com/ytka/gcs-proxy-cloud-run/backends/token"
	"github.com/ytka/gcs-proxy-cloud-run/common"
)

// Setup will be called once at the start of the program.
func Setup() (*token.TokenClient, error) {
	tokenClient, err := token.NewClient()
	if err != nil {
		return nil, err
	}
	return tokenClient, gcs.Setup()
}

// GET will be called in main.go for GET requests
func GET(ctx context.Context, output http.ResponseWriter, input *http.Request, tokenClient *token.TokenClient) {
	bucket, objectName := common.NormalizePath(input.URL.Path)
	tk, err := tokenClient.GetToken(ctx, bucket, objectName)
	if err != nil {
		// TODO: treat 404
		http.Error(output, "can't get token: "+err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(tk)
	gcs.Read(ctx, output, input, LoggingOnly)
	//gcs.ReadWithCache(ctx, output, input, CacheMedia, cacheGetter, LoggingOnly)
}

// HEAD will be called in main.go for HEAD requests
func HEAD(ctx context.Context, output http.ResponseWriter, input *http.Request) {
	gcs.ReadMetadata(ctx, output, input, LoggingOnly)
}

// func POST

// func DELETE

// OPTIONS will be called in main.go for OPTIONS requests
func OPTIONS(ctx context.Context, output http.ResponseWriter, input *http.Request) {
	proxy.SendOptions(ctx, output, input, LoggingOnly)
}

func POST(ctx context.Context, output http.ResponseWriter, input *http.Request, tokenClient *token.TokenClient) {
	bucket, objectName := common.NormalizePath(input.URL.Path)
	token, err := tokenClient.CreateToken(ctx, bucket, objectName)
	if err != nil {
		http.Error(output, "can't create token: "+err.Error(), http.StatusInternalServerError)
	}
	io.WriteString(output, fmt.Sprintf("token: %s", token))
}
