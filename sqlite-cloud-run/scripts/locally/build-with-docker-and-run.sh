#!/bin/bash

set -e

# We probably want to source SECRET-config.sh file first
source SECRET-config.sh

# Build a docker image, with litestrem and everything, and run it locally
# Mostly useful to test the Dockerfile and build scipt locally, without having to wait for Cloud Build

# Copy files to a tmp directory, since dockerfiles are kind of picky about what it will COPY into the image
source  scripts/docker/docker-build-prep.sh
tmpdir=$(gatherFiles) # run a function in the above file 

# Build, using the files we gathered in the temp dir
echo -e "build-with-docker-and-run.sh: Will build with temp dir: $tmpdir"
sudo docker build $tmpdir -t sqlite-test

# Run locally
sudo docker run \
    --rm \
    -ti \
    -v ~/.config/gcloud/application_default_credentials.json:/app/auth \
    -e "GOOGLE_APPLICATION_CREDENTIALS=/app/auth"   \
    -e "T_DB_PATH=${T_DB_PATH}"                     \
    -e "T_REPLICA_URL=${T_REPLICA_URL}"             \
    -e "T_SESSION_KEY=${T_SESSION_KEY}"             \
    -e "T_USER_PASSWORD=${T_USER_PASSWORD}"         \
    -p 8080:8080 sqlite-test:latest