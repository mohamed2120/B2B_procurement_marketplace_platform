FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
COPY shared/ /shared/
COPY services/identity-service/ ./

# Download dependencies
RUN cd /shared && go mod tidy
RUN go mod download
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/identity-service ./cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget wget

WORKDIR /root/

COPY --from=builder /app/identity-service .

EXPOSE 8001

CMD ["./identity-service"]
