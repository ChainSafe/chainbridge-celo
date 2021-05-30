package e2e

import (
	"context"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"math/big"
	"math/rand"
	"time"

	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/core/types"
)

func sendOneWeiWithDelay(sender *client.Client) (*types.Transaction, error) {
	r := rand.Intn(700) + 300
	time.Sleep(time.Duration(r) * time.Millisecond)
	return sendOneWei(sender)
}

func sendOneWei(sender *client.Client) (*types.Transaction, error) {
	err := sender.LockAndUpdateOpts()
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(sender.Opts().Nonce.Uint64(), utils.AliceKp.CommonAddress(), big.NewInt(1), sender.Opts().GasLimit, sender.Opts().GasPrice, sender.Opts().FeeCurrency, sender.Opts().GatewayFeeRecipient, sender.Opts().GatewayFee, nil)

	// Final Step
	signedTx, err := sender.Opts().Signer(types.HomesteadSigner{}, sender.Opts().From, tx)
	if err != nil {
		return nil, err
	}

	err = sender.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}
	sender.UnlockOpts()
	err = utils.WaitForTx(sender, signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
