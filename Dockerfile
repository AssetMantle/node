FROM golang:1.14-buster

# Set up dependencies
ENV PACKAGES curl make git
ENV PATH=/root/.cargo/bin:$PATH

# Set working directory for the build
WORKDIR /usr/local/app

# Install minimum necessary dependencies
RUN apt update && apt install -y $PACKAGES

# Install Rust and wasm32 dependencies
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
RUN rustup default stable \
    && rustup default stable \
    && rustup update stable \
    && rustup target list --installed \
    && rustup target add wasm32-unknown-unknown

# Create appuser
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/app" \
    --shell "/sbin/nologin" \
    --uid "${UID}" \
    "${USER}"
USER 10001

# Add source files
COPY . .

# Build client
RUN make install

# Run persistenceCore by default, omit entrypoint to ease using container with cli
CMD ["assetClient"]
