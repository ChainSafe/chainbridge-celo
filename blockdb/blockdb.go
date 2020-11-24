package blockdb

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-utils/blockstore"
	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/ChainSafe/chainbridge-utils/msg"
)

func NewBlockstoreDB(from string, keystorePath string, insecure bool, bloksorePath string, chainID msg.ChainId, freshStart bool, startblock *big.Int) (*blockstore.Blockstore, error) {

	kpI, err := keystore.KeypairFromAddress(from, keystore.EthChain, keystorePath, insecure)
	if err != nil {
		return nil, err
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	bs, err := setupBlockstore(bloksorePath, chainID, kp, freshStart, startblock)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// checkBlockstore queries the blockstore for the latest known block. If the latest block is
// greater than cfg.startBlock, then cfg.startBlock is replaced with the latest known block.

func setupBlockstore(bloksorePath string, chainID msg.ChainId, kp *secp256k1.Keypair, freshStart bool, startblock *big.Int) (*blockstore.Blockstore, error) {
	bs, err := blockstore.NewBlockstore(bloksorePath, chainID, kp.Address())
	if err != nil {
		return nil, err
	}

	if !freshStart {
		latestBlock, err := bs.TryLoadLatestBlock()
		if err != nil {
			return nil, err
		}

		if latestBlock.Cmp(startblock) == 1 {
			startblock = latestBlock
		}
	}

	return bs, nil
}
