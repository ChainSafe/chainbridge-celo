// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"context"
	"errors"
	"math/big"
	"time"

	celoMsg "github.com/ChainSafe/chainbridge-celo/msg"
	utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/ethereum/go-ethereum/log"
)

// Number of blocks to wait for an finalization event
const ExecuteBlockWatchLimit = 100

// Time between retrying a failed tx
const TxRetryInterval = time.Second * 2

// Time between retrying a failed tx
const TxRetryLimit = 10

var ErrNonceTooLow = errors.New("nonce too low")
var ErrTxUnderpriced = errors.New("replacement transaction underpriced")
var ErrFatalTx = errors.New("submission of transaction failed")
var ErrFatalQuery = errors.New("query of chain state failed")

// proposalIsComplete returns true if the proposal state is either Passed, Transferred or Cancelled
func (w *writer) proposalIsComplete(srcId msg.ChainId, nonce msg.Nonce, dataHash [32]byte) bool {
	prop, err := w.bridgeContract.GetProposal(w.conn.CallOpts(), uint8(srcId), uint64(nonce), dataHash)
	if err != nil {
		w.log.Error("Failed to check proposal existence", "err", err)
		return false
	}
	return prop.Status == PassedStatus || prop.Status == TransferredStatus || prop.Status == CancelledStatus
}

// proposalIsFinalized returns true if the proposal state is Transferred or Cancelled
func (w *writer) proposalIsFinalized(srcId msg.ChainId, nonce msg.Nonce, dataHash [32]byte) bool {
	prop, err := w.bridgeContract.GetProposal(w.conn.CallOpts(), uint8(srcId), uint64(nonce), dataHash)

	if err != nil {
		w.log.Error("Failed to check proposal existence", "err", err)
		return false
	}
	return prop.Status == TransferredStatus || prop.Status == CancelledStatus
}

// hasVoted checks if this relayer has already voted
func (w *writer) hasVoted(srcId msg.ChainId, nonce msg.Nonce, dataHash [32]byte) bool {
	hasVoted, err := w.bridgeContract.HasVotedOnProposal(w.conn.CallOpts(), utils.IDAndNonce(srcId, nonce), dataHash, w.conn.Opts().From)

	if err != nil {
		w.log.Error("Failed to check proposal existance", "err", err)
		return false
	}

	return hasVoted
}

func (w *writer) shouldVote(m msg.Message, dataHash [32]byte) bool {
	// Check if proposal has passed and skip if Passed or Transferred
	if w.proposalIsComplete(m.Source, m.DepositNonce, dataHash) {
		w.log.Info("Proposal complete, not voting", "src", m.Source, "nonce", m.DepositNonce)
	}

	// Check if relayer has previously voted
	if w.hasVoted(m.Source, m.DepositNonce, dataHash) {
		w.log.Info("Relayer has already voted, not voting", "src", m.Source, "source", m.DepositNonce)
		return false
	}

	return true
}

// createErc20Proposal creates an Erc20 proposal.
// Returns true if the proposal is successfully created or is complete
func (w *writer) createErc20Proposal(m msg.Message) bool {
	w.log.Info("Creating erc20 proposal", "src", m.Source, "nonce", m.DepositNonce)

	msgProofOptsInterface := m.Payload[2]

	if msgProofOptsInterface == nil {
		w.log.Error("msgProofOpts cannot be nil")
		return false
	}

	msgProofOpts, ok := msgProofOptsInterface.(*celoMsg.MsgProofOpts)

	if !ok {
		w.log.Error("unable to convert msgProofOptsInterface to *MsgProofOpts")
		return false
	}

	data := ConstructErc20ProposalData(m.Payload[0].([]byte), m.Payload[1].([]byte))
	dataHash := CreateProposalDataHash(data, w.cfg.erc20HandlerContract, msgProofOpts)

	if !w.shouldVote(m, dataHash) {
		return false
	}
	// Capture latest block so when know where to watch from
	latestBlock, err := w.conn.LatestBlock()
	if err != nil {
		w.log.Error("unable to fetch latest block", "err", err)
		return false
	}

	// watch for execution event
	go w.watchThenExecute(m, data, dataHash, latestBlock, msgProofOpts)

	w.voteProposal(m, dataHash)

	return true
}

