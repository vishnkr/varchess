
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY ../go.mod ../go.sum ./
RUN go mod download

COPY . .
RUN go build -o wsserver vc-ws/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/wsserver /app/wsserver
COPY .env .env 
EXPOSE 8080


CMD ["./wsserver"]
