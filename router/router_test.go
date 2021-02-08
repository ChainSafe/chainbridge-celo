// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package router

import (
	"reflect"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-celo/pkg"
)

type mockWriter struct {
	msgs []*pkg.Message
}

func (w *mockWriter) Start() error { return nil }
func (w *mockWriter) Stop() error  { return nil }

func (w *mockWriter) ResolveMessage(msg *pkg.Message) bool {
	w.msgs = append(w.msgs, msg)
	return true
}

func TestRouter(t *testing.T) {
	router := NewRouter()

	ethW := &mockWriter{msgs: make([]*pkg.Message, 0)}
	router.Register(pkg.ChainId(0), ethW)

	ctfgW := &mockWriter{msgs: make([]*pkg.Message, 0)}
	router.Register(pkg.ChainId(1), ctfgW)

	msgEthToCtfg := &pkg.Message{
		Source:      pkg.ChainId(0),
		Destination: pkg.ChainId(1),
	}

	msgCtfgToEth := &pkg.Message{
		Source:      pkg.ChainId(1),
		Destination: pkg.ChainId(0),
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
