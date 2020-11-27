// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package ethtest

import (
	"context"
	"fmt"
	"math/big"

	utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

func WatchEvent(client *utils.Client, bridge common.Address, subStr utils.EventSig) {
	fmt.Printf("Watching for event: %s\n", subStr)
	query := eth.FilterQuery{
		FromBlock: big.NewInt(0),
		Addresses: []common.Address{bridge},
		Topics: [][]common.Hash{
			{subStr.GetTopic()},
		},
	}

	ch := make(chan ethtypes.Log)
	sub, err := client.Client.SubscribeFilterLogs(context.Background(), query, ch)
	if err != nil {
		log.Error().Err(err).Str("event", subStr.GetTopic().Hex()).Msg("Failed to subscribe to event")
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case evt := <-ch:
			fmt.Printf("%s (block: %d): %#v\n", subStr, evt.BlockNumber, evt.Topics)

		case err := <-sub.Err():
			if err != nil {
				log.Error().Err(err).Str("event", subStr.GetTopic().Hex()).Msg("Subscription error")
				return
			}
		}
	}
}
