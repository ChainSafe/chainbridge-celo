# 🌉 <b> Overview </b>

## Summary

An alternative implementation of ChainBridge which supports bridging Celo chains. It includes additional proving mechanisms to provide initial security.

A bridge contract on each chain forms either side of a bridge. Handler contracts allow for customizable behavior upon receiving transactions to and from the bridge. For example locking up an asset on one side and minting a new one on the other. Its highly customizable - you can deploy a handler contract to perform any action you like.

In its current state ChainBridge-Celo operates under a trusted federation model. Deposit events on one chain are detected by a trusted set of off-chain relayers who await finality, submit events to the other chain and vote on submissions to reach acceptance triggering the appropriate handler.

![](./img/system-flow.png)

## Relevant repos

### [ChainBridge](https://github.com/ChainSafe/chainbridge-celo)
This is the core bridging software that Relayers run between chains. Also supports deploying the chainbridge-celo-solidity contracts.

### [chainbridge-solidity](https://github.com/ChainSafe/chainbridge-celo-solidity) 
The Solidity contracts required for chainbridge-celo.
    
### [chainbridge-utils](https://github.com/ChainSafe/chainbridge-celo-blockchain)
A Dockerized Celo-Blockchain node.
