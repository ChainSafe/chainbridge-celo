package config

import (
	"time"
)

type Gas struct {
	DefaultGasLimit int64
	DefaultGasPrice int64
}

type Test struct {
	EndPoint             string
	TestTimeout          time.Duration
	TestRelayerThreshold int64
	TestChainID          uint8
}

type Network struct {
	BlockRetryInterval     time.Duration
	BlockDelay             int64
	BlockRetryLimit        int
	ZeroAddress            string
	ExecuteBlockWatchLimit int
	TxRetryInterval        time.Duration
	TxRetryLimit           int
}

type Configuration struct {
	Gas     Gas
	Network Network
	Test    Test
}
