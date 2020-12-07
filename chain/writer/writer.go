// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package writer

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	"github.com/ChainSafe/chainbridge-celo/chain"
	"github.com/ChainSafe/chainbridge-utils/core"
	metrics "github.com/ChainSafe/chainbridge-utils/metrics/types"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rs/zerolog/log"
)

var _ core.Writer = &writer{}

var PassedStatus uint8 = 2
var TransferredStatus uint8 = 3
var CancelledStatus uint8 = 4
var BlockRetryLimit = 5

type writer struct {
	cfg            *chain.CeloChainConfig
	client         ContractCaller
	bridgeContract *Bridge.Bridge
	stop           <-chan struct{}
	sysErr         chan<- error
	metrics        *metrics.ChainMetrics
}

type ContractCaller interface {
	chain.LogFilterWithLatestBlock
	CallOpts() *bind.CallOpts
	Opts() *bind.TransactOpts
	LockAndUpdateOpts() error
	UnlockOpts()
	WaitForBlock(block *big.Int) error
}

// NewWriter creates and returns writer
func NewWriter(client ContractCaller, cfg *chain.CeloChainConfig, stop <-chan struct{}, sysErr chan<- error, m *metrics.ChainMetrics) *writer {
	return &writer{
		cfg:     cfg,
		client:  client,
		stop:    stop,
		sysErr:  sysErr,
		metrics: m,
	}
}

// setContract adds the bound receiver bridgeContract to the writer
func (w *writer) SetBridge(bridge *Bridge.Bridge) {
	w.bridgeContract = bridge
}

// ResolveMessage handles any given message based on type
// A bool is returned to indicate failure/success
// this should be ignored except for within tests.
func (w *writer) ResolveMessage(m *msg.Message) bool {
	log.Info().Str("type", string(m.Type)).Interface("src", m.Source).Interface("dst", m.Destination).Interface("nonce", m.DepositNonce).Str("rId", m.ResourceId.Hex()).Msg("Attempting to resolve message")
	var dataHash common.Hash
	var data []byte
	var err error
	switch m.Type {
	case msg.FungibleTransfer:
		data, dataHash, err = w.createERC20ProposalDataAndHash(m)
	case msg.NonFungibleTransfer:
		data, dataHash, err = w.createErc721ProposalDataAndHash(m)
	case msg.GenericTransfer:
		data, dataHash, err = w.createGenericDepositProposalDataAndHash(m)
	default:
		log.Error().Str("type", string(m.Type)).Msg("Unknown message type received")
		return false
	}
	if err != nil {
		log.Error().Err(err)
		return false
	}

	if !w.shouldVote(m, dataHash) {
		return false
	}
	// Capture latest block so when know where to watch from
	latestBlock, err := w.client.LatestBlock()
	if err != nil {
		log.Error().Err(err).Msg("unable to fetch latest block")
		return false
	}

	// watch for execution event
	go w.watchThenExecute(m, data, dataHash, latestBlock, msgProofOpts)

	w.voteProposal(m, dataHash)

	return true
}
