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

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.bind.annotation.RequestParam;

@RestController
public class CatalogController {
    private CatalogRepository repository;

    @Autowired
    public CatalogController(CatalogRepository repository) {
        this.repository = repository;
    }

    @GetMapping(value = "/listing/{listingId}")
    public Listing findById(@PathVariable String listingId) {
        Listing found = repository.findById(listingId);
        if (found == null) {
            new ListingNotFoundException(listingId);
        }
        return found;
    }

    @GetMapping(value = "/listing")
    public List<Listing> all() {
        return repository.listAll();
    }

    @PostMapping(value = "/resetcache")
    public Boolean resetCache(@RequestParam(defaultValue = "-1") int age) {
        repository.resetCache(age);
        return true;
    }
}
