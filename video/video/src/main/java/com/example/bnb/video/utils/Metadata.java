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

package com.example.bnb.video.utils;

import com.google.cloud.ServiceOptions;

import java.io.UnsupportedEncodingException;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;

import com.google.cloud.MetadataConfig;

public class Metadata {

    public static String getProjectId() {
        return ServiceOptions.getDefaultProjectId();
    }

    public static String getRegion() {
        String s = MetadataConfig.getAttribute("instance/region");
        // parse region name from "projects/PROJECT_ID/regions/REGION_ID"
        if (s != null) {
            String region = s.substring(s.lastIndexOf('/') + 1);
            if (region.length() > 0) {
                return region;
            }
            return s;
        }
        return "us-central1";
    }

    public static String getIdToken(String audience) {
        try {
            String encodedAudience = URLEncoder.encode(audience, StandardCharsets.UTF_8.toString());
            String s = MetadataConfig
                    .getAttribute("instance/service-accounts/default/identity?audience=" + encodedAudience);
            return s;
        } catch (UnsupportedEncodingException ex) {
            return "";
        }
    }
}
