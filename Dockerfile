FROM golang:1.4.2-cross
MAINTAINER Andreas Koch <andy@allmark.io>

# Install pandoc for RTF conversion
RUN apt-get update && apt-get install -qy pandoc

# Build
ADD . /go
RUN go run make.go -crosscompile

# Data
RUN mkdir /data
ADD . /data

VOLUME ["/data"]

EXPOSE 8080

CMD ["/go/bin/allmark", "serve", "/data"]
