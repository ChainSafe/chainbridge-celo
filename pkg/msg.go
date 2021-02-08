// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package pkg

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ChainId uint8
type TransferType string
type ResourceId [32]byte

func (r ResourceId) Hex() string {
	return fmt.Sprintf("%x", r)
}

type Nonce uint64

func (n Nonce) Big() *big.Int {
	return big.NewInt(int64(n))
}

var FungibleTransfer TransferType = "FungibleTransfer"
var NonFungibleTransfer TransferType = "NonFungibleTransfer"
var GenericTransfer TransferType = "GenericTransfer"

// Message is used as a generic format to communicate between chains
type Message struct {
	Source       ChainId      // Source where message was initiated
	Destination  ChainId      // Destination chain of message
	Type         TransferType // type of bridge transfer
	DepositNonce Nonce        // Nonce for the deposit
	ResourceId   ResourceId
	MPParams     *MerkleProof
	SVParams     *SignatureVerification
	Payload      []interface{} // data associated with event sequence
}

type MerkleProof struct {
	TxRootHash common.Hash // Expected root of trie, in our case should be transactionsRoot from block
	Key        []byte      // RLP encoding of tx index, for the tx we want to prove
	Nodes      []byte      // The actual proof, all the nodes of the trie that between leaf value and root
}

type SignatureVerification struct {
	AggregatePublicKey []byte      // Aggregated public key of block validators
	BlockHash          common.Hash // Hash of block we are proving
	Signature          []byte      // Signature of block we are proving
}

func NewFungibleTransfer(source, dest ChainId, nonce Nonce, resourceId ResourceId, mp *MerkleProof, sv *SignatureVerification, amount *big.Int, recipient []byte) *Message {
	return &Message{
		Source:       source,
		Destination:  dest,
		Type:         FungibleTransfer,
		DepositNonce: nonce,
		ResourceId:   resourceId,
		MPParams:     mp,
		SVParams:     sv,
		Payload: []interface{}{
			amount.Bytes(),
			recipient,
		},
	}
}

func NewNonFungibleTransfer(source, dest ChainId, nonce Nonce, resourceId ResourceId, mp *MerkleProof, sv *SignatureVerification, tokenId *big.Int, recipient, metadata []byte) *Message {
	return &Message{
		Source:       source,
		Destination:  dest,
		Type:         NonFungibleTransfer,
		DepositNonce: nonce,
		ResourceId:   resourceId,
		MPParams:     mp,
		SVParams:     sv,
		Payload: []interface{}{
			tokenId.Bytes(),
			recipient,
			metadata,
		},
	}
}

func NewGenericTransfer(source, dest ChainId, nonce Nonce, resourceId ResourceId, mp *MerkleProof, sv *SignatureVerification, metadata []byte) *Message {
	return &Message{
		Source:       source,
		Destination:  dest,
		Type:         GenericTransfer,
		DepositNonce: nonce,
		ResourceId:   resourceId,
		MPParams:     mp,
		SVParams:     sv,
		Payload: []interface{}{
			metadata,
		},
	}
}
