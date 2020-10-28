package chain

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ChainSafe/chainbridge-celo/connection"
	"github.com/ChainSafe/chainbridge-utils/keystore"
)

var TestEndpoint = "ws://localhost:8545"
var GasLimit = big.NewInt(connection.DefaultGasLimit)
var GasPrice = big.NewInt(connection.DefaultGasPrice)
var GasLimitUint64 = uint64(connection.DefaultGasLimit)
var ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
var TestChainID = uint8(0)
var TestRelayerThreshold = big.NewInt(2)
var TestTimeout = time.Second * 30

var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var BobKp = keystore.TestKeyRing.EthereumKeys[keystore.BobKey]
