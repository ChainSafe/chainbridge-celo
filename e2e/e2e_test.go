//nolint
package e2e

import (
	"context"
	"math/big"
	"testing"
	"time"

	erc20 "github.com/ChainSafe/chainbridge-celo/bindings/ERC20PresetMinterPauser"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	//contracts *DeployedContracts
	client            *client.Client
	client2           *client.Client
	bridgeAddr        common.Address
	erc20HandlerAddr  common.Address
	erc721Addr        common.Address
	genericAddr       common.Address
	erc20ContractAddr common.Address
}

func TestRunE2ETests(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite()    {}
func (s *IntegrationTestSuite) TearDownSuite() {}
func (s *IntegrationTestSuite) SetupTest() {
	chainClient, err := client.NewClient(TestEndpoint, false, utils.AliceKp, big.NewInt(utils.DefaultGasLimit), big.NewInt(utils.DefaultGasPrice))
	if err != nil {
		panic(err)
	}
	s.client = chainClient

	client2, err := client.NewClient(TestEndpoint2, false, utils.AliceKp, big.NewInt(utils.DefaultGasLimit), big.NewInt(utils.DefaultGasPrice))
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

// Deposit hash: 0x42782f963df86c5f31f94d9c610445b72d388bd60f788e2cd8ea4bff17824426
func (s *IntegrationTestSuite) TestDeposit() {
	dstAddr := utils.BobKp.CommonAddress()
	erc20Contract, err := erc20.NewERC20PresetMinterPauser(s.erc20ContractAddr, s.client.Client)
	s.Nil(err)
	erc20Contract2, err := erc20.NewERC20PresetMinterPauser(s.erc20ContractAddr, s.client2.Client)
	s.Nil(err)
	senderBalBefore, err := erc20Contract.BalanceOf(s.client.CallOpts(), utils.AliceKp.CommonAddress())
	s.Nil(err)
	destBalanceBefor, err := erc20Contract2.BalanceOf(s.client2.CallOpts(), dstAddr)
	s.Nil(err)

	amountToDeposit := big.NewInt(1000000)
	resourceID := utils.SliceTo32Bytes(append(common.LeftPadBytes(s.erc20ContractAddr.Bytes(), 31), uint8(5)))
	tx, err := utils.MakeErc20Deposit(s.client, s.bridgeAddr, dstAddr, amountToDeposit, resourceID, 1)
	s.Nil(err)
	receipt, err := utils.WaitAndReturnTxReceipt(s.client, tx)
	s.Nil(err)
	log.Debug().Msg(tx.Hash().String())

	//Wait 30 seconds for relayer vote
	time.Sleep(30 * time.Second)

	lp, err := s.client2.LatestBlock()
	s.Nil(err)

	// wait for vote log event
	query := utils.BuildQuery(s.bridgeAddr, utils.ProposalEvent, receipt.BlockNumber, lp)
	evts, err := s.client2.Client.FilterLogs(context.Background(), query)
	var passedEventFound bool
	for _, evt := range evts {
		status := evt.Topics[3].Big().Uint64()
		if utils.IsPassed(uint8(status)) {
			passedEventFound = true
		}
	}
	s.True(passedEventFound)

	senderBalAfter, err := erc20Contract.BalanceOf(s.client.CallOpts(), utils.AliceKp.CommonAddress())
	s.Nil(err)
	s.Equal(senderBalBefore.Cmp(big.NewInt(0).Add(senderBalAfter, amountToDeposit)), 0)

	//Wait 30 seconds for relayer to execute
	time.Sleep(30 * time.Second)
	lp, err = s.client2.LatestBlock()
	s.Nil(err)
	queryExecute := utils.BuildQuery(s.bridgeAddr, utils.ProposalEvent, receipt.BlockNumber, lp)
	s.Nil(err)
	evts2, err := s.client2.Client.FilterLogs(context.Background(), queryExecute)
	var executedEventFound bool
	for _, evt := range evts2 {
		status := evt.Topics[3].Big().Uint64()
		if utils.IsExecuted(uint8(status)) {
			executedEventFound = true
		}
	}
	s.True(executedEventFound)

	destBalanceAfter, err := erc20Contract2.BalanceOf(s.client2.CallOpts(), dstAddr)
	s.Nil(err)
	//Balance has increased
	s.Equal(destBalanceAfter.Cmp(destBalanceBefor), 1)
}

func (s *IntegrationTestSuite) TestMultipleTransactionsInBlock() {
	eveSender, err := client.NewClient(TestEndpoint, false, utils.EveKp, big.NewInt(utils.DefaultGasLimit), big.NewInt(utils.DefaultGasPrice))
	s.Nil(err)

	dstAddr := utils.BobKp.CommonAddress()
	s.Nil(err)
	erc20Contract, err := erc20.NewERC20PresetMinterPauser(s.erc20ContractAddr, s.client.Client)
	s.Nil(err)
	erc20Contract2, err := erc20.NewERC20PresetMinterPauser(s.erc20ContractAddr, s.client2.Client)
	s.Nil(err)
	senderBalBefore, err := erc20Contract.BalanceOf(s.client.CallOpts(), utils.AliceKp.CommonAddress())
	s.Nil(err)
	destBalanceBefor, err := erc20Contract2.BalanceOf(s.client2.CallOpts(), dstAddr)
	s.Nil(err)

	amountToDeposit := big.NewInt(1000000)

	for i := 0; i <= 20; i += 1 {
		go sendOneWeiWithDelay(eveSender)
	}

	resourceID := utils.SliceTo32Bytes(append(common.LeftPadBytes(s.erc20ContractAddr.Bytes(), 31), uint8(5)))
	tx, err := utils.MakeErc20Deposit(s.client, s.bridgeAddr, dstAddr, amountToDeposit, resourceID, 1)
	s.Nil(err)
	receipt, err := utils.WaitAndReturnTxReceipt(s.client, tx)
	s.Nil(err)
	log.Debug().Msg(tx.Hash().String())

	//Wait 30 seconds for relayer vote
	time.Sleep(30 * time.Second)

	lp, err := s.client2.LatestBlock()
	s.Nil(err)
	// wait for vote log event
	query := utils.BuildQuery(s.bridgeAddr, utils.ProposalEvent, receipt.BlockNumber, lp)
	evts, err := s.client2.Client.FilterLogs(context.Background(), query)
	var passedEventFound bool
	for _, evt := range evts {
		status := evt.Topics[3].Big().Uint64()
		if utils.IsPassed(uint8(status)) {
			passedEventFound = true
		}
	}
	s.True(passedEventFound)

	senderBalAfter, err := erc20Contract.BalanceOf(s.client.CallOpts(), utils.AliceKp.CommonAddress())
	s.Nil(err)
	s.Equal(senderBalBefore.Cmp(big.NewInt(0).Add(senderBalAfter, amountToDeposit)), 0)

	//Wait 30 seconds for relayer to execute
	time.Sleep(30 * time.Second)
	lp, err = s.client2.LatestBlock()
	s.Nil(err)
	queryExecute := utils.BuildQuery(s.bridgeAddr, utils.ProposalEvent, receipt.BlockNumber, lp)
	s.Nil(err)
	evts2, err := s.client2.Client.FilterLogs(context.Background(), queryExecute)
	var executedEventFound bool
	for _, evt := range evts2 {
		status := evt.Topics[3].Big().Uint64()
		if utils.IsExecuted(uint8(status)) {
			executedEventFound = true
		}
	}
	s.True(executedEventFound)

	destBalanceAfter, err := erc20Contract2.BalanceOf(s.client2.CallOpts(), dstAddr)
	s.Nil(err)
	//Balance has increased
	s.Equal(destBalanceAfter.Cmp(destBalanceBefor), 1)
}

//nolint
//func (s *IntegrationTestSuite) TestSimulate() {
//	block := big.NewInt(100)
//	hash := common.HexToHash("0x1991265fe7bfad3cd2cce0cc4e0d4e72e05aa201ae03845df24b985098b6298e")
//	res, err := utils.Simulate(s.client, block, hash, utils.AliceKp.CommonAddress())
//	s.Nil(err)
//	hexres := common.Bytes2Hex(res)
//	log.Info().Msgf("simulate result: %s", hexres)
//
//	rec, err := s.client.TransactionReceipt(context.TODO(), hash)
//	s.Nil(err)
//	log.Debug().Msgf("%+v", rec)
//}
