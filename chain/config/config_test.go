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

	http := "http://localhost:8080"

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
