FROM golang:latest
MAINTAINER Andreas Koch <andy@allmark.io>

# Install pandoc for RTF conversion
RUN apt-get update && apt-get install -qy pandoc

# Build
ADD . /go/src/vendor/github.com/andreaskoch/allmark
WORKDIR /go/src/vendor/github.com/andreaskoch/allmark
RUN make install
RUN make crosscompile

# Data
RUN mkdir /data
ADD . /data

VOLUME ["/data"]

EXPOSE 33001

CMD ["bin/files/allmark", "serve", "/data"]