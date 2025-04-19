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
	"log/slog"
	"os"
)

func SetupLogging() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(group []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.LevelKey:
				a.Key = "severity"
				if level := a.Value.Any().(slog.Level); level == slog.LevelWarn {
					a.Value = slog.StringValue("WARNING")
				}
			case slog.MessageKey:
				a.Key = "message"
			case slog.TimeKey:
				a.Key = "timestamp"
			}
			return a
		},
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(jsonHandler))
}
