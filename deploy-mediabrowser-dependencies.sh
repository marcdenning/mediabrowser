#!/usr/bin/env bash

# Requires Google Cloud SDK to be installed, on PATH, and authenticated
# Requires `jq` to be installed and on PATH
# Input parameters:
# 1. PROJECT_ID _(required)_ - ID of the project to deploy resources to; the project must be created and have billing enabled

PROJECT_ID=$1
SERVICE_ACCOUNT=sa-mediabrowser-prod
gcloud config set project $PROJECT_ID
gcloud services enable \
  cloudapis.googleapis.com \
  cloudbuild.googleapis.com \
  containerregistry.googleapis.com \
  run.googleapis.com \
  secretmanager.googleapis.com \
  storage-api.googleapis.com \
  storage-component.googleapis.com
gcloud iam service-accounts create $SERVICE_ACCOUNT --description="Provides access to Cloud Storage and Secret Manager for mediabrowser."
sleep 5
SERVICE_ACCOUNT_FULL=$(gcloud iam service-accounts list --format="value (email)" | grep $SERVICE_ACCOUNT)
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:$SERVICE_ACCOUNT_FULL" --role="roles/secretmanager.secretAccessor"
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:$SERVICE_ACCOUNT_FULL" --role="roles/storage.objectViewer"
gcloud iam service-accounts keys create --iam-account=$SERVICE_ACCOUNT_FULL ./mediabrowser-key-prod.json
jq ".private_key" ./mediabrowser-key-prod.json | gcloud secrets create mediabrowser-pk --data-file=- --replication-policy=automatic
gsutil mb -b on gs://bucket-$PROJECT_ID-mediabrowser
