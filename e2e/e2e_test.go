//nolint
package e2e

import (
	"context"
	"github.com/ChainSafe/chainbridge-celo/chain/sender"
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	erc20 "github.com/ChainSafe/chainbridge-celo/bindings/ERC20PresetMinterPauser"
	"github.com/ChainSafe/chainbridge-celo/pkg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	//contracts *DeployedContracts
	sender            *sender.Sender
	sender2           *sender.Sender
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
	client, err := sender.NewSender(TestEndpoint, false, AliceKp, big.NewInt(DefaultGasLimit), big.NewInt(DefaultGasPrice))
	if err != nil {
		panic(err)
	}
	s.sender = client

	client2, err := sender.NewSender(TestEndpoint2, false, AliceKp, big.NewInt(DefaultGasLimit), big.NewInt(DefaultGasPrice))
	if err != nil {
		panic(err)
	}
	s.sender2 = client2

	s.bridgeAddr = common.HexToAddress("0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B")
	s.erc20HandlerAddr = common.HexToAddress("0x3167776db165D8eA0f51790CA2bbf44Db5105ADF")
	s.erc721Addr = common.HexToAddress("0x3f709398808af36ADBA86ACC617FeB7F5B7B193E")
	s.genericAddr = common.HexToAddress("0x2B6Ab4b880A45a07d83Cf4d664Df4Ab85705Bc07")
	s.erc20ContractAddr = common.HexToAddress("0x21605f71845f372A9ed84253d2D024B7B10999f4")
}
func (s *IntegrationTestSuite) TearDownTest() {}

// Deposit hash: 0x42782f963df86c5f31f94d9c610445b72d388bd60f788e2cd8ea4bff17824426
func (s *IntegrationTestSuite) TestDeposit() {
	dstAddr := BobKp.CommonAddress()
	bridgeContract, err := Bridge.NewBridge(s.bridgeAddr, s.sender.Client)
	s.Nil(err)
	erc20Contract, err := erc20.NewERC20PresetMinterPauser(s.erc20ContractAddr, s.sender.Client)
	s.Nil(err)
	erc20Contract2, err := erc20.NewERC20PresetMinterPauser(s.erc20ContractAddr, s.sender2.Client)
	s.Nil(err)
	senderBalBefore, err := erc20Contract.BalanceOf(s.sender.CallOpts(), AliceKp.CommonAddress())
	s.Nil(err)
	destBalanceBefor, err := erc20Contract2.BalanceOf(s.sender2.CallOpts(), dstAddr)
	s.Nil(err)

	amountToDeposit := big.NewInt(1000000)
	tx, err := makeErc20Deposit(s.sender, bridgeContract, s.erc20ContractAddr, dstAddr, amountToDeposit)
	s.Nil(err)
	receipt, err := waitAndReturnTxReceipt(s.sender, tx)
	s.Nil(err)
	log.Debug().Msg(tx.Hash().String())

	time.Sleep(20 * time.Second)

	// wait for vote log event
	query := buildQuery(s.bridgeAddr, pkg.ProposalEvent, receipt.BlockNumber, big.NewInt(0).Add(receipt.BlockNumber, big.NewInt(20)))
	evts, err := s.sender2.Client.FilterLogs(context.Background(), query)
	var passedEventFound bool
	for _, evt := range evts {
		status := evt.Topics[3].Big().Uint64()
		if pkg.IsPassed(uint8(status)) {
			passedEventFound = true
		}
	}
	s.True(passedEventFound)

	senderBalAfter, err := erc20Contract.BalanceOf(s.sender.CallOpts(), AliceKp.CommonAddress())
	s.Nil(err)
	s.Equal(senderBalBefore.Cmp(big.NewInt(0).Add(senderBalAfter, amountToDeposit)), 0)

	//wait for execution event
	queryExecute := buildQuery(s.bridgeAddr, pkg.ProposalEvent, big.NewInt(0).Add(receipt.BlockNumber, big.NewInt(1)), big.NewInt(0).Add(receipt.BlockNumber, big.NewInt(20)))
	s.Nil(err)
	evts2, err := s.sender2.Client.FilterLogs(context.Background(), queryExecute)
	var executedEventFound bool
	for _, evt := range evts2 {
		status := evt.Topics[3].Big().Uint64()
		if pkg.IsExecuted(uint8(status)) {
			executedEventFound = true
		}
	}
	s.True(executedEventFound)

	destBalanceAfter, err := erc20Contract2.BalanceOf(s.sender2.CallOpts(), dstAddr)
	s.Nil(err)
	//Balance has increased
	s.Equal(destBalanceAfter.Cmp(destBalanceBefor), 1)
}

