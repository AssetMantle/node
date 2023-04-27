# syntax=docker/dockerfile:latest
FROM golang:1.19-alpine3.17 as build
WORKDIR /go/node
SHELL [ "/bin/sh", "-cex" ]
RUN --mount=type=cache,target=/var/cache/apk/ \
  apk add make ca-certificates build-base git libc-dev gcc linux-headers curl
ARG TARGETARCH
ARG BUILDARCH
RUN if [ "${TARGETARCH}" = "arm64" ] && [ "${BUILDARCH}" != "arm64" ]; then \
    wget -c https://musl.cc/aarch64-linux-musl-cross.tgz -O - | tar -xz --strip-components 1 -C /usr; \
  elif [ "${TARGETARCH}" = "amd64" ] && [ "${BUILDARCH}" != "amd64" ]; then \
    wget -c https://musl.cc/x86_64-linux-musl-cross.tgz -O - | tar -xz --strip-components 1 -C /usr; \
  fi
COPY . .
RUN --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy; \
  export LIBDIR=/lib; \
  if [ "${TARGETARCH}" = "arm64" ]; then \
    ARCH=aarch64; \
    if [ "${BUILDARCH}" != "arm64" ]; then \
      export LIBDIR=/usr/aarch64-linux-musl/lib; \
      mkdir -p $LIBDIR; \
      export CC=aarch64-linux-musl-gcc CXX=aarch64-linux-musl-g++;\
    fi;\
  elif [ "${TARGETARCH}" = "amd64" ]; then \
    ARCH=x86_64; \
    if [ "${BUILDARCH}" != "amd64" ]; then \
      export LIBDIR=/usr/x86_64-linux-musl/lib; \
      mkdir -p $LIBDIR; \
      export CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++; \
    fi; \
  fi; \
  WASM_VERSION=$(go list -m all | grep github.com/CosmWasm/wasmvm | awk '{print $2}'); \
  if [ ! -z "${WASM_VERSION}" ]; then \
    wget -O $LIBDIR/libwasmvm_muslc.a https://github.com/CosmWasm/wasmvm/releases/download/${WASM_VERSION}/libwasmvm_muslc.$ARCH.a; \
  fi; \
  BUILD_TAGS=muslc CGO_ENABLED=1 LDFLAGS='-linkmode external -extldflags "-static"' make install; \
  mantleNode version

FROM alpine:3.17
COPY --from=build /go/bin/mantleNode /usr/bin/mantleNode
ENTRYPOINT ["mantleNode"]
