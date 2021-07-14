# CLI Options

ChainBridge CLI allows you to interact with the on-chain components of ChainBridge. Detailed CLI docs can be found [here](https://github.com/ChainSafe/chainbridge-celo/blob/main/cbcli/README.md).

## Flags

### Global

```zsh
    --help, -h           show help (default: false)
    --version, -v        print the version (default: false)
```

## Commands

### `chainbridge-celo run`
```zsh
   --config value       JSON configuration file
   --verbosity value    Supports levels crit (silent) to trce (trace) (default: “info”)
   --keystore value     Path to keystore directory (default: “./keys”)
   --blockstore value   Specify path for blockstore
   --fresh              Disables loading from blockstore at start. Opts will still be used if specified. (default: false)
   --latest             Overrides blockstore and start block, starts from latest block (default: false)
   --metrics            Enables metric server (default: false)
   --metricsPort value  Port to serve metrics on (default: 8001)
   --leveldb value      sets path to leveldb database
   --testkey value      Applies a predetermined test keystore to the chains.
   --help, -h           show help (default: false)
```

### `chainbridge-celo cli`
```
    --url value                 RPC url of blockchain node (default: "ws://localhost:8545")
    --gasLimit value            gasLimit used in transactions (default: 6721975)
    --gasPrice value            gasPrice used for transactions (default: 20000000000)
    --networkID value           networkID (default: 0)
    --privateKey value          Private key to use
    --jsonWallet value          Encrypted JSON wallet
    --jsonWalletPassword value  Password for encrypted JSON wallet
    --help, -h                  show help (default: false)
```

### `chainbridge-celo deploy`
```zsh
    --help, -h           show help (default: false)
```
