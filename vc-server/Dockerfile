FROM golang:latest

RUN mkdir /app
WORKDIR /app/vc-server

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

CMD ["go build","./vc-server"]
EXPOSE 5000:5000
