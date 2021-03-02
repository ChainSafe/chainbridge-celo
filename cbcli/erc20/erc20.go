package erc20

import "github.com/urfave/cli/v2"

var ERC20CLICMDS = &cli.Command{
	Name: "erc20",
	Subcommands: []*cli.Command{
		mintCMD,
		addMinterCMD,
		approveCMD,
		depositCMD,
		balanceCMD,
		allowanceCMD,
	},
}
