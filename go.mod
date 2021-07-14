module github.com/ChainSafe/chainbridge-celo

go 1.15

replace github.com/celo-org/celo-bls-go => ./celo-bls-go

require (
	github.com/ChainSafe/chainbridge-utils v1.0.6
	github.com/celo-org/celo-blockchain v1.3.2
	github.com/celo-org/celo-bls-go v0.2.4
	github.com/centrifuge/go-substrate-rpc-client v2.0.0-alpha.5+incompatible
	github.com/golang/mock v1.6.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/rs/zerolog v1.23.0
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210305035536-64b5b1c73954
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
)
