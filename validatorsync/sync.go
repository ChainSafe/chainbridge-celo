//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package validatorsync

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	timeToWaitUntilNextBLockAppear = 5
)

type HeaderByNumberGetter interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

func StoreBlockValidators(stopChn <-chan struct{}, errChn chan error, c HeaderByNumberGetter, db *ValidatorsStore, chainID uint8) {
	block, err := db.GetLatestKnownBlock(chainID)
	if err != nil {
		errChn <- fmt.Errorf("error on get latest known block from db: %w", err)
		return
	}
	for {
		select {
		case <-stopChn:
			return
		default:
			header, err := c.HeaderByNumber(context.Background(), block)
			if err != nil {
				if errors.Is(err, ethereum.NotFound) {
					// Block not yet mined, waiting
					time.Sleep(timeToWaitUntilNextBLockAppear * time.Second)
					continue
				}
				errChn <- fmt.Errorf("gettings header by number err: %w", err)
				return
			}
			actualValidators, err := db.GetLatestKnownValidators(chainID)
			if err != nil {
				errChn <- fmt.Errorf("error on get latest known validators from db: %w", err)
				return
			}
			extra, err := types.ExtractIstanbulExtra(header)
			if err != nil {
				errChn <- fmt.Errorf("error on extracting istanbul extra: %w", err)
				return
			}
			b := bytes.NewBuffer(extra.RemovedValidators.Bytes())
			if len(extra.AddedValidators) != 0 || b.Len() > 0 {
				actualValidators, err = applyValidatorsDiff(extra, actualValidators)
				if err != nil {
					errChn <- fmt.Errorf("error applying validators diff: %w", err)
					return
				}
			}
			err = db.SetValidatorsForBlock(block, actualValidators, chainID)
			if err != nil {
				errChn <- fmt.Errorf("error on set validators to db: %w", err)
				return
			}
			block.Add(block, big.NewInt(1))
		}
	}
}
