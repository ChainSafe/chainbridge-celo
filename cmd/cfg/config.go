// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package cfg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ChainSafe/chainbridge-celo/flags"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

const DefaultConfigPath = "./config.json"

type Config struct {
	Chains []RawChainConfig `json:"chains"`
}

// RawChainConfig is parsed directly from the config file and should be using to construct the core.ChainConfig
type RawChainConfig struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Id       string            `json:"id"`       // ChainID
	Endpoint string            `json:"endpoint"` // url for rpc endpoint
	From     string            `json:"from"`     // address of key to use
	Opts     map[string]string `json:"opts"`
}

func NewConfig() *Config {
	return &Config{
		Chains: []RawChainConfig{},
	}
}

func (c *Config) ToJSON(file string) *os.File {
	var (
		newFile *os.File
		err     error
	)

	var raw []byte
	if raw, err = json.Marshal(*c); err != nil {
		log.Error().Err(fmt.Errorf("error marshalling json: %w", err))
		os.Exit(1)
	}

	newFile, err = os.Create(file)
	if err != nil {
		log.Error().Err(fmt.Errorf("error creating config file %w", err))
	}
	_, err = newFile.Write(raw)
	if err != nil {
		log.Error().Err(fmt.Errorf("error writing to config file %w", err))
	}

	if err := newFile.Close(); err != nil {
		log.Error().Err(fmt.Errorf("error closing file %w", err))
	}
	return newFile
}

func (c *Config) validate() error {
	for _, chain := range c.Chains {
		if chain.Type == "" {
			return fmt.Errorf("required field chain.Type empty for chain %s", chain.Id)
		}
		if chain.Endpoint == "" {
			return fmt.Errorf("required field chain.Endpoint empty for chain %s", chain.Id)
		}
		if chain.Name == "" {
			return fmt.Errorf("required field chain.Name empty for chain %s", chain.Id)
		}
		if chain.Id == "" {
			return fmt.Errorf("required field chain.Id empty for chain %s", chain.Id)
		}
		if chain.From == "" {
			return fmt.Errorf("required field chain.From empty for chain %s", chain.Id)
		}
	}
	return nil
}

func GetConfig(ctx *cli.Context) (*Config, error) {
	var fig Config
	path := DefaultConfigPath
	if file := ctx.String(flags.ConfigFileFlag.Name); file != "" {
		path = file
	}
	err := loadConfig(path, &fig)
	if err != nil {
		log.Error().Err(fmt.Errorf("err loading json file: %s", err))
		return &fig, err
	}
	log.Debug().Msgf("Loaded config path: %s", path)
	err = fig.validate()
	if err != nil {
		return nil, err
	}
	return &fig, nil
}

func loadConfig(file string, config *Config) error {
	ext := filepath.Ext(file)
	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	log.Debug().Msgf("Loading configuration path: %s", filepath.Clean(fp))

	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return err
	}

	if ext == ".json" {
		if err = json.NewDecoder(f).Decode(&config); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unrecognized extention: %s", ext)
	}

	return nil
}
