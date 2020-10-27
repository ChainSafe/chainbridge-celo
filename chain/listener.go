// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"context"
	"fmt"
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

var _ Connection = &connection.Connection{}

type Connection interface {
	Connect() error
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
	return nil
}

func (l *listener) close() {
	l.conn.Close()
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
