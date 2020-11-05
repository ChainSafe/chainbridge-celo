// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package msg

import (
	"math/big"
	"github.com/ChainSafe/chainbridge-utils/msg"
)
 
type MsgProofOpts struct {
	Source msg.ChainId
	Dest msg.ChainId
	Nonce msg.Nonce
	Amount *big.Int
	ResourceId msg.ResourceId
	Recipient []byte
	TokenId *big.Int
	Metadata []byte
	//
	RootHash [32]byte
	AggregatePublicKey []byte
	HashedMessage []byte
	Key []byte
	Data []interface{}
	SignatureHeader []byte
	Nodes []byte
	G1 []byte
}

type MessageExtraData struct {
	RootHash [32]byte
	AggregatePublicKey []byte
	HashedMessage []byte
	Key []byte
	SignatureHeader []byte
	Nodes []byte
	G1 []byte
}

 

func NewFungibleTransfer(param MsgProofOpts) msg.Message {
	return msg.Message{
		Source:       param.Source,
		Destination:  param.Dest,
		Type:         msg.FungibleTransfer,
		DepositNonce: param.Nonce,
		ResourceId:   param.ResourceId,
		Payload: []interface{}{
			param.Amount.Bytes(),
			param.Recipient,
			
			&MessageExtraData{
                RootHash: param.RootHash,
				AggregatePublicKey: param.AggregatePublicKey,
				HashedMessage: param.HashedMessage,
				Key: param.Key,
				SignatureHeader: param.SignatureHeader,
				Nodes: param.Nodes,
				G1: param.G1,
			},
		},
	}
}

func NewNonFungibleTransfer(param MsgProofOpts) msg.Message {

	return msg.Message{
		Source:       param.Source,
		Destination:  param.Dest,
		Type:         msg.NonFungibleTransfer,
		DepositNonce: param.Nonce,
		ResourceId:   param.ResourceId,
		Payload: []interface{}{
			param.TokenId.Bytes(),
			param.Recipient,
			param.Metadata,  
			//
			&MessageExtraData{
                RootHash: param.RootHash,
				AggregatePublicKey: param.AggregatePublicKey,
				HashedMessage: param.HashedMessage,
				Key: param.Key,
				SignatureHeader: param.SignatureHeader,
				Nodes: param.Nodes,
				G1: param.G1,
			},
		},
	}
}

func NewGenericTransfer(param MsgProofOpts) msg.Message {
	return msg.Message{
		Source:       param.Source,
		Destination:  param.Dest,
		Type:         msg.GenericTransfer,
		DepositNonce: param.Nonce,
		ResourceId:   param.ResourceId,
		Payload: []interface{}{
			param.Metadata,

			&MessageExtraData{
				RootHash: param.RootHash,
				AggregatePublicKey: param.AggregatePublicKey,
				HashedMessage: param.HashedMessage,
				Key: param.Key,
				SignatureHeader: param.SignatureHeader,
				Nodes: param.Nodes,
				G1: param.G1,
		  },

		},
	}
}


