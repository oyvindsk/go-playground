

# This asumes a script copies the necessary file into the same direcory as this Dockerfile
# Based on: https://github.com/magnuswahlstrand/litestream-cloud-run
#
# https://docs.docker.com/engine/reference/builder/

FROM golang:1.18-buster as builder

WORKDIR /app

# Copy the files needed to build the server (just compile, not run..)
COPY go.* ./
COPY cmd/server ./cmd/server/

# Go build
RUN go build -v -mod=readonly ./cmd/server

# Download the static build of Litestream directly into the path & make it executable.
# This is done in the builder and copied as the chmod doubles the size.
ADD https://github.com/benbjohnson/litestream/releases/download/v0.3.8/litestream-v0.3.8-linux-amd64-static.tar.gz /tmp/litestream.tar.gz
RUN tar -C /app -xzf /tmp/litestream.tar.gz


########################################################################
########################################################################
########################################################################
########################################################################


FROM debian:buster-slim

# Get certificates for tls, so litestream can communicate with gcs
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the files needed to *run* the server and litestream
COPY --from=builder /app/server .
COPY --from=builder /app/litestream .
COPY run.sh .

# Run
RUN chmod +x run.sh
CMD ["/app/run.sh"]