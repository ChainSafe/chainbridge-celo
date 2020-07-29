module github.com/ChainSafe/chainbridge-celo

go 1.13

require (
	github.com/ChainSafe/ChainBridge v1.0.3
	github.com/ChainSafe/log15 v1.0.0
	github.com/ethereum/go-ethereum v1.9.18
)

replace github.com/ethereum/go-ethereum => github.com/celo-org/celo-blockchain v0.0.0-20200612100840-bf2ba25426f9

replace github.com/celo-org/celo-bls-go => github.com/celo-org/celo-bls-go v0.1.7
