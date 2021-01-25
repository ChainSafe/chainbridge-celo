// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package listener

import (
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func GetTransactions() types.Transactions {
	return getTransactions(getStringData())
}

func getTransactions(data string) types.Transactions {
	transactions := types.Transactions{}
	err := json.Unmarshal([]byte(data), &transactions)

	if err != nil {
		log.Fatal(err)
	}

	return transactions
}

func getStringData() string {

	data := string(`[{"nonce": "0xf12a",
			"gasPrice": "0x2540be400",
			"gas": "0x28d0f",
			"to": "0x508d2c2c3584e6e1e2adfa0b7c0823846914bfff",
			"value": "0x0",
			"input": "0x402af46700000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000005fd9033000000000000000000000000000000000000000000000000000000000000000044c494e41000000000000000000000000000000000000000000000000000000006c455448000000000000000000000000000000000000000000000000000000006c425443000000000000000000000000000000000000000000000000000000006c48423130000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000275160bbbab00000000000000000000000000000000000000000000000001fd86bcefbab24c00000000000000000000000000000000000000000000000041c09ed05f31d7d00000000000000000000000000000000000000000000000000000d8b72d434c80000",
			"v": "0x29",
			"r": "0x379b69f6186bebd2044d261d36b588fab4bea2485a1a33bf279e8b5471ef0a14",
			"s": "0x475d8c750cc0c88455d04144f2c84407a212d9657e9ecdbf515e5f555485ca86",
			"hash": "0x6501b4720729abd65c3a0c2207e5d0a13cd3df2c0805b2b008f6b65250518c84"
			},
			{
			"nonce": "0xf12b",
			"gasPrice": "0x2540be400",
			"gas": "0x28d0f",
			"to": "0x508d2c2c3584e6e1e2adfa0b7c0823846914bfff",
			"value": "0x0",
			"input": "0x402af46700000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000005fd903a800000000000000000000000000000000000000000000000000000000000000044c494e41000000000000000000000000000000000000000000000000000000006c455448000000000000000000000000000000000000000000000000000000006c425443000000000000000000000000000000000000000000000000000000006c48423130000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000275160bbbab00000000000000000000000000000000000000000000000001fd63dfc7dd106700000000000000000000000000000000000000000000000041c09ed05f31d7d00000000000000000000000000000000000000000000000000000d8b72d434c80000",
			"v": "0x2a",
			"r": "0xd0b08a726661b71d8b0a24b7757a773c70966bccafae16c8e253d4e905b45de5",
			"s": "0x394225058e50029166f56f9ce7ebfdfd98fe6ef408520568a9248b43cfff1f1e",
			"hash": "0x84a14db17d181e735a8c088f3a16a4bb1eca41da4a2dd57724446c5a8a31f8e7"
			}]`)

	return data
}

func GetTxRoot(transactions types.Transactions) (common.Hash, error) {
	emptyHash := common.HexToHash("")
	emptyRoot := common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	newTrie, err := trie.New(emptyRoot, trie.NewDatabaseWithCache(nil, 0))
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