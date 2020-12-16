// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package listener

import (
	"github.com/ChainSafe/chainbridge-celo/msg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rs/zerolog/log"
)

func (l *listener) handleErc20DepositedEvent(destId msg.ChainId, nonce msg.Nonce) (*msg.Message, error) {
	log.Info().Interface("dest", destId).Interface("nonce", nonce).Msg("Handling fungible deposit event")
	//TODO no call opts. should have From in original chainbridge.
	record, err := l.erc20HandlerContract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
	if err != nil {
		log.Error().Err(err).Msg("Error Unpacking ERC20 Deposit Record")
		return nil, err
	}

	return msg.NewFungibleTransfer(
		l.cfg.ID,
		destId,
		nonce,
		record.Amount,
		record.ResourceID,
		record.DestinationRecipientAddress,
	), nil
}

func (l *listener) handleErc721DepositedEvent(destId msg.ChainId, nonce msg.Nonce) (*msg.Message, error) {
	log.Info().Msg("Handling nonfungible deposit event")
	//TODO no call opts. should have From in original chainbridge.
	record, err := l.erc721HandlerContract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
	if err != nil {
		log.Error().Err(err).Msg("Error Unpacking ERC721 Deposit Record")
		return nil, err
	}

	return msg.NewNonFungibleTransfer(
		l.cfg.ID,
		destId,
		nonce,
		record.ResourceID,
		record.TokenID,
		record.DestinationRecipientAddress,
		record.MetaData,
	), nil
}

func (l *listener) handleGenericDepositedEvent(destId msg.ChainId, nonce msg.Nonce) (*msg.Message, error) {
	log.Info().Msg("Handling generic deposit event")
	//TODO no call opts. should have From in original chainbridge.
	record, err := l.genericHandlerContract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
	if err != nil {
		log.Error().Err(err).Msg("Error Unpacking Generic Deposit Record")
		return nil, err
	}

	return msg.NewGenericTransfer(
		l.cfg.ID,
		destId,
		nonce,
		record.ResourceID,
		record.MetaData[:],
	), nil
}
