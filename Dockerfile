# Build container
FROM golang:bookworm AS build

# Setup environment
RUN mkdir -p /data
WORKDIR /data

# Build the release
COPY . .
RUN make depend/cli
RUN make build/cli

# Extract the release
RUN mkdir -p /out
RUN cp out/skysweeper-server /out/skysweeper-server

# Release container
FROM debian:bookworm

# Add certificates
RUN apt update
RUN apt install -y ca-certificates

# Add the release
COPY --from=build /out/skysweeper-server /usr/local/bin/skysweeper-server

CMD /usr/local/bin/skysweeper-server
