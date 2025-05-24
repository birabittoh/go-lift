# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder

WORKDIR /build

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Transfer source code
COPY *.go ./
COPY templates ./templates
COPY src ./src

# Build
RUN CGO_ENABLED=0 go build -trimpath -o /dist/go-lift

# Test
FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM scratch AS build-release-stage

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY templates ./templates
COPY static ./static
COPY --from=builder /dist .

ENTRYPOINT ["./go-lift"]
