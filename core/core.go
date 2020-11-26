// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package core

import (
	"os"
	"os/signal"
	"syscall"

	defaultRouter "github.com/ChainSafe/chainbridge-celo/router"
	metrics "github.com/ChainSafe/chainbridge-utils/metrics/types"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/rs/zerolog/log"
)

type Chain interface {
	Start() error // Start chain
	SetRouter(*defaultRouter.Router)
	ID() msg.ChainId
	Name() string
	LatestBlock() *metrics.LatestBlock
	Stop()
}

type Core struct {
	Registry []Chain
	sysErr   <-chan error
}

func NewCore(sysErr <-chan error) *Core {
	return &Core{
		Registry: make([]Chain, 0),
		sysErr:   sysErr,
	}
}

// AddChain registers the chain in the Registry and calls Chain.SetRouter()
func (c *Core) AddChain(chain Chain) {
	c.Registry = append(c.Registry, chain)
}

// Start will call all registered chains' Start methods and block forever (or until signal is received)
func (c *Core) Start() {
	for _, chain := range c.Registry {
		err := chain.Start()
		if err != nil {
			log.Error().Interface("chain", chain.ID()).Err(err).Msg("failed to start chain")
			return
		}
		log.Info().Msgf("Started %s chain", chain.Name())
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)

	// Block here and wait for a signal
	select {
	case err := <-c.sysErr:
		log.Error().Err(err).Msg("FATAL ERROR. Shutting down.")
	case <-sigc:
		log.Warn().Msg("Interrupt received, shutting down now.")
	}

	// Signal chains to shutdown
	for _, chain := range c.Registry {
		chain.Stop()
	}
}

func (c *Core) Errors() <-chan error {
	return c.sysErr
}
