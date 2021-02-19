// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package main

import (
	"github.com/ChainSafe/chainbridge-celo/cbcli"
	"os"

	"github.com/ChainSafe/chainbridge-celo/cmd"
	"github.com/ChainSafe/chainbridge-celo/cmdutils/testutils"
	"github.com/ChainSafe/chainbridge-celo/e2e"
	"github.com/ChainSafe/chainbridge-celo/flags"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

var cliFlags = []cli.Flag{
	flags.ConfigFileFlag, // path to config file
	flags.VerbosityFlag,  // logger flag
	flags.KeystorePathFlag,
	flags.BlockstorePathFlag, // seems to be used only in tests
	flags.FreshStartFlag,     // start blocks from scratch. Used on chain initialization
	flags.LatestBlockFlag,    // latest block to start listen from. Used on chain initialization
	flags.MetricsFlag,
	flags.MetricsPort,
	flags.LevelDBPath,
}

//
var generateFlags = []cli.Flag{
	flags.PasswordFlag,
	flags.Secp256k1Flag,
}

//
var devFlags = []cli.Flag{
	flags.TestKeyFlag,
}

var importFlags = []cli.Flag{
	flags.EthereumImportFlag,
	flags.PrivateKeyFlag,
	flags.Secp256k1Flag,
	flags.PasswordFlag,
}

var accountCommand = cli.Command{
	Name:  "accounts",
	Usage: "manage bridge keystore",
	Description: "The accounts command is used to manage the bridge keystore. \n" +
		"\tTo generate a new account (key type generated is determined on the flag passed in): chainbridge accounts generate\n" +
		"\tTo import a keystore file: chainbridge accounts import path/to/file\n" +
		"\tTo import a geth keystore file: chainbridge accounts import --ethereum path/to/file\n" +
		"\tTo import a private key file: chainbridge accounts import --privateKey private_key\n" +
		"\tTo list keys: chainbridge accounts list",
	Subcommands: []*cli.Command{
		{
			Action: wrapHandler(handleGenerateCmd),
			Name:   "generate",
			Usage:  "generate bridge keystore, key type determined by flag",
			Flags:  generateFlags,
			Description: "The generate subcommand is used to generate the bridge keystore.\n" +
				"\tIf no options are specified, a secp256k1 key will be made.",
		},
		{
			Action: wrapHandler(handleImportCmd),
			Name:   "import",
			Usage:  "import bridge keystore",
			Flags:  importFlags,
			Description: "The import subcommand is used to import a keystore for the bridge.\n" +
				"\tA path to the keystore must be provided\n" +
				"\tUse --ethereum to import an ethereum keystore from external sources such as geth\n" +
				"\tUse --privateKey to create a keystore from a provided private key.",
		},
		{
			Action:      wrapHandler(handleListCmd),
			Name:        "list",
			Usage:       "list bridge keystore",
			Description: "The list subcommand is used to list all of the bridge keystores.\n",
		},
	},
}

// TODO: organize to subcommands under test command
var validatorsSyncerCommands = cli.Command{
	Name:   "syncer",
	Action: testutils.Sync,
}

var deployerTestCommands = cli.Command{
	Name:   "deploy",
	Action: e2e.Deploy,
}

var cliCmd = cli.Command{
	Name:        "cli",
	Description: "This CLI supports on-chain interactions with components of ChainBridge",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "url",
			Value: "http://localhost:8545",
			Usage: "RPC url of blockchain node",
		},
		&cli.Int64Flag{
			Name:  "gasLimit",
			Value: 6721975,
			Usage: "gasLimit used in transactions",
		},
		&cli.Int64Flag{
			Name:  "gasPrice",
			Value: 20000000000,
			Usage: "gasPrice used for transactions",
		},
		&cli.Int64Flag{
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
	Subcommands: []*cli.Command{
		{
			Name:        "deploy",
			Description: "This command can be used to deploy all or some of the contracts required for bridging. Selection of contracts can be made by either specifying --all or a subset of flags",
			Action:      cbcli.Deploy,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "bridge",
					Value: "bridge",
					Usage: "deploy bridge",
				},
				&cli.StringFlag{
					Name:  "erc20Handler",
					Value: "erc20Handler",
					Usage: "deploy erc20Handler",
				},
				&cli.StringFlag{
					Name:  "erc721Handler",
					Value: "erc721Handler",
					Usage: "deploy erc721Handler",
				},
				&cli.StringFlag{
					Name:  "genericHandler",
					Value: "genericHandler",
					Usage: "deploy genericHandler",
				},
				&cli.StringFlag{
					Name:  "erc20",
					Value: "erc20",
					Usage: "deploy erc20",
				},
				&cli.StringFlag{
					Name:  "erc721",
					Value: "erc721",
					Usage: "deploy erc721",
				},
				&cli.StringFlag{
					Name:  "all",
					Value: "all",
					Usage: "deploy all contracts",
				},
				&cli.Int64Flag{
					Name:  "relayerThreshold",
					Value: 1,
					Usage: "deploy all contracts",
				},
				&cli.Uint64Flag{
					Name:  "chainId",
					Value: 1,
					Usage: "deploy all contracts",
				},
				&cli.StringSliceFlag{
					Name:  "relayers",
					Value: cli.NewStringSlice(),
					Usage: "deploy all contracts",
				},
				&cli.Int64Flag{
					Name:  "fee",
					Value: 0,
					Usage: "deploy all contracts",
				},
				&cli.StringFlag{
					Name:  "bridgeAddress",
					Value: "",
					Usage: "deploy all contracts",
				},
			},
		},
	},
}

// init initializes CLI
func init() {
	app.Action = cmd.Run
	app.Copyright = "Copyright 2019 ChainSafe Systems Authors"
	app.Name = "chainbridge-celo"
	app.Usage = "ChainBridge-celo"
	app.Authors = []*cli.Author{{Name: "ChainSafe Systems 2020"}}
	app.Version = "0.0.1"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		&accountCommand,
		&validatorsSyncerCommands,
		&deployerTestCommands,
		&cliCmd,
	}

	app.Flags = append(app.Flags, cliFlags...)
	app.Flags = append(app.Flags, devFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Error().Err(err).Msg("Start failed")
		os.Exit(1)
	}
}
