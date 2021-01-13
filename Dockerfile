# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

FROM  golang:1.15.6-alpine AS builder

RUN apk add --no-cache make gcc musl-dev linux-headers git build-base

ADD . /src

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

WORKDIR /src
RUN go mod tidy
RUN go build -o build/chainbridge-celo .

# # final stage
FROM debian:stretch-slim
RUN apt-get -y update && apt-get -y upgrade && apt-get install ca-certificates wget -y
RUN wget -P /usr/local/bin/ https://chainbridge.ams3.digitaloceanspaces.com/subkey-rc6 \
  && mv /usr/local/bin/subkey-rc6 /usr/local/bin/subkey \
  && chmod +x /usr/local/bin/subkey
RUN subkey --version

COPY --from=builder /build ./
RUN chmod +x ./build

ENTRYPOINT ["./build/chainbridge-celo"]
