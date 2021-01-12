package listener

import (
	"github.com/ChainSafe/chainbridge-ethereum-trie/txtrie"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	size   = 1
	dbPath = "./trie-database"
)

func getTrie(txtRootHash common.Hash, transactions types.Transactions) (*txtrie.TxTries, *leveldb.Database, error) {

	tries := txtrie.NewTxTries(size)

	db, err := getDb()

	if err != nil {
		return nil, nil, err
	}

	err = tries.AddNewTrie(txtRootHash, transactions, db)

	if err != nil {
		return nil, nil, err
	}

	return tries, db, nil
}

func getTrieProof(txtRootHash common.Hash, tries *txtrie.TxTries, key uint) ([]byte, error) {

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

func getDb() (*leveldb.Database, error) {

	diskdb, err := leveldb.New(dbPath, 256, 0, "")
	if err != nil {
		return nil, err
	}
	return diskdb, nil
}
