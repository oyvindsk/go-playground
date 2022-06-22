## Based on: https://github.com/magnuswahlstrand/litestream-cloud-run

FROM golang:1.18-buster as builder

WORKDIR /app

COPY go.mod go.sum *.go /app/

RUN go build -mod=readonly -v -o server

# Download the static build of Litestream directly into the path & make it executable.
# This is done in the builder and copied as the chmod doubles the size.
ADD https://github.com/benbjohnson/litestream/releases/download/v0.3.8/litestream-v0.3.8-linux-amd64-static.tar.gz /tmp/litestream.tar.gz
RUN tar -C /app -xzf /tmp/litestream.tar.gz


#########   

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/litestream .
COPY run.sh .
RUN chmod +x run.sh

CMD ["/app/run.sh"]
