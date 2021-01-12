module github.com/ChainSafe/chainbridge-celo

go 1.15

require (
	github.com/ChainSafe/chainbridge-ethereum-trie v0.0.0-20210112124933-e83f68cceafb
	github.com/ChainSafe/chainbridge-utils v1.0.3
	github.com/aristanetworks/goarista v0.0.0-20200609010056-95bcf8053598 // indirect
	github.com/celo-org/celo-bls-go v0.1.4
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/ethereum/go-ethereum v1.9.24
	github.com/golang/mock v1.4.4
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0 // indirect
	github.com/rs/zerolog v1.19.0
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/net v0.0.0-20200822124328-c89045814202
)

replace github.com/celo-org/celo-bls-go => github.com/celo-org/celo-bls-go v0.1.7

replace github.com/ethereum/go-ethereum => github.com/celo-org/celo-blockchain v0.0.0-20200612100840-bf2ba25426f9
