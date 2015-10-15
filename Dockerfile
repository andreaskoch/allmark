FROM golang:latest
MAINTAINER Andreas Koch <andy@allmark.io>

# Install pandoc for RTF conversion
RUN apt-get update && apt-get install -qy pandoc

# Build
ADD . /go
RUN go run make.go -crosscompile
RUN go run make.go -install

# Data
RUN mkdir /data
ADD . /data

VOLUME ["/data"]

CMD ["/go/bin/allmark", "serve", "/data"]
