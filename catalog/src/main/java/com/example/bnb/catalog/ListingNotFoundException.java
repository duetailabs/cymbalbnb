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

package com.example.bnb.catalog;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.NOT_FOUND)
public class ListingNotFoundException extends RuntimeException {

    private String id;

    public ListingNotFoundException() {
        this("unknown");
    }

    public ListingNotFoundException(String id) {
        this(id, null);
    }

    public ListingNotFoundException(String id, Throwable cause) {
        super("Listing with ID='" + id + "' can't be found", cause);
        this.id = id;
    }

    public String getId() {
        return id;
    }
}