func (s *IntegrationTestSuite) TestMultipleTransactionsInBlock() {
	eveSender, err := sender.NewSender(TestEndpoint, false, EveKp, big.NewInt(DefaultGasLimit), big.NewInt(DefaultGasPrice))
	s.Nil(err)

	dstAddr := BobKp.CommonAddress()
	bridgeContract, err := Bridge.NewBridge(s.bridgeAddr, s.sender.Client)
	s.Nil(err)
	erc20Contract, err := erc20.NewERC20PresetMinterPauser(s.erc20ContractAddr, s.sender.Client)
	s.Nil(err)
	erc20Contract2, err := erc20.NewERC20PresetMinterPauser(s.erc20ContractAddr, s.sender2.Client)
	s.Nil(err)
	senderBalBefore, err := erc20Contract.BalanceOf(s.sender.CallOpts(), AliceKp.CommonAddress())
	s.Nil(err)
	destBalanceBefor, err := erc20Contract2.BalanceOf(s.sender2.CallOpts(), dstAddr)
	s.Nil(err)

	amountToDeposit := big.NewInt(1000000)

	for i := 0; i <= 20; i += 1 {
		go sendOneWeiWithDelay(eveSender)
	}
	tx, err := makeErc20Deposit(s.sender, bridgeContract, s.erc20ContractAddr, dstAddr, amountToDeposit)
	s.Nil(err)
	receipt, err := waitAndReturnTxReceipt(s.sender, tx)
	s.Nil(err)
	log.Debug().Msg(tx.Hash().String())

	time.Sleep(20 * time.Second)

	// wait for vote log event
	query := buildQuery(s.bridgeAddr, pkg.ProposalEvent, receipt.BlockNumber, big.NewInt(0).Add(receipt.BlockNumber, big.NewInt(20)))
	evts, err := s.sender2.Client.FilterLogs(context.Background(), query)
	var passedEventFound bool
	for _, evt := range evts {
		status := evt.Topics[3].Big().Uint64()
		if pkg.IsPassed(uint8(status)) {
			passedEventFound = true
		}
	}
	s.True(passedEventFound)

	senderBalAfter, err := erc20Contract.BalanceOf(s.sender.CallOpts(), AliceKp.CommonAddress())
	s.Nil(err)
	s.Equal(senderBalBefore.Cmp(big.NewInt(0).Add(senderBalAfter, amountToDeposit)), 0)

	//wait for execution event
	queryExecute := buildQuery(s.bridgeAddr, pkg.ProposalEvent, big.NewInt(0).Add(receipt.BlockNumber, big.NewInt(1)), big.NewInt(0).Add(receipt.BlockNumber, big.NewInt(20)))
	s.Nil(err)
	evts2, err := s.sender2.Client.FilterLogs(context.Background(), queryExecute)
	var executedEventFound bool
	for _, evt := range evts2 {
		status := evt.Topics[3].Big().Uint64()
		if pkg.IsExecuted(uint8(status)) {
			executedEventFound = true
		}
	}
	s.True(executedEventFound)

	destBalanceAfter, err := erc20Contract2.BalanceOf(s.sender2.CallOpts(), dstAddr)
	s.Nil(err)
	//Balance has increased
	s.Equal(destBalanceAfter.Cmp(destBalanceBefor), 1)
}
