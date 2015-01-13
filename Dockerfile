FROM golang:1.4
MAINTAINER Andreas Koch <andy@allmark.io>

# Build
ADD . /go
RUN go run make.go -install

# Data
RUN mkdir /data
ADD . /data

VOLUME ["/data"]

EXPOSE 8080

CMD ["/go/bin/allmark", "serve", "/data"]
