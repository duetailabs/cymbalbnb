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

package com.example.bnb.video.generator;

import com.example.bnb.video.generator.types.*;
import com.example.bnb.video.utils.Metadata;

import java.io.Serializable;
import java.nio.charset.StandardCharsets;
import java.util.List;

import org.springframework.http.MediaType;
import org.springframework.web.reactive.function.client.WebClient;

public class VideoGenerationClient {
    private final static String AIPLATFORM_URL = "https://%s-aiplatform.googleapis.com/";
    private final static String VIDEO_GENERATION_URI = "/v1/projects/%s/locations/%s/publishers/google/models/%s:";
    private final static String GENERATE_OP = "predictLongRunning";
    private final static String FETCH_OP = "fetchPredictOperation";
    private final String videoGenerationUri;
    private final String baseUri;

    public VideoGenerationClient(String projectId, String location, String model) {
        this.baseUri = AIPLATFORM_URL.formatted(location);
        this.videoGenerationUri = VIDEO_GENERATION_URI.formatted(location, projectId, location, model);
    }

    public Operation generateVideos(String prompt, Image image, GenerateVideoConfig config) {
        var webClient = newClient();
        var body = new PredictLongRunningRequest(
                java.util.Collections.singletonList(new Instance(prompt, image)),
                config);
        return webClient.post()
                .uri(videoGenerationUri + GENERATE_OP)
                .bodyValue(body)
                .retrieve()
                .bodyToMono(Operation.class)
                .block();
    }

    public Operation getStatus(String opName) {
        var webClient = newClient();
        var body = new Operation(opName, null, null);
        return webClient.post()
                .uri(videoGenerationUri + FETCH_OP)
                .bodyValue(body)
                .retrieve()
                .bodyToMono(Operation.class)
                .block();
    }

    private WebClient newClient() {
        return WebClient.builder()
                .baseUrl(baseUri)
                .defaultHeaders(httpHeaders -> {
                    httpHeaders.setAccept(java.util.Collections.singletonList(MediaType.APPLICATION_JSON));
                    httpHeaders.setAcceptCharset(java.util.Collections.singletonList(StandardCharsets.UTF_8));
                    httpHeaders.setContentType(MediaType.APPLICATION_JSON);
                    httpHeaders.setBearerAuth(Metadata.getIdToken(videoGenerationUri));
                }).build();
    }

    public record PredictLongRunningRequest(
            List<Instance> instances,
            GenerateVideoConfig parameters) implements Serializable {
    }

    public record Instance(String prompt, Image image) implements Serializable {
    }
}
