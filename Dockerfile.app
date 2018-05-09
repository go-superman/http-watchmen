FROM alpine:latest
MAINTAINER huoyinghui "huoyinghui@apkpure.com"
WORKDIR /app
RUN apk update && apk add curl bash tree tzdata \
    && cp -r -f /usr/share/zoneinfo/Hongkong /etc/localtime
COPY ./build/linux-amd64/ /app/
CMD ["/app/http-watchmen"]
