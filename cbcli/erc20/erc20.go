package erc20

import "github.com/urfave/cli/v2"

var ERC20CLICMDS = &cli.Command{
	Name: "admin",
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:  "decimals",
			Usage: "erc20Token decimals",
		},
	},
	Subcommands: []*cli.Command{
		mintCMD,
		addMinterCMD,
		approveCMD,
		depositCMD,
		balanceCMD,
		allowanceCMD,
	},
}
