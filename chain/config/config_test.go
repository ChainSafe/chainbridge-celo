// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package config

import (
	"flag"
	"strconv"
	"testing"

	"github.com/ChainSafe/chainbridge-celo/cmd/cfg"
	"github.com/ChainSafe/chainbridge-celo/msg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

func TestParseConfig(t *testing.T) {

	var chainID int = 3
	chainIDStr := strconv.Itoa(chainID)

	var maxGasPrice int = 21000
	maxGasPriceStr := strconv.Itoa(maxGasPrice)

	var gasLimit int = 31000
	gasLimitStr := strconv.Itoa(gasLimit)

	var startBlock int = 3000000
	startBlockStr := strconv.Itoa(startBlock)

	http := "true"

	_type := "test"
	_name := "test"
	_endpoint := "http://localhost:8080"
	fromAddress := "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa201"
	bridge := "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202"
	erc20Handler := "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa203"
	erc721Handler := "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa204"
	genericHandler := "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa205"

	_keystorePath := "./keystore"
	_config := "./config"
	_blockstore := "blockstore"
	_latest := true
	_fresh := true

	rCon := &cfg.RawChainConfig{
		Name:     _name,
		Type:     _type,
		Id:       chainIDStr,
		Endpoint: _endpoint,
		From:     fromAddress,
		Opts: map[string]string{
			"bridge":         bridge,
			"erc20Handler":   erc20Handler,
			"erc721Handler":  erc721Handler,
			"genericHandler": genericHandler,
			"maxGasPrice":    maxGasPriceStr,
			"gasLimit":       gasLimitStr,
			"http":           http,
			"startBlock":     startBlockStr,
		},
	}

	set := flag.NewFlagSet("test", 0)

	set.String("config", _config, "config")
	set.String("keystore", _keystorePath, "keystore")
	set.String("blockstore", _blockstore, "blockstore")
	set.Bool("fresh", _fresh, "fresh")
	set.Bool("latest", _latest, "fresh")

	ctx := cli.NewContext(nil, set, nil)

	config, err := ParseChainConfig(rCon, ctx)

	if err != nil {
		t.Error(err)
	}

	if config.ID != msg.ChainId(chainID) {
		t.Errorf("expected %v got %v ", chainID, config.ID)
	}

	if config.MaxGasPrice.Int64() != int64(maxGasPrice) {
		t.Errorf("expected %v got %v ", maxGasPrice, config.MaxGasPrice)
	}

	if config.GasLimit.Int64() != int64(gasLimit) {
		t.Errorf("expected %v got %v ", gasLimit, config.GasLimit)
	}

	if config.StartBlock.Int64() != int64(startBlock) {
		t.Errorf("expected %v got %v ", startBlock, config.StartBlock)
	}

	if config.BridgeContract != common.HexToAddress(bridge) {
		t.Errorf("expected %v got %v ", common.HexToAddress(bridge), config.BridgeContract)
	}

	if config.Erc20HandlerContract != common.HexToAddress(erc20Handler) {
		t.Errorf("expected %v got %v ", common.HexToAddress(erc20Handler), config.Erc20HandlerContract)
	}

	if config.Erc721HandlerContract != common.HexToAddress(erc721Handler) {
		t.Errorf("expected %v got %v ", common.HexToAddress(erc721Handler), config.Erc721HandlerContract)
	}

	if config.GenericHandlerContract != common.HexToAddress(genericHandler) {
		t.Errorf("expected %v got %v ", common.HexToAddress(genericHandler), config.GenericHandlerContract)
	}

	if config.From != fromAddress {
		t.Errorf("expected %v got %v ", common.HexToAddress(fromAddress), config.From)
	}

	if config.From != fromAddress {
		t.Errorf("expected %v got %v ", common.HexToAddress(fromAddress), config.From)
	}

	if config.KeystorePath != _keystorePath {
		t.Errorf("expected %v got %v ", _keystorePath, config.KeystorePath)
	}

	if config.BlockstorePath != _blockstore {
		t.Errorf("expected %v got %v ", _blockstore, config.BlockstorePath)
	}

	if config.FreshStart != _fresh {
		t.Errorf("expected %v got %v ", _fresh, config.FreshStart)
	}

	if config.LatestBlock != _latest {
		t.Errorf("expected %v got %v ", _latest, config.LatestBlock)
	}

}

func TestParseConfigInvalidChainID(t *testing.T) {

	_type := "test"
	_name := "test"
	_endpoint := "http://localhost:8080"
	fromAddress := "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa201"

	rCon := &cfg.RawChainConfig{
		Name:     _name,
		Type:     _type,
		Id:       "2X",
		Endpoint: _endpoint,
		From:     fromAddress,
		Opts:     map[string]string{},
	}

	set := flag.NewFlagSet("test", 0)

	ctx := cli.NewContext(nil, set, nil)

	_, err := ParseChainConfig(rCon, ctx)

	if err == nil {
		t.Error("expected invalid syntax error got nil ")
	}

}

func TestParseConfigNoBridgeContract(t *testing.T) {

	var chainID int = 3
	chainIDStr := strconv.Itoa(chainID)

	var maxGasPrice int = 21000
	maxGasPriceStr := strconv.Itoa(maxGasPrice)

	var gasLimit int = 31000
	gasLimitStr := strconv.Itoa(gasLimit)

	var startBlock int = 3000000
	startBlockStr := strconv.Itoa(startBlock)

	http := "http://localhost:8080"

	_type := "test"
	_name := "test"
	_endpoint := "http://localhost:8080"
	fromAddress := "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa201"

	erc20Handler := "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa203"
	erc721Handler := "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa204"
	genericHandler := "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa205"

	rCon := &cfg.RawChainConfig{
		Name:     _name,
		Type:     _type,
		Id:       chainIDStr,
		Endpoint: _endpoint,
		From:     fromAddress,
		Opts: map[string]string{
			"bridge":         "",
			"erc20Handler":   erc20Handler,
			"erc721Handler":  erc721Handler,
			"genericHandler": genericHandler,
			"maxGasPrice":    maxGasPriceStr,
			"gasLimit":       gasLimitStr,
			"http":           http,
			"startBlock":     startBlockStr,
		},
	}

	set := flag.NewFlagSet("test", 0)

	ctx := cli.NewContext(nil, set, nil)

	_, err := ParseChainConfig(rCon, ctx)

	if err == nil {
		t.Error("expected error  got error=nil ")
	}

}

func TestParseConfigInvalidMaxGasPrice(t *testing.T) {

	maxGasPriceStr := "2x"

	rCon := &cfg.RawChainConfig{
		Name:     "test",
		Type:     "test",
		Id:       "3",
		Endpoint: "http://localhost:8080",
		From:     "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
		Opts: map[string]string{
			"bridge":         "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"erc20Handler":   "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"erc721Handler":  "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"genericHandler": "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"maxGasPrice":    maxGasPriceStr,
			"gasLimit":       "3000",
			"http":           "http://localhost:8080",
			"startBlock":     "1",
		},
	}

	set := flag.NewFlagSet("test", 0)

	ctx := cli.NewContext(nil, set, nil)

	_, err := ParseChainConfig(rCon, ctx)

	if err == nil {
		t.Error("expected invalid maxGasPrice error , got error=nil")
	}

}

func TestParseConfigInvalidGasLimit(t *testing.T) {

	gasLimit := "2x"

	rCon := &cfg.RawChainConfig{
		Name:     "test",
		Type:     "test",
		Id:       "3",
		Endpoint: "http://localhost:8080",
		From:     "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
		Opts: map[string]string{
			"bridge":         "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"erc20Handler":   "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"erc721Handler":  "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"genericHandler": "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"maxGasPrice":    "3000",
			"gasLimit":       gasLimit,
			"http":           "http://localhost:8080",
			"startBlock":     "1",
		},
	}

	set := flag.NewFlagSet("test", 0)

	ctx := cli.NewContext(nil, set, nil)

	_, err := ParseChainConfig(rCon, ctx)

	if err == nil {
		t.Error("expected invalid gasLimit error , got error=nil")
	}

}

func TestParseConfigInvalidStartBlock(t *testing.T) {

	startBlock := "2x"

	rCon := &cfg.RawChainConfig{
		Name:     "test",
		Type:     "test",
		Id:       "3",
		Endpoint: "http://localhost:8080",
		From:     "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
		Opts: map[string]string{
			"bridge":         "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"erc20Handler":   "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"erc721Handler":  "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"genericHandler": "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"maxGasPrice":    "3000",
			"gasLimit":       "210000",
			"http":           "http://localhost:8080",
			"startBlock":     startBlock,
		},
	}

	set := flag.NewFlagSet("test", 0)

	ctx := cli.NewContext(nil, set, nil)

	_, err := ParseChainConfig(rCon, ctx)

	if err == nil {
		t.Error("expected invalid startBlock error , got error=nil")
	}

}

func TestParseConfigHttpFalseIsSet(t *testing.T) {

	http := "false"
	rCon := &cfg.RawChainConfig{
		Name:     "test",
		Type:     "test",
		Id:       "3",
		Endpoint: "http://localhost:8080",
		From:     "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
		Opts: map[string]string{
			"bridge":         "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"erc20Handler":   "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"erc721Handler":  "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"genericHandler": "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"maxGasPrice":    "30000",
			"gasLimit":       "30000",
			"http":           http,
			"startBlock":     "30000",
		},
	}

	set := flag.NewFlagSet("test", 0)

	ctx := cli.NewContext(nil, set, nil)

	config, err := ParseChainConfig(rCon, ctx)

	if err != nil {
		t.Error(err)
	}

	if config.Http != false {
		t.Errorf("expected Http %v got %v ", false, config.Http)
	}

}

func TestParseConfigInsecureKeyIsSetTrue(t *testing.T) {

	rCon := &cfg.RawChainConfig{
		Name:     "test",
		Type:     "test",
		Id:       "3",
		Endpoint: "http://localhost:8080",
		From:     "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
		Opts: map[string]string{
			"bridge":         "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"erc20Handler":   "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"erc721Handler":  "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"genericHandler": "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"maxGasPrice":    "30000",
			"gasLimit":       "30000",
			"http":           "false",
			"startBlock":     "30000",
		},
	}

	set := flag.NewFlagSet("test", 0)

	set.String("testkey", "testkey", "testkey")

	ctx := cli.NewContext(nil, set, nil)

	config, err := ParseChainConfig(rCon, ctx)

	if err != nil {
		t.Error(err)
	}

	if config.Insecure == false {
		t.Errorf("expected Insecure %v got %v ", true, config.Insecure)
	}

}

func TestParseConfigInsecureKeyIsSetFalse(t *testing.T) {

	http := "false"
	rCon := &cfg.RawChainConfig{
		Name:     "test",
		Type:     "test",
		Id:       "3",
		Endpoint: "http://localhost:8080",
		From:     "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
		Opts: map[string]string{
			"bridge":         "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"erc20Handler":   "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"erc721Handler":  "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"genericHandler": "0x18DfB0f9B4138d70d3EFe504A4D716D483Cfa202",
			"maxGasPrice":    "30000",
			"gasLimit":       "30000",
			"http":           http,
			"startBlock":     "30000",
		},
	}

	set := flag.NewFlagSet("test", 0)

	ctx := cli.NewContext(nil, set, nil)

	config, err := ParseChainConfig(rCon, ctx)

	if err != nil {
		t.Error(err)
	}

	if config.Insecure == true {
		t.Errorf("expected Insecure %v got %v ", false, config.Insecure)
	}

}
