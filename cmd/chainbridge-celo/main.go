// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
/*
Provides the command-line interface for the chainbridge-celo application.
For configuration and CLI commands see the README: https://github.com/ChainSafe/chainbridge-celo
*/
package main

import (
	"os"

	log "github.com/ChainSafe/log15"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

var cliFlags = []cli.Flag{
	//TODO
}

var devFlags = []cli.Flag{
	// TODO
}

var accountCommand = cli.Command{
	// TODO
}

// init initializes CLI
func init() {
	app.Action = run
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
		log.Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	// TODO - Implement run method
	return nil
}
