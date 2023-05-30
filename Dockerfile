FROM golang:alpine as builder

LABEL maintainer="Insomnia Development <insomniadevlabs@gmail.com>"

COPY . $GOPATH/src/github.com/insomniadev/collective-db
WORKDIR $GOPATH/src/github.com/insomniadev/collective-db

RUN apk add git
# RUN go get -u github.com/golang/dep/cmd/dep;export GOOS=linux && export CGO_ENABLED=0; dep ensure
RUN go build .
RUN ls -al
RUN pwd

FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/insomniadev/collective-db/collective-db /collective-db
RUN apk add curl
RUN pwd && ls -al
WORKDIR /

CMD ["./collective-db"]
EXPOSE 10000
EXPOSE 9091
VOLUME /data