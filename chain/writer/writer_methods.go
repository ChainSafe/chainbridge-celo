// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package writer

import (
	"context"
	"errors"
	"math/big"
	"time"

	celoMsg "github.com/ChainSafe/chainbridge-celo/msg"
	utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
	"github.com/ChainSafe/chainbridge-utils/msg"
	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
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
//TODO: unsderstand CallOpt deeply
func (w *writer) proposalIsComplete(srcId msg.ChainId, nonce msg.Nonce, dataHash [32]byte) bool {
	prop, err := w.bridgeContract.GetProposal(w.client.CallOpts(), uint8(srcId), uint64(nonce), dataHash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check proposal existence")
		return false
	}
	return prop.Status == PassedStatus || prop.Status == TransferredStatus || prop.Status == CancelledStatus
}

// proposalIsFinalized returns true if the proposal state is Transferred or Cancelled
func (w *writer) proposalIsFinalized(srcId msg.ChainId, nonce msg.Nonce, dataHash [32]byte) bool {
	prop, err := w.bridgeContract.GetProposal(w.client.CallOpts(), uint8(srcId), uint64(nonce), dataHash)

	if err != nil {
		log.Error().Err(err).Msg("Failed to check proposal existence")
		return false
	}
	return prop.Status == TransferredStatus || prop.Status == CancelledStatus
}

// hasVoted checks if this relayer has already voted
func (w *writer) hasVoted(srcId msg.ChainId, nonce msg.Nonce, dataHash [32]byte) bool {
	hasVoted, err := w.bridgeContract.HasVotedOnProposal(w.client.CallOpts(), utils.IDAndNonce(srcId, nonce), dataHash, w.client.Opts().From)

	if err != nil {
		log.Error().Err(err).Msg("Failed to check proposal existance")
		return false
	}

	return hasVoted
}

func (w *writer) shouldVote(m msg.Message, dataHash [32]byte) bool {
	// Check if proposal has passed and skip if Passed or Transferred
	if w.proposalIsComplete(m.Source, m.DepositNonce, dataHash) {
		log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Proposal complete, not voting")
	}

	// Check if relayer has previously voted
	if w.hasVoted(m.Source, m.DepositNonce, dataHash) {
		log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Relayer has already voted, not voting")
		return false
	}

	return true
}

// createErc20Proposal creates an Erc20 proposal.
// Returns true if the proposal is successfully created or is complete
func (w *writer) createErc20Proposal(m msg.Message) bool {
	log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Creating erc20 proposal")

	msgProofOptsInterface := m.Payload[2]

	if msgProofOptsInterface == nil {
		log.Error().Msg("msgProofOpts cannot be nil")
		return false
	}

	msgProofOpts, ok := msgProofOptsInterface.(*celoMsg.MsgProofOpts)

	if !ok {
		log.Error().Msg("unable to convert msgProofOptsInterface to *MsgProofOpts")
		return false
	}

	data := ConstructErc20ProposalData(m.Payload[0].([]byte), m.Payload[1].([]byte))
	dataHash := CreateProposalDataHash(data, w.cfg.Erc20HandlerContract, msgProofOpts)

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

// createErc721Proposal creates an Erc721 proposal.
// Returns true if the proposal is succesfully created or is complete
func (w *writer) createErc721Proposal(m msg.Message) bool {
	log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Creating erc721 proposal")

	msgProofOptsInterface := m.Payload[3]

	if msgProofOptsInterface == nil {
		log.Error().Msg("msgProofOpts cannot be nil")
		return false
	}

	msgProofOpts, ok := msgProofOptsInterface.(*celoMsg.MsgProofOpts)

	if !ok {
		log.Error().Msg("unable to convert msgProofOptsInterface to *MsgProofOpts")
		return false
	}

	data := ConstructErc721ProposalData(m.Payload[0].([]byte), m.Payload[1].([]byte), m.Payload[2].([]byte))
	dataHash := CreateProposalDataHash(data, w.cfg.Erc721HandlerContract, msgProofOpts)

	if !w.shouldVote(m, dataHash) {
		return false
	}

	// Capture latest block so we know where to watch from
	latestBlock, err := w.client.LatestBlock()
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch latest block")
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
	log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Creating generic proposal")

	metadata, ok := m.Payload[0].([]byte)

	if !ok {
		log.Error().Msg("Unable to convert metadata to []byte")
		return false
	}

	msgProofOptsInterface := m.Payload[1]

	if msgProofOptsInterface == nil {
		log.Error().Msg("msgProofOpts cannot be nil")
		return false
	}

	msgProofOpts, ok := msgProofOptsInterface.(*celoMsg.MsgProofOpts)

	if !ok {
		log.Error().Msg("unable to convert msgProofOptsInterface to *msgProofOpts")
		return false
	}

	data := ConstructGenericProposalData(metadata)
	dataHash := CreateProposalDataHash(data, w.cfg.GenericHandlerContract, msgProofOpts)

	if !w.shouldVote(m, dataHash) {
		return false
	}

	// Capture latest block so when know where to watch from
	latestBlock, err := w.client.LatestBlock()
	if err != nil {
		log.Error().Err(err).Msg("Unable to fetch latest block")
		return false
	}

	// watch for execution event
	go w.watchThenExecute(m, data, dataHash, latestBlock, msgProofOpts)

	w.voteProposal(m, dataHash)

	return true
}

// watchThenExecute watches for the latest block and executes once the matching finalized event is found
func (w *writer) watchThenExecute(m msg.Message, data []byte, dataHash [32]byte, latestBlock *big.Int, msgProofOpts *celoMsg.MsgProofOpts) {
	log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Watching for finalization event")

	// watching for the latest block, querying and matching the finalized event will be retried up to ExecuteBlockWatchLimit times
	for i := 0; i < ExecuteBlockWatchLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			// watch for the lastest block, retry up to BlockRetryLimit times
			for waitRetrys := 0; waitRetrys < BlockRetryLimit; waitRetrys++ {
				err := w.client.WaitForBlock(latestBlock)
				if err != nil {
					log.Error().Err(err).Msg("Waiting for block failed")
					// Exit if retries exceeded
					if waitRetrys+1 == BlockRetryLimit {
						log.Error().Err(err).Msg("Waiting for block retries exceeded, shutting down")
						w.sysErr <- ErrFatalQuery
						return
					}
				} else {
					break
				}
			}

			// query for logs
			query := buildQuery(w.cfg.BridgeContract, utils.ProposalEvent, latestBlock, latestBlock)
			evts, err := w.client.FilterLogs(context.Background(), query)
			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch logs")
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
					log.Trace().Interface("src", sourceId).Interface("nonce", depositNonce).Msg("Ignoring event")
				}
			}
			log.Trace().Interface("block", latestBlock).Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("No finalization event found in current block")
			latestBlock = latestBlock.Add(latestBlock, big.NewInt(1))
		}
	}
	log.Warn().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Block watch limit exceeded, skipping execution")
}

