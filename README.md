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

## Getting Started

1. Run `go build .` from the project directory to download packages and build the application.
2. _Optional:_ Execute the steps in [Google Cloud Setup](#google-cloud-setup) to configure your own Google Cloud project.
3. Set environment variables for your own Google Cloud project. See below 
4. Spin up the application by executing the `mediabrowser` executable with your environment variables configured.

## Configuration

### Environment Variables

* `GOOGLE_APPLICATION_CREDENTIALS` _(required)_ - path to Google Cloud service account key file (JSON file); the service account must have Cloud Storage object reader and Secret Manager Secret Accessor roles
  * For local development, you may [use your own account credentials](https://cloud.google.com/docs/authentication/application-default-credentials#personal) by authenticating with the `gcloud` CLI. 
* `BUCKET_NAME` _(required)_ - name of the Google Cloud Storage bucket to browse
* `SERVICE_ACCOUNT_NAME` _(required)_ - name of the IAM service account used to access Cloud Storage and Secret Manager
* `PK_SECRET_NAME` _(required)_ - resource ID of the secret storing the service account's private key
* `PORT` _(optional)_ - set the port to expose the HTTP server on
* `WEB_USERNAME` _(optional)_ - username to expect for HTTP Basic Authentication
* `WEB_PASSWORD` _(optional)_ - password to expect for HTTP Basic Authentication

### Google Cloud Setup

1. Create a project for your app.
2. Install the [Google Cloud SDK](https://cloud.google.com/sdk) locally as well as the `jq` program. Then, you may execute the deployment script [`deploy-mediabrowser-dependencies.sh`](/deploy-mediabrowser-dependencies.sh) with the name of your project as the first argument.
   This script performs the following operations:
    1. Enables the Cloud APIs necessary for the app.
    2. Create a Cloud Storage bucket to serve media. No special permissions are required for the bucket.
    3. Create a service account to use for the app. Securely store the credentials file. Extract the private key. The service account needs permissions for:
        1. Reading Cloud Storage objects
        2. Reading Secret Manager payloads
    4. Store the private key as a secret in Secret Manager.
3. Make sure that Cloud Build and Cloud Run are enabled (they are enabled by the deployment script), then execute a deployment of the app substituting the environment variables as needed:
    ```bash
    gcloud run deploy mediabrowser --platform managed --region $REGION --image gcr.io/$PROJECT_ID/mediabrowser:latest --allow-unauthenticated --service-account $_SERVICE_ACCOUNT_EMAIL --set-env-vars BUCKET_NAME=$BUCKET_NAME,WEB_USERNAME=$WEB_USERNAME,WEB_PASSWORD=$WEB_PASSWORD,PK_SECRET_NAME=$PK_SECRET_NAME
    ```
    Alternatively, you can set up a Cloud Build pipeline using the [`cloudbuild.yaml`](/cloudbuild.yaml) file included in the repo and trigger the pipeline.
    _Note:_ within the Cloud Run environment, the `GOOGLE_APPLICATION_CREDENTIALS` variable does not need to be specified - the specified service account is exposed to the runtime automatically.
4. Now, you should be able to upload media into your Cloud Storage bucket and navigate to the Cloud Run app URL to browse content.

## Technology

* [Go](https://golang.org/)
* [Google Cloud Storage](https://cloud.google.com/storage/docs/)
* [Google Cloud Secret Manager](https://cloud.google.com/secret-manager/docs/)
* [Google Cloud Run](https://cloud.google.com/run/docs/)
* [Google Cloud Build](https://cloud.google.com/cloud-build/docs/)

## Credit

Favicon Video File by Vicons Design from the Noun Project
