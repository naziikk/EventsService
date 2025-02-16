FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o redis-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/redis-service /usr/local/bin/redis-service

EXPOSE 8010

CMD ["/usr/local/bin/redis-service"]