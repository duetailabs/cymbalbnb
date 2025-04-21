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
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"maps"
	"net/http"
	"strconv"
	"time"

	"github.com/minherz/bnb-demo/frontend/utils"
)

type FrontendServer struct {
	config *Config
}

var (
	footerSetting = map[string]any{
		"currentYear": time.Now().Year(),
		"service":     "unknown",
		"location":    "unknown",
	}
	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"renderMoney": utils.RenderMoney,
		"maxImages":   utils.RangeMax[Image],
		"twoWords":    utils.FirstTwoWords,
	}).ParseGlob("_templates/*.html"))
)

func NewServer(config *Config) http.Handler {
	footerSetting["service"] = config.ServiceName
	footerSetting["location"] = config.Location

	srv := &FrontendServer{config: config}
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.StaticPath))))
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "User-agent: *\nDisallow: /") })
	mux.HandleFunc("/_alive", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "OK") })
	mux.HandleFunc("/_ready", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "OK") })
	mux.HandleFunc("/", srv.Default)
	mux.HandleFunc("/listing/{id}", srv.Listing)
	mux.HandleFunc("/loadgen", srv.Loadgen)

	return mux
}

func (s *FrontendServer) Default(w http.ResponseWriter, r *http.Request) {
	listings, err := listings(r.Context(), s.config.CatalogServiceURI)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, err)
		return
	}

	for _, listing := range listings {
		if listing.Price > 250.0 {
			listing.featured = true
		}
	}

	slog.Info("service listings", slog.String("listing", fmt.Sprintf("%+v", listings)))
	data := map[string]interface{}{
		"sessionID":   utils.Session(r.Context()),
		"requestID":   utils.Request(r.Context()),
		"listings":    listings,
		"titleSuffix": " - All Listings",
	}
	maps.Copy(data, footerSetting)
	err = templates.ExecuteTemplate(w, "home", data)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, err)
	}
}

func (s *FrontendServer) Listing(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		RenderError(w, http.StatusBadRequest, errors.New("listing id is not specified"))
		return
	}
	listing, err := listing(r.Context(), s.config.CatalogServiceURI, id)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, fmt.Errorf("listing for id# %s does not exist", id))
		return
	}
	slog.Info("service listing", slog.String("id", id), slog.String("listing", fmt.Sprintf("%+v", listing)))
	data := map[string]interface{}{
		"sessionID":   utils.Session(r.Context()),
		"requestID":   utils.Request(r.Context()),
		"listing":     listing,
		"titleSuffix": " - ",
	}
	maps.Copy(data, footerSetting)
	if err := templates.ExecuteTemplate(w, "listing", data); err != nil {
		RenderError(w, http.StatusInternalServerError, err)
	}
}

func (s *FrontendServer) Loadgen(w http.ResponseWriter, r *http.Request) {
	delayms := int64(3000)
	argms := r.URL.Query().Get("delayms")
	args := r.URL.Query().Get("delay")
	if args != "" {
		sec, err := strconv.ParseInt(args, 10, 64)
		if err == nil {
			delayms = sec * 1000
		}
	} else if argms != "" {
		ms, err := strconv.ParseInt(argms, 10, 64)
		if err == nil {
			delayms = ms
		}
	}
	defer customizedDelay(time.Now(), delayms)
	s.Default(w, r)
}

func RenderError(w http.ResponseWriter, httpCode int, err error) {
	msg := fmt.Sprintf("%+v", err)
	w.WriteHeader(httpCode)
	slog.Error("render error to client", slog.String("error", msg), slog.Int("code", httpCode))
	if err2 := templates.ExecuteTemplate(w, "error", map[string]any{
		"error":       msg,
		"status_code": httpCode,
		"status":      http.StatusText(httpCode),
	}); err2 != nil {
		slog.Error("failed to render error template", slog.Any("error", err2))
	}
}

func customizedDelay(start time.Time, delay int64) {
	diff := time.Since(start)
	slog.Debug("Custimize delay time", slog.Any("start", start), slog.Any("end", time.Now()), slog.Int64("delay", delay), slog.Int64("execution_time", diff.Milliseconds()))
	if diff.Milliseconds() < delay {
		slog.Debug("Additional delay", slog.Any("time", time.Duration(delay-diff.Milliseconds())*time.Millisecond))
		time.Sleep(time.Duration(delay-diff.Milliseconds()) * time.Millisecond)
	}
}
