// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package listener

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	"github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/chain/config"
	"github.com/ChainSafe/chainbridge-celo/msg"
	"github.com/ChainSafe/chainbridge-celo/shared/ethereum"
	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
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
	bridgeContract         *Bridge.Bridge // instance of bound bridge contract
	erc20HandlerContract   *ERC20Handler.ERC20Handler
	erc721HandlerContract  *ERC721Handler.ERC721Handler
	genericHandlerContract *GenericHandler.GenericHandler
	blockstore             Blockstorer
	stop                   <-chan struct{}
	sysErr                 chan<- error // Reports fatal error to core
	syncer                 BlockSyncer
	//latestBlock            *metrics.LatestBlock
	//metrics                *metrics.ChainMetrics
	client client.LogFilterWithLatestBlock
}

type BlockSyncer interface {
	Sync(latestBlock *big.Int) error
}

type IRouter interface {
	Send(msg *msg.Message) error
}
type Blockstorer interface {
	StoreBlock(*big.Int) error
}

func NewListener(cfg *config.CeloChainConfig, client client.LogFilterWithLatestBlock, bs Blockstorer, stop <-chan struct{}, sysErr chan<- error, syncer BlockSyncer, router IRouter) *listener {
	return &listener{
		cfg:        cfg,
		blockstore: bs,
		stop:       stop,
		sysErr:     sysErr,
		syncer:     syncer,
		router:     router,
		client:     client,
	}
}

func (l *listener) SetContracts(bridge *Bridge.Bridge, erc20Handler *ERC20Handler.ERC20Handler, erc721Handler *ERC721Handler.ERC721Handler, genericHandler *GenericHandler.GenericHandler) {
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
				log.Debug().Str("target", currentBlock.String()).Str("latest", latestBlock.String()).Msg("Block not ready, will retry")
				time.Sleep(BlockRetryInterval)
				continue
			}

			err = l.syncer.Sync(currentBlock)
			if err != nil {
				log.Error().Str("block", currentBlock.String()).Err(err).Msg("Failed to sync validators for block")
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

// TODO: Proof construction.
func (l *listener) getDepositEventsAndProofsForBlock(latestBlock *big.Int) error {
	log.Debug().Str("block", latestBlock.String()).Msg("Querying block for deposit events")
	query := buildQuery(l.cfg.BridgeContract, utils.Deposit, latestBlock, latestBlock)

	// querying for logs
	logs, err := l.client.FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("unable to Filter Logs: %w", err)
	}

	// read through the log events and handle their deposit event if handler is recognized
	for _, eventLog := range logs {
		var m *msg.Message
		destId := msg.ChainId(eventLog.Topics[1].Big().Uint64())
		rId := msg.ResourceId(eventLog.Topics[2])
		nonce := msg.Nonce(eventLog.Topics[3].Big().Uint64())

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

		err = l.router.Send(m)
		if err != nil {
			log.Error().Err(err).Msg("subscription error: failed to route message")
		}
	}
	return nil
}

//TODO removenolint
//nolint
//COMMENTED SINCE CURRENTLTY UNUSED. SEEMS TO BE USED FOR BLOCK PROOF BUILDING
//func (l *listener) getBlockHashFromTransactionHash(txHash ethcommon.Hash) (blockHash ethcommon.Hash, err error) {
//
//	receipt, err := l.conn.Client().TransactionReceipt(context.Background(), txHash)
//	if err != nil {
//		return txHash, fmt.Errorf("unable to get BlockHash: %w", err)
//	}
//	return receipt.BlockHash, nil
//}
//
////TODO removenolint
////nolint
//func (l *listener) getTransactionsFromBlockHash(blockHash ethcommon.Hash) (txHashes []ethcommon.Hash, txRoot ethcommon.Hash, err error) {
//	block, err := l.conn.Client().BlockByHash(context.Background(), blockHash)
//	if err != nil {
//		return nil, ethcommon.Hash{}, fmt.Errorf("unable to get BlockHash: %w", err)
//	}
//
//	var transactionHashes []ethcommon.Hash
//
//	transactions := block.Transactions()
//	for _, transaction := range transactions {
//		transactionHashes = append(transactionHashes, transaction.Hash())
//	}
//
//	return transactionHashes, block.Root(), nil
//}
//
//nolint
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
