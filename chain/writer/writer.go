// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package writer

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	"github.com/ChainSafe/chainbridge-celo/chain"
	"github.com/ChainSafe/chainbridge-celo/msg"
	metrics "github.com/ChainSafe/chainbridge-utils/metrics/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rs/zerolog/log"
)

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
func (w *writer) ResolveMessage(m msg.Message) bool {
	log.Info().Str("type", string(m.Type)).Interface("src", m.Source).Interface("dst", m.Destination).Interface("nonce", m.DepositNonce).Str("rId", m.ResourceId.Hex()).Msg("Attempting to resolve message")
	switch m.Type {
	case msg.FungibleTransfer:
		return w.createErc20Proposal(m)
	case msg.NonFungibleTransfer:
		return w.createErc721Proposal(m)
	case msg.GenericTransfer:
		return w.createGenericDepositProposal(m)
	default:
		log.Error().Str("type", string(m.Type)).Msg("Unknown message type received")
		return false
	}
}
