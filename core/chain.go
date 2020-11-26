// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package core

import (
	metrics "github.com/ChainSafe/chainbridge-utils/metrics/types"
	"github.com/ChainSafe/chainbridge-utils/msg"
)

type Chain interface {
	Start() error // Start chain
	SetRouter(*Router)
	Id() msg.ChainId
	Name() string
	LatestBlock() *metrics.LatestBlock
	Stop()
}
