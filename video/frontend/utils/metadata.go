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

package utils

import (
	"context"
	"net/url"
	"strings"

	"cloud.google.com/go/compute/metadata"
)

var (
	accessToken   string
	instanceId    string
	projectID     string
	projectNumber string
	region        string
)

func ProjectID(ctx context.Context) (string, error) {
	if projectID != "" {
		return projectID, nil
	}
	projectID = GetEnv("DEBUG_PROJECT_ID", "")
	if projectID != "" {
		return projectID, nil
	}
	var err error
	projectID, err = metadata.ProjectIDWithContext(ctx)
	return projectID, err
}

func ProjectNumber(ctx context.Context) (string, error) {
	if projectNumber != "" {
		return projectNumber, nil
	}
	projectNumber = GetEnv("DEBUG_PROJECT_NUMBER", "")
	if projectNumber != "" {
		return projectNumber, nil
	}
	var err error
	projectNumber, err = metadata.GetWithContext(ctx, "project/numeric-project-id")
	return projectID, err
}

func Region(ctx context.Context) (string, error) {
	if region != "" {
		return region, nil
	}
	region = GetEnv("DEBUG_REGION", "")
	if region != "" {
		return region, nil
	}
	var err error
	region, err = metadata.GetWithContext(ctx, "instance/region")
	// parse region from fully qualified name projects/<projNum>/regions/<region>
	if pos := strings.LastIndex(region, "/"); err == nil && pos >= 0 {
		region = region[pos+1:]
	}
	return region, err
}

func InstanceId(ctx context.Context) (string, error) {
	if instanceId != "" {
		return instanceId, nil
	}
	instanceId = GetEnv("DEBUG_INSTANCE_ID", "")
	if instanceId != "" {
		return instanceId, nil
	}
	var err error
	instanceId, err = metadata.GetWithContext(ctx, "instance/id")
	return instanceId, err
}

func IDToken(ctx context.Context, aud string) (string, error) {
	if accessToken != "" {
		return accessToken, nil
	}
	accessToken = GetEnv("DEBUG_ACCESS_TOKEN", "")
	if accessToken != "" {
		return accessToken, nil
	}
	var err error
	accessToken, err = metadata.GetWithContext(ctx, "instance/service-accounts/default/identity?audience="+url.QueryEscape(aud))
	return accessToken, err
}

func ServiceName() string {
	return GetEnv("K_SERVICE", "Unknown")
}

func RevisionName() string {
	return GetEnv("K_REVISION", "Unknown")
}
