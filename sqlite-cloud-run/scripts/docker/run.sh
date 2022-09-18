#!/bin/bash

set -e

echo "ENV"
env

cd /app

echo "LS:"
ls -la

echo ./litestream restore -o $T_DB_PATH $T_REPLICA_URL
./litestream restore -if-replica-exists -o $T_DB_PATH $T_REPLICA_URL

echo ./litestream replicate --exec ./server $T_DB_PATH $T_REPLICA_URL
./litestream replicate --exec ./server $T_DB_PATH $T_REPLICA_URL