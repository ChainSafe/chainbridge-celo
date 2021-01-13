// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package writer

import (
	"bytes"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/msg"
	"github.com/ethereum/go-ethereum/common"
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
	data := bytes.Buffer{}
	metadataLen := big.NewInt(int64(len(metadata))).Bytes()
	data.Write(common.LeftPadBytes(metadataLen, 32)) // length of metadata (uint256)
	data.Write(metadata)
	return data.Bytes()
}

// ConstructErc721ProposalData returns the bytes to construct a proposal suitable for Erc721
func ConstructErc721ProposalData(tokenId []byte, recipient []byte, metadata []byte) []byte {
	data := bytes.Buffer{}
	data.Write(common.LeftPadBytes(tokenId, 32))

	recipientLen := big.NewInt(int64(len(recipient))).Bytes()
	data.Write(common.LeftPadBytes(recipientLen, 32))
	data.Write(recipient)

	metadataLen := big.NewInt(int64(len(metadata))).Bytes()
	data.Write(common.LeftPadBytes(metadataLen, 32))
	data.Write(metadata)
	return data.Bytes()
}

// CreateProposalDataHash constructs and returns proposal data hash
// https://github.com/ChainSafe/chainbridge-celo-solidity/blob/1fae9c66a07139c277b03a09877414024867a8d9/contracts/Bridge.sol#L452-L454
func CreateProposalDataHash(data []byte, handler common.Address, mp *msg.MerkleProof, sv *msg.SignatureVerification) common.Hash {
	b := bytes.NewBuffer(data)
	b.Write(handler.Bytes())
	b.Write(mp.TxRootHash[:])
	b.Write(mp.Key)
	b.Write(mp.Nodes)
	b.Write(sv.AggregatePublicKey)
	b.Write(sv.BlockHash[:])
	b.Write(sv.Signature)
	return crypto.Keccak256Hash(b.Bytes())
}
