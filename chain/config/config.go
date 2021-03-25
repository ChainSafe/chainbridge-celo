// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
	"math/big"
	"strconv"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/cmd/cfg"
	"github.com/ChainSafe/chainbridge-celo/flags"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000
const DefaultGasMultiplier = 1

type CeloChainConfig struct {
	ID                     utils.ChainId // ChainID
	Name                   string        // Human-readable chain name
	Endpoint               string        // url for rpc endpoint
	From                   string        // address of key to use // TODO: name should be changed
	KeystorePath           string        // Location of keyfiles
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
	EpochSize              uint64 // Size of chain epoch. eg. The number of blocks after which to checkpoint and reset the pending votes
	GasMultiplier          *big.Float
}

func (cfg *CeloChainConfig) EnsureContractsHaveBytecode(conn *client.Client) error {
	err := conn.EnsureHasBytecode(cfg.BridgeContract)
	if err != nil {
		return err
	}
	err = conn.EnsureHasBytecode(cfg.Erc20HandlerContract)
	if err != nil {
		return err
	}
	err = conn.EnsureHasBytecode(cfg.GenericHandlerContract)
	if err != nil {
		return err
	}
	err = conn.EnsureHasBytecode(cfg.Erc721HandlerContract)
	if err != nil {
		return err
	}
	return nil
}

// parseChainConfig uses a core.ChainConfig to construct a corresponding Config
func ParseChainConfig(rawCfg *cfg.RawChainConfig, ctx *cli.Context) (*CeloChainConfig, error) {
	var ks string
	var insecure bool
	if key := ctx.String(flags.TestKeyFlag.Name); key != "" {
		ks = key
		insecure = true
	} else {
		if ksPath := ctx.String(flags.KeystorePathFlag.Name); ksPath != "" {
			ks = ksPath
		}
	}
	chainId, err := strconv.Atoi(rawCfg.Id)
	if err != nil {
		return nil, err
	}

	config := &CeloChainConfig{
		Name:                   rawCfg.Name,
		ID:                     utils.ChainId(chainId),
		Endpoint:               rawCfg.Endpoint,
		From:                   rawCfg.From,
		KeystorePath:           ks,
		BlockstorePath:         ctx.String(flags.BlockstorePathFlag.Name),
		FreshStart:             ctx.Bool(flags.FreshStartFlag.Name),
		LatestBlock:            ctx.Bool(flags.LatestBlockFlag.Name),
		BridgeContract:         common.Address{},
		Erc20HandlerContract:   common.Address{},
		Erc721HandlerContract:  common.Address{},
		GenericHandlerContract: common.Address{},
		GasLimit:               big.NewInt(DefaultGasLimit),
		MaxGasPrice:            big.NewInt(DefaultGasPrice),
		Http:                   false,
		StartBlock:             big.NewInt(0),
		Insecure:               insecure,
	}

	epochSize, ok := rawCfg.Opts["epochSize"]
	if !ok {
		return nil, errors.New("Set epochSize")
	}
	epochSizeUint, err := strconv.ParseUint(epochSize, 10, 64)
	if err != nil {
		return nil, err
	}
	config.EpochSize = epochSizeUint

	if contract, ok := rawCfg.Opts["bridge"]; ok && contract != "" {
		config.BridgeContract = common.HexToAddress(contract)
	} else {
		return nil, errors.New("must provide opts.bridge field for ethereum config")
	}

	config.Erc20HandlerContract = common.HexToAddress(rawCfg.Opts["erc20Handler"])

	config.Erc721HandlerContract = common.HexToAddress(rawCfg.Opts["erc721Handler"])

	config.GenericHandlerContract = common.HexToAddress(rawCfg.Opts["genericHandler"])

	if gasPrice, ok := rawCfg.Opts["maxGasPrice"]; ok {
		price := big.NewInt(0)
		_, pass := price.SetString(gasPrice, 10)
		if pass {
			config.MaxGasPrice = price
			delete(rawCfg.Opts, "maxGasPrice")
		} else {
			return nil, errors.New("unable to parse max gas price")
		}
	}

	if gasLimit, ok := rawCfg.Opts["gasLimit"]; ok {
		limit := big.NewInt(0)
		_, pass := limit.SetString(gasLimit, 10)
		if pass {
			config.GasLimit = limit
		} else {
			return nil, errors.New("unable to parse gas limit")
		}
	}

	if gasMultiplier, ok := rawCfg.Opts["gasMultiplier"]; ok {
		multiplier := big.NewFloat(DefaultGasMultiplier)
		_, pass := multiplier.SetString(gasMultiplier)
		if pass {
			config.GasMultiplier = multiplier
		} else {
			return nil, errors.New("unable to parse gasMultiplier to float")
		}
	}

	if HTTP, ok := rawCfg.Opts["http"]; ok && HTTP == "true" {
		config.Http = true
	} else if HTTP, ok := rawCfg.Opts["http"]; ok && HTTP == "false" {
		config.Http = false
	}

	if startBlock, ok := rawCfg.Opts["startBlock"]; ok && startBlock != "" {
		block := big.NewInt(0)
		_, pass := block.SetString(startBlock, 10)
		if pass {
			config.StartBlock = block
		} else {
			return nil, errors.New("unable to parse start block")
		}
	}
	return config, nil
}
