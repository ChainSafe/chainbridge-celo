package erc721

import "github.com/urfave/cli/v2"

var ERC721CLICMDS = &cli.Command{
	Name: "erc721",
	Subcommands: []*cli.Command{
		mintCMD,
		ownerCMD,
		addMinterCMD,
		approveCMD,
		depositCMD,
	},
}
