# Installation

## Building from Source

To build `chainbridge-celo` in `./build`.
```
make build
```

**or**

Use`go install` to add `chainbridge-celo` to your GOBIN.

```
make install
```

## Docker

The official ChainBridge-Celo Docker image can be found [here](https://hub.docker.com/r/chainsafe/chainbridge-celo).

To build the Docker image locally run: 

```
docker build -t chainsafe/chainbridge-celo .
```

To start ChainBridge-Celo:

```
docker run -v $(pwd)/config.json:/config.json --network host chainsafe/chainbridge-celo run --testkey alice --fresh --leveldb ./lvldb
```
