# chainbridge-celo

An alternative implementation of ChainBridge which supports bridging Celo chains. It includes additional proving mechanisms to provide initial security.

## Developers

See [developers.md](/docs/developers.md).

# ChainSafe Security Policy

## Reporting a Security Bug

We take all security issues seriously, if you believe you have found a security issue within a ChainSafe
project please notify us immediately. If an issue is confirmed, we will take all necessary precautions 
to ensure a statement and patch release is made in a timely manner.

Please email us a description of the flaw and any related information (e.g. reproduction steps, version) to
[security at chainsafe dot io](mailto:security@chainsafe.io).


# Configuration

> Note: TOML configs have been deprecated in favour of JSON

A chain configurations take this form:

```
{
    "name": "eth",                      // Human-readable name
    "type": "ethereum",                 // Chain type (eg. "ethereum" or "substrate")
    "id": "0",                          // Chain ID
    "endpoint": "ws://<host>:<port>",   // Node endpoint
    "from": "0xff93...",                // On-chain address of relayer
    "opts": {},                         // Chain-specific configuration options (see below)
}
```

See `config.json.example` for an example configuration.

### Celo Options

Celo chains support the following additional options:

```
{
    "bridge": "0x12345...",          // Address of the bridge contract (required)
    "erc20Handler": "0x1234...",     // Address of erc20 handler (required)
    "erc721Handler": "0x1234...",    // Address of erc721 handler (required)
    "genericHandler": "0x1234...",   // Address of generic handler (required)
    "maxGasPrice": "0x1234",         // Gas price for transactions (default: 20000000000)
    "gasLimit": "0x1234",            // Gas limit for transactions (default: 6721975)
    "http": "true",                  // Whether the chain connection is ws or http (default: false)
    "startBlock": "1234",            // The block to start processing events from (default: 0)
    "blockConfirmations": "10",      // Number of blocks to wait before processing a block
    "epochSize": "12"                // Size of chain epoch. eg. The number of blocks after which to checkpoint and reset the pending votes
    "gasMultiplier": "1.25", 		 // Multiplies the gas price by the supplied value (default: 1)
}
```


## Building

`make build`: Builds `chainbridge` in `./build`.

**or**

`make install`: Uses `go install` to add `chainbridge-celo` to your GOBIN.

## Docker
To build the Docker image locally run:

```
docker build -t chainsafe/chainbridge-celo .
```

To start ChainBridge:

```
docker run -v ./config.json:/config.json chainsafe/chainbridge-celo
```

## CLI
ChainBridge CLI allows you to interact with the on-chain components of ChainBridge. Detailed CLI docs are [here](cbcli/README.md).
