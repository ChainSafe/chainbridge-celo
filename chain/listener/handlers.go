// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package listener

import (
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rs/zerolog/log"
	"github.com/status-im/keycard-go/hexutils"
)

func (l *listener) handleErc20DepositedEvent(destId utils.ChainId, nonce utils.Nonce) (*utils.Message, error) {
	record, err := l.erc20HandlerContract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
	if err != nil {
		log.Error().Err(err).Msg("Error Unpacking ERC20 Deposit Record")
		return nil, err
	}

	log.Info().Interface("dest", destId).Interface("nonce", nonce).Str("resourceID", hexutils.BytesToHex(record.ResourceID[:])).Msg("Handling fungible deposit event")
	return utils.NewFungibleTransfer(
		l.cfg.ID,
		destId,
		nonce,
		record.ResourceID,
		nil,
		nil,
		record.Amount,
		record.DestinationRecipientAddress,
	), nil
}

func (l *listener) handleErc721DepositedEvent(destId utils.ChainId, nonce utils.Nonce) (*utils.Message, error) {
	//TODO no call opts. should have From in original chainbridge.
	record, err := l.erc721HandlerContract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
	if err != nil {
		log.Error().Err(err).Msg("Error Unpacking ERC721 Deposit Record")
		return nil, err
	}
	log.Info().Interface("dest", destId).Interface("nonce", nonce).Str("resourceID", hexutils.BytesToHex(record.ResourceID[:])).Msg("Handling nonfungible deposit event")
	return utils.NewNonFungibleTransfer(
		l.cfg.ID,
		destId,
		nonce,
		record.ResourceID,
		nil,
		nil,
		record.TokenID,
		record.DestinationRecipientAddress,
		record.MetaData,
	), nil
}

func (l *listener) handleGenericDepositedEvent(destId utils.ChainId, nonce utils.Nonce) (*utils.Message, error) {
	record, err := l.genericHandlerContract.GetDepositRecord(&bind.CallOpts{}, uint64(nonce), uint8(destId))
	if err != nil {
		log.Error().Err(err).Msg("Error Unpacking Generic Deposit Record")
		return nil, err
	}
	log.Info().Interface("dest", destId).Interface("nonce", nonce).Str("resourceID", hexutils.BytesToHex(record.ResourceID[:])).Msg("Handling generic deposit event")
	return utils.NewGenericTransfer(
		l.cfg.ID,
		destId,
		nonce,
		record.ResourceID,
		nil,
		nil,
		record.MetaData[:],
	), nil
}
