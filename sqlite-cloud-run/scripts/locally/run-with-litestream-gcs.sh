#!/bin/bash

set -e

DB_PATH="database-files/foo.db"
REPLICA_URL="gcs://oyvindsk-sqlite-test-litestream"

export DB_PATH

go build ./cmd/server

echo "deleting database $DB_PATH"
rm -rf $DB_PATH 

echo ./litestream restore -o $DB_PATH $REPLICA_URL
litestream restore -if-replica-exists -o $DB_PATH $REPLICA_URL

echo ./litestream replicate --exec ./server $DB_PATH $REPLICA_URL
litestream replicate --exec ./server $DB_PATH $REPLICA_URL