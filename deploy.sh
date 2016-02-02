#!/bin/bash

set -e

echo "Authenticating with Google cloud services"
codeship_google authenticate

echo "Setting default zone $DEFAULT_ZONE"
gcloud config set compute/zone $DEFAULT_ZONE

echo "Tagging the Docker machine for Google Container Registry push"
docker tag -f go-eventsource $GOOGLE_CONTAINER_NAME

echo "Pushing to Google Container Registry: $GOOGLE_CONTAINER_NAME"
gcloud docker push $GOOGLE_CONTAINER_NAME
