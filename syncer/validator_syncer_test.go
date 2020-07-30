package syncer

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-celo/connection"
	"github.com/ChainSafe/chainbridge-utils/keystore"
	"github.com/ChainSafe/log15"
)

var TestEndpoint = "ws://localhost:8545"
var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var GasLimit = big.NewInt(connection.DefaultGasLimit)
var GasPrice = big.NewInt(connection.DefaultGasPrice)

func createTestConnection(t *testing.T) *connection.Connection {
	conn := connection.NewConnection(TestEndpoint, false, AliceKp, log15.Root(), GasLimit, GasPrice)
	return conn
}

func TestValidatorSyncer_Sync(t *testing.T) {
	conn := createTestConnection(t)
	vsyncer := ValidatorSyncer{conn: conn}
	err := vsyncer.start()
	if err != nil {
		t.Fatal(err)
	}
	// TODO: validators which is the first skipped arg needs to be tested
	_, err = vsyncer.Sync(0)
	if err != nil {
		t.Fatal(err)
	}
	vsyncer.close()
}