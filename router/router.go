// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package router

import (
	"fmt"
	"sync"

	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/rs/zerolog/log"
)

// Writer consumes a message and makes the requried on-chain interactions.
type MessageResolver interface {
	ResolveMessage(message msg.Message) bool
}

// BaseRouter forwards messages from their source to their destination
type BaseRouter struct {
	registry map[msg.ChainId]MessageResolver
	lock     *sync.RWMutex
}

func NewRouter() *BaseRouter {
	return &BaseRouter{
		registry: make(map[msg.ChainId]MessageResolver),
		lock:     &sync.RWMutex{},
	}
}

// Send passes a message to the destination Writer if it exists
func (r *BaseRouter) Send(msg msg.Message) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	log.Trace().Interface("src", msg.Source).Interface("dest", msg.Destination).Interface("nonce", msg.DepositNonce).Interface("rId", msg.ResourceId.Hex()).Msg("Routing message")
	w := r.registry[msg.Destination]
	if w == nil {
		return fmt.Errorf("unknown destination chainId: %d", msg.Destination)
	}

	go w.ResolveMessage(msg)
	return nil
}

// Register registers a Writer with a ChainId which BaseRouter.Send can then use to propagate messages
func (r *BaseRouter) Register(id msg.ChainId, w MessageResolver) {
	r.lock.Lock()
	defer r.lock.Unlock()
	log.Debug().Interface("id", id).Msg("Registering new chain in router")
	r.registry[id] = w
}
