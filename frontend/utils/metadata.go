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
	if projectID == "" {
		var err error
		if projectID, err = metadata.ProjectIDWithContext(ctx); err != nil {
			return "", err
		}
	}
	return projectID, nil
}

func ProjectNumber(ctx context.Context) (string, error) {
	if projectNumber == "" {
		var err error
		if projectNumber, err = metadata.GetWithContext(ctx, "project/numeric-project-id"); err != nil {
			return "", err
		}
	}
	return projectNumber, nil
}

func Region(ctx context.Context) (string, error) {
	if region == "" {
		var err error
		if region, err = metadata.GetWithContext(ctx, "instance/region"); err != nil {
			return "", err
		}
		// parse region from fully qualified name projects/<projNum>/regions/<region>
		if pos := strings.LastIndex(region, "/"); pos >= 0 {
			region = region[pos+1:]
		}
	}
	return region, nil
}

func InstanceId(ctx context.Context) (string, error) {
	if instanceId == "" {
		var err error
		if instanceId, err = metadata.GetWithContext(ctx, "instance/id"); err != nil {
			return "", err
		}
	}
	return instanceId, nil
}

func IDToken(ctx context.Context, aud string) (string, error) {
	if accessToken == "" {
		var err error
		if accessToken, err = metadata.GetWithContext(ctx, "instance/service-accounts/default/identity?audience="+url.QueryEscape(aud)); err != nil {
			return "", err
		}
	}
	return accessToken, nil
}

func ServiceName() string {
	return GetEnv("K_SERVICE", "Unknown")
}

func RevisionName() string {
	return GetEnv("K_REVISION", "Unknown")
}
