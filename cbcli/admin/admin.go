package admin

import "github.com/urfave/cli/v2"

var AdminCLICMDS = &cli.Command{
	Name: "admin",
	Subcommands: []*cli.Command{
		isRelayerCMD,
		addRelayerCMD,
		removeRelayerCMD,
		setTresholdCMD,
		pauseCMD,
		unpauseCMD,
		setFeeCMD,
		withdrawCMD,
		addAdminCMD,
		removeAdminCMD,
	},
}