// voteProposal submits a vote proposal
// a vote proposal will try to be submitted up to the TxRetryLimit times
func (w *writer) voteProposal(m msg.Message, dataHash [32]byte) {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			err := w.client.LockAndUpdateOpts()
			if err != nil {
				log.Error().Err(err).Msg("Failed to update tx opts")
				continue
			}

			tx, err := w.bridgeContract.VoteProposal(
				w.client.Opts(),
				uint8(m.Source),
				uint64(m.DepositNonce),
				m.ResourceId,
				dataHash,
			)
			w.client.UnlockOpts()

			if err == nil {
				log.Info().Str("tx", tx.Hash().Hex()).Interface("src", m.Source).Interface("depositNonce", m.DepositNonce).Msg("Submitted proposal vote")
				if w.metrics != nil {
					w.metrics.VotesSubmitted.Inc()
				}
				return
			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				log.Debug().Msg("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				log.Warn().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Voting failed")
				time.Sleep(TxRetryInterval)
			}

			// Verify proposal is still open for voting, otherwise no need to retry
			if w.proposalIsComplete(m.Source, m.DepositNonce, dataHash) {
				log.Info().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Proposal voting complete on chain")
				return
			}
		}
	}
	log.Error().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Submission of Vote transaction failed")
	w.sysErr <- ErrFatalTx
}

// executeProposal executes the proposal
func (w *writer) executeProposal(m msg.Message, data []byte, dataHash [32]byte, msgProofOpts *celoMsg.MsgProofOpts) {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			err := w.client.LockAndUpdateOpts()
			if err != nil {
				log.Error().Err(err).Msg("Failed to update nonce")
				return
			}

			tx, err := w.bridgeContract.ExecuteProposal(
				w.client.Opts(),
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
			w.client.UnlockOpts()

			if err == nil {
				log.Info().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal execution")
				return
			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				log.Error().Err(err).Msg("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				log.Error().Err(err).Msg("Execution failed, proposal may already be complete")
				time.Sleep(TxRetryInterval)
			}

			// Verify proposal is still open for execution, tx will fail if we aren't the first to execute,
			// but there is no need to retry
			if w.proposalIsFinalized(m.Source, m.DepositNonce, dataHash) {
				log.Info().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Proposal finalized on chain")
				return
			}
		}
	}
	log.Error().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Submission of Execute transaction failed")
	w.sysErr <- ErrFatalTx
}

// buildQuery constructs a query for the bridgeContract by hashing sig to get the event topic
func buildQuery(contract ethcommon.Address, sig utils.EventSig, startBlock *big.Int, endBlock *big.Int) eth.FilterQuery {
	query := eth.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Addresses: []ethcommon.Address{contract},
		Topics: [][]ethcommon.Hash{
			{sig.GetTopic()},
		},
	}
	return query
}
