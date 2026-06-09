# Build executable using the golang official image
FROM golang:1.26.4 AS builder

WORKDIR /app

ENV GOPRIVATE=github.com/meetgeekai/*

# Configures Git to use specific credentials
ARG AUTH_USER
ARG AUTH_TOKEN
RUN git config --global url."https://${AUTH_USER}:${AUTH_TOKEN}@github.com".insteadOf "https://github.com"

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod tidy

COPY . .

# Build
RUN GOOS=linux go build -o /release .

# Minimize image size footprint by running the executable in another image
FROM debian:bookworm-slim AS runner

# Install dependencies
RUN apt-get update && \
    apt-get --no-install-recommends install ca-certificates -y && \
    apt-get clean && \
    groupadd -g 3000 podgroup && useradd -g podgroup -u 3000 poduser

WORKDIR /build/

# Copy compiled binary into the runner image
COPY --from=builder /release .

# Change permissions
RUN chown -R poduser:podgroup /build/release

# Change user to the nonroot user
USER poduser

# Additional settings inside the container
EXPOSE 8000

# Run code
ENTRYPOINT [ "/build/release" ]
