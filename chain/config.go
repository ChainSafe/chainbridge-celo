// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"

	utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
	"github.com/ChainSafe/chainbridge-utils/core"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/ethereum/go-ethereum/common"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000

type CeloChainConfig struct {
	ID                     msg.ChainId // ChainID
	Name                   string      // Human-readable chain name
	Endpoint               string      // url for rpc endpoint
	From                   string      // address of key to use
	KeystorePath           string      // Location of keyfiles
	BlockstorePath         string
	FreshStart             bool // Disables loading from blockstore at start
	BridgeContract         common.Address
	Erc20HandlerContract   common.Address
	Erc721HandlerContract  common.Address
	GenericHandlerContract common.Address
	GasLimit               *big.Int
	MaxGasPrice            *big.Int
	Http                   bool // Config for type of connection
	StartBlock             *big.Int
	LatestBlock            bool
	Insecure               bool
}

// parseChainConfig uses a core.ChainConfig to construct a corresponding Config
func ParseChainConfig(chainCfg *core.ChainConfig) (*CeloChainConfig, error) {

	config := &CeloChainConfig{
		Name:                   chainCfg.Name,
		ID:                     chainCfg.Id,
		Endpoint:               chainCfg.Endpoint,
		From:                   chainCfg.From,
		KeystorePath:           chainCfg.KeystorePath,
		BlockstorePath:         chainCfg.BlockstorePath,
		FreshStart:             chainCfg.FreshStart,
		BridgeContract:         utils.ZeroAddress,
		Erc20HandlerContract:   utils.ZeroAddress,
		Erc721HandlerContract:  utils.ZeroAddress,
		GenericHandlerContract: utils.ZeroAddress,
		GasLimit:               big.NewInt(DefaultGasLimit),
		MaxGasPrice:            big.NewInt(DefaultGasPrice),
		Http:                   false,
		StartBlock:             big.NewInt(0),
		LatestBlock:            chainCfg.LatestBlock,
		Insecure:               chainCfg.Insecure,
	}

	if contract, ok := chainCfg.Opts["bridge"]; ok && contract != "" {
		config.BridgeContract = common.HexToAddress(contract)
		delete(chainCfg.Opts, "bridge")
	} else {
		return nil, errors.New("must provide opts.bridge field for ethereum config")
	}

	config.Erc20HandlerContract = common.HexToAddress(chainCfg.Opts["erc20Handler"])
	delete(chainCfg.Opts, "erc20Handler")

	config.Erc721HandlerContract = common.HexToAddress(chainCfg.Opts["erc721Handler"])
	delete(chainCfg.Opts, "erc721Handler")

	config.GenericHandlerContract = common.HexToAddress(chainCfg.Opts["genericHandler"])
	delete(chainCfg.Opts, "genericHandler")

	if gasPrice, ok := chainCfg.Opts["maxGasPrice"]; ok {
		price := big.NewInt(0)
		_, pass := price.SetString(gasPrice, 10)
		if pass {
			config.MaxGasPrice = price
			delete(chainCfg.Opts, "maxGasPrice")
		} else {
			return nil, errors.New("unable to parse max gas price")
		}
	}

	if gasLimit, ok := chainCfg.Opts["gasLimit"]; ok {
		limit := big.NewInt(0)
		_, pass := limit.SetString(gasLimit, 10)
		if pass {
			config.GasLimit = limit
			delete(chainCfg.Opts, "gasLimit")
		} else {
			return nil, errors.New("unable to parse gas limit")
		}
	}

	if HTTP, ok := chainCfg.Opts["http"]; ok && HTTP == "true" {
		config.Http = true
		delete(chainCfg.Opts, "http")
	} else if HTTP, ok := chainCfg.Opts["http"]; ok && HTTP == "false" {
		config.Http = false
		delete(chainCfg.Opts, "http")
	}

	if startBlock, ok := chainCfg.Opts["startBlock"]; ok && startBlock != "" {
		block := big.NewInt(0)
		_, pass := block.SetString(startBlock, 10)
		if pass {
			config.StartBlock = block
			delete(chainCfg.Opts, "startBlock")
		} else {
			return nil, errors.New("unable to parse start block")
		}
	}

	if len(chainCfg.Opts) != 0 {
		return nil, fmt.Errorf("unknown Opts Encountered: %#v", chainCfg.Opts)
	}

	return config, nil
}
