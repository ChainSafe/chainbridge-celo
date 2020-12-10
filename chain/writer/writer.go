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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

var ProposalStatusPassed uint8 = 2
var ProposalStatusTransferred uint8 = 3
var ProposalStatusCancelled uint8 = 4
var BlockRetryLimit = 5

type writer struct {
	cfg            *chain.CeloChainConfig
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
	ExecuteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, data []byte, resourceID [32]byte, signatureHeader []byte, aggregatePublicKey []byte, g1 []byte, hashedMessage []byte, rootHash [32]byte, key []byte, nodes []byte) (*types.Transaction, error)
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
func (w *writer) SetBridge(bridge Bridger) {
	w.bridgeContract = bridge
}

func buildMsgOpts()

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
		data, err = w.createERC20ProposalDataAndHash(m)
	case msg.NonFungibleTransfer:
		data, err = w.createErc721ProposalDataAndHash(m)
	case msg.GenericTransfer:
		data, err = w.createGenericDepositProposalDataAndHash(m)
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
	go w.watchThenExecute(m, data, dataHash, latestBlock)

	w.voteProposal(m, dataHash)

	return true
}
