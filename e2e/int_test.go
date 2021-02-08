package e2e

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"math/big"
	"testing"
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
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	//contracts *DeployedContracts
	client            *Client
	client2           *Client
	bridgeAddr        common.Address
	erc20HandlerAddr  common.Address
	erc721Addr        common.Address
	genericAddr       common.Address
	erc20ContractAddr common.Address
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite()    {}
func (s *IntegrationTestSuite) TearDownSuite() {}
func (s *IntegrationTestSuite) SetupTest() {
	client, err := NewClient(TestEndpoint, AliceKp)
	if err != nil {
		panic(err)
	}
	s.client = client

	client2, err := NewClient(TestEndpoint2, AliceKp)
	if err != nil {
		panic(err)
	}
	s.client2 = client2

	s.bridgeAddr = common.HexToAddress("0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B")
	s.erc20HandlerAddr = common.HexToAddress("0x3167776db165D8eA0f51790CA2bbf44Db5105ADF")
	s.erc721Addr = common.HexToAddress("0x3f709398808af36ADBA86ACC617FeB7F5B7B193E")
	s.genericAddr = common.HexToAddress("0x2B6Ab4b880A45a07d83Cf4d664Df4Ab85705Bc07")
	s.erc20ContractAddr = common.HexToAddress("0x21605f71845f372A9ed84253d2D024B7B10999f4")
}
func (s *IntegrationTestSuite) TearDownTest() {}

func (s *IntegrationTestSuite) TestDeposit() {
	bridgeContract, err := Bridge.NewBridge(s.bridgeAddr, s.client.Client)
	s.Nil(err)
	erc20Contract, err := erc20.NewERC20PresetMinterPauser(s.erc20ContractAddr, s.client.Client)
	s.Nil(err)

	balBefore, err := erc20Contract.BalanceOf(s.client.CallOpts, AliceKp.CommonAddress())
	s.Nil(err)

	amountToDeposit := big.NewInt(1000000)
	tx, err := makeErc20Deposit(s.client, bridgeContract, s.erc20ContractAddr, BobKp.CommonAddress(), amountToDeposit)
	s.Nil(err)

	receipt, err := waitAndReturnTxReceipt(s.client, tx)
	s.Nil(err)

	// wait for vote log event
	query := buildQuery(s.bridgeAddr, pkg.ProposalEvent, receipt.BlockNumber, receipt.BlockNumber)
	evts, err := s.client2.Client.FilterLogs(context.Background(), query)
	for _, evt := range evts {
		status := evt.Topics[3].Big().Uint64()
		s.True(pkg.IsPassed(uint8(status)))
	}

	balAfter, err := erc20Contract.BalanceOf(s.client.CallOpts, AliceKp.CommonAddress())
	s.Nil(err)
	s.Equal(balBefore.Cmp(big.NewInt(0).Add(balAfter, amountToDeposit)), 0)

	// wait for execution event
}

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

func makeErc20Deposit(client *Client, bridge *Bridge.Bridge, erc20ContractAddr, dest common.Address, amount *big.Int) (*types.Transaction, error) {
	data := constructErc20DepositData(dest.Bytes(), amount)
	err := client.LockNonceAndUpdate()
	if err != nil {
		return nil, err
	}

	src := pkg.ChainId(5)
	resourceID := pkg.SliceTo32Bytes(append(common.LeftPadBytes(erc20ContractAddr.Bytes(), 31), uint8(src)))
	tx, err := bridge.Deposit(client.Opts, 1, resourceID, data)
	if err != nil {
		return nil, err
	}
	client.UnlockNonce()
	return tx, nil
}

func constructErc20DepositData(destRecipient []byte, amount *big.Int) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(amount, 32)...)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...)
	data = append(data, destRecipient...)
	return data
}

func simulate(client *Client, block *big.Int, txHash common.Hash, from common.Address) ([]byte, error) {
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
	log.Debug().Msgf("HEX:   %s", hexutils.BytesToHex(res))
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
func waitAndReturnTxReceipt(client *Client, tx *types.Transaction) (*types.Receipt, error) {
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
