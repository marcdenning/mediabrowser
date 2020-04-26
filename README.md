# Media Browser

> Web-based media browser written in Go.

## Overview

Media Browser is a directory browser for Google Cloud Storage.
Similar to common web servers' directory and file browsing capabilities, it allows user to browse a storage bucket via a simple web interface using hyperlinks.

Once an individual file is selected, a Signed URL is generated for the object in Cloud Storage and the server redirects the user to the Signed URL.
Since the intention is browsing media, this feature is intended to reduce the runtime of the server itself and allow video media to be streamed directly from the storage bucket.
This kind of optimization makes Media Browser a perfect fit for serverless platforms like Cloud Run.

## Configuration

* `GOOGLE_APPLICATION_CREDENTIALS` _(required)_ - path to Google Cloud service account key file (JSON file); the service account must have 
* `BUCKET_NAME` _(required)_ - name of the Google Cloud Storage bucket to browse
* `WEB_USERNAME` _(optional)_ - username to expect for HTTP Basic Authentication
* `WEB_PASSWORD` _(optional)_ - password to expect for HTTP Basic Authentication

## Technology

* [Go](https://golang.org/)
* [Google Cloud Storage](https://cloud.google.com/storage/docs/)
* [Google Cloud Run](https://cloud.google.com/run/docs/)
