package admin

import (
	"flag"
	"github.com/urfave/cli/v2"
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

const TestEndpoint = "ws://localhost:8545"
const TestEndpoint2 = "ws://localhost:8547"

type AdminCLIIntegrationTestSuite struct {
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
	suite.Run(t, new(AdminCLIIntegrationTestSuite))
}

func (s *AdminCLIIntegrationTestSuite) SetupSuite()    {}
func (s *AdminCLIIntegrationTestSuite) TearDownSuite() {}
func (s *AdminCLIIntegrationTestSuite) SetupTest() {
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

func (s *AdminCLIIntegrationTestSuite) TestIsRelayer() {
	set := flag.NewFlagSet("isRelayer test flagset", 0)
	set.String("url", TestEndpoint, "")
	set.Uint64("gasLimit", utils.DefaultGasLimit, "")
	set.Uint64("gasPrice", utils.DefaultGasPrice, "")
	set.String("bridge", "0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B", "")
	set.String("relayer", utils.AliceKp.Address(), "")
	ctx := cli.NewContext(nil, set, nil)
	err := isRelayer(ctx)
	s.Nil(err)
}

func (s *AdminCLIIntegrationTestSuite) TestIsRelayerWrongAddress() {
	set := flag.NewFlagSet("isRelayer test flagset", 0)
	set.String("url", TestEndpoint, "")
	set.Uint64("gasLimit", utils.DefaultGasLimit, "")
	set.Uint64("gasPrice", utils.DefaultGasPrice, "")
	set.String("bridge", "123", "")
	set.String("relayer", utils.AliceKp.Address(), "")
	ctx := cli.NewContext(nil, set, nil)
	err := isRelayer(ctx)
	s.NotNil(err)
}
