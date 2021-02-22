package cliutils

import (
	"github.com/ChainSafe/chainbridge-celo/utils"
	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/urfave/cli/v2"
)

func DefineSender(cctx *cli.Context) (*secp256k1.Keypair, error) {
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
