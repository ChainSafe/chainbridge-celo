// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package msg

import (
	"math/big"
	"github.com/ChainSafe/chainbridge-utils/msg"
)
 

type MsgProofOpts struct {
	RootHash [32]byte
	AggregatePublicKey []byte
	HashedMessage []byte
	Key []byte
	SignatureHeader []byte
	Nodes []byte
	G1 []byte
}

 

func NewFungibleTransfer(source, dest msg.ChainId, nonce msg.Nonce, amount *big.Int, resourceId msg.ResourceId, recipient []byte, msgProofOpts *MsgProofOpts) msg.Message {
	return msg.Message{
		Source:       source,
		Destination:  dest,
		Type:         msg.FungibleTransfer,
		DepositNonce: nonce,
		ResourceId:   resourceId,
		Payload: []interface{}{
			amount.Bytes(),
			recipient,
			msgProofOpts,
		},
	}
}

func NewNonFungibleTransfer(source, dest msg.ChainId, nonce msg.Nonce, resourceId msg.ResourceId, tokenId *big.Int, recipient, metadata []byte,  msgProofOpts *MsgProofOpts) msg.Message {
	return msg.Message{
		Source:       source,
		Destination:  dest,
		Type:         msg.NonFungibleTransfer,
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

func NewGenericTransfer(source, dest msg.ChainId, nonce msg.Nonce, resourceId msg.ResourceId, metadata []byte, msgProofOpts *MsgProofOpts) msg.Message {
	return msg.Message{
		Source:       source,
		Destination:  dest,
		Type:         msg.GenericTransfer,
		DepositNonce: nonce,
		ResourceId:   resourceId,
		Payload: []interface{}{
			metadata,
			msgProofOpts,
		},
	}
}
