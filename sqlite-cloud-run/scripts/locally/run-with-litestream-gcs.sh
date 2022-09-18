#!/bin/bash

set -e

# We probably want to source SECRET-config.sh file first
source SECRET-config.sh

go build ./cmd/server

echo "deleting database $T_DB_PATH"
rm -rf $T_DB_PATH 

echo ./litestream restore -o $T_DB_PATH $T_REPLICA_URL
litestream restore -if-replica-exists -o $T_DB_PATH $T_REPLICA_URL

echo ./litestream replicate --exec ./server $T_DB_PATH $T_REPLICA_URL
litestream replicate --exec ./server $T_DB_PATH $T_REPLICA_URL