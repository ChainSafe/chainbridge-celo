// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"fmt"
	"math/big"
	"testing"

	connection "github.com/ChainSafe/chainbridge-celo/connection"
	utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000

var TestLogger = newTestLogger("test")

var TestChainId = msg.ChainId(0)

var aliceTestConfig = createConfig("alice", nil, nil)

func createConfig(name string, startBlock *big.Int, contracts *utils.DeployedContracts) *Config {
	cfg := &Config{
		name:                   name,
		id:                     0,
		endpoint:               TestEndpoint,
		from:                   name,
		keystorePath:           "",
		blockstorePath:         "",
		freshStart:             true,
		bridgeContract:         common.Address{},
		erc20HandlerContract:   common.Address{},
		erc721HandlerContract:  common.Address{},
		genericHandlerContract: common.Address{},
		gasLimit:               big.NewInt(DefaultGasLimit),
		maxGasPrice:            big.NewInt(DefaultGasPrice),
		http:                   false,
		startBlock:             startBlock,
	}

	if contracts != nil {
		cfg.bridgeContract = contracts.BridgeAddress
		cfg.erc20HandlerContract = contracts.ERC20HandlerAddress
		cfg.erc721HandlerContract = contracts.ERC721HandlerAddress
		cfg.genericHandlerContract = contracts.GenericHandlerAddress
	}

	return cfg
}

func newTestLogger(name string) log15.Logger {
	tLog := log15.New("chain", name)
	tLog.SetHandler(log15.LvlFilterHandler(log15.LvlInfo, tLog.GetHandler()))
	return tLog
}

func newLocalConnection(t *testing.T, cfg *Config) *connection.Connection {
	kp := keystore.TestKeyRing.EthereumKeys[cfg.from]
	conn := connection.NewConnection(TestEndpoint, false, kp, TestLogger, big.NewInt(DefaultGasLimit))
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}

	return conn
}

func deployTestContracts(t *testing.T, client *utils.Client, id msg.ChainId) *utils.DeployedContracts {
	contracts, err := utils.DeployContracts(
		client,
		uint8(id),
		TestRelayerThreshold,
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("=======================================================")
	fmt.Printf("Bridge: %s\n", contracts.BridgeAddress.Hex())
	fmt.Printf("Erc20Handler: %s\n", contracts.ERC20HandlerAddress.Hex())
	fmt.Printf("ERC721Handler: %s\n", contracts.ERC721HandlerAddress.Hex())
	fmt.Printf("GenericHandler: %s\n", contracts.GenericHandlerAddress.Hex())
	fmt.Println("========================================================")

	return contracts
}
