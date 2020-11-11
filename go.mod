module github.com/ChainSafe/chainbridge-celo

go 1.13

require (
	github.com/ChainSafe/chainbridge-utils v1.0.3
	github.com/ChainSafe/log15 v1.0.0
	github.com/aristanetworks/goarista v0.0.0-20200609010056-95bcf8053598 // indirect
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/celo-org/celo-bls-go v0.1.4
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/ethereum/go-ethereum v1.9.18
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/pierrec/xxHash v0.1.5 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.4.1
	github.com/rs/cors v1.7.0 // indirect
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
)

replace github.com/ethereum/go-ethereum => github.com/celo-org/celo-blockchain v0.0.0-20200612100840-bf2ba25426f9

replace github.com/celo-org/celo-bls-go => github.com/celo-org/celo-bls-go v0.1.7
