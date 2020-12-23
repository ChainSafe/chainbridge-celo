package trie

import (
	"errors"

	"github.com/ChainSafe/chainbridge-celo/chain/config"
	"github.com/ChainSafe/chainbridge-ethereum-trie/txtrie"
	"github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	size   = 1
	dbPath = "./trie-database"
)

type CeloTrie struct {
}

func NewCeloTrie() *CeloTrie {

	return &CeloTrie{}
}

func (t *CeloTrie) GetProof(txtRootHash common.Hash, transactions types.Transactions, key uint) ([]byte, error) {

	tries := txtrie.NewTxTries(size)

	db, err := getDb()

	if err != nil {
		return nil, err
	}

	err = tries.AddNewTrie(txtRootHash, transactions, db)

	if err != nil {
		return nil, err
	}

	keyRlp, err := rlp.EncodeToBytes(key)

	if err != nil {
		return nil, err
	}

	proof, err := tries.RetrieveEncodedProof(txtRootHash, keyRlp)

	if err != nil {
		return nil, err
	}

	return proof, nil

}

func getKeypair(cfg *config.CeloChainConfig) (*secp256k1.Keypair, error) {

	kpI, err := keystore.KeypairFromAddress(cfg.From, keystore.EthChain, cfg.KeystorePath, cfg.Insecure)

	if err != nil {
		return nil, err
	}

	keypair, ok := kpI.(*secp256k1.Keypair)

	if !ok {
		return nil, errors.New("failed to convert kpI to *secp256k1.Keypair")
	}

	return keypair, nil
}

func getDb() (*leveldb.Database, error) {

	diskdb, err := leveldb.New(dbPath, 256, 0, "")
	if err != nil {
		return nil, err
	}
	return diskdb, nil
}
