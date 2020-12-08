package blockdb

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-utils/blockstore"
	"github.com/ChainSafe/chainbridge-utils/msg"
)

func NewBlockStoreDB(relayerAddress string, blockstorePath string, chainID msg.ChainId, freshStart bool, startBlock *big.Int) (*blockstore.Blockstore, error) {
	bs, err := blockstore.NewBlockstore(blockstorePath, chainID, relayerAddress)
	if err != nil {
		return nil, err
	}
	if !freshStart {
		latestBlock, err := bs.TryLoadLatestBlock()
		if err != nil {
			return nil, err
		}
		if latestBlock.Cmp(startBlock) == 1 {
			*startBlock = *latestBlock
		}
	}
	return bs, nil
}
