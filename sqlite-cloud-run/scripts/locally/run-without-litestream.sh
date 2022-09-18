#!/bin/bash

set -e

# We probably want to source SECRET-config.sh file first
source SECRET-config.sh


go run ./cmd/server/*.go