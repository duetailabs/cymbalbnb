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
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

var getenvFunc func(string) string

func Setup(getenv func(string) string) {
	getenvFunc = getenv
}

func StringOnly(s string, err error) string {
	if err != nil {
		slog.Warn("StringOnly call hides error", slog.Any("error", err))
		return "unknown"
	}
	return s
}

func GetStringParam(flagName, evName, defaultValue string) string {
	v := *flag.String(flagName, "", "flag "+flagName)
	if v != "" {
		return v
	}
	v = getenvFunc(evName)
	if v != "" {
		return v
	}
	return defaultValue
}

func GetEnv(name, defaultValue string) string {
	v := getenvFunc(name)
	if v != "" {
		return v
	}
	return defaultValue
}

func NewGUID() string {
	v, _ := uuid.NewRandom()
	return v.String()
}

func RestCall(ctx context.Context, uri string, verb string, payload []byte) ([]byte, error) {
	r, err := http.NewRequestWithContext(ctx, verb, uri, bytes.NewReader(payload))
	if err != nil {
		return []byte{}, err
	}
	t, err := IDToken(ctx, uri)
	if err != nil {
		return []byte{}, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		data, _ := io.ReadAll(res.Body)
		return []byte{}, errors.New("rest call failed with " + strconv.Itoa(res.StatusCode) + ": " + string(data))
	}
	return io.ReadAll(res.Body)
}

func FormatServiceName(ctx context.Context, n, d string) string {
	if n == "" {
		return d
	}
	_, err := url.ParseRequestURI(n)
	if err == nil {
		return n
	}
	return fmt.Sprintf("https://%s-%s.%s.run.app", n, StringOnly(ProjectNumber(ctx)), StringOnly(Region(ctx)))
}
