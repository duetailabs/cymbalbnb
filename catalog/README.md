# Catalog service

Catalog service is implemented in Java using Spring Boot framework. It exposes two APIs:

- `listing` -- to return the full list of BnB locations.
- `listing/{id}` -- to return information for a listing by ID.
- `resetcache` -- to reset cache and force reading BnB locations from the database.

The service reads data from Cloud SQL Postgres database named "bnbcatalog".
The service does not initialize or validates database or its schema.
The service assumes that the default user "postgres" has read access to the data of the database.

## Invoking the service

The service requires authentication i.e. a valid ID token of the identity that is authorize to invoke this Cloud Run service.
To call the service from another service, grant the `roles/run.invoker` role to the service account assigned to the _calling_ service.
Use the following CLI command to reset service cache using `gcloud` when the user you logged in has the `roles/run.invoker` role on the project that hosting the catalog service:

```shell
SERVICE_URL=$(gcloud run services list --format='value(URL)' --filter='SERVICE:"bnb-catalog"'
curl -v -X GET https://{$SERVICE_URL}/resetcache -H "Authorization: bearer $(gcloud auth print-identity-token)" -H "Content-Type: application/json"
```

## Deployment configuration

The service is designed to run on Cloud Run. When deploying the service the following environment variables should be set up to the main container of the Cloud Run service:

- `STATIC_BUCKET_ID` -- defines the name of the bucket (without "gs://" prefix!) where images and videos of the BnB locations are stored.
- `SQL_INSTANCE_NAME` -- the Cloud SQL instance name of the PostgreSQL instance.
- `DB_NAME` -- the name of the database that stores the listings data. If you follow the instructions below this value should be `bnbcatalog`.
- `DB_PWD` -- the password to the `postgres` user of the database. It is recommended to store the password in the Secret Manager and inject the value using [Cloud Run secrets](https://cloud.google.com/run/docs/configuring/services/secrets)

The service assumes that the Postgres DB name and schema are fixed and provisioned. Run the following commands in psql CLI to create the database and the schema:

1. Create a database:

   ```sql
   CREATE DATABASE bnbcatalog;
   ```

1. Connect to database (you will need a password):

   ```shell
   \c bnbcatalog
   ```

1. Create the `listing` schema to store BnB locations:

   ```sql
   CREATE TABLE Listing (
       id VARCHAR PRIMARY KEY,
       name VARCHAR,
       description TEXT,
       price FLOAT,
       categories VARCHAR[],
       front_picture_uri VARCHAR,
       images JSONB[]
   );
   ```
   
