package validator_syncer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/consensus/istanbul"
)

type validatorsDB struct {
	epochSize         *big.Int
	db                map[*big.Int][]*istanbul.ValidatorData
	latestServedBlock *big.Int
	latestValidators  []*istanbul.ValidatorData
}

func NewValidatorsDB() *validatorsDB {
	return &validatorsDB{epochSize: big.NewInt(12), db: make(map[*big.Int][]*istanbul.ValidatorData)}
}

func (db *validatorsDB) SetValidatorsForBlock(blockNumber *big.Int, vals []*istanbul.ValidatorData) {
	db.db[blockNumber] = vals
	db.latestServedBlock = blockNumber
	db.latestValidators = vals
}

func (db *validatorsDB) getValidatorsForBlock(blockNumber *big.Int) []*istanbul.ValidatorData {
	return db.db[blockNumber]
}

func (db *validatorsDB) GetLatestBlock() *big.Int {
	if db.latestServedBlock == nil {
		return big.NewInt(0)
	} else {
		return db.latestServedBlock
	}
}

func (db *validatorsDB) GetLatestValidators() []*istanbul.ValidatorData {
	if db.latestValidators != nil {
		return db.latestValidators
	}
	return make([]*istanbul.ValidatorData, 0)
}
