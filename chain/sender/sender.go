// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package sender

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000

var BlockRetryInterval = time.Second * 5

type Sender struct {
	*ethclient.Client
	endpoint    string
	http        bool
	kp          *secp256k1.Keypair
	gasLimit    *big.Int
	maxGasPrice *big.Int
	opts        *bind.TransactOpts
	callOpts    *bind.CallOpts
	nonce       uint64
	nonceLock   sync.Mutex
	optsLock    sync.Mutex
	stop        chan int // All routines should exit when this channel is closed
}

type LogFilterWithLatestBlock interface {
	FilterLogs(ctx context.Context, q eth.FilterQuery) ([]types.Log, error)
	LatestBlock() (*big.Int, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

// NewConnection returns an uninitialized connection, must call Sender.Connect() before using.
func NewSender(endpoint string, http bool, kp *secp256k1.Keypair, gasLimit *big.Int, gasPrice *big.Int) (*Sender, error) {
	c := &Sender{
		endpoint:    endpoint,
		http:        http,
		kp:          kp,
		maxGasPrice: gasPrice,
		gasLimit:    gasLimit,
		stop:        make(chan int),
	}
	if err := c.Connect(); err != nil {
		return nil, err
	}
	return c, nil
}

// Connect starts the ethereum WS connection
func (c *Sender) Connect() error {
	log.Info().Str("url", c.endpoint).Msg("Connecting to ethereum chain...")
	var rpcClient *rpc.Client
	var err error
	// Start http or ws client
	if c.http {
		rpcClient, err = rpc.DialHTTP(c.endpoint)
	} else {
		rpcClient, err = rpc.DialWebsocket(context.Background(), c.endpoint, "/ws")
	}
	if err != nil {
		return err
	}
	c.Client = ethclient.NewClient(rpcClient)

	// Construct tx opts, call opts, and nonce mechanism
	opts, _, err := c.newTransactOpts(big.NewInt(0), c.gasLimit, c.maxGasPrice)
	if err != nil {
		return err
	}
	c.opts = opts
	c.nonce = 0
	c.callOpts = &bind.CallOpts{From: c.kp.CommonAddress()}
	return nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (c *Sender) newTransactOpts(value, gasLimit, gasPrice *big.Int) (*bind.TransactOpts, uint64, error) {
	privateKey := c.kp.PrivateKey()
	address := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := c.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, 0, err
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = uint64(gasLimit.Int64())
	auth.GasPrice = gasPrice
	auth.Context = context.Background()

	return auth, nonce, nil
}

func (c *Sender) Keypair() *secp256k1.Keypair {
	return c.kp
}

func (c *Sender) Opts() *bind.TransactOpts {
	return c.opts
}

func (c *Sender) CallOpts() *bind.CallOpts {
	return c.callOpts
}

func (c *Sender) LockAndUpdateNonce() error {
	c.nonceLock.Lock()
	nonce, err := c.PendingNonceAt(context.Background(), c.opts.From)
	if err != nil {
		c.nonceLock.Unlock()
		return err
	}
	c.opts.Nonce.SetUint64(nonce)
	return nil
}

func (c *Sender) UnlockNonce() {
	c.nonceLock.Unlock()
}

func (c *Sender) UnlockOpts() {
	c.optsLock.Unlock()
}

// LatestBlock returns the latest block from the current chain
func (c *Sender) LatestBlock() (*big.Int, error) {
	header, err := c.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

// EnsureHasBytecode asserts if contract code exists at the specified address
func (c *Sender) EnsureHasBytecode(addr ethcommon.Address) error {
	code, err := c.CodeAt(context.Background(), addr, nil)
	if err != nil {
		return err
	}

	if len(code) == 0 {
		return fmt.Errorf("no bytecode found at %s", addr.Hex())
	}
	return nil
}

// WaitForBlock will poll for the block number until the current block is equal or greater than
func (c *Sender) WaitForBlock(block *big.Int) error {
	for {
		select {
		case <-c.stop:
			return errors.New("connection terminated")
		default:
			currBlock, err := c.LatestBlock()
			if err != nil {
				return err
			}

			// Equal or greater than target
			if currBlock.Cmp(block) >= 0 {
				return nil
			}
			log.Trace().Interface("target", block).Interface("current", currBlock).Msg("Block not ready, waiting")
			time.Sleep(BlockRetryInterval)
			continue
		}
	}
}

// LockAndUpdateOpts acquires a lock on the opts before updating the nonce
// and gas price.
func (c *Sender) LockAndUpdateOpts() error {
	c.optsLock.Lock()

	gasPrice, err := c.SafeEstimateGas(context.TODO())
	if err != nil {
		return err
	}
	c.opts.GasPrice = gasPrice

	nonce, err := c.PendingNonceAt(context.Background(), c.opts.From)
	if err != nil {
		c.optsLock.Unlock()
		return err
	}
	c.opts.Nonce.SetUint64(nonce)
	return nil
}

func (c *Sender) SafeEstimateGas(ctx context.Context) (*big.Int, error) {
	gasPrice, err := c.SuggestGasPrice(context.TODO())
	if err != nil {
		return nil, err
	}

	// Check we aren't exceeding our limit

	if gasPrice.Cmp(c.maxGasPrice) == 1 {
		return c.maxGasPrice, nil
	} else {
		return gasPrice, nil
	}
}

// Close terminates the client connection and stops any running routines
func (c *Sender) Close() {
	if c.Client != nil {
		c.Client.Close()
	}
}