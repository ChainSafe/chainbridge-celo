// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package writer

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/chainbridge-celo/pkg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestCreateProposalDataHash(t *testing.T) {

	handler := common.HexToAddress("0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa201")
	data := []byte{}

	mp := &pkg.MerkleProof{
		TxRootHash: common.BytesToHash([]byte{1, 2, 3}),
		Key:        []byte{12, 3, 4},
		Nodes:      []byte{1, 2, 3},
	}

	sv := &pkg.SignatureVerification{
		AggregatePublicKey: []byte{1, 2, 3},
		BlockHash:          common.BytesToHash([]byte{1, 2, 3}),
		Signature:          []byte{2, 3, 4},
	}

	b := bytes.NewBuffer(data)
	b.Write(handler.Bytes())
	b.Write(mp.TxRootHash[:])
	b.Write(mp.Key)
	b.Write(mp.Nodes)
	b.Write(sv.AggregatePublicKey)
	b.Write(sv.BlockHash[:])
	b.Write(sv.Signature)
	testResult := crypto.Keccak256Hash(b.Bytes())

	result := CreateProposalDataHash(data, handler, mp, sv)

	if result != testResult {
		t.Errorf("expected %v got %v", testResult, result)
	}

}
