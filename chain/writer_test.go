// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package chain  

import (
	"testing"
	ethtest "github.com/ChainSafe/ChainBridge/shared/ethereum/testing"
	utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
)

func createWriters(t *testing.T, client *utils.Client, contracts *utils.DeployedContracts) (*writer, *writer, func(), func(), chan error, chan error) {
	latestBlock := ethtest.GetLatestBlock(t, client)
	errA := make(chan error)
	writerA, stopA := createTestWriter(t, createConfig("bob", latestBlock, contracts), errA)
	errB := make(chan error)
	writerB, stopB := createTestWriter(t, createConfig("charlie", latestBlock, contracts), errB)
	return writerA, writerB, stopA, stopB, errA, errB
}

func createTestWriter(t *testing.T, cfg *Config, errs chan<- error) (*writer, func()) {
 
	conn := newLocalConnection(t, cfg)
	stop := make(chan int)
	writer := NewWriter(conn, cfg, newTestLogger(cfg.name), stop, errs, nil)

	bridge, err := Bridge.NewBridge(cfg.bridgeContract, conn.Client())
	if err != nil {
		t.Fatal(err)
	}

	writer.setContract(bridge)

	err = writer.start()
	if err != nil {
		t.Fatal(err)
	}
	return writer, func() { close(stop) }

}

func routeMessageAndWait(t *testing.T, client *utils.Client, alice, bob *writer, m msg.Message, aliceErr, bobErr chan error) {
	// Watch for executed event
	query := eth.FilterQuery{
		FromBlock: big.NewInt(0),
		Addresses: []common.Address{alice.cfg.bridgeContract},
		Topics: [][]common.Hash{
			{utils.ProposalEvent.GetTopic()},
		},
	}

	ch := make(chan ethtypes.Log)
	sub, err := client.Client.SubscribeFilterLogs(context.Background(), query, ch)
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe()

	if err != nil {
		log15.Error("Failed to subscribe to finalization event", "err", err)
	}

	// Alice processes the message, then waits to execute
	if ok := alice.ResolveMessage(m); !ok {
		t.Fatal("Alice failed to resolve the message")
	}

	// Now Bob receives the same message and also waits to execute
	if ok := bob.ResolveMessage(m); !ok {
		t.Fatal("Bob failed to resolve the message")
	}

	for {
		select {
		case evt := <-ch:
			sourceId := evt.Topics[1].Big().Uint64()
			depositNonce := evt.Topics[2].Big().Uint64()
			status := uint8(evt.Topics[3].Big().Uint64())

			if m.Source == msg.ChainId(sourceId) &&
				uint64(m.DepositNonce) == depositNonce &&
				utils.IsExecuted(status) {
				return
			}

		case err = <-sub.Err():
			if err != nil {
				t.Fatal(err)
			}
		case err = <-aliceErr:
			t.Fatalf("Fatal error: %s", err)
		case err = <-bobErr:
			t.Fatalf("Fatal error: %s", err)
		case <-time.After(TestTimeout):
			t.Fatal("test timed out")
		}
	}
}

func TestWriter_start_stop(t *testing.T) {
	conn := newLocalConnection(t, aliceTestConfig)
	defer conn.Close()

	stop := make(chan int)
	writer := NewWriter(conn, aliceTestConfig, TestLogger, stop, nil, nil)

	err := writer.start()
	if err != nil {
		t.Fatal(err)
	}

	// Initiate shutdown
	close(stop)
}