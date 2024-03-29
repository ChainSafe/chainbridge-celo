// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package writer

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/chain/config"
	"github.com/ChainSafe/chainbridge-celo/utils"
	metrics "github.com/ChainSafe/chainbridge-utils/metrics/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

var ProposalNotPassedStatus uint8 = 1
var ProposalStatusPassed uint8 = 2
var ProposalStatusTransferred uint8 = 3
var ProposalStatusCancelled uint8 = 4
var BlockRetryLimit = 5

type writer struct {
	cfg            *config.CeloChainConfig
	client         ContractCaller
	bridgeContract Bridger
	stop           <-chan struct{}
	sysErr         chan<- error
	metrics        *metrics.ChainMetrics
}

type Bridger interface {
	GetProposal(opts *bind.CallOpts, originChainID uint8, depositNonce uint64, dataHash [32]byte) (Bridge.BridgeProposal, error)
	HasVotedOnProposal(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte, arg2 common.Address) (bool, error)
	VoteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, resourceID [32]byte, dataHash [32]byte) (*types.Transaction, error)
	ExecuteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, data []byte, resourceID [32]byte, signatureHeader []byte, aggregatePublicKey []byte, hashedMessage [32]byte, rootHash [32]byte, key []byte, nodes []byte) (*types.Transaction, error)
}

type ContractCaller interface {
	client.LogFilterWithLatestBlock
	CallOpts() *bind.CallOpts
	Opts() *bind.TransactOpts
	LockAndUpdateOpts() error
	UnlockOpts()
	WaitForBlock(block *big.Int) error
}

// NewWriter creates and returns writer
func NewWriter(client ContractCaller, cfg *config.CeloChainConfig, stop <-chan struct{}, sysErr chan<- error, m *metrics.ChainMetrics) *writer {
	return &writer{
		cfg:     cfg,
		client:  client,
		stop:    stop,
		sysErr:  sysErr,
		metrics: m,
	}
}

// setContract adds the bound receiver bridgeContract to the writer
func (w *writer) SetBridge(bridge Bridger) {
	w.bridgeContract = bridge
}

// ResolveMessage handles any given message based on type
// A bool is returned to indicate failure/success
// this should be ignored except for within tests.
func (w *writer) ResolveMessage(m *utils.Message) bool {
	log.Info().Str("type", string(m.Type)).Interface("src", m.Source).Interface("dst", m.Destination).Interface("nonce", m.DepositNonce).Str("rId", m.ResourceId.Hex()).Msg("Attempting to resolve message")
	var data []byte
	var handlerContract common.Address
	var err error
	switch m.Type {
	case utils.FungibleTransfer:
		data, err = w.createERC20ProposalData(m)
		handlerContract = w.cfg.Erc20HandlerContract
	case utils.NonFungibleTransfer:
		data, err = w.createErc721ProposalData(m)
		handlerContract = w.cfg.Erc721HandlerContract
	case utils.GenericTransfer:
		data, err = w.createGenericDepositProposalData(m)
		handlerContract = w.cfg.GenericHandlerContract
	default:
		log.Error().Str("type", string(m.Type)).Msg("Unknown message type received")
		return false
	}
	if err != nil {
		log.Error().Err(err)
		return false
	}
	dataHash := CreateProposalDataHash(data, handlerContract, m.MPParams, m.SVParams)

	if !w.shouldVote(m, dataHash) {
		if w.proposalIsPassed(m.Source, m.DepositNonce, dataHash) {
			// We should not vote for this proposal but it is ready to be executed
			w.executeProposal(m, data, dataHash)
			return true
		} else {
			return false
		}
	}
	// Capture latest block so when know where to watch from
	latestBlock, err := w.client.LatestBlock()
	if err != nil {
		log.Error().Err(err).Msg("unable to fetch latest block")
		return false
	}

	// watch for execution event
	go w.watchThenExecute(m, data, dataHash, latestBlock)

	w.voteProposal(m, dataHash)

	return true
}
