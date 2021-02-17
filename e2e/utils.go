package e2e

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ChainSafe/chainbridge-celo/chain/sender"
	"math/big"
	"math/rand"
	"time"

	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	erc20 "github.com/ChainSafe/chainbridge-celo/bindings/ERC20PresetMinterPauser"
	"github.com/ChainSafe/chainbridge-celo/utils"
	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/status-im/keycard-go/hexutils"
)

func makeErc20Deposit(client *sender.Sender, bridge *Bridge.Bridge, erc20ContractAddr, dest common.Address, amount *big.Int) (*types.Transaction, error) {
	data := constructErc20DepositData(dest.Bytes(), amount)
	err := client.LockAndUpdateOpts()
	if err != nil {
		return nil, err
	}

	src := utils.ChainId(5)
	resourceID := utils.SliceTo32Bytes(append(common.LeftPadBytes(erc20ContractAddr.Bytes(), 31), uint8(src)))
	tx, err := bridge.Deposit(client.Opts(), 1, resourceID, data)
	if err != nil {
		return nil, err
	}
	client.UnlockOpts()
	return tx, nil
}

func constructErc20DepositData(destRecipient []byte, amount *big.Int) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(amount, 32)...)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...)
	data = append(data, destRecipient...)
	return data
}

//nolint
func simulate(client *sender.Sender, block *big.Int, txHash common.Hash, from common.Address) ([]byte, error) {
	tx, _, err := client.Client.TransactionByHash(context.TODO(), txHash)
	if err != nil {
		return nil, err
	}
	msg := eth.CallMsg{
		From:                from,
		To:                  tx.To(),
		Gas:                 tx.Gas(),
		FeeCurrency:         tx.FeeCurrency(),
		GatewayFeeRecipient: tx.GatewayFeeRecipient(),
		GatewayFee:          tx.GatewayFee(),
		GasPrice:            tx.GasPrice(),
		Value:               tx.Value(),
		Data:                tx.Data(),
	}
	res, err := client.Client.CallContract(context.TODO(), msg, block)
	if err != nil {
		return nil, err
	}
	bs, err := hex.DecodeString(hexutils.BytesToHex(res))
	if err != nil {
		panic(err)
	}
	log.Debug().Msg(string(bs))
	return nil, nil
}

func buildQuery(contract common.Address, sig utils.EventSig, startBlock *big.Int, endBlock *big.Int) eth.FilterQuery {
	query := eth.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Addresses: []common.Address{contract},
		Topics: [][]common.Hash{
			{sig.GetTopic()},
		},
	}
	return query
}

// WaitForTx will query the chain at ExpectedBlockTime intervals, until a receipt is returned.
// Returns an error if the tx failed.
func WaitForTx(client *sender.Sender, tx *types.Transaction) error {
	retry := 10
	for retry > 0 {
		receipt, err := client.Client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			retry--
			time.Sleep(ExpectedBlockTime)
			continue
		}

		if receipt.Status != 1 {
			return fmt.Errorf("transaction failed on chain")
		}
		return nil
	}
	return nil
}

// WaitForTx will query the chain at ExpectedBlockTime intervals, until a receipt is returned.
// Returns an error if the tx failed.
func waitAndReturnTxReceipt(client *sender.Sender, tx *types.Transaction) (*types.Receipt, error) {
	retry := 10
	for retry > 0 {
		receipt, err := client.Client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			retry--
			time.Sleep(ExpectedBlockTime)
			continue
		}
		if receipt.Status != 1 {
			return receipt, fmt.Errorf("transaction failed on chain")
		}
		return receipt, nil
	}
	return nil, errors.New("Tx do not appear")
}

//nolint
func transfer(client *sender.Sender, erc20 *erc20.ERC20PresetMinterPauser, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	err := client.LockAndUpdateOpts()
	if err != nil {
		return nil, err
	}

	tx, err := erc20.Transfer(client.Opts(), recipient, amount)
	if err != nil {
		return nil, err
	}
	client.UnlockOpts()
	return tx, nil
}

func sendOneWeiWithDelay(sender *sender.Sender) (*types.Transaction, error) {
	r := rand.Intn(700) + 300
	time.Sleep(time.Duration(r) * time.Millisecond)
	return sendOneWei(sender)
}

func sendOneWei(sender *sender.Sender) (*types.Transaction, error) {
	err := sender.LockAndUpdateOpts()
	if err != nil {
		return nil, err
	}
	tx := types.NewTransaction(sender.Opts().Nonce.Uint64(), AliceKp.CommonAddress(), big.NewInt(1), sender.Opts().GasLimit, sender.Opts().GasPrice, sender.Opts().FeeCurrency, sender.Opts().GatewayFeeRecipient, sender.Opts().GatewayFee, nil)

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
	err = WaitForTx(sender, signedTx)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}
