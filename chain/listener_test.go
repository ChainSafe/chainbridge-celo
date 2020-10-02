// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	"github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
	"github.com/ChainSafe/chainbridge-celo/connection"
	"github.com/ChainSafe/chainbridge-celo/shared/ethereum"

	log "github.com/ChainSafe/log15"
	"github.com/ChainSafe/chainbridge-utils/blockstore"
	"github.com/ChainSafe/chainbridge-utils/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var GasLimitUint64 = uint64(connection.DefaultGasLimit)
var ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
var TestChainID = uint8(0)
var TestRelayerThreshold = big.NewInt(2)

func createTestListener(t *testing.T) *listener {
	newConfig := Config{}
	conn := connection.NewConnection(TestEndpoint, false, AliceKp, log.Root(), GasLimit, GasPrice)
	vsyncer := ValidatorSyncer{conn: conn}
	stop := make(chan int)
	errs := make(chan error)

	l := NewListener(conn, &newConfig, log.Root(), &blockstore.EmptyStore{}, stop, errs, vsyncer)

	client, err := utils.NewClient(TestEndpoint, AliceKp)
	if err != nil {
		t.Fatal(err)
	}
	contracts, err := utils.DeployContracts(
		client,
		TestChainID,
		TestRelayerThreshold,
	)
	if err != nil {
		t.Fatal(err)
	}

	bridgeContract, err := Bridge.NewBridge(contracts.BridgeAddress, conn.Client())
	if err != nil {
		t.Fatal(err)
	}
	erc20HandlerContract, err := ERC20Handler.NewERC20Handler(contracts.ERC20HandlerAddress, conn.Client())
	if err != nil {
		t.Fatal(err)
	}
	erc721HandlerContract, err := ERC721Handler.NewERC721Handler(contracts.ERC721HandlerAddress, conn.Client())
	if err != nil {
		t.Fatal(err)
	}
	genericHandlerContract, err := GenericHandler.NewGenericHandler(contracts.GenericHandlerAddress, conn.Client())
	if err != nil {
		t.Fatal(err)
	}
	l.setContracts(bridgeContract, erc20HandlerContract, erc721HandlerContract, genericHandlerContract)

	router := core.NewRouter(log.Root())
	l.setRouter(router)

	return l
}

// creating and sending a new transaction
func newTransaction(t *testing.T, l *listener) common.Hash {

	// Creating a new transaction
	nonce := l.conn.Opts().Nonce
	tx := types.NewTransaction(nonce.Uint64(), ZeroAddress, big.NewInt(0), GasLimitUint64, GasPrice, nil, nil, nil, nil)

	chainId, err := l.conn.Client().ChainID(context.Background())
	signer := types.NewEIP155Signer(chainId)
	if err != nil {
		t.Fatal(err)
	}

	// Sign the transaction
	signedTx, err := types.SignTx(tx, signer, AliceKp.PrivateKey()) //SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error)
	if err != nil {
		t.Fatal(err)
	}

	// Send the transaction for execution
	err = l.conn.Client().SendTransaction(context.Background(), signedTx)
	if err != nil {
		t.Fatal(err)
	}
	return signedTx.Hash()
}

// WaitForTx will query the chain at ExpectedBlockTime intervals, until a receipt is returned.
// Returns an error if the tx failed.
func waitForTx(client *ethclient.Client, hash common.Hash) (*types.Receipt, error) {
	retry := 10
	for retry > 0 {
		receipt, err := client.TransactionReceipt(context.Background(), hash)
		if err != nil {
			retry--
			time.Sleep(ExpectedBlockTime)
			continue
		}

		if receipt.Status != 1 {
			return nil, fmt.Errorf("transaction failed on chain")
		}
		return receipt, nil
	}
	return nil, fmt.Errorf("transaction after retries failed")
}

func TestListener_start_stop(t *testing.T) {
	l := createTestListener(t)

	err := l.start()
	if err != nil {
		t.Fatal(err)
	}

	// Initiate shutdown
	l.close()
}

// Testing transaction Block hash
func TestListener_BlockHashFromTransactionHash(t *testing.T) {

	l := createTestListener(t)
	err := l.start()
	if err != nil {
		t.Fatal(err)
	}

	// Create and submit a new transaction and return the signed transaction hash
	txHash := newTransaction(t, l)

	receipt, err := waitForTx(l.conn.Client(), txHash)
	if err != nil {
		t.Fatal(err)
	}

	blockHash, err := l.getBlockHashFromTransactionHash(txHash)
	if err != nil {
		t.Fatal(err)
	}

	// Confirm that the receipt blockhash and the block's blockhash are the same
	if blockHash != receipt.BlockHash {
		t.Fatalf("block hash are not equal, expected: %x, %x", receipt.BlockHash, blockHash)
	}
}

func TestListener_TransactionsFromBlockHash(t *testing.T) {
	l := createTestListener(t)
	err := l.start()
	if err != nil {
		t.Fatal(err)
	}

	// Create and submit a new transaction
	txHash := newTransaction(t, l)

	// Get receipt from hash
	receipt, err := waitForTx(l.conn.Client(), txHash)
	if err != nil {
		t.Fatal(err)
	}

	// Get txHashes and txroot from blockHash
	txs, _, err := l.getTransactionsFromBlockHash(receipt.BlockHash)
	if err != nil {
		t.Fatal(err)
	}

	ok := false
	for _, tx := range txs {
		if tx == txHash {
			ok = true
		}
	}

	if !ok {
		t.Fatalf("expected %x to be in %x", txHash, txs)
	}
}
