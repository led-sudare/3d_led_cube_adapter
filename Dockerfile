FROM golang:1.12-alpine3.10

RUN apk update && \
    apk upgrade && \
    apk add --no-cache \
    gcc \
    libc-dev \
    git \
    czmq-dev \
    libzmq \
    libsodium

ENV GO111MODULE=on
EXPOSE 9001 5563

RUN mkdir /go/src/3d_led_cube_adapter
WORKDIR /go/src/3d_led_cube_adapter

ENTRYPOINT ["./docker-entrypoint.sh"]
