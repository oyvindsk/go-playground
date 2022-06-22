#!/bin/bash

set -e

IMAGE_URL="europe-north1-docker.pkg.dev/sqlite-test-353918/sqlite-test/foo:"`date '+%Y-%m-%d--%H-%M-%S'`

# Build with Cloud Run and store the image in Artifact Registry
gcloud \
    --project sqlite-test-353918 \
    builds submit \
    --region=europe-north1  \
    --tag $IMAGE_URL \
    .

# Deploy the image from Artifact Registry
gcloud \
    --project sqlite-test-353918 \
    beta run deploy \
    test-1 \
    --region=europe-north1 \
    --allow-unauthenticated \
    --max-instances=1 \
    --set-env-vars='DB_PATH=foo.db,REPLICA_URL=gcs://oyvindsk-sqlite-test-litestream' \
    --execution-environment gen2 \
    --image=$IMAGE_URL