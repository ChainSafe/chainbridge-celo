package config

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
