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

import javax.sql.DataSource;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.cloud.MetadataConfig;
import com.google.cloud.ServiceOptions;
import com.google.common.base.Strings;
import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;

import java.io.IOException;

public class CloudSQLConnectionPoolFactory {
    private static final Logger LOGGER = LoggerFactory.getLogger(CloudSQLConnectionPoolFactory.class);
    private static final String DEFAULT_DB_USER = "postgres";
    private static final String DEFAULT_DB_NAME = "bnbcatalog";
    private static final String INSTANCE_CONNECTION_NAME = "%s:%s:%s";

    public static DataSource createConnectionPool() throws IOException {
        String projectId = ServiceOptions.getDefaultProjectId();
        String dbName = System.getenv().getOrDefault("DB_NAME", DEFAULT_DB_NAME);
        String secret = System.getenv("DB_PWD");
        String instanceName = System.getenv("SQL_INSTANCE_NAME");

        if (Strings.isNullOrEmpty(projectId) || Strings.isNullOrEmpty(secret) || Strings.isNullOrEmpty(instanceName)) {
            throw new IOException("Missing information about project ID, instance name or password");
        }
        String instanceConnectionName = String.format(INSTANCE_CONNECTION_NAME, projectId, getRegion(), instanceName);
        LOGGER.debug(
                "@@@@@@@@@@@@@@@@@@@@@@@ DB instance connection is initialized to '" + instanceConnectionName + "'");

        // create a new configuration and set the database credentials
        HikariConfig config = new HikariConfig();
        config.setJdbcUrl(String.format("jdbc:postgresql:///%s", dbName));
        config.setMinimumIdle(3);
        config.setMaximumPoolSize(10);
        config.setConnectionTimeout(3000);
        config.setIdleTimeout(60000);
        config.setMaxLifetime(60000);
        config.setUsername(DEFAULT_DB_USER);
        config.setPassword(secret);
        config.addDataSourceProperty("socketFactory", "com.google.cloud.sql.postgres.SocketFactory");
        config.addDataSourceProperty("cloudSqlInstance", instanceConnectionName);

        // Initialize the connection pool using the configuration object.
        return new HikariDataSource(config);
    }

    private static String getRegion() {
        String s = MetadataConfig.getAttribute("instance/region");
        // parse region name from "projects/PROJECT_ID/regions/REGION_ID"
        if (s != null) {
            String region = s.substring(s.lastIndexOf('/') + 1);
            if (region.length() > 0) {
                return region;
            }
        }
        return "us-central1";
    }
}