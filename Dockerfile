# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy application source
COPY . .

# Compile binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/quantito-engine ./cmd/api-quantito-engine

# Final runtime image
FROM alpine:3.19

LABEL AUTHOR="Leandro Gutierrez <leandrogutierrez148@gmail.com>" \
      image="quantito-engine:latest"

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/quantito-engine /app/quantito-engine

EXPOSE 8080

CMD ["/app/quantito-engine"]