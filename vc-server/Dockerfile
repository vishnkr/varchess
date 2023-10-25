FROM golang:1.19-alpine
ARG PB_VERSION=0.19.0
RUN apk add --no-cache \
    unzip \
    ca-certificates
WORKDIR /app/vc-server
ADD https://github.com/pocketbase/pocketbase/releases/download/v${PB_VERSION}/pocketbase_${PB_VERSION}_linux_amd64.zip /tmp/pb.zip
RUN unzip /tmp/pb.zip -d /app/vc-server/pb/
COPY go.mod .
COPY go.sum .
#COPY ./pb/pb_migrations /app/vc-server/pb/pb_migrations
RUN go mod download
COPY . .

EXPOSE 5000 8090
RUN chmod +x entry.sh
CMD ["./entry.sh"]