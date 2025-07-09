# CymbalB&B Frontend

This document provides a detailed overview of the Go-based frontend web server for the CymbalB&B application, located in the `frontend/` directory.

## Architecture Overview

The frontend is a web server written in Go. Its primary responsibilities are:
- Serving the HTML pages for the user interface, including the home page with all listings and detailed pages for each listing.
- Communicating with the backend catalog service via a REST API to fetch listing data.
- Handling user sessions and request tracking through middleware.
- Rendering dynamic content using Go's standard `html/template` package.

It is designed to be deployed as a containerized application, as defined in the `Dockerfile`.

## Folder Structure

The `frontend/` directory is structured as follows:

```
frontend/
├── main.go          # Main application entrypoint
├── handlers.go      # HTTP request handlers
├── config.go        # Application configuration management
├��─ listing.go       # Data models and backend communication
├── middleware.go    # Request/session middleware
├── telemetry.go     # Structured logging setup
├── go.mod           # Go module dependencies
├── Dockerfile       # Container build definition
├── _templates/      # HTML templates
│   ├── home.html
│   └── listing.html
│   └── ...
├── _static/         # Static assets (CSS, images)
│   └── styles/
│       └── styles.css
└── utils/           # Utility and helper functions
    ├── contextvalues.go
    ├── metadata.go
    ├── renderers.go
    └── utils.go
```

## File-by-File Breakdown

### Core Application Logic

*   **`main.go`**: The entry point for the web server. It initializes the configuration, sets up the server and logging, chains the middleware, and handles graceful shutdown on interrupt signals.
*   **`handlers.go`**: Defines the HTTP handlers for the application's routes.
    *   **`/` (Default)**: Fetches all listings from the catalog service and renders the `home.html` template.
    *   **`/listing/{id}`**: Fetches a single listing by its ID and renders the `listing.html` template.
    *   **`/loadgen`**: A utility endpoint for generating load, which calls the default handler with a customizable delay.
    *   **Dependencies**: `listing.go` (to fetch data), `utils/` (for template functions and context values), `_templates/` (for rendering).
*   **`config.go`**: Manages the application's configuration. It reads settings from environment variables and uses the Google Cloud metadata service to discover information like Project ID and the catalog service URI.
*   **`middleware.go`**: Implements HTTP middleware.
    *   **`SessionIDMiddleware`**: Manages a session cookie (`cymbal-bnb-session-id`) to track user sessions across requests.
    *   **`RequestIDMiddleware`**: Injects a unique request ID into the context of each incoming request for traceability.
    *   **Dependencies**: `utils/contextvalues.go`, `github.com/google/uuid`.

### Backend Communication & Data Models

*   **`listing.go`**:
    *   **Function**: Defines the data structures (`Listing`, `Image`, `ListingCategory`) that mirror the JSON response from the backend catalog service. It contains the core logic for making HTTP GET requests to the backend to fetch all listings or a single listing by ID.
    *   **Local Debugging**: Includes a feature flag (`LOCAL_DEBUGGING`) that, when enabled, serves a hardcoded list of `DebugListings` instead of calling the real backend service.
    *   **Dependencies**: `utils/utils.go` (for `RestCall`).

### Utilities (`utils/` package)

*   **`utils.go`**: Contains generic helper functions. The most critical is `RestCall`, which creates an authenticated HTTP request to a Google Cloud service, adding the necessary `Authorization` header with an ID token fetched from the metadata server.
*   **`metadata.go`**: A client for the Google Cloud metadata service. It fetches environment details like Project ID, Project Number, Region, and ID tokens required to securely call other Google Cloud services.
*   **`renderers.go`**: Provides helper functions that are passed to the HTML templates, such as `RenderMoney` to format currency and `FirstTwoWords` to truncate strings.
*   **`contextvalues.go`**: Defines the keys and accessor functions for retrieving the session and request IDs from the request context.

### Presentation Layer

*   **`_templates/`**: This directory holds all the HTML templates. They use Go's template syntax (`{{ . }}`) to render dynamic data passed from the handlers.
    *   `header.html` and `footer.html` are common partials included in other pages.
    *   `home.html` displays the grid of all listings.
    *   `listing.html` shows the detailed view for a single listing.
    *   `error.html` is a generic page for displaying server errors.
*   **`_static/`**: Contains all static assets, primarily the `styles.css` stylesheet and application icons and images. These are served directly by the web server.

### Other Files

*   **`go.mod` / `go.sum`**: Standard Go module files that define the project's dependencies.
*   **`Dockerfile`**: A multi-stage Dockerfile that builds the Go application into a minimal, production-ready container image from a `scratch` base.
*   **`telemetry.go`**: Configures the `slog` structured logger to output JSON-formatted logs with severity levels and field names that are automatically recognized by Google Cloud Logging.

## Libraries and Dependencies

The frontend service relies on the following key Go modules, as defined in `go.mod`:

*   **`cloud.google.com/go/compute/metadata`**: The official Google Cloud client for interacting with the metadata server to fetch instance and project information.
*   **`github.com/google/uuid`**: Used to generate new UUIDs for session and request IDs.
*   **Standard Library**: The application makes extensive use of Go's standard library, including `net/http` (for the web server), `html/template` (for rendering), `log/slog` (for logging), and `encoding/json` (for processing API responses).
