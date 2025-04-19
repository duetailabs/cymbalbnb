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

import java.io.Serializable;
import java.util.Arrays;
import java.util.Collection;

import com.google.common.collect.ImmutableList;

public class Listing implements Serializable {
    private String id;
    private String name, description, location, frontPictureUri;
    private Collection<ListingCategory> categories;
    private Float price;
    private Collection<Image> images;

    public Listing() {
    }

    public Listing(String id, String name, String location, String description, ListingCategory category, float price,
            String frontPicture, Image... images) {
        setId(id);
        setName(name);
        setLocation(location);
        setDescription(description);
        setCategories(ImmutableList.of(category));
        setFrontPicture(frontPicture);
        setImages(Arrays.asList(images));
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getLocation() {
        return name;
    }

    public void setLocation(String location) {
        this.location = location;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Float getPrice() {
        return price;
    }

    public void setPrice(Float price) {
        this.price = price;
    }

    public Collection<ListingCategory> getCategories() {
        return categories;
    }

    public void setCategories(Collection<ListingCategory> categories) {
        this.categories = categories;
    }

    public String getFrontPicture() {
        return frontPictureUri;
    }

    public void setFrontPicture(String frontPicture) {
        this.frontPictureUri = frontPicture;
    }

    public Collection<Image> getImages() {
        return images;
    }

    public void setImages(Collection<Image> images) {
        this.images = images;
    }
}
