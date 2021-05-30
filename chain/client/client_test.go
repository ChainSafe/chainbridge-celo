package client

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/celo-org/celo-blockchain/accounts/abi/bind"
)

func Test_ClientOpts(t *testing.T) {
	kp, err := secp256k1.GenerateKeypair()
	if err != nil {
		t.Fatal(fmt.Errorf("keyp pair generation error %w", err))
	}
	chainClient := &Client{
		endpoint:    "",
		http:        true,
		kp:          kp,
		maxGasPrice: big.NewInt(1),
		gasLimit:    big.NewInt(1),
		stop:        make(chan int),
		opts:        bind.NewKeyedTransactor(kp.PrivateKey()),
	}
	chainClient.opts.Nonce = big.NewInt(int64(0))
	chainClient.opts.GasLimit = uint64(1)
	chainClient.opts.GasPrice = big.NewInt(0)
	chainClient.opts.Context = context.Background()
	chainClient.opts.Value = big.NewInt(0)

	newOpts := chainClient.OptsCopyWithArgs(OptsWithValue(big.NewInt(10)))

	if newOpts.Value.Cmp(big.NewInt(10)) != 0 {
		t.Fatal(fmt.Sprintf("New opts value is not equal to required value, but : %s", newOpts.Value.String()))
	}
	// Client opts should be the same
	if chainClient.opts.Value.Cmp(big.NewInt(0)) != 0 {
		t.Fatal(fmt.Sprintf("Original opts value is not equal to required value, but : %s", chainClient.opts.Value.String()))
	}
}

func Test_ClientWithArgs(t *testing.T) {
	kp, err := secp256k1.GenerateKeypair()
	if err != nil {
		t.Fatal(fmt.Errorf("keyp pair generation error %w", err))
	}
	chainClient := &Client{
		endpoint:    "",
		http:        true,
		kp:          kp,
		maxGasPrice: big.NewInt(1),
		gasLimit:    big.NewInt(1),
		stop:        make(chan int),
		opts:        bind.NewKeyedTransactor(kp.PrivateKey()),
	}

	chainClient.opts.Value = big.NewInt(0)

	chainClient.ClientWithArgs(ClientWithValue(big.NewInt(10)))

	if chainClient.opts.Value.Cmp(big.NewInt(10)) != 0 {
		t.Fatal()
	}
}
