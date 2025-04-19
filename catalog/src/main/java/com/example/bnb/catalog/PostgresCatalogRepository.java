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

import java.io.IOException;
import java.sql.Array;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.time.Duration;
import java.time.Instant;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import javax.sql.DataSource;

import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Repository;
import org.slf4j.Logger;

import com.google.common.collect.ImmutableList;
import com.google.common.collect.ImmutableMap;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;

@Repository
public class PostgresCatalogRepository implements CatalogRepository {
    private static final Logger LOGGER = LoggerFactory.getLogger(CatalogRepository.class);
    private static final int DEFAULT_AGE_SECS = 60 * 60 * 24; // default age 1 day
    private final DataSource connPool;
    private final String uriPrefix;
    private ImmutableMap<String, Listing> cachedListing;
    private Instant cachedTimestamp;
    private Duration maxCacheAge = Duration.ofSeconds(DEFAULT_AGE_SECS); 

    public PostgresCatalogRepository() throws IOException {
        try {
            this.connPool = CloudSQLConnectionPoolFactory.createConnectionPool();
        } catch (IOException ex) {
            LOGGER.error("Failed to initialize connection pool: " + ex.toString());
            throw ex;
        }
        this.cachedTimestamp = Instant.MIN;
        String bucketId = System.getenv().getOrDefault("STATIC_BUCKET_ID", System.getenv("K_SERVICE"));
        this.uriPrefix = "https://storage.googleapis.com/" + bucketId;
        LOGGER.debug("@@@@@@@@@@@@@@@@@@@@@@@ Prefix is initialized to '" + this.uriPrefix + "'");
    }

    private void refreshCached() {
        Instant now = Instant.now();
        if (Duration.between(cachedTimestamp, now).compareTo(maxCacheAge) > 0) {
            Map<String, Listing> data = new HashMap<String, Listing>();
            try {
                ResultSet rs = connPool.getConnection()
                        .prepareStatement(
                                "select id,name,location,description,price,categories,front_picture_uri,images from listing")
                        .executeQuery();
                // print the results to the console
                while (rs.next()) {
                    Listing entry = new Listing();
                    entry.setId(rs.getString("id"));
                    entry.setName(rs.getString("name"));
                    entry.setLocation(rs.getString("location"));
                    entry.setDescription(rs.getString("description"));
                    entry.setPrice(rs.getFloat("price"));
                    Array ar = rs.getArray("categories");
                    String[] strings = (String[]) ar.getArray();
                    entry.setCategories(toCategories(strings));
                    entry.setFrontPicture(uriPrefix + rs.getString("front_picture_uri"));
                    ar = rs.getArray("images");
                    strings = (String[]) ar.getArray();
                    entry.setImages(toImages(strings));
                    data.put(entry.getId(), entry);
                }
            } catch (SQLException ex) {
                LOGGER.error("Failed to refresh cache: " + ex.toString());
            }
            LOGGER.debug("@@@@@@@@@@@@@@@@@@@@@@@ Cache refreshed with " + data.size() + " records.");
            cachedListing = ImmutableMap.copyOf(data);
            cachedTimestamp = now;
        }
    }

    private Collection<ListingCategory> toCategories(String[] categories) {
        ImmutableList.Builder<ListingCategory> builder = ImmutableList.builder();
        for (String s : categories) {
            try {
                builder.add(ListingCategory.valueOf(s));
            } catch (IllegalArgumentException ex) {
                LOGGER.error("Failed to convert '" + s + "' to ListingCategory: " + ex.toString());
            }
        }
        return builder.build();
    }

    private Collection<Image> toImages(String[] images) {
        ImmutableList.Builder<Image> builder = ImmutableList.builder();
        for (String s : images) {
            JsonObject jo = (JsonObject) JsonParser.parseString(s);
            Image image = new Image(uriPrefix + jo.get("uri").getAsString(), jo.get("label").getAsString());
            LOGGER.debug("@@@@@@@@@@@@@@@@@@@@@@@ Converting '" + s + "' to Image object with uri='" + image.getURI()
                    + "' label='" + image.getLabel() + "'");
            builder.add(image);
        }
        return builder.build();
    }

    @Override
    public Listing findById(String id) {
        refreshCached();
        return cachedListing.get(id);
    }

    @Override
    public List<Listing> listAll() {
        refreshCached();
        return cachedListing.values().asList();
    }

    @Override
    public void setMaxCacheAge(long seconds) {
        this.maxCacheAge = Duration.ofSeconds(seconds);
    }

    @Override
    public void resetCache(int seconds) {
        if (seconds > 0) {
            maxCacheAge = Duration.ofSeconds(seconds);
        }
        cachedListing = ImmutableMap.of();
        cachedTimestamp = Instant.MIN;
        LOGGER.debug("@@@@@@@@@@@@@@@@@@@@@@@ cache was reset to " + Integer.toString(seconds));
    }
}
