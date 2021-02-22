package cbcli

import "github.com/urfave/cli/v2"

var CLICMD = cli.Command{
	Name:        "cli",
	Description: "This CLI supports on-chain interactions with components of ChainBridge",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "url",
			Value: "http://localhost:8545",
			Usage: "RPC url of blockchain node",
		},
		&cli.Uint64Flag{
			Name:  "gasLimit",
			Value: 6721975,
			Usage: "gasLimit used in transactions",
		},
		&cli.Uint64Flag{
			Name:  "gasPrice",
			Value: 20000000000,
			Usage: "gasPrice used for transactions",
		},
		&cli.Uint64Flag{
			Name:  "networkID",
			Value: 0,
			Usage: "networkID",
		},
		&cli.StringFlag{
			Name:  "privateKey",
			Value: "",
			Usage: "Private key to use",
		},
		&cli.PathFlag{
			Name:  "jsonWallet",
			Value: "",
			Usage: "Encrypted JSON wallet",
		},
		&cli.StringFlag{
			Name:  "jsonWalletPassword",
			Value: "",
			Usage: "Password for encrypted JSON wallet",
		},
	},
	Subcommands: []*cli.Command{deployCMD, registerResourceCMD, registerGenericResourceCMD, setBurnCMD, cancelProposalCMD, queryProposalCMD, queryResourceCMD},
}
