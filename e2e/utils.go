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
	"github.com/ChainSafe/chainbridge-celo/pkg"
	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/status-im/keycard-go/hexutils"
)

//
//{"level":"trace","src":1,"dest":1,"nonce":1,"rId":"000000000000000000000021605f71845f372a9ed84253d2d024b7b10999f405","time":"2021-02-05T12:51:01+02:00","message":"Routing message"}
//{"level":"info","type":"FungibleTransfer","src":1,"dst":1,"nonce":1,"rId":"000000000000000000000021605f71845f372a9ed84253d2d024b7b10999f405","time":"2021-02-05T12:51:01+02:00","message":"Attempting to resolve message"}
//{"level":"info","src":1,"nonce":1,"time":"2021-02-05T12:51:01+02:00","message":"Creating erc20 proposal"}
//{"level":"info","src":1,"nonce":1,"time":"2021-02-05T12:51:01+02:00","message":"Watching for finalization event"}
//{"level":"trace","block":737,"src":1,"nonce":1,"time":"2021-02-05T12:51:01+02:00","message":"No finalization event found in current block"}
//{"level":"trace","target":738,"current":737,"time":"2021-02-05T12:51:01+02:00","message":"Block not ready, waiting"}
//{"level":"info","tx":"0xe8fdc387d8a51bfe9c62d501aedd9f270647bc83895c1a683a0673b7539db834","src":1,"depositNonce":1,"time":"2021-02-05T12:51:01+02:00","message":"Submitted proposal vote"}
//{"level":"trace","src":1,"nonce":1,"time":"2021-02-05T12:51:06+02:00","message":"Ignoring event"}
//{"level":"info","source":1,"dest":1,"nonce":1,"tx":"0xbd7a6e74c3c57bde06464bc3997397d5bbc8f44095ed4654486e06b17e838d26","time":"2021-02-05T12:51:06+02:00","message":"Submitted proposal execution"}

func makeErc20Deposit(client *sender.Sender, bridge *Bridge.Bridge, erc20ContractAddr, dest common.Address, amount *big.Int) (*types.Transaction, error) {
	data := constructErc20DepositData(dest.Bytes(), amount)
	err := client.LockAndUpdateOpts()
	if err != nil {
		return nil, err
	}

	src := pkg.ChainId(5)
	resourceID := pkg.SliceTo32Bytes(append(common.LeftPadBytes(erc20ContractAddr.Bytes(), 31), uint8(src)))
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

func buildQuery(contract common.Address, sig pkg.EventSig, startBlock *big.Int, endBlock *big.Int) eth.FilterQuery {
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
