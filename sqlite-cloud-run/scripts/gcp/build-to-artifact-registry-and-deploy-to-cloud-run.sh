#!/bin/bash

set -e

# We probably want to source SECRET-config.sh file first
source SECRET-config.sh


# Copy files to a tmp directory so we can "submit" that to gcp build
# also used by local Docker builds
source  scripts/docker/docker-build-prep.sh
tmpdir=$(gatherFiles) # run a function in the above file 

    
# Build with Cloud Run and store the image in Artifact Registry
gcloud \
    --project ${T_PROJECT} \
    builds submit \
    --region=${T_REGION}  \
    --tag $T_IMAGE_URL \
    $tmpdir

# Deploy the image from Artifact Registry
#   --no-cpu-throttling: Might be needed to complte litestream sync to gcs?
#       https://cloud.google.com/run/docs/configuring/cpu-allocation#command-line
gcloud \
    --project ${T_PROJECT} \
    beta run deploy \
    ${T_SERVICE_NAME} \
    --region=${T_REGION}    \
    --allow-unauthenticated \
    --max-instances=1 \
    --no-cpu-throttling \
    --set-env-vars="T_DB_PATH=${T_DB_PATH},T_REPLICA_URL=${T_REPLICA_URL},T_SESSION_KEY=${T_SESSION_KEY},T_USER_PASSWORD=${T_USER_PASSWORD}" \
    --execution-environment gen2 \
    --image=$T_IMAGE_URL

    # TODO: a better way than --set-env-vars ? 