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
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv); err != nil {
		slog.Error("invalid run termination", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context, getenv func(string) string) error {
	SetupLogging()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	config := NewConfig(ctx, getenv)
	srv := NewServer(config)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: ChainMiddleware(srv, SessionIDMiddleware, RequestIDMiddleware),
	}
	go func() {
		slog.InfoContext(ctx, "starting web server", slog.String("address", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "error listening and serving", slog.String("error", err.Error()))
		}
	}()

	// wait for interrupt and gracefully stop web server after 5 sec
	<-ctx.Done()
	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 5*time.Second)
	defer cancel()
	return httpServer.Shutdown(shutdownCtx)
}
