// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package ethtest

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	msg "github.com/ChainSafe/chainbridge-celo/msg"
	utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func CreateErc20Deposit(
	t *testing.T,
	contract *Bridge.Bridge,
	client *utils.Client,
	rId msg.ResourceId,
	destRecipient common.Address,
	destId msg.ChainId,
	amount *big.Int,
) {
	data := utils.ConstructErc20DepositData(destRecipient.Bytes(), amount)

	err := client.LockNonceAndUpdate()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := contract.Deposit(
		client.Opts,
		uint8(destId),
		rId,
		data,
	); err != nil {
		t.Fatal(err)
	}

	client.UnlockNonce()
}

func CreateErc721Deposit(
	t *testing.T,
	bridge *Bridge.Bridge,
	client *utils.Client,
	rId msg.ResourceId,
	destRecipient common.Address,
	destId msg.ChainId,
	tokenId *big.Int,
) {
	data := utils.ConstructErc721DepositData(tokenId, destRecipient.Bytes())

	err := client.LockNonceAndUpdate()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := bridge.Deposit(
		client.Opts,
		uint8(destId),
		rId,
		data,
	); err != nil {
		t.Fatal(err)
	}

	client.UnlockNonce()
}

func CreateGenericDeposit(
	t *testing.T,
	bridge *Bridge.Bridge,
	client *utils.Client,
	rId msg.ResourceId,
	destId msg.ChainId,
	hash []byte,
) {
	data := utils.ConstructGenericDepositData(hash)

	err := client.LockNonceAndUpdate()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := bridge.Deposit(
		client.Opts,
		uint8(destId),
		rId,
		data,
	); err != nil {
		t.Fatal(err)
	}

	client.UnlockNonce()
}
