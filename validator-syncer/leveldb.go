package validator_syncer

import (
	"bytes"
	"encoding/gob"
	"math/big"

	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	latestKnowBlockKey      = "latestKnownBlock"
	latestKnowValidatorsKey = "latestKnownValidators"
)

func NewSyncerDB(path string) (*SyncerDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &SyncerDB{db: db}, nil
}

type SyncerDB struct {
	db *leveldb.DB
}

func (db *SyncerDB) setLatestKnownBlock(block *big.Int) error {
	err := db.db.Put([]byte(latestKnowBlockKey), block.Bytes(), nil)
	if err != nil {
		return err
	}
	return nil
}
func (db *SyncerDB) setLatestKnownValidators(validators []*istanbul.ValidatorData) error {
	b := &bytes.Buffer{}
	enc := gob.NewEncoder(b)
	err := enc.Encode(validators)
	if err != nil {
		return err
	}
	err = db.db.Put([]byte(latestKnowValidatorsKey), b.Bytes(), nil)
	if err != nil {
		return err
	}
	return nil
}

func (db *SyncerDB) GetLatestKnownBlock() (*big.Int, error) {
	data, err := db.db.Get([]byte(latestKnowBlockKey), nil)
	if err != nil {
		return nil, err
	}
	v := big.NewInt(0)
	v.SetBytes(data)
	return v, nil
}

func (db *SyncerDB) GetLatestKnownValidators() ([]*istanbul.ValidatorData, error) {
	res, err := db.db.Get([]byte(latestKnowValidatorsKey), nil)
	if err != nil {
		return nil, err
	}
	b := &bytes.Buffer{}
	b.Write(res)
	dec := gob.NewDecoder(b)
	dataArr := make([]*istanbul.ValidatorData, 0)
	err = dec.Decode(&dataArr)
	if err != nil {
		return nil, err
	}
	return dataArr, nil
}

// Atomically sets block and validators as related KV to underlying DB backend
func (db *SyncerDB) SetValidatorsForBlock(block *big.Int, validators []*istanbul.ValidatorData) error {
	byteValidators := &bytes.Buffer{}
	enc := gob.NewEncoder(byteValidators)
	err := enc.Encode(validators)
	if err != nil {
		return err
	}
	tx, err := db.db.OpenTransaction()
	if err != nil {
		return err
	}
	err = tx.Put(block.Bytes(), byteValidators.Bytes(), nil)
	if err != nil {
		tx.Discard()
		return err
	}
	err = db.setLatestKnownBlockWithTransaction(block, tx)
	if err != nil {
		tx.Discard()
		return err
	}
	err = db.setLatestKnownValidatorsWithTransaction(validators, tx)
	if err != nil {
		tx.Discard()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Discard()
		return err
	}
	return nil
}

func (db *SyncerDB) GetValidatorsForBLock(block *big.Int) ([]*istanbul.ValidatorData, error) {
	res, err := db.db.Get(block.Bytes(), nil)
	if err != nil {
		return nil, err
	}
	b := &bytes.Buffer{}
	b.Write(res)
	dec := gob.NewDecoder(b)
	dataArr := make([]*istanbul.ValidatorData, 0)
	err = dec.Decode(&dataArr)
	if err != nil {
		return nil, err
	}
	return dataArr, nil
}

func (db *SyncerDB) setLatestKnownBlockWithTransaction(block *big.Int, transaction *leveldb.Transaction) error {
	err := transaction.Put([]byte(latestKnowBlockKey), block.Bytes(), nil)
	if err != nil {
		return err
	}
	return nil
}

func (db *SyncerDB) setLatestKnownValidatorsWithTransaction(validators []*istanbul.ValidatorData, transaction *leveldb.Transaction) error {
	b := &bytes.Buffer{}
	enc := gob.NewEncoder(b)
	err := enc.Encode(validators)
	if err != nil {
		return err
	}
	err = transaction.Put([]byte(latestKnowValidatorsKey), b.Bytes(), nil)
	if err != nil {
		return err
	}
	return nil
}

// Closes connection to underlying DB backend
func (db *SyncerDB) Close() {
	db.db.Close()
}
