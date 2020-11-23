package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/ChainSafe/chainbridge-celo/cmd"

)

var app = cli.NewApp()

var cliFlags = []cli.Flag{}

var devFlags = []cli.Flag{}

var accountCommand = cli.Command{}

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
		//log.Error(err.Error())
		os.Exit(1)
	}
}