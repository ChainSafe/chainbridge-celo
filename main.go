package main

import (
	"os"

	"github.com/ChainSafe/chainbridge-celo/cmd"
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
	Description: "The accounts command is used to manage the bridge keystore.\n" +
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
