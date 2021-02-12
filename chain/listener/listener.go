// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package listener

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridge-celo/chain/config"
	"github.com/ChainSafe/chainbridge-celo/chain/sender"
	"github.com/ChainSafe/chainbridge-celo/pkg"
	"github.com/ChainSafe/chainbridge-celo/txtrie"
	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog/log"
)

var BlockDelay = big.NewInt(10)
var BlockRetryInterval = time.Second * 5
var ErrFatalPolling = errors.New("listener block polling failed")
var ExpectedBlockTime = time.Second
var BlockRetryLimit = 5

type listener struct {
	cfg                    *config.CeloChainConfig
	router                 IRouter
	bridgeContract         IBridge // instance of bound bridge contract
	erc20HandlerContract   IERC20Handler
	erc721HandlerContract  IERC721Handler
	genericHandlerContract IGenericHandler
	blockstore             Blockstorer
	stop                   <-chan struct{}
	sysErr                 chan<- error // Reports fatal error to core
	//latestBlock            *metrics.LatestBlock
	//metrics                *metrics.ChainMetrics
	client   sender.LogFilterWithLatestBlock
	valsAggr ValidatorsAggregator
}

type IRouter interface {
	Send(msg *pkg.Message) error
}
type Blockstorer interface {
	StoreBlock(*big.Int) error
}

type ValidatorsAggregator interface {
	GetAPKForBlock(block *big.Int, chainID uint8, epochSize uint64) ([]byte, error)
}

func NewListener(cfg *config.CeloChainConfig, client sender.LogFilterWithLatestBlock, bs Blockstorer, stop <-chan struct{}, sysErr chan<- error, router IRouter, valsAggr ValidatorsAggregator) *listener {
	return &listener{
		cfg:        cfg,
		blockstore: bs,
		stop:       stop,
		sysErr:     sysErr,
		router:     router,
		client:     client,
		valsAggr:   valsAggr,
	}
}

func (l *listener) SetContracts(bridge IBridge, erc20Handler IERC20Handler, erc721Handler IERC721Handler, genericHandler IGenericHandler) {
	l.bridgeContract = bridge
	l.erc20HandlerContract = erc20Handler
	l.erc721HandlerContract = erc721Handler
	l.genericHandlerContract = genericHandler
}

func (l *listener) StartPollingBlocks() error {
	log.Debug().Msg("Starting listener...")

	go func() {
		err := l.pollBlocks()
		if err != nil {
			log.Error().Err(err).Msg("Polling blocks failed")
		}
	}()

	return nil
}

// TODO this is metrics latest block, naming mess
//func (l *listener) LatestBlock() *metrics.LatestBlock {
//	return l.latestBlock
//}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.cfg.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before continuing to the next block.
func (l *listener) pollBlocks() error {
	log.Info().Msg("Polling Blocks...")
	var currentBlock = l.cfg.StartBlock
	var retry = BlockRetryLimit
	for {
		select {
		case <-l.stop:
			return errors.New("polling terminated")
		default:
			// No more retries, goto next block
			if retry == 0 {
				log.Error().Msg("Polling failed, retries exceeded")
				l.sysErr <- ErrFatalPolling
				return nil
			}

			latestBlock, err := l.client.LatestBlock()
			if err != nil {
				log.Error().Err(err).Str("block", currentBlock.String()).Msg("Unable to get latest block")
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Sleep if the difference is less than BlockDelay; (latest - current) < BlockDelay
			if big.NewInt(0).Sub(latestBlock, currentBlock).Cmp(BlockDelay) == -1 {
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Parse out events
			err = l.getDepositEventsAndProofsForBlock(currentBlock)
			if err != nil {
				log.Error().Str("block", currentBlock.String()).Err(err).Msg("Failed to get events for block")
				retry--
				continue
			}

			// Write to block store. Not a critical operation, no need to retry
			err = l.blockstore.StoreBlock(currentBlock)
			if err != nil {
				log.Error().Str("block", currentBlock.String()).Err(err).Msg("Failed to write latest block to blockstore")
			}

			//if l.metrics != nil {
			//	l.metrics.BlocksProcessed.Inc()
			//	l.metrics.LatestProcessedBlock.Set(float64(latestBlock.Int64()))
			//}
			//
			//l.latestBlock.Height = big.NewInt(0).Set(latestBlock)
			//l.latestBlock.LastUpdated = time.Now()

			// Goto next block and reset retry counter
			currentBlock.Add(currentBlock, big.NewInt(1))
			retry = BlockRetryLimit
		}
	}
}

func (l *listener) getDepositEventsAndProofsForBlock(latestBlock *big.Int) error {
	query := buildQuery(l.cfg.BridgeContract, pkg.Deposit, latestBlock, latestBlock)

	// querying for logs
	logs, err := l.client.FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("unable to Filter Logs: %w", err)
	}
	if len(logs) == 0 {
		return nil
	}

	blockData, err := l.client.BlockByNumber(context.Background(), latestBlock)
	if err != nil {
		return err
	}

	trie, err := txtrie.CreateNewTrie(blockData.TxHash(), blockData.Transactions())
	if err != nil {
		return err
	}
	// read through the log events and handle their deposit event if handler is recognized
	for _, eventLog := range logs {

		var m *pkg.Message
		destId := pkg.ChainId(eventLog.Topics[1].Big().Uint64())
		rId := pkg.ResourceId(eventLog.Topics[2])
		nonce := pkg.Nonce(eventLog.Topics[3].Big().Uint64())

		addr, err := l.bridgeContract.ResourceIDToHandlerAddress(&bind.CallOpts{}, rId)
		if err != nil {
			return fmt.Errorf("failed to get handler from resource ID %x, reason: %w", rId, err)
		}

		if addr == l.cfg.Erc20HandlerContract {
			m, err = l.handleErc20DepositedEvent(destId, nonce)
		} else if addr == l.cfg.Erc721HandlerContract {
			m, err = l.handleErc721DepositedEvent(destId, nonce)
		} else if addr == l.cfg.GenericHandlerContract {
			m, err = l.handleGenericDepositedEvent(destId, nonce)
		} else {
			log.Error().Err(err).Str("handler", addr.Hex()).Msg("event has unrecognized handler")
			return nil
		}
		if err != nil {
			return err
		}
		apk, err := l.valsAggr.GetAPKForBlock(latestBlock, uint8(l.cfg.ID), l.cfg.EpochSize)
		if err != nil {
			return err

		}
		keyRlp, err := rlp.EncodeToBytes(eventLog.TxIndex)
		if err != nil {
			return fmt.Errorf("encoding TxIndex to rlp: %w", err)
		}
		proof, key, err := txtrie.RetrieveProof(trie, keyRlp)
		if err != nil {
			return err
		}

		m.SVParams = &pkg.SignatureVerification{AggregatePublicKey: apk, BlockHash: blockData.Header().Hash(), Signature: blockData.EpochSnarkData().Signature}
		m.MPParams = &pkg.MerkleProof{TxRootHash: pkg.SliceTo32Bytes(blockData.TxHash().Bytes()), Nodes: proof, Key: key}
		err = l.router.Send(m)

		if err != nil {
			log.Error().Err(err).Msg("subscription error: failed to route message")
		}
	}
	return nil
}

// buildQuery constructs a query for the bridgeContract by hashing sig to get the event topic
func buildQuery(contract ethcommon.Address, sig pkg.EventSig, startBlock *big.Int, endBlock *big.Int) eth.FilterQuery {
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
