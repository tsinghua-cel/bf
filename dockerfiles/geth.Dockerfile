# syntax = docker/dockerfile:1-experimental
FROM golang:1.20-alpine AS build

# Install dependencies
RUN apk update && \
    apk upgrade && \
    apk add --no-cache bash git openssh make

ENV http_proxy http://192.168.1.8:7890
ENV https_proxy http://192.168.1.8:7890

WORKDIR /build

RUN git clone -b v1.13.5 --single-branch https://github.com/ethereum/go-ethereum

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd /build/go-ethereum && make geth && cp /build/go-ethereum/build/bin/geth /geth 

FROM alpine

WORKDIR /root

COPY  --from=build /geth /usr/bin/geth
COPY ./entrypoint/execution.sh /usr/local/bin/execution.sh
RUN chmod u+x /usr/local/bin/execution.sh

ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]
