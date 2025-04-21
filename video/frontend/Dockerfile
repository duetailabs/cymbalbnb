# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.23.4-alpine@sha256:c23339199a08b0e12032856908589a6d41a0dab141b8b3b21f156fc571a3f1d3 AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -mod=readonly -installsuffix 'static' -v -o frontend .

# Use empty image for a lean production container.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM scratch

WORKDIR /app

# Copy certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /src/frontend ./frontend
COPY ./_templates ./_templates
COPY ./_static ./_static

ENTRYPOINT ["/app/frontend"]