// createErc721Proposal creates an Erc721 proposal.
// Returns true if the proposal is succesfully created or is complete
func (w *writer) createErc721Proposal(m msg.Message) bool {
	w.log.Info("Creating erc721 proposal", "src", m.Source, "nonce", m.DepositNonce)

	msgProofOptsInterface := m.Payload[3]

	if msgProofOptsInterface == nil {
		w.log.Error("msgProofOpts cannot be nil")
		return false
	}

	msgProofOpts, ok := msgProofOptsInterface.(*celoMsg.MsgProofOpts)

	if !ok {
		w.log.Error("unable to convert msgProofOptsInterface to *MsgProofOpts")
		return false
	}

	data := ConstructErc721ProposalData(m.Payload[0].([]byte), m.Payload[1].([]byte), m.Payload[2].([]byte))
	dataHash := CreateProposalDataHash(data, w.cfg.erc721HandlerContract, msgProofOpts)

	if !w.shouldVote(m, dataHash) {
		return false
	}

	// Capture latest block so we know where to watch from
	latestBlock, err := w.conn.LatestBlock()
	if err != nil {
		w.log.Error("Unable to fetch latest block", "err", err)
		return false
	}

	// watch for execution event
	go w.watchThenExecute(m, data, dataHash, latestBlock, msgProofOpts)

	w.voteProposal(m, dataHash)

	return true

}

// createGenericDepositProposal creates a generic proposal
// returns true if the proposal is complete or is succesfully created
func (w *writer) createGenericDepositProposal(m msg.Message) bool {
	w.log.Info("Creating generic proposal", "src", m.Source, "nonce", m.DepositNonce)

	metadata, ok := m.Payload[0].([]byte)

	if !ok {
		w.log.Error("Unable to convert metadata to []byte")
		return false
	}

	msgProofOptsInterface := m.Payload[1]

	if msgProofOptsInterface == nil {
		w.log.Error("msgProofOpts cannot be nil")
		return false
	}

	msgProofOpts, ok := msgProofOptsInterface.(*celoMsg.MsgProofOpts)

	if !ok {
		w.log.Error("unable to convert msgProofOptsInterface to *msgProofOpts")
		return false
	}

	data := ConstructGenericProposalData(metadata)
	dataHash := CreateProposalDataHash(data, w.cfg.genericHandlerContract, msgProofOpts)

	if !w.shouldVote(m, dataHash) {
		return false
	}

	// Capture latest block so when know where to watch from
	latestBlock, err := w.conn.LatestBlock()
	if err != nil {
		w.log.Error("Unable to fetch latest block", "err", err)
		return false
	}

	// watch for execution event
	go w.watchThenExecute(m, data, dataHash, latestBlock, msgProofOpts)

	w.voteProposal(m, dataHash)

	return true
}

// watchThenExecute watches for the latest block and executes once the matching finalized event is found
func (w *writer) watchThenExecute(m msg.Message, data []byte, dataHash [32]byte, latestBlock *big.Int, msgProofOpts *celoMsg.MsgProofOpts) {
	w.log.Info("Watching for finalization event", "src", m.Source, "nonce", m.DepositNonce)

	// watching for the latest block, querying and matching the finalized event will be retried up to ExecuteBlockWatchLimit times
	for i := 0; i < ExecuteBlockWatchLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			// watch for the lastest block, retry up to BlockRetryLimit times
			for waitRetrys := 0; waitRetrys < BlockRetryLimit; waitRetrys++ {
				err := w.conn.WaitForBlock(latestBlock)
				if err != nil {
					w.log.Error("Waiting for block failed", "err", err)
					// Exit if retries exceeded
					if waitRetrys+1 == BlockRetryLimit {
						w.log.Error("Waiting for block retries exceeded, shutting down")
						w.sysErr <- ErrFatalQuery
						return
					}
				} else {
					break
				}
			}

			// query for logs
			query := buildQuery(w.cfg.bridgeContract, utils.ProposalEvent, latestBlock, latestBlock)
			evts, err := w.conn.Client().FilterLogs(context.Background(), query)
			if err != nil {
				w.log.Error("Failed to fetch logs", "err", err)
				return
			}

			// execute the proposal once we find the matching finalized event
			for _, evt := range evts {
				sourceId := evt.Topics[1].Big().Uint64()
				depositNonce := evt.Topics[2].Big().Uint64()
				status := evt.Topics[3].Big().Uint64()

				if m.Source == msg.ChainId(sourceId) &&
					m.DepositNonce.Big().Uint64() == depositNonce &&
					utils.IsFinalized(uint8(status)) {
					w.executeProposal(m, data, dataHash, msgProofOpts)
					return
				} else {
					w.log.Trace("Ignoring event", "src", sourceId, "nonce", depositNonce)
				}
			}
			w.log.Trace("No finalization event found in current block", "block", latestBlock, "src", m.Source, "nonce", m.DepositNonce)
			latestBlock = latestBlock.Add(latestBlock, big.NewInt(1))
		}
	}
	log.Warn("Block watch limit exceeded, skipping execution", "source", m.Source, "dest", m.Destination, "nonce", m.DepositNonce)
}

