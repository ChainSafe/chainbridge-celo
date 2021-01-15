package listener

import (
	"github.com/ChainSafe/chainbridge-ethereum-trie/txtrie"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func getTrie(txtRootHash common.Hash, transactions types.Transactions) (*txtrie.TxTries, error) {

	tries := txtrie.NewTxTries()

	err := tries.CreateNewTrie(txtRootHash, transactions)

	if err != nil {
		return nil, err
	}

	return tries, nil
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
