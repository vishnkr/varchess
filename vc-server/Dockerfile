FROM golang:latest

RUN mkdir /app
WORKDIR /app/vc-server

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

CMD ["make","build_run"]
EXPOSE 5000
