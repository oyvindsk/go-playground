#!/bin/bash

set -e

DB_PATH="database-files/foo.db"
export DB_PATH

go run ./cmd/server/*.go