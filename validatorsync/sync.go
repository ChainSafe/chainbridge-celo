//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package validatorsync

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

const (
	timeToWaitUntilNextBlockAppear = 5
)

type HeaderByNumberGetter interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

func StoreBlockValidators(stopChn <-chan struct{}, errChn chan error, c HeaderByNumberGetter, db *ValidatorsStore, chainID uint8, epochSize uint64) {
	// If DB is empty will return 0 (first epoch by itself)
	var prevValidators []*istanbul.ValidatorData
	var currentValidators []*istanbul.ValidatorData
	block, err := db.GetLatestKnownEpochLastBlock(chainID)
	if err != nil {
		errChn <- fmt.Errorf("error on get latest known block from db: %w", err)
		return
	}
	if block.Cmp(big.NewInt(0)) == 0 {
		// If block is zero, initial validators should be empty array
		log.Info().Msg("Syncing validators from zero block")
		prevValidators = make([]*istanbul.ValidatorData, 0)
	} else {
		prevValidators, err = db.GetValidatorsForBlock(block, chainID)
		if err != nil {
			errChn <- fmt.Errorf("error on get latest known validators from db: %w", err)
			return
		}
		// We already know validators for that block so moving to next one
		block.Add(block, big.NewInt(0).SetUint64(epochSize))
		log.Info().Msg(fmt.Sprintf("Syncing validators from %s block", block.String()))
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
					time.Sleep(timeToWaitUntilNextBlockAppear * time.Second)
					continue
				}
				errChn <- fmt.Errorf("gettings header by number err: %w", err)
				return
			}
			extra, err := types.ExtractIstanbulExtra(header)
			if err != nil {
				errChn <- fmt.Errorf("error on extracting istanbul extra: %w", err)
				return
			}
			b := bytes.NewBuffer(extra.RemovedValidators.Bytes())

			currentValidators = append(make([]*istanbul.ValidatorData, 0), prevValidators...)

			if len(extra.AddedValidators) != 0 || b.Len() > 0 {
				log.Debug().Str("block", block.String()).Msg("New validators data")
				currentValidators, err = applyValidatorsDiff(extra, prevValidators)
				if err != nil {
					errChn <- fmt.Errorf("error applying validators diff: %w", err)
					return
				}
			}
			// Zero block is first and last block of first epoch. So zero block should be set with its own diff validators
			if block.Cmp(big.NewInt(0)) == 0 {
				err = db.SetValidatorsForBlock(block, currentValidators, chainID)
				if err != nil {
					errChn <- fmt.Errorf("error on set validators to db: %w", err)
					return
				}
			} else {
				// If block is not zero, then it is last block of epoch, so
				err = db.SetValidatorsForBlock(block, currentValidators, chainID)
				if err != nil {
					errChn <- fmt.Errorf("error on set validators to db: %w", err)
					return
				}
			}
			// Current validators for next epoch, will be set for next last epoch block and applied with its diff
			prevValidators = append(make([]*istanbul.ValidatorData, 0), currentValidators...)
			block.Add(block, big.NewInt(0).SetUint64(epochSize))
		}
	}
}
