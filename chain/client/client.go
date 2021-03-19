// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package client

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

type Client struct {
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

// NewConnection returns an uninitialized connection, must call Client.Connect() before using.
func NewClient(endpoint string, http bool, kp *secp256k1.Keypair, gasLimit *big.Int, gasPrice *big.Int) (*Client, error) {
	c := &Client{
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
func (c *Client) Connect() error {
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
	opts, _, err := c.newTransactOpts(c.gasLimit, c.maxGasPrice)
	if err != nil {
		return err
	}
	c.opts = opts
	c.nonce = 0
	c.callOpts = &bind.CallOpts{From: c.kp.CommonAddress()}
	return nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (c *Client) newTransactOpts(gasLimit, gasPrice *big.Int) (*bind.TransactOpts, uint64, error) {
	privateKey := c.kp.PrivateKey()
	address := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := c.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, 0, err
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = uint64(gasLimit.Int64())
	auth.GasPrice = gasPrice
	auth.Context = context.Background()

	return auth, nonce, nil
}

func (c *Client) Keypair() *secp256k1.Keypair {
	return c.kp
}

func (c *Client) Opts() *bind.TransactOpts {
	return c.opts
}

func (c *Client) OptsCopyWithArgs(opts ...func(*bind.TransactOpts)) *bind.TransactOpts {
	copyOfOpts := *c.opts
	for _, opt := range opts {
		opt(&copyOfOpts)
	}
	return &copyOfOpts
}

func OptsWithValue(value *big.Int) func(*bind.TransactOpts) {
	return func(opts *bind.TransactOpts) {
		opts.Value = value
	}
}

func (c *Client) ClientWithArgs(args ...func(*Client)) *Client {
	for _, arg := range args {
		arg(c)
	}
	return c
}

func TheClientWithValue(value *big.Int) func(*Client) {
	return func(c *Client) {
		c.opts.Value = value
	}
}

func (c *Client) CallOpts() *bind.CallOpts {
	return c.callOpts
}

func (c *Client) UnlockNonce() {
	c.nonceLock.Unlock()
}

func (c *Client) UnlockOpts() {
	c.optsLock.Unlock()
}

// LatestBlock returns the latest block from the current chain
func (c *Client) LatestBlock() (*big.Int, error) {
	header, err := c.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

// EnsureHasBytecode asserts if contract code exists at the specified address
func (c *Client) EnsureHasBytecode(addr ethcommon.Address) error {
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
func (c *Client) WaitForBlock(block *big.Int) error {
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
func (c *Client) LockAndUpdateOpts() error {
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

func (c *Client) LockAndUpdateNonce() error {
	c.nonceLock.Lock()
	nonce, err := c.PendingNonceAt(context.Background(), c.opts.From)
	if err != nil {
		c.nonceLock.Unlock()
		return err
	}
	c.opts.Nonce.SetUint64(nonce)
	return nil
}

func (c *Client) SafeEstimateGas(ctx context.Context) (*big.Int, error) {
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
func (c *Client) Close() {
	if c.Client != nil {
		c.Client.Close()
	}
}
