package config

import (
	"math/big"
	"time"
)

type Gas struct {
	DefaultGasLimit int
	DefaultGasPrice int
}

type Network struct {
	BlockRetryInterval int
	BlockDelay         int
	BlockRetryLimit    int
}

type Configuration struct {
	Gas Gas
	Network: Network
}

// var BlockDelay = big.NewInt(10)
// var BlockRetryInterval = time.Second * 5
// var BlockRetryLimit = 5
