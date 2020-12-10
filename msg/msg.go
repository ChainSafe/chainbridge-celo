// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package msg

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"math/big"
)

type ChainId msg.ChainId
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
	Payload      []interface{} // data associated with event sequence
}

type MsgProofOpts struct {
	RootHash           [32]byte
	AggregatePublicKey []byte
	HashedMessage      [32]byte
	Key                []byte
	SignatureHeader    []byte
	Nodes              []byte
	G1                 []byte
}

func NewFungibleTransfer(source, dest ChainId, nonce Nonce, amount *big.Int, resourceId ResourceId, recipient []byte, msgProofOpts *MsgProofOpts) Message {
	return Message{
		Source:       source,
		Destination:  dest,
		Type:         FungibleTransfer,
		DepositNonce: nonce,
		ResourceId:   resourceId,
		Payload: []interface{}{
			amount.Bytes(),
			recipient,
			msgProofOpts,
		},
	}
}

func NewNonFungibleTransfer(source, dest ChainId, nonce Nonce, resourceId ResourceId, tokenId *big.Int, recipient, metadata []byte, msgProofOpts *MsgProofOpts) Message {
	return Message{
		Source:       source,
		Destination:  dest,
		Type:         NonFungibleTransfer,
		DepositNonce: nonce,
		ResourceId:   resourceId,
		Payload: []interface{}{
			tokenId.Bytes(),
			recipient,
			metadata,
			msgProofOpts,
		},
	}
}

func NewGenericTransfer(source, dest ChainId, nonce Nonce, resourceId ResourceId, metadata []byte, msgProofOpts *MsgProofOpts) Message {
	return Message{
		Source:       source,
		Destination:  dest,
		Type:         GenericTransfer,
		DepositNonce: nonce,
		ResourceId:   resourceId,
		Payload: []interface{}{
			metadata,
			msgProofOpts,
		},
	}
}
