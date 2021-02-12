// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package txtrie

import (
	"bytes"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/rlp"
	ethtrie "github.com/ethereum/go-ethereum/trie"
	"github.com/rs/zerolog/log"
	"github.com/status-im/keycard-go/hexutils"
)

var (
	// from https://github.com/ethereum/go-ethereum/blob/bcb308745010675671991522ad2a9e811938d7fb/trie/trie.go#L32
	emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)

func CreateNewTrie(root common.Hash, transactions types.Transactions) (*ethtrie.Trie, error) {
	if transactions == nil {
		return nil, errors.New("transactions cannot be nil")
	}
	db := memorydb.New()
	trie, err := ethtrie.New(emptyRoot, ethtrie.NewDatabase(db))
	if err != nil {
		return nil, err
	}
	for i, tx := range transactions {
		key, err := rlp.EncodeToBytes(uint(i))
		if err != nil {
			return nil, err
		}
		value, err := rlp.EncodeToBytes(tx)
		if err != nil {
			return nil, err
		}
		trie.Update(key, value)
	}
	if trie.Hash().Hex() != root.Hex() {
		return nil, errors.New("transaction roots don't match")
	}
	return trie, nil
}

func RetrieveProof(trie *ethtrie.Trie, key []byte) ([]byte, []byte, error) {
	//RetrieveProof(trie, key)
	iterator := trie.NodeIterator(key)
	it2 := ethtrie.NewIterator(iterator)

	proof := make([][][]byte, 1)
	for it2.Next() {
		log.Debug().Msgf("KEY %s", hexutils.BytesToHex(it2.Key))
		value := it2.Prove()
		log.Debug().Msgf("VALUE %s  %s", hexutils.BytesToHex(value[0]), hexutils.BytesToHex(value[0]))
		proof[0] = value
		proof123 := make([][][]byte, 0)
		for _, v := range value {
			n := make([][]byte, 0, 17)
			err := rlp.DecodeBytes(v, &n)
			if err != nil {
				panic(err)
			}
			proof123 = append(proof123, n)
		}
		buf := &bytes.Buffer{}
		rlp.Encode(buf, proof123)
		log.Debug().Msgf("ITERATOR LEAF PATH %s", hexutils.BytesToHex(keybytesToHex(iterator.LeafKey())))
		key := keybytesToHex(iterator.LeafKey())
		key = key[:len(key)-1]
		return buf.Bytes(), key, nil
	}
	//log.Debug().Msgf("IS THE LEAF %v", iterator.Leaf())
	//for iterator.Next(true) {
	//	if iterator.Leaf() {
	//		log.Debug().Msgf("ITERATOR LEAF PATH %s", hexutils.BytesToHex(iterator.LeafKey()))
	//		log.Debug().Msgf("ITERATOR LEAF PROOF %s", hexutils.BytesToHex(iterator.LeafProof()[0]))
	//	}
	//log.Debug().Msgf("IS THE LEAF %v", iterator.Leaf())
	//}
	//log.Debug().Msgf("IS THE LEAF %v", iterator.Leaf())
	//log.Debug().Msgf("%s", hexutils.BytesToHex(iterator.LeafProof()[0]))
	return iterator.Path(), nil, nil
}

func keybytesToHex(str []byte) []byte {
	l := len(str)*2 + 1
	var nibbles = make([]byte, l)
	for i, b := range str {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	nibbles[l-1] = 16
	return nibbles
}

//func RetrieveNewProof(trie *ethtrie.Trie, root common.Hash, key []byte) ([]byte, error) {
//	proofDB, err := RetrieveProof(trie, key)
//	if err != nil {
//		return nil, err
//	}
//	iterator := proofDB.NewIterator()
//
//	iterator.Next()
//
//	buf := &bytes.Buffer{}
//	byteByteArr := make([][][]byte, 1)
//
//	n := make([][]byte, 2)
//	err = rlp.DecodeBytes(iterator.Value(), &n)
//	if err != nil {
//		return nil, err
//	}
//	iterator.
//
//	byteByteArr[0] = n
//
//	err = rlp.Encode(buf, byteByteArr)
//	if err != nil {
//		return nil, err
//	}
//	return buf.Bytes(), nil
//}

// VerifyProof verifies merkle proof on path key against the provided root
func VerifyProof(root common.Hash, key []byte, proof ethdb.KeyValueStore) (bool, error) {
	exists, _, err := ethtrie.VerifyProof(root, key, proof)

	if err != nil {
		return false, err
	}

	return exists != nil, nil
}
