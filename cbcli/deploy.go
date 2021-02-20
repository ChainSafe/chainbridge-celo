package cbcli

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var ErrNoDeploymentFalgsProvided = errors.New("provide at least one deployment flag. For help use help.")

func defineSender(cctx *cli.Context) (*secp256k1.Keypair, error) {
	privateKey := cctx.String("privateKey")
	//jsonWallet := cctx.Path("jsonWallet")
	//jsonWalletPassword := cctx.String("jsonWalletPassword")

	if privateKey != "" {
		kp, err := secp256k1.NewKeypairFromString(privateKey)
		if err != nil {
			return nil, err
		}
		return kp, nil
	}
	//if jsonWallet != "" {
	//
	//}
	return utils.AliceKp, nil
}

func Deploy(cctx *cli.Context) error {
	url := cctx.String("url")
	gasLimit := cctx.Int64("gasLimit")
	gasPrice := cctx.Int64("gasPrice")
	//networkID := cctx.String("networkID")

	sender, err := defineSender(cctx)
	if err != nil {
		return err
	}
	chainID := cctx.Uint64("chainId")
	relayerThreshold := cctx.Int64("relayerThreshold")

	var relayerAddresses []common.Address
	relayerAddressesString := cctx.StringSlice("relayers")
	if len(relayerAddresses) == 0 {
		relayerAddresses = utils.DefaultRelayerAddresses
	} else {
		relayerAddresses = make([]common.Address, len(relayerAddresses))
		for i, addr := range relayerAddressesString {
			relayerAddresses[i] = common.HexToAddress(addr)
		}
	}
	var bridgeAddress common.Address
	bridgeAddressString := cctx.String("bridgeAddress")
	if common.IsHexAddress(bridgeAddressString) {
		bridgeAddress = common.HexToAddress(bridgeAddressString)
	}

	deployments := make([]string, 0)
	if cctx.Bool("all") {
		deployments = append(deployments, []string{"bridge", "erc20Handler", "erc721Handler", "genericHandler", "erc20", "erc721"}...)
	} else {
		if cctx.Bool("bridge") {
			deployments = append(deployments, "bridge")
		}
		if cctx.Bool("erc20Handler") {
			deployments = append(deployments, "erc20Handler")
		}
		if cctx.Bool("erc721Handler") {
			deployments = append(deployments, "erc721Handler")
		}
		if cctx.Bool("genericHandler") {
			deployments = append(deployments, "genericHandler")
		}
		if cctx.Bool("erc20") {
			deployments = append(deployments, "erc20")
		}
		if cctx.Bool("erc721") {
			deployments = append(deployments, "erc721")
		}
	}
	if len(deployments) == 0 {
		return ErrNoDeploymentFalgsProvided
	}
	ethClient, err := client.NewClient(url, false, sender, big.NewInt(gasLimit), big.NewInt(gasPrice))
	if err != nil {
		return err
	}
	deployedContracts := make(map[string]string)
	for _, v := range deployments {
		switch v {
		case "bridge":
			bridgeAddress, err = utils.DeployBridge(ethClient, uint8(chainID), relayerAddresses, big.NewInt(relayerThreshold))
			if err != nil {
				return err
			}
			deployedContracts["bridge"] = bridgeAddress.String()
		case "erc20Handler":
			if bridgeAddress.String() == "" {
				return errors.New("bridge flag or bridgeAddress param should be set for contracts deployments")
			}
			erc20HandlerAddr, err := utils.DeployERC20Handler(ethClient, bridgeAddress)
			deployedContracts["erc20Handler"] = erc20HandlerAddr.String()
			if err != nil {
				return err
			}
		case "erc721Handler":
			if bridgeAddress.String() == "" {
				return errors.New("bridge flag or bridgeAddress param should be set for contracts deployments")
			}
			erc721HandlerAddr, err := utils.DeployERC721Handler(ethClient, bridgeAddress)
			deployedContracts["erc721Handler"] = erc721HandlerAddr.String()
			if err != nil {
				return err
			}
		case "genericHandler":
			if bridgeAddress.String() == "" {
				return errors.New("bridge flag or bridgeAddress param should be set for contracts deployments")
			}
			genericHandlerAddr, err := utils.DeployGenericHandler(ethClient, bridgeAddress)
			deployedContracts["genericHandler"] = genericHandlerAddr.String()
			if err != nil {
				return err
			}
		case "erc20":
			erc20Token, err := utils.DeployERC20Token(ethClient)
			deployedContracts["erc20Token"] = erc20Token.String()
			if err != nil {
				return err
			}
		case "erc721":
			erc721Token, err := utils.DeployERC721Token(ethClient)
			deployedContracts["erc721Token"] = erc721Token.String()
			if err != nil {
				return err
			}
		}
	}
	fmt.Printf("%+v", deployedContracts)
	return nil
}
