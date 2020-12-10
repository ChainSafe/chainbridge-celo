// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package writer

import (
	"bytes"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/msg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

// constructErc20ProposalData returns the bytes to construct a proposal suitable for Erc20
func ConstructErc20ProposalData(amount []byte, recipient []byte) []byte {
	b := bytes.Buffer{}
	b.Write(common.LeftPadBytes(amount, 32)) // amount (uint256)
	recipientLen := big.NewInt(int64(len(recipient))).Bytes()
	b.Write(common.LeftPadBytes(recipientLen, 32))
	b.Write(recipient)
	return b.Bytes()
}

// constructGenericProposalData returns the bytes to construct a generic proposal
func ConstructGenericProposalData(metadata []byte) []byte {
	var data []byte

	metadataLen := big.NewInt(int64(len(metadata)))
	data = append(data, math.PaddedBigBytes(metadataLen, 32)...) // length of metadata (uint256)
	data = append(data, metadata...)                             // metadata ([]byte)
	return data
}

// ConstructErc721ProposalData returns the bytes to construct a proposal suitable for Erc721
func ConstructErc721ProposalData(tokenId []byte, recipient []byte, metadata []byte) []byte {
	var data []byte
	data = append(data, common.LeftPadBytes(tokenId, 32)...) // tokenId ([]byte)

	recipientLen := big.NewInt(int64(len(recipient))).Bytes()
	data = append(data, common.LeftPadBytes(recipientLen, 32)...) // length of recipient
	data = append(data, recipient...)                             // recipient ([]byte)

	metadataLen := big.NewInt(int64(len(metadata))).Bytes()
	data = append(data, common.LeftPadBytes(metadataLen, 32)...) // length of metadata (uint256)
	data = append(data, metadata...)                             // metadata ([]byte)
	return data
}

// CreateProposalDataHash constructs and returns proposal data hash
// https://github.com/ChainSafe/chainbridge-celo-solidity/blob/1fae9c66a07139c277b03a09877414024867a8d9/contracts/Bridge.sol#L452-L454
func CreateProposalDataHash(data []byte, handler common.Address, msgProofOpts *msg.MsgProofOpts) common.Hash {
	b := bytes.NewBuffer(data)
	b.Write(handler.Bytes())
	b.Write(msgProofOpts.RootHash[:])
	b.Write(msgProofOpts.Key)
	b.Write(msgProofOpts.Nodes)
	b.Write(msgProofOpts.AggregatePublicKey)
	b.Write(msgProofOpts.HashedMessage[:])
	b.Write(msgProofOpts.SignatureHeader)
	return crypto.Keccak256Hash(b.Bytes())
}
