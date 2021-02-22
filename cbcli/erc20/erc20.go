package erc20

import "github.com/urfave/cli/v2"

var ERC20CLICMDS = &cli.Command{
	Name: "admin",
	Subcommands: []*cli.Command{
		mintCMD,
	},
}
