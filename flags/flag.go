// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package flags

import (
	"github.com/rs/zerolog"

	"github.com/urfave/cli/v2"
)

// Env vars
const DefaultKeystorePath = "./keys"

var (
	HealthBlockTimeout = "BLOCK_TIMEOUT"
)

var (
	ConfigFileFlag = &cli.StringFlag{
		Name:  "config",
		Usage: "JSON configuration file",
	}

	VerbosityFlag = &cli.StringFlag{
		Name:  "verbosity",
		Usage: "Supports levels crit (silent) to trce (trace)",
		Value: zerolog.InfoLevel.String(),
	}

	KeystorePathFlag = &cli.StringFlag{
		Name:  "keystore",
		Usage: "Path to keystore directory",
		Value: DefaultKeystorePath,
	}

	BlockstorePathFlag = &cli.StringFlag{
		Name:  "blockstore",
		Usage: "Specify path for blockstore",
		Value: "", // Empty will use home dir
	}

	FreshStartFlag = &cli.BoolFlag{
		Name:  "fresh",
		Usage: "Disables loading from blockstore at start. Opts will still be used if specified.",
	}

	LatestBlockFlag = &cli.BoolFlag{
		Name:  "latest",
		Usage: "Overrides blockstore and start block, starts from latest block",
	}

	LevelDBPath = &cli.StringFlag{
		Name:     "leveldb",
		Usage:    "sets path to leveldb database",
		Required: true,
	}
)

// Metrics flags
var (
	MetricsFlag = &cli.BoolFlag{
		Name:  "metrics",
		Usage: "Enables metric server",
	}

	MetricsPort = &cli.IntFlag{
		Name:  "metricsPort",
		Usage: "Port to serve metrics on",
		Value: 8001,
	}
)

// Generate subcommand flags
var (
	PasswordFlag = &cli.StringFlag{
		Name:  "password",
		Usage: "Password used to encrypt the keystore. Used with --generate, --import, or --unlock",
	}
	Secp256k1Flag = &cli.BoolFlag{
		Name:  "secp256k1",
		Usage: "Specify account/key type as secp256k1.",
	}
)

var (
	EthereumImportFlag = &cli.BoolFlag{
		Name:  "ethereum",
		Usage: "Import an existing ethereum keystore, such as from geth.",
	}
	PrivateKeyFlag = &cli.StringFlag{
		Name:  "privateKey",
		Usage: "Import a hex representation of a private key into a keystore.",
	}
)

// Test Setting Flags
var (
	TestKeyFlag = &cli.StringFlag{
		Name:  "testkey",
		Usage: "Applies a predetermined test keystore to the chains.",
	}
)
