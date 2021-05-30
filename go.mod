module github.com/ChainSafe/chainbridge-celo

go 1.15

replace github.com/celo-org/celo-bls-go => github.com/celo-org/celo-bls-go v0.1.7

replace github.com/ethereum/go-ethereum => github.com/celo-org/celo-blockchain v1.3.2

require (
	github.com/ChainSafe/chainbridge-utils v1.0.6
	github.com/celo-org/celo-bls-go v0.2.4
	github.com/ethereum/go-ethereum v1.9.17
	github.com/golang/mock v1.4.4
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20190923125748-758128399b1d
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
)
