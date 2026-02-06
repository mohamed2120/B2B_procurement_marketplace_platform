FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY shared/ /shared/
COPY services/virtual-warehouse-service/ ./

RUN cd /shared && go mod tidy
RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/virtual-warehouse-service ./cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

COPY --from=builder /app/virtual-warehouse-service .

EXPOSE 8011

CMD ["./virtual-warehouse-service"]
