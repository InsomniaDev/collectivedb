FROM golang:alpine as builder

LABEL maintainer="Insomnia Development <insomniadevlabs@gmail.com>"

COPY . $GOPATH/src/github.com/insomniadev/collectivedb
WORKDIR $GOPATH/src/github.com/insomniadev/collectivedb

RUN apk add --no-cache git
# RUN go get -u github.com/golang/dep/cmd/dep;export GOOS=linux && export CGO_ENABLED=0; dep ensure
RUN GOOS=linux GOARCH=amd64 go build .

FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/insomniadev/collectivedb/collectivedb /collectivedb

ENTRYPOINT [ "/collectivedb" ]
EXPOSE 10000
EXPOSE 9091
VOLUME /data