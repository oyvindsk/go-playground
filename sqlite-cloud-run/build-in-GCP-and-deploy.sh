#!/bin/bash

# Build from source and deploy the result to Cloud Run,
# in one big swoop,
# kind of like App Engine Standard
# Results can be seen in Cloud Build and Artifact Registry
gcloud \
    --project sqlite-test-353918 \
    run deploy test-1 \
    --region=europe-north1 \
    --allow-unauthenticated \
    --max-instances=1 \
    --source .

# TODO:
#   --concurrency=CONCURRENCY 


# You can also deploy an docker image, that can create locally or in Cloud Build, or some other way
# To build in Cloud Build, with a buildpack:
#   gcloud --project sqlite-test-353918 builds submit --pack image=europe-north1-docker.pkg.dev/sqlite-test-353918/sqlite-test/foo
# image must point to an existing repo in Artifact Registry

# To deploy Artifact Registry or Container Registry, run deploy like above, but substitute `--source` with:
#   --image=europe-north1-docker.pkg.dev/sqlite-test-353918/sqlite-test/foo:latest

