#!/bin/bash

set -e

sudo docker build . -t sqlite-test

sudo docker run --rm -ti -v ~/.config/gcloud/application_default_credentials.json:/app/auth -e "GOOGLE_APPLICATION_CREDENTIALS=/app/auth" -e "DB_PATH=foo.db" -e "REPLICA_URL=gcs://oyvindsk-sqlite-test-litestream" -p 8080:8080 sqlite-test:latest

# $ sudo docker run --rm -ti -e "DB_PATH=foo.db" -e "REPLICA_URL=gcs://oyvindsk-sqlite-test-litestream" sqlite-test:latest 