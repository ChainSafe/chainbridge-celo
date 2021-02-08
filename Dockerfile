# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

FROM golang:1.14-alpine as builder

RUN apk add --no-cache gcc musl-dev linux-headers

ADD . /src

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

WORKDIR /src
RUN go mod download
RUN go build -o build/chainbridge-celo .

# # final stage
FROM alpine
RUN apk update && apk add ca-certificates binutils && rm -rf /var/cache/apk/*

COPY --from=builder /build ./
RUN chmod +x ./build

ENTRYPOINT ["./build/chainbridge-celo"]
