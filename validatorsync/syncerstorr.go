package validatorsync

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	latestKnowBlockKey      = "latestKnownBlock"
	latestKnowValidatorsKey = "latestKnownValidators"
)

func NewSyncerStorr(db *leveldb.DB) *SyncerStorr {
	return &SyncerStorr{db: db}
}

type SyncerStorr struct {
	db *leveldb.DB
}

func (db *SyncerStorr) setLatestKnownBlock(block *big.Int, chainID uint8) error {
	key := new(bytes.Buffer)
	err := binary.Write(key, binary.BigEndian, chainID)
	if err != nil {
		return err
	}
	key.WriteString(latestKnowBlockKey)
	err = db.db.Put(key.Bytes(), block.Bytes(), nil)
	if err != nil {
		return err
	}
	return nil
}
func (db *SyncerStorr) setLatestKnownValidators(validators []*istanbul.ValidatorData, chainID uint8) error {
	b := &bytes.Buffer{}
	enc := gob.NewEncoder(b)
	err := enc.Encode(validators)
	if err != nil {
		return err
	}
	key := new(bytes.Buffer)
	err = binary.Write(key, binary.BigEndian, chainID)
	if err != nil {
		return err
	}
	key.WriteString(latestKnowValidatorsKey)
	err = db.db.Put(key.Bytes(), b.Bytes(), nil)
	if err != nil {
		return err
	}
	return nil
}

func (db *SyncerStorr) GetLatestKnownBlock(chainID uint8) (*big.Int, error) {
	key := new(bytes.Buffer)
	err := binary.Write(key, binary.BigEndian, chainID)
	if err != nil {
		return nil, err
	}
	key.WriteString(latestKnowBlockKey)
	data, err := db.db.Get(key.Bytes(), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return big.NewInt(0), nil
		}
		return nil, err
	}
	v := big.NewInt(0)
	v.SetBytes(data)
	return v, nil
}

func (db *SyncerStorr) GetLatestKnownValidators(chainID uint8) ([]*istanbul.ValidatorData, error) {
	key := new(bytes.Buffer)
	err := binary.Write(key, binary.BigEndian, chainID)
	if err != nil {
		return nil, err
	}
	key.WriteString(latestKnowValidatorsKey)
	res, err := db.db.Get(key.Bytes(), nil)
	if err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
			return make([]*istanbul.ValidatorData, 0), nil
		}
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
func (db *SyncerStorr) SetValidatorsForBlock(block *big.Int, validators []*istanbul.ValidatorData, chainID uint8) error {
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
	key := new(bytes.Buffer)
	err = binary.Write(key, binary.BigEndian, chainID)
	if err != nil {
		return err
	}
	key.Write(block.Bytes())
	err = tx.Put(key.Bytes(), byteValidators.Bytes(), nil)
	if err != nil {
		tx.Discard()
		return err
	}
	err = db.setLatestKnownBlockWithTransaction(block, chainID, tx)
	if err != nil {
		tx.Discard()
		return err
	}
	err = db.setLatestKnownValidatorsWithTransaction(validators, chainID, tx)
	if err != nil {
		tx.Discard()
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Discard()
		return err
	}
	log.Info().Int64("block", block.Int64()).Msgf("New validators set for block")
	return nil
}

func (db *SyncerStorr) GetValidatorsForBLock(block *big.Int, chainID uint8) ([]*istanbul.ValidatorData, error) {
	key := new(bytes.Buffer)
	err := binary.Write(key, binary.BigEndian, chainID)
	if err != nil {
		return nil, err
	}
	key.Write(block.Bytes())
	res, err := db.db.Get(key.Bytes(), nil)
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

func (db *SyncerStorr) setLatestKnownBlockWithTransaction(block *big.Int, chainID uint8, transaction *leveldb.Transaction) error {
	key := new(bytes.Buffer)
	err := binary.Write(key, binary.BigEndian, chainID)
	if err != nil {
		return err
	}
	key.WriteString(latestKnowBlockKey)
	err = transaction.Put(key.Bytes(), block.Bytes(), nil)
	if err != nil {
		return err
	}
	return nil
}

func (db *SyncerStorr) setLatestKnownValidatorsWithTransaction(validators []*istanbul.ValidatorData, chainID uint8, transaction *leveldb.Transaction) error {
	b := &bytes.Buffer{}
	enc := gob.NewEncoder(b)
	err := enc.Encode(validators)
	if err != nil {
		return err
	}
	key := new(bytes.Buffer)
	err = binary.Write(key, binary.BigEndian, chainID)
	if err != nil {
		return err
	}
	key.WriteString(latestKnowValidatorsKey)
	err = transaction.Put(key.Bytes(), b.Bytes(), nil)
	if err != nil {
		return err
	}
	return nil
}

// Closes connection to underlying DB backend
func (db *SyncerStorr) Close() error {
	if err := db.db.Close(); err != nil {
		return err
	}
	return nil
}