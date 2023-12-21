# https://www.visioncortex.org/vtracer/
# https://github.com/sfomuseum/go-sfomuseum-coloringbook
# https://xmlgraphics.apache.org/batik/

FROM rust:alpine AS rusttools

RUN cargo install vtracer

FROM golang:1.21-alpine as gotools

RUN mkdir /build
COPY . /build/go-coloringbook

RUN apk update && apk upgrade \
    && cd /build/go-coloringbook \
    && go build -mod vendor -ldflags="-s -w" -o /usr/local/bin/outline cmd/outline/main.go \    
    && cd \
    && rm -rf build
    
FROM alpine

RUN mkdir /usr/local/src

RUN apk update && apk upgrade \
    && apk add openjdk21-jre \
    && cd /usr/local/src \
    && wget -O batik-bin-1.17.tar.gz 'https://www.apache.org/dyn/closer.cgi?filename=/xmlgraphics/batik/binaries/batik-bin-1.17.tar.gz&action=download' \
    && tar -xvzf batik-bin-1.17.tar.gz

COPY --from=rusttools /usr/local/cargo/bin/vtracer /usr/local/bin/vtracer
COPY --from=gotools /usr/local/bin/outline /usr/local/bin/outline