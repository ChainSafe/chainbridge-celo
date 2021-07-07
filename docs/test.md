# Testing

## Testing Environments

There are two Docker enviroments provided for testing.


### Without Elections

`make docker` will start two independant Celo chains with RPC endpoints on localhost ports `8545` and `8546`. 

This can be useful for simple tests that don't depend on validator set changes.

### With Elections

`make elections` starts three nodes running a singe chain and starts a script which runs the elections process. They
 expose RPC endpoints on localhost ports `8545`, `8645` and `8745`.

To provide deterministic validator sets, the elections script will submit votes that cause the set to repeatedly add and then remove one validator group:

Even Epochs (blockNumber % epochLength == 0):
```
Charlie
Dave
```
Odd Epochs (blockNumber % epochLength == 0):
```
Charlie
Dave
Eve
```
