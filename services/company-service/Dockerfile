FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY shared/ /shared/
COPY services/company-service/ ./

RUN cd /shared && go mod tidy
RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/company-service ./cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

COPY --from=builder /app/company-service .

EXPOSE 8002

CMD ["./company-service"]
