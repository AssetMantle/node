# syntax=docker/dockerfile:latest
FROM golang:1.14-buster as build

# Set up dependencies
ENV PACKAGES curl make git
ENV PATH=/root/.cargo/bin:$PATH

# Set working directory for the build
WORKDIR /usr/local/app

# Install minimum necessary dependencies
RUN --mount=type=cache,target=/var/cache/apt \
  apt update && apt install -y $PACKAGES

# Install Rust and wasm32 dependencies
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

# Add source files
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod download

# Build
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  make install

FROM ubuntu
COPY --from=build '/go/pkg/mod/github.com/!cosm!wasm/go-cosmwasm@v0.10.0/api/libgo_cosmwasm.so' /lib/x86_64-linux-gnu/libgo_cosmwasm.so
COPY --from=build /go/bin/assetClient /usr/bin/assetClient
COPY --from=build /go/bin/assetNode /usr/bin/assetNode
ENTRYPOINT ["assetClient"]
