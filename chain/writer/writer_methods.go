// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package writer

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridge-celo/utils"
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
func (w *writer) proposalIsComplete(srcId utils.ChainId, nonce utils.Nonce, dataHash ethcommon.Hash) bool {
	prop, err := w.bridgeContract.GetProposal(w.client.CallOpts(), uint8(srcId), uint64(nonce), dataHash)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check proposal existence")
		return false
	}
	return prop.Status == ProposalStatusPassed || prop.Status == ProposalStatusTransferred || prop.Status == ProposalStatusCancelled
}

// proposalIsFinalized returns true if the proposal state is Transferred or Cancelled
func (w *writer) proposalIsFinalized(srcId utils.ChainId, nonce utils.Nonce, dataHash ethcommon.Hash) bool {
	prop, err := w.bridgeContract.GetProposal(w.client.CallOpts(), uint8(srcId), uint64(nonce), dataHash)

	if err != nil {
		log.Error().Err(err).Msg("Failed to check proposal existence")
		return false
	}
	return prop.Status == ProposalStatusTransferred || prop.Status == ProposalStatusCancelled
}

// hasVoted checks if this relayer has already voted
func (w *writer) hasVoted(srcId utils.ChainId, nonce utils.Nonce, dataHash ethcommon.Hash) bool {
	hasVoted, err := w.bridgeContract.HasVotedOnProposal(w.client.CallOpts(), idAndNonce(srcId, nonce), dataHash, w.client.Opts().From)

	if err != nil {
		log.Error().Err(err).Msg("Failed to check proposal existence")
		return false
	}
	return hasVoted
}

func (w *writer) proposalIsPassed(srcId utils.ChainId, nonce utils.Nonce, dataHash [32]byte) bool {
	prop, err := w.bridgeContract.GetProposal(w.client.CallOpts(), uint8(srcId), uint64(nonce), dataHash)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to check proposal existence")
		return false
	}
	return prop.Status == ProposalStatusPassed
}

func idAndNonce(srcId utils.ChainId, nonce utils.Nonce) *big.Int {
	var data []byte
	data = append(data, nonce.Big().Bytes()...)
	data = append(data, uint8(srcId))
	return big.NewInt(0).SetBytes(data)
}

func (w *writer) shouldVote(m *utils.Message, dataHash ethcommon.Hash) bool {
	// Check if proposal has passed and skip if Passed or Transferred
	if w.proposalIsComplete(m.Source, m.DepositNonce, dataHash) {
		log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Proposal complete, not voting")
		return false
	}

	// Check if relayer has previously voted
	if w.hasVoted(m.Source, m.DepositNonce, dataHash) {
		log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Relayer has already voted, not voting")
		return false
	}
	return true
}

func (w *writer) createERC20ProposalData(m *utils.Message) ([]byte, error) {
	log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Creating erc20 proposal")
	if len(m.Payload) != 2 {
		return nil, errors.New("malformed payload. Len  of payload should be 2")
	}
	amount, ok := m.Payload[0].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads amount format")
	}

	recipient, ok := m.Payload[1].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads recipient format")
	}
	data := ConstructErc20ProposalData(amount, recipient)
	return data, nil
}

func (w *writer) createErc721ProposalData(m *utils.Message) ([]byte, error) {
	log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Creating erc721 proposal")
	if len(m.Payload) != 3 {
		return nil, errors.New("malformed payload. Len  of payload should be 3")
	}
	tokenID, ok := m.Payload[0].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads tokenID format")
	}
	recipient, ok := m.Payload[1].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads recipient format")
	}
	metadata, ok := m.Payload[2].([]byte)
	if !ok {
		return nil, errors.New("wrong payloads metadata format")
	}
	return ConstructErc721ProposalData(tokenID, recipient, metadata), nil
}

func (w *writer) createGenericDepositProposalData(m *utils.Message) ([]byte, error) {
	log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Creating generic proposal")
	if len(m.Payload) != 1 {
		return nil, errors.New("malformed payload. Len  of payload should be 1")
	}
	metadata, ok := m.Payload[0].([]byte)
	if !ok {
		return nil, errors.New("unable to convert metadata to []byte")
	}
	return ConstructGenericProposalData(metadata), nil
}

// watchThenExecute watches for the latest block and executes once the matching finalized event is found
func (w *writer) watchThenExecute(m *utils.Message, data []byte, dataHash ethcommon.Hash, latestBlock *big.Int) {
	log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Watching for finalization event")

	// watching for the latest block, querying and matching the finalized event will be retried up to ExecuteBlockWatchLimit times
	for i := 0; i < ExecuteBlockWatchLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			// watch for the lastest block, retry up to BlockRetryLimit times
			for waitRetrys := 0; waitRetrys <= BlockRetryLimit; waitRetrys++ {
				err := w.client.WaitForBlock(latestBlock)
				if err != nil {
					log.Error().Err(err).Msg("Waiting for block failed")
					// Exit if retries exceeded
					if waitRetrys == BlockRetryLimit {
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

				if m.Source == utils.ChainId(sourceId) &&
					m.DepositNonce.Big().Uint64() == depositNonce &&
					utils.IsPassed(uint8(status)) {
					w.executeProposal(m, data, dataHash)
					return
				} else {
					log.Trace().Interface("src", sourceId).Interface("nonce", depositNonce).Uint64("status", status).Msg("Ignoring event")
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
func (w *writer) voteProposal(m *utils.Message, dataHash ethcommon.Hash) {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			// Checking first does proposal complete? If so, we do not need to vote for it
			if w.proposalIsComplete(m.Source, m.DepositNonce, dataHash) {
				log.Info().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Proposal voting complete on chain")
				return
			}
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
			if err != nil {
				if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
					log.Debug().Msg("Nonce too low, will retry")
					time.Sleep(TxRetryInterval)
					continue
				} else {
					log.Warn().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Voting failed")
					time.Sleep(TxRetryInterval)
					continue
				}
			}
			log.Info().Str("tx", tx.Hash().Hex()).Interface("src", m.Source).Interface("depositNonce", m.DepositNonce).Msg("Submitted proposal vote")
			return
		}
	}
	log.Error().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Submission of Vote transaction failed")
	w.sysErr <- ErrFatalTx
}

// executeProposal executes the proposal
func (w *writer) executeProposal(m *utils.Message, data []byte, dataHash ethcommon.Hash) {
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
				m.SVParams.Signature,
				m.SVParams.AggregatePublicKey,
				// TODO: Remove once G1 has been removed from contracts
				[]byte{},
				m.SVParams.BlockHash,
				m.MPParams.TxRootHash,
				m.MPParams.Key,
				m.MPParams.Nodes,
			)
			w.client.UnlockOpts()

			if err == nil {
				log.Info().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal execution")
				return
			}
			if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				log.Error().Err(err).Msg("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				log.Error().Err(err).Msg("Execution failed, proposal may already be complete")
				time.Sleep(TxRetryInterval)
			}
			// Checking proposal status one more time (Since it could be execute by some other bridge). If it is finalized then we do not need to retry
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
