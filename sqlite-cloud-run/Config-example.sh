#!/usr/bin/env bash

# You probably want to rename this to SECRET-config.sh, since it's sourced from some scripts

####################################################################################
## Tavla config:
####################################################################################

# Path to the sqlite file(s):
export T_DB_PATH="database-files/tavla.db"

# Litestream replicate url,
# typically a s3 or gcs bucket + path, see litestream docs
export T_REPLICA_URL="gcs://.. .."

# Key used to sign (?) the session cookie.
# Do we really need this if we store the sessions in the backend?
# "Should be 32 or 64 bytes", see Gorrila Securecookie GenerateRandomKey()
export T_SESSION_KEY='<SecretKey>' #  FIXME 

# The password the users need to access
export T_USER_PASSWORD="foo bar 123"


####################################################################################
## Google Cloud and Docker config:  
####################################################################################
export T_PROJECT=.. ..                     # GCP project 
export T_REGION=europe-north1              # GCP Region
export T_SERVICE_NAME=sqlite-test2         # Service Name in Cloud Run

# Cloud Artifact Artifact Repository
# !! Must be created in the console - TODO: Put this in some script? 
export T_ARTIFACT_REPO=sqlite-test

# Docker image tag. Used when building the Docker image locally and when deploying to Cloud Run
export T_IMAGE_URL="${T_REGION}-docker.pkg.dev/${T_PROJECT}/${T_ARTIFACT_REPO}/${T_SERVICE_NAME}:"`date '+%Y-%m-%d--%H-%M-%S'`
