// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package e2e

import (
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
	"math/big"
)

const TestEndpoint = "ws://localhost:8545"
const TestEndpoint2 = "ws://localhost:8547"

func Deploy(ctx *cli.Context) error {
	ethClient, err := client.NewClient(TestEndpoint, false, utils.AliceKp, big.NewInt(utils.DefaultGasLimit), big.NewInt(utils.DefaultGasPrice))
	if err != nil {
		return err
	}

	dpc, err := utils.DeployContracts(ethClient, 1, big.NewInt(1), utils.DefaultRelayerAddresses, "test", "TST", big.NewInt(0))
	if err != nil {
		return err
	}
	src := utils.ChainId(5)
	resourceID := utils.SliceTo32Bytes(append(common.LeftPadBytes(dpc.ERC20TokenAddress.Bytes(), 31), uint8(src)))
	err = utils.RegisterResource(ethClient, dpc.BridgeAddress, dpc.ERC20HandlerAddress, resourceID, dpc.ERC20TokenAddress)
	if err != nil {
		return err
	}

	tenTokens := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))
	err = utils.MintTokens(ethClient, dpc.ERC20TokenAddress, tenTokens)
	if err != nil {
		return err
	}
	err = utils.Erc20Approve(ethClient, dpc.ERC20TokenAddress, dpc.ERC20HandlerAddress, tenTokens)
	if err != nil {
		return err
	}

	err = utils.Erc20AddMinter(ethClient, dpc.ERC20TokenAddress, dpc.ERC20HandlerAddress)
	if err != nil {
		return err
	}

	err = utils.SetBurnable(ethClient, dpc.BridgeAddress, dpc.ERC20HandlerAddress, dpc.ERC20TokenAddress)
	if err != nil {
		return err
	}

	return nil
}
