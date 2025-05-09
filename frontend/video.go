// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/minherz/bnb-demo/frontend/utils"
)

type videoRequest struct {
	ListingID     string   `json:"id"`
	ImageURIs     []string `json:"imageUris"`
	StorageBucket string   `json:"storageBucket"`
}

type videoResponse struct {
	ListingID string `json:"id"`
	URI       string `json:"uri"`
}

func video(ctx context.Context, serviceURI string, listing Listing) (string, error) {
	uri, err := url.JoinPath(serviceURI, "/video/"+listing.Id)
	if err != nil {
		return "", err
	}

	imageURIs := []string{}
	for _, image := range listing.Images {
		imageURIs = append(imageURIs, image.URI)
	}
	// all URIs are in the format https://storage.googleapis.com/<BUCKET>/<PATH_TO_FILE>
	var storageBucket, filePath string
	fmt.Scanf("https://storage.googleapis.com/%s/%s", &storageBucket, &filePath)

	video := videoRequest{
		ListingID:     listing.Id,
		ImageURIs:     imageURIs,
		StorageBucket: storageBucket,
	}
	payload, err := json.Marshal(video)
	if err != nil {
		return "", err
	}
	data, err := utils.RestCall(ctx, uri, http.MethodGet, payload)
	if err != nil {
		return "", err
	}
	var v videoResponse
	if err := json.Unmarshal(data, &v); err != nil {
		return "", err
	}
	return v.URI, err
}
