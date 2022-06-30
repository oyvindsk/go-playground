#!/bin/bash

set -e

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
    -e "GOOGLE_APPLICATION_CREDENTIALS=/app/auth" \
    -e "DB_PATH=database-files/foo.db" \
    -e "REPLICA_URL=gcs://oyvindsk-sqlite-test-litestream" \
    -p 8080:8080 sqlite-test:latest