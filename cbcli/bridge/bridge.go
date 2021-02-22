package bridge

import "github.com/urfave/cli/v2"

var BridgeCLICMDS = &cli.Command{
	Name: "bridge",
	Subcommands: []*cli.Command{
		registerResourceCMD,
		registerGenericResourceCMD,
		setBurnCMD,
		cancelProposalCMD,
		queryProposalCMD,
		queryResourceCMD,
	},
}
