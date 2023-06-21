FROM golang:1.18-alpine

WORKDIR /app/vc-server

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN go build -o vc-server ./cmd/server

EXPOSE 5000

CMD ["./vc-server"]