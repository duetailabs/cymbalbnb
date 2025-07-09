# CymbalB&B Catalog Backend

This document provides a detailed overview of the Java Spring Boot backend service for the CymbalB&B application, located in the `catalog/` directory.

## Architecture Overview

The backend is a monolithic Spring Boot application that serves as the catalog service for the frontend. It exposes a REST API to provide information about B&B listings. It connects to a PostgreSQL database on Google Cloud SQL for data persistence and implements an in-memory caching layer to improve performance.

## Folder Structure

The `catalog/` directory is structured as a standard Maven project:

```
catalog/
├── pom.xml      # Maven build and dependency configuration
└── src/
    └── main/
        ├── java/com/example/bnb/catalog/ # Java source code
        │   ├── CatalogService.java       # Main application class
        │   ├── CatalogController.java    # REST API controller
        │   ├── CatalogRepository.java    # Data access interface
        │   ├── PostgresCatalogRepository.java # PostgreSQL implementation of the repository
        │   ├── CloudSQLConnectionPoolFactory.java # Creates DB connection pool
        │   ├── Listing.java              # Data model for a listing
        │   ├── Image.java                # Data model for an image
        │   ├── ListingCategory.java      # Enum for listing categories
        │   ├── ListingNotFoundException.java # Custom exception for 404 errors
        │   └── LoggingEventGoogleCloudEncoder.java # Custom log encoder for Google Cloud
        └── resources/
            ├── application.properties    # Spring Boot configuration
            └── logback.xml               # Logging framework configuration
```

## File-by-File Breakdown

### Core Application Logic

*   **`CatalogService.java`**: This is the main entry point for the Spring Boot application. The `main` method bootstraps the application, setting the server port from the `PORT` environment variable (defaulting to `8080`).

*   **`CatalogController.java`**: This class defines the public REST API endpoints.
    *   **`GET /listing`**: Returns a list of all B&B listings.
    *   **`GET /listing/{listingId}`**: Returns a single listing by its unique ID.
    *   **`POST /resetcache`**: Clears the in-memory cache, forcing subsequent requests to fetch fresh data from the database.
    *   **Dependencies**: It depends on the `CatalogRepository` interface to fetch data.

### Data Access Layer

*   **`CatalogRepository.java`**: This is an interface that defines the contract for data access operations. It decouples the controller from the specific data source implementation.

*   **`PostgresCatalogRepository.java`**: This is the primary implementation of `CatalogRepository`.
    *   **Function**: It connects to the PostgreSQL database to retrieve listing data.
    *   **Caching**: It features a time-based in-memory cache (`ImmutableMap`) to store listings. Data is fetched from the database only if the cache is empty or has expired (default TTL is 24 hours). The `resetCache` method invalidates this cache.
    *   **Dependencies**:
        *   `CloudSQLConnectionPoolFactory`: To create the database connection pool.
        *   `Image.java`, `Listing.java`, `ListingCategory.java`: To construct the data objects from the database result set.
        *   `com.google.gson`: To parse JSON data stored in the database for image details.
        *   `com.google.common`: For immutable collections.

*   **`CloudSQLConnectionPoolFactory.java`**: A utility class responsible for creating a JDBC `DataSource`.
    *   **Function**: It configures and initializes a `HikariCP` connection pool to connect securely to a Google Cloud SQL for PostgreSQL instance using the Cloud SQL Socket Factory.
    *   **Configuration**: It reads database connection details (instance name, DB name, password) from environment variables.

### Data Models & Exceptions

*   **`Listing.java`**: A POJO (Plain Old Java Object) representing a single B&B listing, including its ID, name, price, images, etc.
*   **`Image.java`**: A POJO representing an image associated with a listing, containing a URI and a label.
*   **`ListingCategory.java`**: An `enum` that defines the possible categories for a listing (e.g., `HOUSE`, `APARTMENT`).
*   **`ListingNotFoundException.java`**: A custom `RuntimeException` that is thrown when a requested listing ID does not exist. It is annotated with `@ResponseStatus(HttpStatus.NOT_FOUND)`, which causes Spring to return a 404 HTTP status code.

### Configuration and Utilities

*   **`pom.xml`**: The Maven Project Object Model file. It defines project metadata, dependencies, and build configurations.
*   **`application.properties`**: The Spring Boot configuration file. It is used to disable the default database auto-configuration, as a custom `DataSource` is created in `CloudSQLConnectionPoolFactory`.
*   **`logback.xml`**: The configuration file for the Logback logging framework. It sets up a custom encoder.
*   **`LoggingEventGoogleCloudEncoder.java`**: A custom Logback `Encoder` that formats log messages into a structured JSON format. This is optimized for ingestion and analysis in Google Cloud Logging.
    *   **Dependencies**: `ch.qos.logback`, `com.google.gson`, `org.joda.time`.

## Libraries and Dependencies

The backend service relies on the following key libraries, as defined in `pom.xml`:

*   **Spring Boot Starter Web**: Core framework for building the web application.
*   **Google Cloud SQL Postgres Socket Factory**: For secure, authorized connections to Cloud SQL.
*   **PostgreSQL JDBC Driver**: The standard driver for connecting to a PostgreSQL database.
*   **HikariCP**: A high-performance JDBC connection pool.
*   **Google Cloud Libraries BOM**: Manages versions for Google Cloud client libraries.
*   **Google Guava**: Provides core utility libraries, used here for immutable collections.
*   **Google Gson**: A library for converting Java Objects into their JSON representation and vice-versa.
*   **Joda-Time**: A library for date and time manipulation, used in the custom logging encoder.
*   **Logback**: A logging framework for Java.