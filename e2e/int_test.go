package e2e

import (
	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	erc20 "github.com/ChainSafe/chainbridge-celo/bindings/ERC20PresetMinterPauser"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
	"math/big"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite
	//contracts *DeployedContracts
	client *Client
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite()    {}
func (s *IntegrationTestSuite) TearDownSuite() {}
func (s *IntegrationTestSuite) SetupTest() {
	//
	client, err := NewClient(TestEndpoint, AliceKp)
	if err != nil {
		panic(err)
	}
	//
	s.client = client
	//dpc, err := DeployContracts(client, 1, big.NewInt(1))
	//if err != nil {
	//	panic(err)
	//}
	//s.contracts = dpc
	//log.Debug().Msgf("Bridge %s erc20 /n \n  handler %s /n erc721 handler %s /n\n generic handler %s /n", dpc.BridgeAddress, dpc.ERC20HandlerAddress, dpc.ERC721HandlerAddress, dpc.GenericHandlerAddress)
	//RegisterResource(client, dpc.BridgeAddress, dpc.ERC20HandlerAddress, )
	//MintTokens()
}

//"message":"Bridge 0xc279648CE5cAa25B9bA753dAb0Dfef44A069BaF4 "
//"\r\nerc20 handler 0x84b141Aada70e2B0C3Ec25d24E81032328ea1b1A "
//"\r\nerc721 handler 0x771ce6cda91d18eB76029d17aCcd85834F5C0303 "
//"\r\ngeneric handler 0x209373B5047f1F8a619977D3C9bCB0A6b10126c1 "
//"\r\nerc20Contract 0x39863e3eDB5255dB93bBf8E76c12578357dBe6c7"

// no elections addresses
// "Bridge 0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B
// erc20 handler 0x3167776db165D8eA0f51790CA2bbf44Db5105ADF \r\ne
// rc721 handler 0x3f709398808af36ADBA86ACC617FeB7F5B7B193E \r\n
// generic handler 0x2B6Ab4b880A45a07d83Cf4d664Df4Ab85705Bc07 \r\n
// erc20Contract 0x21605f71845f372A9ed84253d2D024B7B10999f4"}
func (s *IntegrationTestSuite) TearDownTest() {}

func (s *IntegrationTestSuite) TestGetBalance() {
	//addr := common.HexToAddress("0xF4314cb9046bECe6AA54bb9533155434d0c76909")
	//bal, err := s.client.Client.BalanceAt(context.TODO(), addr, nil)
	//s.Nil(err)
	//log.Debug().Msgf("%s", bal.String())
	//bal, err = s.client.Client.BalanceAt(context.TODO(), AliceKp.CommonAddress(), nil)
	//s.Nil(err)
	//log.Debug().Msgf("%s %s", bal.String(), AliceKp.Address())

	erc20ContractAddr := common.HexToAddress("0x21605f71845f372A9ed84253d2D024B7B10999f4")
	erc20Contract, err := erc20.NewERC20PresetMinterPauser(erc20ContractAddr, s.client.Client)
	s.Nil(err)
	balT, err := erc20Contract.BalanceOf(s.client.CallOpts, AliceKp.CommonAddress())
	s.Nil(err)

	log.Debug().Msgf("BALANCE: %s", balT.String())
}

func (s *IntegrationTestSuite) TestDeposit() {
	bridgeContractAddr := common.HexToAddress("0xc279648CE5cAa25B9bA753dAb0Dfef44A069BaF4")
	bridgeContract, err := Bridge.NewBridge(bridgeContractAddr, s.client.Client)
	erc20ContractAddr := common.HexToAddress("0x39863e3eDB5255dB93bBf8E76c12578357dBe6c7")
	erc20Contract, err := erc20.NewERC20PresetMinterPauser(erc20ContractAddr, s.client.Client)
	s.Nil(err)

	data := ConstructErc20DepositData(AliceKp.CommonAddress().Bytes(), big.NewInt(120000000))

	err = s.client.LockNonceAndUpdate()
	s.Nil(err)

	src := msg.ChainId(5)
	resourceId := msg.ResourceIdFromSlice(append(common.LeftPadBytes(erc20ContractAddr.Bytes(), 31), uint8(src)))

	tx, err := bridgeContract.Deposit(s.client.Opts, 1, resourceId, data)
	s.Nil(err)
	s.NotNil(tx)

	err = WaitForTx(s.client, tx)
	s.Nil(err)

	balT, err := erc20Contract.BalanceOf(s.client.CallOpts, AliceKp.CommonAddress())
	s.Nil(err)
	log.Debug().Msgf("BALANCE: %s", balT.String())
	s.client.UnlockNonce()

}

func (s *IntegrationTestSuite) TestNonce() {
	n, err := s.client.Client.NonceAt(context.TODO(), AliceKp.CommonAddress(), nil)
	s.Nil(err)
	log.Debug().Msgf("NONCE: %v", n)
	log.Debug().Msg(AliceKp.Address())
	//10
}

func ConstructErc20DepositData(destRecipient []byte, amount *big.Int) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(amount, 32)...)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...)
	data = append(data, destRecipient...)
	return data
}
