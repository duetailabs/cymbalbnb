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

package com.example.bnb.video;

import java.io.IOException;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import com.example.bnb.video.generator.VideoGenerationClient;
import com.example.bnb.video.generator.types.GenerateVideoConfig;
import com.example.bnb.video.generator.types.Image;
import com.example.bnb.video.generator.types.Video;
import com.example.bnb.video.utils.Metadata;
import com.google.cloud.vertexai.VertexAI;
import com.google.cloud.vertexai.api.Content;
import com.google.cloud.vertexai.api.FileData;
import com.google.cloud.vertexai.api.Part;
import com.google.cloud.vertexai.generativeai.ContentMaker;
import com.google.cloud.vertexai.generativeai.GenerativeModel;
import com.google.cloud.vertexai.generativeai.ResponseHandler;

import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;

@RestController
public class VideoController {
        private static final Logger LOGGER = LoggerFactory.getLogger(VideoController.class);
        private static final String PROMPT_GENERATION_INSTRUCTIONS = """
                        generate instructions for the Veo 2 video from image gen AI model.
                        start instructions with 'create a short 3D clip in realistic style'.
                        use the description of the property to generate step-by-step instructions.
                        the instructions begin with a general view of the image.
                        then the instructions should describe zoom in and walkthrough the property.
                        instructions should create a 8 second clip with 16:9 aspect ratio.
                        do not include explanations how instructions were generated.
                        do not include title.
                        do not use markdown.
                        description: %s""";
        private VertexAI vertexAI;
        private GenerativeModel model;
        private VideoGenerationClient videoClient;

        @PostConstruct
        public void init() throws IOException {
                var projectId = Metadata.getProjectId();
                var region = Metadata.getRegion();
                vertexAI = new VertexAI(projectId, region);
                model = new GenerativeModel("gemini-2.0-flash-001", vertexAI);
                videoClient = new VideoGenerationClient(projectId, region, "veo-2.0-generate-001");
        }

        @PostMapping(path = "/newvideo", consumes = MediaType.APPLICATION_JSON_VALUE, produces = MediaType.APPLICATION_JSON_VALUE)
        public ResponseEntity<Video> create(@RequestBody CreateParameters payload) {
                if (payload.imageUris() == null || payload.imageUris().length == 0
                                || payload.id().isEmpty() || payload.storageBucket().isEmpty()) {
                        LOGGER.atWarn()
                                        .addKeyValue("listing_id", payload.id())
                                        .addKeyValue("images", payload.imageUris())
                                        .addKeyValue("gcs_bucket", payload.storageBucket())
                                        .log("invalid request parameter(s)");
                        return new ResponseEntity<>(HttpStatus.BAD_REQUEST);
                }
                LOGGER.atDebug()
                                .addKeyValue("listing_id", payload.id())
                                .addKeyValue("images", payload.imageUris())
                                .addKeyValue("gcs_bucket", payload.storageBucket())
                                .log("new video request");
                try {
                        // https://cloud.google.com/vertex-ai/generative-ai/docs/multimodal/image-understanding
                        var imageData = FileData.newBuilder()
                                        .setFileUri(payload.imageUris()[0]).setMimeType(MediaType.IMAGE_JPEG_VALUE)
                                        .build();
                        var content = Content.newBuilder()
                                        .addParts(Part.newBuilder().setText("describe what is in the image").build())
                                        .addParts(Part.newBuilder().setFileData(imageData).build())
                                        .setRole("user")
                                        .build();
                        var response = model.generateContent(content);
                        var description = ResponseHandler.getText(response);
                        LOGGER.atDebug()
                                        .addKeyValue("description", description)
                                        .addKeyValue("image", payload.imageUris()[0])
                                        .log("image description generated");
                        // generate Veo2 prompt
                        content = ContentMaker
                                        .fromString(PROMPT_GENERATION_INSTRUCTIONS.formatted(description));
                        response = model.generateContent(content);
                        description = ResponseHandler.getText(response);
                        LOGGER.atDebug()
                                        .addKeyValue("prompt", description)
                                        .log("prompt for creating video generated");
                        // https://cloud.google.com/vertex-ai/generative-ai/docs/model-reference/veo-video-generation
                        var outputUri = "%s/videos/%s/".formatted(payload.storageBucket(), payload.id());
                        var op = videoClient.generateVideos(
                                        description,
                                        new Image(payload.imageUris()[0], MediaType.IMAGE_JPEG_VALUE),
                                        new GenerateVideoConfig(1, outputUri));
                        while (op.done() == null || !op.done()) {
                                op = videoClient.getStatus(op.name());
                        }
                        List<Video> videos = op.response().videos();
                        if (videos.size() > 0) {
                                return ResponseEntity.ok(videos.get(0));
                        }
                        return new ResponseEntity<>(HttpStatus.I_AM_A_TEAPOT);
                } catch (IOException ex) {
                        LOGGER.atError()
                                        .addKeyValue("error", ex)
                                        .log("fail to generate video");
                        return new ResponseEntity<>(HttpStatus.INTERNAL_SERVER_ERROR);
                }
        }

        @PreDestroy
        public void destroy() {
                vertexAI.close();
        }

        
}
