// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	"github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
	"github.com/ChainSafe/chainbridge-celo/connection"
	utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
	ethtest "github.com/ChainSafe/chainbridge-celo/shared/ethereum/testing"

	"github.com/ChainSafe/chainbridge-utils/blockstore"
	"github.com/ChainSafe/chainbridge-utils/core"
	"github.com/ChainSafe/chainbridge-utils/msg"
	log "github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type mockWriter struct {
	msgs chan msg.Message
}

func (w *mockWriter) ResolveMessage(msg msg.Message) bool {
	w.msgs <- msg
	return true
}

func createTestListener(t *testing.T, client *utils.Client, stop <-chan int, errs chan<- error) *listener {
	contracts, err := utils.DeployContracts(
		client,
		TestChainID,
		TestRelayerThreshold,
	)
	if err != nil {
		t.Fatal(err)
	}
	conn := connection.NewConnection(TestEndpoint, false, AliceKp, log.Root(), GasLimit)
	err = conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	vsyncer := ValidatorSyncer{conn: conn}

	latestBlock, err := conn.LatestBlock()
	if err != nil {
		t.Fatal(err)
	}

	newConfig := Config{}
	newConfig.startBlock = latestBlock
	newConfig.bridgeContract = contracts.BridgeAddress
	newConfig.erc20HandlerContract = contracts.ERC20HandlerAddress
	newConfig.erc721HandlerContract = contracts.ERC721HandlerAddress
	newConfig.genericHandlerContract = contracts.GenericHandlerAddress

	l := NewListener(conn, &newConfig, log.Root(), &blockstore.EmptyStore{}, stop, errs, vsyncer)

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
	err := l.conn.LockAndUpdateNonce()
	if err != nil {
		t.Fatal(err)
	}
	defer l.conn.UnlockNonce()
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
	client, err := utils.NewClient(TestEndpoint, AliceKp)
	if err != nil {
		t.Fatal(err)
	}
	stop := make(chan int)
	l := createTestListener(t, client, stop, make(chan error))

	err = l.start()
	if err != nil {
		t.Fatal(err)
	}

	// Initiate shutdown
	close(stop)
	l.conn.Close()
}

