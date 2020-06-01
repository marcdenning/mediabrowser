# Media Browser

> Web-based media browser written in Go.

## Overview

Media Browser is a directory browser for Google Cloud Storage.
Similar to common web servers' directory and file browsing capabilities, it allows user to browse a storage bucket via a simple web interface using hyperlinks.

Once a user selects an individual file, the server redirects the user to the Signed URL for the object in Cloud Storage.
Since the intention of this app is browsing media, this feature reduces the runtime of the server itself and allow video media to be streamed directly from the storage bucket.
This kind of optimization makes Media Browser a perfect fit for serverless platforms like Cloud Run.

The interface and HTTP responses are modeled after the Apache HTTPD directory browsing feature.
This enables better compatibility with clients that are looking for that particular layout (ex. Kodi).

## Configuration

### Environment Variables

* `GOOGLE_APPLICATION_CREDENTIALS` _(required)_ - path to Google Cloud service account key file (JSON file); the service account must have Cloud Storage object reader and Secret Manager Secret Accessor roles
* `BUCKET_NAME` _(required)_ - name of the Google Cloud Storage bucket to browse
* `SERVICE_ACCOUNT_NAME` _(required)_ - name of the IAM service account used to access Cloud Storage and Secret Manager
* `PK_SECRET_NAME` _(required)_ - resource ID of the secret storing the service account's private key
* `WEB_USERNAME` _(optional)_ - username to expect for HTTP Basic Authentication
* `WEB_PASSWORD` _(optional)_ - password to expect for HTTP Basic Authentication

### Google Cloud Setup

1. Create a project for your app.
2. Enable Cloud Storage and create a bucket to serve media. No special permissions are required for the bucket.
3. Create a service account to use for the app. Securely store the credentials file. Extract the private key. The service account needs permissions for:
    1. Reading Cloud Storage objects
    2. Reading Secret Manager payloads
4. Enable Secret Manager and store the private key as a secret.
5. Enable Cloud Build and Cloud Run, then execute a deployment of the app substituting the environment variables as needed:
    ```bash
    gcloud run deploy mediabrowser --platform managed --region $REGION --image gcr.io/$PROJECT_ID/mediabrowser:latest --allow-unauthenticated --set-env-vars BUCKET_NAME=$BUCKET_NAME,WEB_USERNAME=$WEB_USERNAME,WEB_PASSWORD=$WEB_PASSWORD,PK_SECRET_NAME=$PK_SECRET_NAME
    ```
    Alternatively, you can set up a Cloud Build pipeline using the [`cloudbuild.yaml`](/cloudbuild.yaml) file included in the repo and trigger the pipeline.
    _Note:_ within the Cloud Run environment, the `GOOGLE_APPLICATION_CREDENTIALS` variable does not need to be specified - an implicit account is used.
6. Now, you should be able to upload media into your Cloud Storage bucket and navigate to the Cloud Run app URL to browse content.

## Technology

* [Go](https://golang.org/)
* [Google Cloud Storage](https://cloud.google.com/storage/docs/)
* [Google Cloud Secret Manager](https://cloud.google.com/secret-manager/docs/)
* [Google Cloud Run](https://cloud.google.com/run/docs/)
* [Google Cloud Build](https://cloud.google.com/cloud-build/docs/)

## Credit

Favicon Video File by Vicons Design from the Noun Project
