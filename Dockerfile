FROM ubuntu:18.04 as builder
RUN apt update \
 && apt install -y \
    curl \
 && apt clean \
 && apt autoremove -y \
 && rm -rf /var/lib/apt/lists/* \
 && curl -LSs https://dl.google.com/go/go1.16.2.linux-amd64.tar.gz -o go.tar.gz \
 && tar -xf go.tar.gz \
 && rm -v go.tar.gz \
 && mv go /usr/local/
ENV PATH=${PATH}:/usr/local/go/bin
WORKDIR /app
RUN apt update \
 && apt install -y \
    gcc \
    libblkid-dev \
 && apt clean \
 && apt autoremove -y \
 && rm -rf /var/lib/apt/lists/*


ADD . /src
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org
WORKDIR /src
RUN go mod download
RUN GOOS=linux CGO_ENABLED=1 go build -o build/chainbridge-celo .

RUN chmod +x ./build/chainbridge-celo

ENTRYPOINT ["/src/build/chainbridge-celo"]

# # final stage
FROM ubuntu:18.04

COPY --from=builder /src/build ./
RUN chmod +x ./chainbridge-celo

ENTRYPOINT ["./chainbridge-celo"]