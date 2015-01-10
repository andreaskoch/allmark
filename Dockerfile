FROM golang:1.4
MAINTAINER Andreas Koch <andy@allmark.io>

RUN mkdir /data

ADD make.go /go/make.go
ADD Makefile /go/Makefile
ADD src /go/src

RUN go run make.go -install

ADD documentation /data
VOLUME ["/data"]

EXPOSE 8080

CMD ["/go/bin/allmark", "serve", "/data"]