#!/bin/bash

set -e

echo "ENV"
env

cd /app

echo "LS1: " $PWD

ls -la

echo "1"
echo ./litestream restore -o $DB_PATH $REPLICA_URL
./litestream restore -if-replica-exists -o $DB_PATH $REPLICA_URL

echo "2"
echo ./litestream replicate --exec ./server $DB_PATH $REPLICA_URL
./litestream replicate --exec ./server $DB_PATH $REPLICA_URL