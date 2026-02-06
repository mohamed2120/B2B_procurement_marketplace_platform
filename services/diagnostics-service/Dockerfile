FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies for building
RUN apk add --no-cache git wget

# Copy go mod files
COPY shared/go.mod shared/go.sum* ./shared/
COPY services/diagnostics-service/go.mod services/diagnostics-service/go.sum* ./services/diagnostics-service/

# Download dependencies
WORKDIR /app/services/diagnostics-service
RUN go mod download

# Copy source code
COPY shared/ /app/shared/
COPY services/diagnostics-service/ /app/services/diagnostics-service/

# Build
WORKDIR /app/services/diagnostics-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/diagnostics-service ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

COPY --from=builder /app/bin/diagnostics-service .

EXPOSE 8013

CMD ["./diagnostics-service"]
