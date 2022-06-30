#!/bin/bash

set -e

echo "ENV"
env

cd /app

echo "LS:"
ls -la

echo ./litestream restore -o $DB_PATH $REPLICA_URL
./litestream restore -if-replica-exists -o $DB_PATH $REPLICA_URL

echo ./litestream replicate --exec ./server $DB_PATH $REPLICA_URL
./litestream replicate --exec ./server $DB_PATH $REPLICA_URL