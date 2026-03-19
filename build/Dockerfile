# Start by building the application.
FROM golang:1.25 AS builder

LABEL contact=tombenke@gmail.com

# Install Task binary to BINDIR=/usr/local/bin
ENV BINDIR=/usr/local/bin

# Download and install Taskfile
RUN /bin/bash -o pipefail -c 'curl -sL https://taskfile.dev/install.sh | /bin/bash'

WORKDIR /workdir
COPY go.mod go.sum ./

ARG REPO_PULL
RUN git config --global url."https://$REPO_PULL:x-oauth-basic@github.com/".insteadOf "https://github.com/" && \
    go mod download

ENV CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2 -fstack-protector-all"
ENV CGO_ENABLED=0

# Build the binaries
COPY cmd cmd
COPY internal internal
COPY .git .git

COPY Taskfile.yml ./
RUN task build

FROM alpine:3.23

# app's parameter
ENV PRINT_CONFIG=false
ENV LOG_LEVEL=info
ENV LOG_FORMAT=json
ENV LIVENESS_CHECK_PATH=/live
ENV READINESS_CHECK_PATH=/ready
ENV HEALTHCHECK_PORT=8080
EXPOSE 8080

# Copy assets and binaries
COPY --from=builder /workdir/server /mcp-time-server

WORKDIR /
# Start the app
ENTRYPOINT ["/mcp-time-server",  "--transport", "sse", "--addr", "0.0.0.0:8080" ]
