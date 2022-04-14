# Simple usage with a mounted data directory:
# > docker build -t AssetMantle .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.AssetMantle:/root/.AssetMantle AssetMantle mantleNode init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.AssetMantle:/root/.AssetMantle AssetMantle mantleNode start
FROM golang:1.17-alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /go/src/github.com/AssetMantle/node

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
  make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/mantleNode /usr/bin/mantleNode

COPY ./.scripts/single-node.sh .

EXPOSE 26657

ENTRYPOINT [ "./single-node.sh" ]
# NOTE: to run this image, docker run -d -p 26657:26657 ./single-node.sh {{chain_id}} {{genesis_account}}
