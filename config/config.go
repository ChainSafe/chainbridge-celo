package config

type Gas struct {
	DefaultGasLimit int64
	DefaultGasPrice int64
}

type Network struct {
	BlockRetryInterval     int64
	BlockDelay             int64
	BlockRetryLimit        int
	ZeroAddress            string
	ExecuteBlockWatchLimit int
	TxRetryInterval        int64
	TxRetryLimit           int
}

type Test struct {
	EndPoint             string
	TestTimeout          int64
	TestRelayerThreshold int64
	TestChainID          uint8
}

type Configuration struct {
	Gas     Gas
	Network Network
	Test    Test
}
