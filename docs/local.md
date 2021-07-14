# Running Locally
## Prerequisites

- Docker, docker-compose
- chainbridge-celo `v0.0.1` binary (see [README](https://github.com/chainsafe/chainbridge-celo))

# Steps To Get Started
1. [`Start Local Chains`](#start-local-chains)
2. [`Deploy Contracts`](#deploy-contracts)
3. [`Run Relayer`](#run-relayer)

## Start Local Chains

The easiest way to get started is to use the [docker-compose-chains.yml](https://github.com/ChainSafe/chainbridge-celo/blob/main/docker-compose-chains.yml), shown below.

These instructions will start two independent instances of Celo Blockchain nodes (Celo uses a custom fork of [go-ethereum](https://github.com/ethereum/go-ethereum) under the hood in [celo-blockchain](https://github.com/celo-org/celo-blockchain)):
```yaml
# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only
version: ‘3’
services:
  geth1:
    image: “chainsafe/chainbridge-celo-blockchain:20200720215414-7a61816”
    container_name: celo-geth1
    ports:
      - “8545:8545"
  geth2:
    image: “chainsafe/chainbridge-celo-blockchain:20200720215414-7a61816”
    container_name: celo-geth2
    ports:
      - “8546:8545"
```

```zsh
docker-compose -f docker-compose-chains.yml up -V
```

(Use `-V` to always start with new chains. These instructions depend on deterministic Ethereum addresses, which are used as defaults implicitly by some of these commands. Avoid re-deploying the contracts without restarting both chains, or ensure to specify all the required parameters.)

## On-Chain Setup (Celo Blockchain)

### Deploy Contracts

To deploy the contracts on to the two chains, run the following:

```zsh
chainbridge-celo deploy
```

After running, the expected output looks like this:

```zsh
{"level":"info","url":"ws://localhost:8545","time":"2021-05-21T15:30:34-07:00","message":"Connecting to ethereum chain..."}
{“level”:“info”,“time”:“2021-05-21T15:30:34-07:00",“message”:“Waiting for transaction with hash 0x2bc822511ee3e829cc41ffc2031f37ad4402c69c9e3ec2b9f4471297c23c4a47"}
{“level”:“info”,“time”:“2021-05-21T15:30:36-07:00",“message”:“Waiting for transaction with hash 0x0c3efa9d2755166336a0a19a52e422941b7703a2a1806d28bb34d8051348a139"}
{“level”:“info”,“time”:“2021-05-21T15:30:38-07:00",“message”:“Waiting for transaction with hash 0x1cd8726d2a00e62e9b48c6cfd42ad72c1e39611b009a87b99ee7c3240f0b6979"}
{“level”:“info”,“time”:“2021-05-21T15:30:40-07:00",“message”:“Waiting for transaction with hash 0xea7273998eedebc435b8a1253b8bb6ba265f9a98549e10d3922ad4d590d07a33"}
{“level”:“info”,“time”:“2021-05-21T15:30:42-07:00",“message”:“Waiting for transaction with hash 0xebcd64eb5fdaf9d296632f6e8ea215a63a6593e41e6e519c7683adc665d658fd”}
{“level”:“debug”,“time”:“2021-05-21T15:30:44-07:00",“message”:“Bridge 0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B \r\nerc20 handler 0x3167776db165D8eA0f51790CA2bbf44Db5105ADF \r\nerc721 handler 0x3f709398808af36ADBA86ACC617FeB7F5B7B193E \r\ngeneric handler 0x2B6Ab4b880A45a07d83Cf4d664Df4Ab85705Bc07 \r\nerc20Contract 0x21605f71845f372A9ed84253d2D024B7B10999f4"}
{“level”:“info”,“time”:“2021-05-21T15:30:44-07:00",“message”:“Waiting for transaction with hash 0x9b4530b91162203265a7078b8ae264d181e2d7e1ca1a7af63340ea281de4d506"}
{“level”:“info”,“time”:“2021-05-21T15:30:46-07:00",“message”:“Waiting for transaction with hash 0xc0748fb7feb57f233024c5aaea4153784a3d70c21c1ed7be6ed13e36d8886bba”}
{“level”:“info”,“time”:“2021-05-21T15:30:48-07:00",“message”:“Waiting for transaction with hash 0x03f153850c11c1c4ba7195a4454d02122278ed8e2113446ff411b880ed9f7bd1"}
{“level”:“info”,“time”:“2021-05-21T15:30:50-07:00",“message”:“Waiting for transaction with hash 0xbb9959bf88a6e8d4185bbe222742d93089e229cf5dee32fa40c403f13030ef02"}
{“level”:“info”,“time”:“2021-05-21T15:30:52-07:00",“message”:“Waiting for transaction with hash 0xa9d193dc545f4684e5b5ec59ce28a769e360ec6f9eb3478ecf73d453b12274ff”}
```

## Run Relayer

Steps to run a relayer:
1. Create `config.json` using the below as a sample template
2. Start relayer as a binary using the default "Alice" key
```zsh
chainbridge-celo run --config config.json --testkey alice --fresh --leveldb ./lvldb
```

OR

If you prefer Docker, see the following:

1. Build an image first
```zsh
docker build -t chainsafe/chainbridge-celo .
```
2. Start the relayer as a docker container:
```zsh
docker run -v $(pwd)/config.json:/config.json --network host chainsafe/chainbridge-celo run --testkey alice --fresh --leveldb ./lvldb
```

Sample `config.json`:
```json
{
  "chains": [
    {
      "name": "eth",
      "type": "ethereum",
      "id": "1",
      "endpoint": "ws://localhost:8545",
      "from": "0xff93B45308FD417dF303D6515aB04D9e89a750Ca",
      "opts": {
        "bridge": "0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B",
        "erc20Handler": "0x3167776db165D8eA0f51790CA2bbf44Db5105ADF",
        "erc721Handler": "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E",
        "genericHandler": "0x2B6Ab4b880A45a07d83Cf4d664Df4Ab85705Bc07",
        "gasLimit": "1000000",
        "maxGasPrice": "20000000",
        "epochSize": "12",
        "blockConfirmations": "10"
      }
    },
    {
      "name": "eth",
      "type": "ethereum",
      "id": "1",
      "endpoint": "ws://localhost:8545",
      "from": "0xff93B45308FD417dF303D6515aB04D9e89a750Ca",
      "opts": {
        "bridge": "0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B",
        "erc20Handler": "0x3167776db165D8eA0f51790CA2bbf44Db5105ADF",
        "erc721Handler": "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E",
        "genericHandler": "0x2B6Ab4b880A45a07d83Cf4d664Df4Ab85705Bc07",
        "gasLimit": "1000000",
        "maxGasPrice": "20000000",
        "epochSize": "12",
        "blockConfirmations": "10"
      }
    }
  ]
}
```
- This is an example config file for a single relayer ("Alice") using the contracts we've deployed.
