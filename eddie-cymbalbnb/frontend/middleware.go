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
	"log/slog"
	"net/http"

	"github.com/minherz/bnb-demo/frontend/utils"
)

type MiddlewareHandler func(http.Handler) http.Handler

const (
	sessionCookie = "cymbal-bnb-session-id"
)

func requestSessionID(r *http.Request) string {
	c, err := r.Cookie(sessionCookie)
	if err == nil {
		return c.Value
	}
	if err != http.ErrNoCookie {
		slog.Error("failed to read cookie", slog.Any("error", err))
	}
	return ""
}

func ChainMiddleware(h http.Handler, chain ...MiddlewareHandler) http.Handler {
	if len(chain) < 1 {
		return h
	}
	wrapper := h
	for i := len(chain) - 1; i >= 0; i-- {
		wrapper = chain[i](wrapper)
	}
	return wrapper
}

func SessionIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sid := requestSessionID(r)
		if sid == "" {
			sid = utils.NewGUID()
			http.SetCookie(w, &http.Cookie{
				Name:   sessionCookie,
				Value:  sid,
				MaxAge: 60 * 60,
			})
		}
		ctx := context.WithValue(r.Context(), utils.SessionKey, sid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := utils.NewGUID()
		ctx := context.WithValue(r.Context(), utils.RequestKey, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
