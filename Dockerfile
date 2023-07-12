# FROM golang:alpine as builder


# COPY . $GOPATH/src/github.com/insomniadev/collectivedb
# WORKDIR $GOPATH/src/github.com/insomniadev/collectivedb

# RUN apk add --no-cache git
# RUN GOOS=linux GOARCH=amd64 go build .

FROM alpine

LABEL maintainer="Insomnia Development <insomniadevlabs@gmail.com>"

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /go/src/github.com/insomniadev/collectivedb/collectivedb /collectivedb

COPY collectivedb collectivedb
ENV COLLECTIVE_DATA_DIRECTORY="/data/"

ENTRYPOINT [ "/collectivedb" ]
EXPOSE 9090
EXPOSE 9091
VOLUME /data