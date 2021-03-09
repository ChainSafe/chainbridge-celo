package cliutils

import (
	"fmt"
	"github.com/pkg/errors"

	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

func DefineSender(cctx *cli.Context) (*secp256k1.Keypair, error) {
	privateKey := cctx.String("privateKey")
	if privateKey != "" {
		kp, err := secp256k1.NewKeypairFromString(privateKey)
		if err != nil {
			return nil, err
		}
		return kp, nil
	}
	return utils.AliceKp, nil
}

func DefineBridgeAddress(cctx *cli.Context) (common.Address, error) {
	bridge := cctx.String("bridge")
	if !common.IsHexAddress(bridge) {
		return common.Address{}, errors.New(fmt.Sprintf("invalid bridge address %s", bridge))
	}
	return common.HexToAddress(bridge), nil
}
