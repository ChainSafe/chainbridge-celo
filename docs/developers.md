# Developers

## Testing Environments

There are two Docker enviroments provided for testing.

They both start two independant Celo chains with RPC endpoints on `http://localhost:8545` and `http://localhost:8546`.

### Without Elections

`make docker` will simply run the chains with a static single validator. This can be useful for simple tests that don't depend on validator set changes.

### With Elections

`make elections` starts the chains and starts the elections process for each.

To provide deterministic validator sets, the elections script will submit votes that cause the set to repeatedly add and then remove one validator group:

(Assumption: Epochs start from 0)

Even Epochs:
```
Charlie
Dave
```
Odd Epochs:
```
Charlie
Dave
Eve
```