// Testing transaction Block hash
func TestListener_BlockHashFromTransactionHash(t *testing.T) {
	client, err := utils.NewClient(TestEndpoint, AliceKp)
	if err != nil {
		t.Fatal(err)
	}
	l := createTestListener(t, client, make(chan int), make(chan error))
	err = l.start()
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
	client, err := utils.NewClient(TestEndpoint, AliceKp)
	if err != nil {
		t.Fatal(err)
	}
	l := createTestListener(t, client, make(chan int), make(chan error))
	err = l.start()
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

func TestListener_Erc20DepositedEvent(t *testing.T) {
	client, err := utils.NewClient(TestEndpoint, AliceKp)
	if err != nil {
		t.Fatal(err)
	}
	errs := make(chan error)
	l := createTestListener(t, client, make(chan int), errs)
	err = l.start()
	if err != nil {
		t.Fatal(err)
	}
	router := l.router
	writer := &mockWriter{msgs: make(chan msg.Message)}
	router.Listen(msg.ChainId(1), writer)

	// For debugging
	go ethtest.WatchEvent(client, l.cfg.bridgeContract, utils.Deposit)

	erc20Contract := ethtest.DeployMintApproveErc20(t, client, l.cfg.erc20HandlerContract, big.NewInt(100))

	amount := big.NewInt(10)
	src := msg.ChainId(0)
	dst := msg.ChainId(1)
	resourceId := msg.ResourceIdFromSlice(append(common.LeftPadBytes(erc20Contract.Bytes(), 31), uint8(src)))
	recipient := ethcrypto.PubkeyToAddress(BobKp.PrivateKey().PublicKey)

	ethtest.RegisterResource(t, client, l.cfg.bridgeContract, l.cfg.erc20HandlerContract, resourceId, erc20Contract)

	expectedMessage := msg.NewFungibleTransfer(
		src,
		dst,
		1,
		amount,
		resourceId,
		common.HexToAddress(BobKp.Address()).Bytes(),
	)
	// Create an ERC20 Deposit
	ethtest.CreateErc20Deposit(
		t,
		l.bridgeContract,
		client,
		resourceId,

		recipient,
		dst,
		amount,
	)

	verifyMessage(t, writer, expectedMessage, errs)

	// Create second deposit, verify nonce change
	expectedMessage = msg.NewFungibleTransfer(
		src,
		dst,
		2,
		amount,
		resourceId,
		common.HexToAddress(BobKp.Address()).Bytes(),
	)
	ethtest.CreateErc20Deposit(
		t,
		l.bridgeContract,
		client,
		resourceId,

		recipient,
		dst,
		amount,
	)

	verifyMessage(t, writer, expectedMessage, errs)
}

func TestListener_Erc721DepositedEvent(t *testing.T) {
	client, err := utils.NewClient(TestEndpoint, AliceKp)
	if err != nil {
		t.Fatal(err)
	}
	errs := make(chan error)
	l := createTestListener(t, client, make(chan int), errs)
	err = l.start()
	if err != nil {
		t.Fatal(err)
	}
	router := l.router
	writer := &mockWriter{msgs: make(chan msg.Message)}
	router.Listen(msg.ChainId(1), writer)

	// For debugging
	go ethtest.WatchEvent(client, l.cfg.bridgeContract, utils.Deposit)

	tokenId := big.NewInt(99)

	erc721Contract := ethtest.Erc721Deploy(t, client)
	ethtest.Erc721Mint(t, client, erc721Contract, tokenId, []byte{})
	ethtest.Erc721Approve(t, client, erc721Contract, l.cfg.erc721HandlerContract, tokenId)
	log.Info("Deployed erc721, minted and approved handler", "handler", l.cfg.erc721HandlerContract, "contract", erc721Contract, "tokenId", tokenId.Bytes())
	ethtest.Erc721AssertOwner(t, client, erc721Contract, tokenId, client.Opts.From)
	src := msg.ChainId(0)
	dst := msg.ChainId(1)
	resourceId := msg.ResourceIdFromSlice(append(common.LeftPadBytes(erc721Contract.Bytes(), 31), uint8(src)))
	recipient := BobKp.CommonAddress()

	ethtest.RegisterResource(t, client, l.cfg.bridgeContract, l.cfg.erc721HandlerContract, resourceId, erc721Contract)

	expectedMessage := msg.NewNonFungibleTransfer(
		src,
		dst,
		1,
		resourceId,
		tokenId,
		recipient.Bytes(),
		[]byte{},
	)

	// Create an ERC20 Deposit
	ethtest.CreateErc721Deposit(
		t,
		l.bridgeContract,
		client,
		resourceId,

		recipient,
		dst,
		tokenId,
	)

	verifyMessage(t, writer, expectedMessage, errs)
}

func TestListener_GenericDepositedEvent(t *testing.T) {
	client, err := utils.NewClient(TestEndpoint, AliceKp)
	if err != nil {
		t.Fatal(err)
	}
	errs := make(chan error)
	l := createTestListener(t, client, make(chan int), errs)
	err = l.start()
	if err != nil {
		t.Fatal(err)
	}
	router := l.router
	writer := &mockWriter{msgs: make(chan msg.Message)}
	router.Listen(msg.ChainId(1), writer)

	// For debugging
	go ethtest.WatchEvent(client, l.cfg.bridgeContract, utils.Deposit)

	src := msg.ChainId(0)
	dst := msg.ChainId(1)
	hash := utils.Hash(common.LeftPadBytes([]byte{1}, 32))
	resourceId := msg.ResourceIdFromSlice(append(common.LeftPadBytes([]byte{1}, 31), uint8(src)))
	depositSig := utils.CreateFunctionSignature("")
	executeSig := utils.CreateFunctionSignature("store()")
	ethtest.RegisterGenericResource(t, client, l.cfg.bridgeContract, l.cfg.genericHandlerContract, resourceId, utils.ZeroAddress, depositSig, executeSig)

	expectedMessage := msg.NewGenericTransfer(
		src,
		dst,
		1,
		resourceId,
		hash[:],
	)

	// Create an ERC20 Deposit
	ethtest.CreateGenericDeposit(
		t,
		l.bridgeContract,
		client,
		resourceId,

		dst,
		hash[:],
	)

	verifyMessage(t, writer, expectedMessage, errs)
}

func verifyMessage(t *testing.T, w *mockWriter, expected msg.Message, errs chan error) {
	// Verify message
	select {
	case m := <-w.msgs:
		err := compareMessage(expected, m)
		if err != nil {
			t.Fatal(err)
		}
	case err := <-errs:
		t.Fatalf("Fatal error: %s", err)
	case <-time.After(TestTimeout):
		t.Fatalf("test timed out")
	}
}

func compareMessage(expected, actual msg.Message) error {
	if !reflect.DeepEqual(expected, actual) {
		if !reflect.DeepEqual(expected.Source, actual.Source) {
			return fmt.Errorf("Source doesn't match. \n\tExpected: %#v\n\tGot: %#v\n", expected.Source, actual.Source)
		} else if !reflect.DeepEqual(expected.Destination, actual.Destination) {
			return fmt.Errorf("Destination doesn't match. \n\tExpected: %#v\n\tGot: %#v\n", expected.Destination, actual.Destination)
		} else if !reflect.DeepEqual(expected.DepositNonce, actual.DepositNonce) {
			return fmt.Errorf("Deposit nonce doesn't match. \n\tExpected: %#v\n\tGot: %#v\n", expected.DepositNonce, actual.DepositNonce)
		} else if !reflect.DeepEqual(expected.Payload, actual.Payload) {
			return fmt.Errorf("Payload doesn't match. \n\tExpected: %#v\n\tGot: %#v\n", expected.Payload, actual.Payload)
		}
	}
	return nil
}
