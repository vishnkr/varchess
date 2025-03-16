
FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o core_server vc-core/main.go

FROM alpine:latest
WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/core_server /app/core_server
COPY vc-core/.env .env 
EXPOSE 5000

CMD ["./core_server"]
