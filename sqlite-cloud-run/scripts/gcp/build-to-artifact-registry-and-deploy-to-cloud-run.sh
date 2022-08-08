#!/bin/bash

set -e

IMAGE_URL="europe-north1-docker.pkg.dev/sqlite-test-353918/sqlite-test/foo:"`date '+%Y-%m-%d--%H-%M-%S'`

# Copy files to a tmp directory so we can "submit" that to gcp build
# also used by local Docker builds
source  scripts/docker/docker-build-prep.sh
tmpdir=$(gatherFiles) # run a function in the above file 

    
# Build with Cloud Run and store the image in Artifact Registry
gcloud \
    --project sqlite-test-353918 \
    builds submit \
    --region=europe-north1  \
    --tag $IMAGE_URL \
    $tmpdir

# Deploy the image from Artifact Registry
#   --no-cpu-throttling: Might be needed to complte litestream sync to gcs?
#       https://cloud.google.com/run/docs/configuring/cpu-allocation#command-line
gcloud \
    --project sqlite-test-353918 \
    beta run deploy \
    test-1 \
    --region=europe-north1 \
    --allow-unauthenticated \
    --max-instances=1 \
    --no-cpu-throttling \
    --set-env-vars='DB_PATH=database-files/foo.db,REPLICA_URL=gcs://oyvindsk-sqlite-test-litestream' \
    --execution-environment gen2 \
    --image=$IMAGE_URL