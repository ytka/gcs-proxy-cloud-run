#!/usr/bin/env bash
# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
set -euf -o pipefail

function usage(){
    echo >&2
    echo "Usage: $0 image_name project" >&2
    echo "Builds this service container in Google Container Registry (gcr.io)." >&2
    echo >&2
    echo "image_name is optional. default value is gcs-streaming-proxy." >&2
    echo >&2
    echo "project is optional; your gcloud config project will be used if it" >&2
    echo "is not provided." >&2
    echo >&2
}

IMAGE_NAME="${1:-stattc-ip}"
PROJECT="${2:-$(gcloud config get-value project 2>/dev/null)}"
TAG=gcr.io/"${PROJECT}"/"${IMAGE_NAME}"

if [[ -z "$PROJECT" ]]; then
    echo >&2 "ERROR: Could not determine project. Please specify it explicitly."
    usage
    exit 2
fi

# quick and dirty way to catch if the user asks for help, like --help
# downside: you can't tag the image as *help or just "-h"
if [[ "${IMAGE_NAME}" == *help ]] || [[ "${IMAGE_NAME}" == "-h" ]]; then
    usage
    exit
fi

gcloud --project="${PROJECT}" builds submit --tag "${TAG}"

echo Container image built:
echo "${TAG}" 