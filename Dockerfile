# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder
WORKDIR /src

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cert-tls-check main.go

FROM gcr.io/distroless/static
COPY --from=builder /src/cert-tls-check /cert-tls-check
USER nonroot:nonroot
ENTRYPOINT ["/cert-tls-check"]
