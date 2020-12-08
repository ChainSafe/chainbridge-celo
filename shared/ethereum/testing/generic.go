// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package ethtest

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/chainbridge-celo/msg"
	utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

func AssertGenericResourceAddress(t *testing.T, client *utils.Client, handler common.Address, rId msg.ResourceId, expected common.Address) {
	actual, err := utils.GetGenericResourceAddress(client, handler, rId)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actual.Bytes(), expected.Bytes()) {
		t.Fatalf("Generic resoruce mismatch for ID %x. Expected address: %x Got: %x", rId, expected, actual)
	}
	log.Info().Str("handler", handler.Hex()).Str("rId", rId.Hex()).Str("contract", actual.Hex()).Msg("Asserted generic resource ID")
}
