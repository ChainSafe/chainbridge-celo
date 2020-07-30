package syncer

import (
	"context"
	"math/big"

	"github.com/ChainSafe/chainbridge-celo/connection"
)

type ValidatorSyncer struct {
	conn *connection.Connection
}

func (v *ValidatorSyncer) Sync(num uint64) error {
	_, err := v.conn.Client().BlockByNumber(context.Background(), new(big.Int).SetUint64(num))
	if err != nil {
		return err
	}

	return nil
}

func (v *ValidatorSyncer) start() error {
	err := v.conn.Connect()
	if err != nil {
		return err
	}
	return nil
}

func (v *ValidatorSyncer) close() {
	v.conn.Close()
}
