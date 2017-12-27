FROM golang:latest AS build-env
ADD . $GOPATH/src/http-watchmen
WORKDIR $GOPATH/src/http-watchmen

ENV CGO_ENABLED=0
RUN mkdir -p /output && go build -o /output/http-watchmen

FROM alpine:latest
MAINTAINER huoyinghui "huoyinghui@apkpure.com"
WORKDIR /app
RUN apk update && apk add curl bash tree tzdata && \
    cp -r -f /usr/share/zoneinfo/Hongkong /etc/localtime

RUN apk add --update ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=build-env /output/http-watchmen /app/
CMD ["/app/http-watchmen"]
