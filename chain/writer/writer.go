// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package writer

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	"github.com/ChainSafe/chainbridge-celo/chain"
	"github.com/ChainSafe/chainbridge-utils/core"
	metrics "github.com/ChainSafe/chainbridge-utils/metrics/types"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

var _ core.Writer = &writer{}

var PassedStatus uint8 = 2
var TransferredStatus uint8 = 3
var CancelledStatus uint8 = 4

type writer struct {
	cfg            *chain.CeloChainConfig
	conn           ConnectionWriter
	bridgeContract *Bridge.Bridge
	stop           <-chan int
	sysErr         chan<- error
	metrics        *metrics.ChainMetrics
}

type ConnectionWriter interface {
	LatestBlock() (*big.Int, error)
	LockAndUpdateOpts() error
	Opts() *bind.TransactOpts
	UnlockOpts()
	CallOpts() *bind.CallOpts
	WaitForBlock(block *big.Int) error
	Client() *ethclient.Client
}

// NewWriter creates and returns writer
func NewWriter(conn ConnectionWriter, cfg *chain.CeloChainConfig, stop <-chan int, sysErr chan<- error, m *metrics.ChainMetrics) *writer {
	return &writer{
		cfg:     cfg,
		conn:    conn,
		stop:    stop,
		sysErr:  sysErr,
		metrics: m,
	}
}

func (w *writer) Start() error {
	log.Debug().Msg("Starting celo writer...")
	return nil
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