// voteProposal submits a vote proposal
// a vote proposal will try to be submitted up to the TxRetryLimit times
func (w *writer) voteProposal(m msg.Message, dataHash [32]byte) {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			err := w.conn.LockAndUpdateOpts()
			if err != nil {
				w.log.Error("Failed to update tx opts", "err", err)
				continue
			}

			tx, err := w.bridgeContract.VoteProposal(
				w.conn.Opts(),
				uint8(m.Source),
				uint64(m.DepositNonce),
				m.ResourceId,
				dataHash,
			)
			w.conn.UnlockOpts()

			if err == nil {
				w.log.Info("Submitted proposal vote", "tx", tx.Hash(), "src", m.Source, "depositNonce", m.DepositNonce)
				if w.metrics != nil {
					w.metrics.VotesSubmitted.Inc()
				}
				return
			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				w.log.Debug("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				w.log.Warn("Voting failed", "source", m.Source, "dest", m.Destination, "depositNonce", m.DepositNonce, "err", err)
				time.Sleep(TxRetryInterval)
			}

			// Verify proposal is still open for voting, otherwise no need to retry
			if w.proposalIsComplete(m.Source, m.DepositNonce, dataHash) {
				w.log.Info("Proposal voting complete on chain", "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce)
				return
			}
		}
	}
	w.log.Error("Submission of Vote transaction failed", "source", m.Source, "dest", m.Destination, "depositNonce", m.DepositNonce)
	w.sysErr <- ErrFatalTx
}

// executeProposal executes the proposal
func (w *writer) executeProposal(m msg.Message, data []byte, dataHash [32]byte, msgProofOpts *celoMsg.MsgProofOpts) {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			err := w.conn.LockAndUpdateOpts()
			if err != nil {
				w.log.Error("Failed to update nonce", "err", err)
				return
			}

			tx, err := w.bridgeContract.ExecuteProposal(
				w.conn.Opts(),
				uint8(m.Source),
				uint64(m.DepositNonce),
				data,
				m.ResourceId,
				//
				msgProofOpts.SignatureHeader,
				msgProofOpts.AggregatePublicKey,
				msgProofOpts.G1,
				msgProofOpts.HashedMessage,
				msgProofOpts.RootHash,
				msgProofOpts.Key,
				msgProofOpts.Nodes,
			)
			w.conn.UnlockOpts()

			if err == nil {
				w.log.Info("Submitted proposal execution", "tx", tx.Hash(), "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce)
				return
			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				w.log.Error("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				w.log.Warn("Execution failed, proposal may already be complete", "err", err)
				time.Sleep(TxRetryInterval)
			}

			// Verify proposal is still open for execution, tx will fail if we aren't the first to execute,
			// but there is no need to retry
			if w.proposalIsFinalized(m.Source, m.DepositNonce, dataHash) {
				w.log.Info("Proposal finalized on chain", "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce)
				return
			}
		}
	}
	w.log.Error("Submission of Execute transaction failed", "source", m.Source, "dest", m.Destination, "depositNonce", m.DepositNonce)
	w.sysErr <- ErrFatalTx
}