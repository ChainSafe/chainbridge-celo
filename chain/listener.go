// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridge-celo/connection"
	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	"github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
	"github.com/ChainSafe/chainbridge-utils/blockstore"
	"github.com/ChainSafe/chainbridge-utils/core"
	log "github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var BlockDelay = big.NewInt(10)
var BlockRetryInterval = time.Second * 5
var BlockRetryLimit = 5
var ErrFatalPolling = errors.New("listener block polling failed")

var _ Connection = &connection.Connection{}

type Connection interface {
	Connect() error
	LatestBlock() (*big.Int, error)
	Close()
	Client() *ethclient.Client
	Opts() *bind.TransactOpts
}

var ExpectedBlockTime = time.Second

type listener struct {
	cfg                    Config
	conn                   Connection
	router                 *core.Router
	bridgeContract         *Bridge.Bridge // instance of bound bridge contract
	erc20HandlerContract   *ERC20Handler.ERC20Handler
	erc721HandlerContract  *ERC721Handler.ERC721Handler
	genericHandlerContract *GenericHandler.GenericHandler
	log                    log.Logger
	blockstore             blockstore.Blockstorer
	stop                   <-chan int
	sysErr                 chan<- error // Reports fatal error to core
	syncer                 ValidatorSyncer
}

func NewListener(conn Connection, cfg *Config, log log.Logger, bs blockstore.Blockstorer, stop <-chan int, sysErr chan<- error, s ValidatorSyncer) *listener {
	return &listener{
		cfg:         *cfg,
		conn:        conn,
		log:         log,
		blockstore:  bs,
		stop:        stop,
		sysErr:      sysErr,
		syncer:      s,
	}
}

func (l *listener) setContracts(bridge *Bridge.Bridge, erc20Handler *ERC20Handler.ERC20Handler, erc721Handler *ERC721Handler.ERC721Handler, genericHandler *GenericHandler.GenericHandler) {
	l.bridgeContract = bridge
	l.erc20HandlerContract = erc20Handler
	l.erc721HandlerContract = erc721Handler
	l.genericHandlerContract = genericHandler
}

func (l *listener) setRouter(r *core.Router) {
	l.router = r
}

func (l *listener) start() error {
	l.log.Debug("Starting listener...")

	err := l.conn.Connect()
	if err != nil {
		return err
	}

	go func() {
		err := l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

func (l *listener) close() {
	l.conn.Close()
}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.cfg.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before continuing to the next block.
func (l *listener) pollBlocks() error {
	l.log.Info("Polling Blocks...")
	var currentBlock = l.cfg.startBlock
	var retry = BlockRetryLimit
	for {
		select {
		case <-l.stop:
			return errors.New("polling terminated")
		default:
			// No more retries, goto next block
			if retry == 0 {
				l.log.Error("Polling failed, retries exceeded")
				l.sysErr <- ErrFatalPolling
				return nil
			}

			latestBlock, err := l.conn.LatestBlock()
			if err != nil {
				l.log.Error("Unable to get latest block", "block", currentBlock, "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Sleep if the difference is less than BlockDelay; (latest - current) < BlockDelay
			if big.NewInt(0).Sub(latestBlock, currentBlock).Cmp(BlockDelay) == -1 {
				l.log.Debug("Block not ready, will retry", "target", currentBlock, "latest", latestBlock)
				time.Sleep(BlockRetryInterval)
				continue
			}

			err = l.syncer.Sync(currentBlock)
			if err != nil {
				l.log.Error("Failed to sync validators for block", "block", currentBlock, "err", err)
				continue
			}

			// Parse out events
			err = l.getDepositEventsAndProofsForBlock(currentBlock)
			if err != nil {
				l.log.Error("Failed to get events for block", "block", currentBlock, "err", err)
				retry--
				continue
			}

			// Write to block store. Not a critical operation, no need to retry
			err = l.blockstore.StoreBlock(currentBlock)
			if err != nil {
				l.log.Error("Failed to write latest block to blockstore", "block", currentBlock, "err", err)
			}

			// Goto next block and reset retry counter
			currentBlock.Add(currentBlock, big.NewInt(1))
			retry = BlockRetryLimit
		}
	}
}

func (l *listener) getDepositEventsAndProofsForBlock(latestBlock *big.Int) error {
	l.log.Debug("Querying block for deposit events", "block", latestBlock)

	return nil
}

func (l *listener) getBlockHashFromTransactionHash(txHash common.Hash) (blockHash common.Hash, err error) {

	receipt, err := l.conn.Client().TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return txHash, fmt.Errorf("unable to get BlockHash: %s", err)
	}
	return receipt.BlockHash, nil
}

func (l *listener) getTransactionsFromBlockHash(blockHash common.Hash) (txHashes []common.Hash, txRoot common.Hash, err error) {
	block, err := l.conn.Client().BlockByHash(context.Background(), blockHash)
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("unable to get BlockHash: %s", err)
	}

	var transactionHashes []common.Hash

	transactions := block.Transactions()
	for _, transaction := range transactions {
		transactionHashes = append(transactionHashes, transaction.Hash())
	}

	return transactionHashes, block.Root(), nil
}
