FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY shared/ /shared/
COPY services/billing-service/ ./

RUN cd /shared && go mod tidy
RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/billing-service ./cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

COPY --from=builder /app/billing-service .

EXPOSE 8010

CMD ["./billing-service"]
