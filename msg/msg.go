// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package msg

import (
	"fmt"
	"math/big"
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

type MsgProofOpts struct {
	source ChainId
	dest ChainId
	nonce Nonce
	amount *big.Int
	resourceId ResourceId
	recipient []byte
	rootHash [32]byte
	aggregatePublicKey []byte
	hashedMessage []byte
	key []byte
	data []interface{}
	signatureHeader []byte
	nodes []byte
	g1 []byte
	tokenId *big.Int
	metadata []byte
}

// Message is used as a generic format to communicate between chains
type Message struct {
	Source             ChainId // Source where message was initiated
	Destination        ChainId // Destination chain of message
	ResourceId         ResourceId
	Type               TransferType  // type of bridge transfer
	DepositNonce       Nonce         // Nonce for the deposit
	Payload            []interface{} // data associated with event sequence
	RootHash           [32]byte
	AggregatePublicKey []byte
	HashedMessage      []byte
	Key                []byte
	Data               []interface{}
	SignatureHeader    []byte
	Nodes              []byte
	G1                 []byte
}

func NewFungibleTransfer(param MsgProofOpts) Message {
	return Message{
		Source:       param.source,
		Destination:  param.dest,
		Type:         FungibleTransfer,
		DepositNonce: param.nonce,
		ResourceId:   param.resourceId,
		Payload: []interface{}{
			param.amount.Bytes(),
			param.recipient,
		},
		RootHash:           param.rootHash,
		AggregatePublicKey: param.aggregatePublicKey,
		HashedMessage:      param.hashedMessage,
		Key:                param.key,
		Data:               param.data,
		SignatureHeader:    param.signatureHeader,
		Nodes:              param.nodes,
		G1:                 param.g1,
	}
}

func NewNonFungibleTransfer(param MsgProofOpts) Message {

	return Message{
		Source:       param.source,
		Destination:  param.dest,
		Type:         NonFungibleTransfer,
		DepositNonce: param.nonce,
		ResourceId:   param.resourceId,
		Payload: []interface{}{
			param.tokenId.Bytes(),
			param.recipient,
			param.metadata,
		},
		RootHash:           param.rootHash,
		AggregatePublicKey: param.aggregatePublicKey,
		HashedMessage:      param.hashedMessage,
		Key:                param.key,
		Data:               param.data,
		SignatureHeader:    param.signatureHeader,
		Nodes:              param.nodes,
		G1:                 param.g1,
	}
}

func NewGenericTransfer(param MsgProofOpts) Message {
	return Message{
		Source:       param.source,
		Destination:  param.dest,
		Type:         GenericTransfer,
		DepositNonce: param.nonce,
		ResourceId:   param.resourceId,
		Payload: []interface{}{
			param.metadata,
		},
		RootHash:           param.rootHash,
		AggregatePublicKey: param.aggregatePublicKey,
		HashedMessage:      param.hashedMessage,
		Key:                param.key,
		Data:               param.data,
		SignatureHeader:    param.signatureHeader,
		Nodes:              param.nodes,
		G1:                 param.g1,
	}
}

func ResourceIdFromSlice(in []byte) ResourceId {
	var res ResourceId
	copy(res[:], in)
	return res
}
