// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package txtrie

import (
	"fmt"
	"github.com/status-im/keycard-go/hexutils"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

var (
	emptyHash = common.HexToHash("")
)

func computeEthReferenceTrieHash(transactions types.Transactions) (common.Hash, error) {
	newTrie, err := trie.New(emptyRoot, trie.NewDatabase(nil))
	if err != nil {
		return emptyHash, err
	}

	for i, tx := range transactions {

		key, err := rlp.EncodeToBytes(uint(i))
		if err != nil {
			return emptyHash, err
		}

		value, err := rlp.EncodeToBytes(tx)
		if err != nil {
			return emptyHash, err
		}

		err = newTrie.TryUpdate(key, value)
		if err != nil {
			return emptyHash, err
		}
	}

	return newTrie.Hash(), nil
}

func TestAddEmptyTrie(t *testing.T) {
	emptyTransactions := make([]*types.Transaction, 0)
	_, err := CreateNewTrie(emptyRoot, types.Transactions(emptyTransactions))
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddSingleTrieUpdate(t *testing.T) {
	vals := GetTransactions1()
	root, err := computeEthReferenceTrieHash(vals)
	if err != nil {
		t.Fatal(err)
	}

	trie, err := CreateNewTrie(root, types.Transactions(vals))
	if err != nil {
		t.Fatal(err)
	}
	keyRlp, err := rlp.EncodeToBytes(0)
	proof, key, err := RetrieveProof(trie, keyRlp)

	if proof == nil {
		t.Fatal("proof is nil")
	}

	if hexutils.BytesToHex(key) != "0001" {
		t.Fatal(fmt.Sprintf("wrong RLP key is %s", hexutils.BytesToHex(key)))
	}

}
