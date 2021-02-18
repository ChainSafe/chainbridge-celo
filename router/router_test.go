// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package router

import (
	"reflect"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-celo/utils"
)

type mockWriter struct {
	msgs []*utils.Message
}

func (w *mockWriter) Start() error { return nil }
func (w *mockWriter) Stop() error  { return nil }

func (w *mockWriter) ResolveMessage(msg *utils.Message) bool {
	w.msgs = append(w.msgs, msg)
	return true
}

func TestRouter(t *testing.T) {
	router := NewRouter()

	ethW := &mockWriter{msgs: make([]*utils.Message, 0)}
	router.Register(utils.ChainId(0), ethW)

	ctfgW := &mockWriter{msgs: make([]*utils.Message, 0)}
	router.Register(utils.ChainId(1), ctfgW)

	msgEthToCtfg := &utils.Message{
		Source:      utils.ChainId(0),
		Destination: utils.ChainId(1),
	}

	msgCtfgToEth := &utils.Message{
		Source:      utils.ChainId(1),
		Destination: utils.ChainId(0),
	}

	err := router.Send(msgCtfgToEth)
	if err != nil {
		t.Fatal(err)
	}
	err = router.Send(msgEthToCtfg)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)

	if !reflect.DeepEqual(*ethW.msgs[0], *msgCtfgToEth) {
		t.Error("Unexpected message")
	}

	if !reflect.DeepEqual(*ctfgW.msgs[0], *msgEthToCtfg) {
		t.Error("Unexpected message")
	}
}
