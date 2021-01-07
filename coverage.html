
<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<title>chainbridge-celo: Go Coverage Report</title>
		<style>
			body {
				background: black;
				color: rgb(80, 80, 80);
			}
			body, pre, #legend span {
				font-family: Menlo, monospace;
				font-weight: bold;
			}
			#topbar {
				background: black;
				position: fixed;
				top: 0; left: 0; right: 0;
				height: 42px;
				border-bottom: 1px solid rgb(80, 80, 80);
			}
			#content {
				margin-top: 50px;
			}
			#nav, #legend {
				float: left;
				margin-left: 10px;
			}
			#legend {
				margin-top: 12px;
			}
			#nav {
				margin-top: 10px;
			}
			#legend span {
				margin: 0 5px;
			}
			.cov0 { color: rgb(192, 0, 0) }
.cov1 { color: rgb(128, 128, 128) }
.cov2 { color: rgb(116, 140, 131) }
.cov3 { color: rgb(104, 152, 134) }
.cov4 { color: rgb(92, 164, 137) }
.cov5 { color: rgb(80, 176, 140) }
.cov6 { color: rgb(68, 188, 143) }
.cov7 { color: rgb(56, 200, 146) }
.cov8 { color: rgb(44, 212, 149) }
.cov9 { color: rgb(32, 224, 152) }
.cov10 { color: rgb(20, 236, 155) }

		</style>
	</head>
	<body>
		<div id="topbar">
			<div id="nav">
				<select id="files">
				
				<option value="file0">github.com/ChainSafe/chainbridge-celo/account.go (64.7%)</option>
				
				<option value="file1">github.com/ChainSafe/chainbridge-celo/chain/chain.go (0.0%)</option>
				
				<option value="file2">github.com/ChainSafe/chainbridge-celo/chain/client/client.go (0.0%)</option>
				
				<option value="file3">github.com/ChainSafe/chainbridge-celo/chain/config/config.go (59.3%)</option>
				
				<option value="file4">github.com/ChainSafe/chainbridge-celo/chain/listener/events.go (100.0%)</option>
				
				<option value="file5">github.com/ChainSafe/chainbridge-celo/chain/listener/listener.go (70.3%)</option>
				
				<option value="file6">github.com/ChainSafe/chainbridge-celo/chain/validator/validator_syncer.go (0.0%)</option>
				
				<option value="file7">github.com/ChainSafe/chainbridge-celo/chain/writer/proposal_data.go (69.0%)</option>
				
				<option value="file8">github.com/ChainSafe/chainbridge-celo/chain/writer/writer.go (32.1%)</option>
				
				<option value="file9">github.com/ChainSafe/chainbridge-celo/chain/writer/writer_methods.go (77.7%)</option>
				
				<option value="file10">github.com/ChainSafe/chainbridge-celo/cmd/cfg/config.go (73.6%)</option>
				
				<option value="file11">github.com/ChainSafe/chainbridge-celo/main.go (76.9%)</option>
				
				<option value="file12">github.com/ChainSafe/chainbridge-celo/router/router.go (92.3%)</option>
				
				</select>
			</div>
			<div id="legend">
				<span>not tracked</span>
			
				<span class="cov0">not covered</span>
				<span class="cov8">covered</span>
			
			</div>
		</div>
		<div id="content">
		
		<pre class="file" id="file0" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "os"
        "path/filepath"

        "github.com/ChainSafe/chainbridge-celo/flags"
        "github.com/ChainSafe/chainbridge-utils/crypto"
        "github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
        "github.com/ChainSafe/chainbridge-utils/keystore"
        gokeystore "github.com/ethereum/go-ethereum/accounts/keystore"
        "github.com/rs/zerolog/log"
        "github.com/urfave/cli/v2"
)

//dataHandler is a struct which wraps any extra data our CMD functions need that cannot be passed through parameters
type dataHandler struct {
        datadir string
}

//wrapHandler takes in a Cmd function (all declared below) and wraps
//it in the correct signature for the Cli Commands
func wrapHandler(hdl func(*cli.Context, *dataHandler) error) cli.ActionFunc <span class="cov8" title="1">{

        return func(ctx *cli.Context) error </span><span class="cov0" title="0">{
                // TODO: check logger1
                //err := startLogger(ctx)
                //if err != nil {
                //        return err
                //}

                datadir, err := getDataDir(ctx)
                if err != nil </span><span class="cov0" title="0">{
                        return fmt.Errorf("failed to access the datadir: %w", err)
                }</span>

                <span class="cov0" title="0">return hdl(ctx, &amp;dataHandler{datadir: datadir})</span>
        }
}

//handleGenerateCmd generates a keystore for the accounts
func handleGenerateCmd(ctx *cli.Context, dHandler *dataHandler) error <span class="cov8" title="1">{

        log.Info().Msg("Generating keypair...")

        keytype := crypto.Secp256k1Type

        // check if --password is set
        var password []byte = nil
        if pwdflag := ctx.String(flags.PasswordFlag.Name); pwdflag != "" </span><span class="cov8" title="1">{
                password = []byte(pwdflag)
        }</span>

        <span class="cov8" title="1">_, err := generateKeypair(keytype, dHandler.datadir, password)
        if err != nil </span><span class="cov0" title="0">{
                return fmt.Errorf("failed to generate key: %w", err)
        }</span>
        <span class="cov8" title="1">return nil</span>
}

//handleImportCmd imports external keystores into the bridge
func handleImportCmd(ctx *cli.Context, dHandler *dataHandler) error <span class="cov8" title="1">{
        log.Info().Msg("Importing key...")
        var err error

        // check if --ed25519 or --sr25519 is set
        keytype := crypto.Secp256k1Type

        if ctx.Bool(flags.EthereumImportFlag.Name) </span><span class="cov0" title="0">{
                if keyimport := ctx.Args().First(); keyimport != "" </span><span class="cov0" title="0">{
                        // check if --password is set
                        var password []byte = nil
                        if pwdflag := ctx.String(flags.PasswordFlag.Name); pwdflag != "" </span><span class="cov0" title="0">{
                                password = []byte(pwdflag)
                        }</span>
                        <span class="cov0" title="0">_, err = importEthKey(keyimport, dHandler.datadir, password, nil)</span>
                } else<span class="cov0" title="0"> {
                        return fmt.Errorf("Must provide a key to import.")
                }</span>
        } else<span class="cov8" title="1"> if privkeyflag := ctx.String(flags.PrivateKeyFlag.Name); privkeyflag != "" </span><span class="cov8" title="1">{
                // check if --password is set
                var password []byte = nil
                if pwdflag := ctx.String(flags.PasswordFlag.Name); pwdflag != "" </span><span class="cov8" title="1">{
                        password = []byte(pwdflag)
                }</span>

                <span class="cov8" title="1">_, err = importPrivKey(ctx, keytype, dHandler.datadir, privkeyflag, password)</span>
        } else<span class="cov0" title="0"> {
                if keyimport := ctx.Args().First(); keyimport != "" </span><span class="cov0" title="0">{
                        _, err = importKey(keyimport, dHandler.datadir)
                }</span> else<span class="cov0" title="0"> {
                        return fmt.Errorf("Must provide a key to import.")
                }</span>
        }

        <span class="cov8" title="1">if err != nil </span><span class="cov0" title="0">{
                return fmt.Errorf("failed to import key: %w", err)
        }</span>

        <span class="cov8" title="1">return nil</span>
}

//handleListCmd lists all accounts currently in the bridge
func handleListCmd(ctx *cli.Context, dHandler *dataHandler) error <span class="cov8" title="1">{

        _, err := listKeys(dHandler.datadir)
        if err != nil </span><span class="cov0" title="0">{
                return fmt.Errorf("failed to list keys: %w", err)
        }</span>

        <span class="cov8" title="1">return nil</span>
}

// getDataDir obtains the path to the keystore and returns it as a string
func getDataDir(ctx *cli.Context) (string, error) <span class="cov8" title="1">{
        // key directory is datadir/keystore/
        if dir := ctx.String(flags.KeystorePathFlag.Name); dir != "" </span><span class="cov8" title="1">{
                datadir, err := filepath.Abs(dir)
                if err != nil </span><span class="cov0" title="0">{
                        return "", err
                }</span>
                <span class="cov8" title="1">log.Trace().Msg(fmt.Sprintf("Using keystore dir: %s", datadir))
                return datadir, nil</span>
        }
        <span class="cov8" title="1">return "", fmt.Errorf("datadir flag not supplied")</span>
}

//importPrivKey imports a private key into a keypair
func importPrivKey(ctx *cli.Context, keytype, datadir, key string, password []byte) (string, error) <span class="cov8" title="1">{
        if password == nil </span><span class="cov0" title="0">{
                password = keystore.GetPassword("Enter password to encrypt keystore file:")
        }</span>
        <span class="cov8" title="1">keystorepath, err := keystoreDir(datadir)

        if keytype == "" </span><span class="cov8" title="1">{
                log.Info().Str("type", keytype).Msg("Using default key type")
                keytype = crypto.Secp256k1Type
        }</span>

        <span class="cov8" title="1">var kp crypto.Keypair

        if keytype == crypto.Secp256k1Type </span><span class="cov8" title="1">{
                // Hex must not have leading 0x
                if key[0:2] == "0x" </span><span class="cov0" title="0">{
                        kp, err = secp256k1.NewKeypairFromString(key[2:])
                }</span> else<span class="cov8" title="1"> {
                        kp, err = secp256k1.NewKeypairFromString(key)
                }</span>

                <span class="cov8" title="1">if err != nil </span><span class="cov0" title="0">{
                        return "", fmt.Errorf("could not generate secp256k1 keypair from given string: %w", err)
                }</span>
        }
        <span class="cov8" title="1">fp, err := filepath.Abs(keystorepath + "/" + kp.Address() + ".key")
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("invalid filepath: %w", err)
        }</span>

        <span class="cov8" title="1">file, err := os.OpenFile(filepath.Clean(fp), os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0600)
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("Unable to Open File: %w", err)
        }</span>

        <span class="cov8" title="1">defer func() </span><span class="cov8" title="1">{
                err = file.Close()
                if err != nil </span><span class="cov0" title="0">{
                        log.Error().Msg("import private key: could not close keystore file")
                }</span>
        }()

        <span class="cov8" title="1">err = keystore.EncryptAndWriteToFile(file, kp, password)
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("could not write key to file: %w", err)
        }</span>

        <span class="cov8" title="1">log.Info().Str("address", kp.Address()).Str("file", fp).Msg("private key imported")
        return fp, nil</span>

}

//importEthKey takes an ethereum keystore and converts it to our keystore format
func importEthKey(filename, datadir string, password, newPassword []byte) (string, error) <span class="cov8" title="1">{
        keystorepath, err := keystoreDir(datadir)
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("could not get keystore directory: %w", err)
        }</span>

        <span class="cov8" title="1">importdata, err := ioutil.ReadFile(filepath.Clean(filename))
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("could not read import file: %w", err)
        }</span>

        <span class="cov8" title="1">if password == nil </span><span class="cov0" title="0">{
                password = keystore.GetPassword("Enter password to decrypt keystore file:")
        }</span>

        <span class="cov8" title="1">key, err := gokeystore.DecryptKey(importdata, string(password))
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("Unable to decrypt file: %w", err)
        }</span>

        <span class="cov8" title="1">kp := secp256k1.NewKeypair(*key.PrivateKey)

        fp, err := filepath.Abs(keystorepath + "/" + kp.Address() + ".key")
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("invalid filepath: %w", err)
        }</span>

        <span class="cov8" title="1">file, err := os.OpenFile(filepath.Clean(fp), os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0600)
        if err != nil </span><span class="cov0" title="0">{
                return "", err
        }</span>

        <span class="cov8" title="1">defer func() </span><span class="cov8" title="1">{
                err = file.Close()
                if err != nil </span><span class="cov0" title="0">{
                        log.Error().Msg("generate keypair: could not close keystore file")
                }</span>
        }()

        <span class="cov8" title="1">if newPassword == nil </span><span class="cov0" title="0">{
                newPassword = keystore.GetPassword("Enter password to encrypt new keystore file:")
        }</span>

        <span class="cov8" title="1">err = keystore.EncryptAndWriteToFile(file, kp, newPassword)
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("could not write key to file: %w", err)
        }</span>

        <span class="cov8" title="1">log.Info().Str("address", kp.Address()).Str("file", fp).Msg("ETH key imported")
        return fp, nil</span>

}

// importKey imports a key specified by its filename to datadir/keystore/
// it saves it under the filename "[publickey].key"
// it returns the absolute path of the imported key file
func importKey(filename, datadir string) (string, error) <span class="cov8" title="1">{
        keystorepath, err := keystoreDir(datadir)
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("could not get keystore directory: %w", err)
        }</span>

        <span class="cov8" title="1">importdata, err := ioutil.ReadFile(filepath.Clean(filename))
        if err != nil </span><span class="cov8" title="1">{
                return "", fmt.Errorf("could not read import file: %w", err)
        }</span>

        <span class="cov0" title="0">ksjson := new(keystore.EncryptedKeystore)
        err = json.Unmarshal(importdata, ksjson)
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("could not read file contents: %w", err)
        }</span>

        <span class="cov0" title="0">keystorefile, err := filepath.Abs(keystorepath + "/" + ksjson.Address[2:] + ".key")
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("could not create keystore file path: %w", err)
        }</span>

        <span class="cov0" title="0">err = ioutil.WriteFile(keystorefile, importdata, 0600)
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("could not write to keystore directory: %w", err)
        }</span>

        <span class="cov0" title="0">log.Info().Str("address", ksjson.Address).Str("file", keystorefile).Msg("successfully imported key")
        return keystorefile, nil</span>
}

// listKeys lists all the keys in the datadir/keystore/ directory and returns them as a list of filepaths
func listKeys(datadir string) ([]string, error) <span class="cov8" title="1">{
        keys, err := getKeyFiles(datadir)
        if err != nil </span><span class="cov0" title="0">{
                return nil, err
        }</span>

        <span class="cov8" title="1">fmt.Printf("=== Found %d keys ===\n", len(keys))
        for i, key := range keys </span><span class="cov8" title="1">{
                fmt.Printf("[%d] %s\n", i, key)
        }</span>

        <span class="cov8" title="1">return keys, nil</span>
}

// getKeyFiles returns the filenames of all the keys in the datadir's keystore
func getKeyFiles(datadir string) ([]string, error) <span class="cov8" title="1">{
        keystorepath, err := keystoreDir(datadir)
        if err != nil </span><span class="cov0" title="0">{
                return nil, fmt.Errorf("could not get keystore directory: %w", err)
        }</span>

        <span class="cov8" title="1">files, err := ioutil.ReadDir(keystorepath)
        if err != nil </span><span class="cov0" title="0">{
                return nil, fmt.Errorf("could not read keystore dir: %w", err)
        }</span>

        <span class="cov8" title="1">keys := []string{}

        for _, f := range files </span><span class="cov8" title="1">{
                ext := filepath.Ext(f.Name())
                if ext == ".key" </span><span class="cov8" title="1">{
                        keys = append(keys, f.Name())
                }</span>
        }

        <span class="cov8" title="1">return keys, nil</span>
}

// generateKeypair create a new keypair with the corresponding type and saves it to datadir/keystore/[public key].key
// in json format encrypted using the specified password
// it returns the resulting filepath of the new key
func generateKeypair(keytype, datadir string, password []byte) (string, error) <span class="cov8" title="1">{
        if password == nil </span><span class="cov0" title="0">{
                password = keystore.GetPassword("Enter password to encrypt keystore file:")
        }</span>

        <span class="cov8" title="1">if keytype == "" </span><span class="cov8" title="1">{
                log.Info().Str("type", keytype).Msg("Using default key type")
                keytype = crypto.Secp256k1Type
        }</span>

        <span class="cov8" title="1">var kp crypto.Keypair
        var err error

        if keytype == crypto.Secp256k1Type </span><span class="cov8" title="1">{
                // generate secp256k1 keys
                kp, err = secp256k1.GenerateKeypair()
                if err != nil </span><span class="cov0" title="0">{
                        return "", fmt.Errorf("could not generate secp256k1 keypair: %w", err)
                }</span>
        } else<span class="cov0" title="0"> {
                return "", fmt.Errorf("invalid key type: %s", keytype)
        }</span>

        <span class="cov8" title="1">keystorepath, err := keystoreDir(datadir)
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("could not get keystore directory: %w", err)
        }</span>

        <span class="cov8" title="1">fp, err := filepath.Abs(keystorepath + "/" + kp.Address() + ".key")
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("invalid filepath: %w", err)
        }</span>

        <span class="cov8" title="1">file, err := os.OpenFile(filepath.Clean(fp), os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0600)
        if err != nil </span><span class="cov0" title="0">{
                return "", err
        }</span>

        <span class="cov8" title="1">defer func() </span><span class="cov8" title="1">{
                err = file.Close()
                if err != nil </span><span class="cov0" title="0">{
                        log.Error().Msg("generate keypair: could not close keystore file")
                }</span>
        }()

        <span class="cov8" title="1">err = keystore.EncryptAndWriteToFile(file, kp, password)
        if err != nil </span><span class="cov0" title="0">{
                return "", fmt.Errorf("could not write key to file: %w", err)
        }</span>

        <span class="cov8" title="1">log.Info().Str("address", kp.Address()).Str("type", keytype).Str("file", fp).Msg("key generated")
        return fp, nil</span>
}

// keystoreDir returnns the absolute filepath of the keystore directory given a datadir
// by default, it is ./keys/
// otherwise, it is datadir/keys/
func keystoreDir(keyPath string) (keystorepath string, err error) <span class="cov8" title="1">{
        // datadir specified, return datadir/keys as absolute path
        if keyPath != "" </span><span class="cov8" title="1">{
                keystorepath, err = filepath.Abs(keyPath)
                if err != nil </span><span class="cov0" title="0">{
                        return "", err
                }</span>
        } else<span class="cov0" title="0"> {
                // datadir not specified, use default
                keyPath = flags.DefaultKeystorePath

                keystorepath, err = filepath.Abs(keyPath)
                if err != nil </span><span class="cov0" title="0">{
                        return "", fmt.Errorf("could not create keystore file path: %w", err)
                }</span>
        }

        // if datadir does not exist, create it
        <span class="cov8" title="1">if _, err = os.Stat(keyPath); os.IsNotExist(err) </span><span class="cov8" title="1">{
                err = os.Mkdir(keyPath, os.ModePerm)
                if err != nil </span><span class="cov0" title="0">{
                        return "", err
                }</span>
        }

        // if datadir/keystore does not exist, create it
        <span class="cov8" title="1">if _, err = os.Stat(keystorepath); os.IsNotExist(err) </span><span class="cov0" title="0">{
                err = os.Mkdir(keystorepath, os.ModePerm)
                if err != nil </span><span class="cov0" title="0">{
                        return "", err
                }</span>
        }

        <span class="cov8" title="1">return keystorepath, nil</span>
}
</pre>
		
		<pre class="file" id="file1" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package chain

import (
        "fmt"

        bridgeHandler "github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
        erc20Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
        erc721Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
        "github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
        "github.com/ChainSafe/chainbridge-celo/chain/client"
        "github.com/ChainSafe/chainbridge-celo/chain/config"
        listener "github.com/ChainSafe/chainbridge-celo/chain/listener"
        "github.com/ChainSafe/chainbridge-celo/chain/writer"
        "github.com/ChainSafe/chainbridge-celo/msg"
        "github.com/pkg/errors"
        "github.com/rs/zerolog/log"
)

// checkBlockstore queries the blockstore for the latest known block. If the latest block is
// greater than cfg.startBlock, then cfg.startBlock is replaced with the latest known block.
type Listener interface {
        StartPollingBlocks() error
        SetContracts(bridge listener.IBridge, erc20Handler listener.IERC20Handler, erc721Handler listener.IERC721Handler, genericHandler listener.IGenericHandler)
        //LatestBlock() *metrics.LatestBlock
}

type Writer interface {
        SetBridge(bridge writer.Bridger)
}

type Chain struct {
        cfg      *config.CeloChainConfig // The config of the chain
        listener Listener                // The listener of this chain
        writer   Writer                  // The writer of the chain
        client   *client.Client
        stopChn  &lt;-chan struct{}
}

func InitializeChain(cc *config.CeloChainConfig, c *client.Client, listener Listener, writer Writer, stopChn &lt;-chan struct{}) (*Chain, error) <span class="cov0" title="0">{

        bridgeContract, err := bridgeHandler.NewBridge(cc.BridgeContract, c)
        if err != nil </span><span class="cov0" title="0">{
                return nil, err
        }</span>

        <span class="cov0" title="0">chainId, err := bridgeContract.ChainID(c.CallOpts())
        if err != nil </span><span class="cov0" title="0">{
                return nil, err
        }</span>
        <span class="cov0" title="0">if chainId != uint8(cc.ID) </span><span class="cov0" title="0">{
                return nil, errors.New(fmt.Sprintf("chainId (%d) and configuration chainId (%d) do not match", chainId, cc.ID))
        }</span>

        <span class="cov0" title="0">erc20HandlerContract, err := erc20Handler.NewERC20Handler(cc.Erc20HandlerContract, c)
        if err != nil </span><span class="cov0" title="0">{
                return nil, err
        }</span>
        <span class="cov0" title="0">erc721HandlerContract, err := erc721Handler.NewERC721Handler(cc.Erc721HandlerContract, c)
        if err != nil </span><span class="cov0" title="0">{
                return nil, err
        }</span>
        <span class="cov0" title="0">genericHandlerContract, err := GenericHandler.NewGenericHandler(cc.GenericHandlerContract, c)
        if err != nil </span><span class="cov0" title="0">{
                return nil, err
        }</span>
        <span class="cov0" title="0">if cc.LatestBlock </span><span class="cov0" title="0">{
                curr, err := c.LatestBlock()
                if err != nil </span><span class="cov0" title="0">{
                        return nil, err
                }</span>
                <span class="cov0" title="0">cc.StartBlock = curr</span>
        }
        <span class="cov0" title="0">listener.SetContracts(bridgeContract, erc20HandlerContract, erc721HandlerContract, genericHandlerContract)
        writer.SetBridge(bridgeContract)
        return &amp;Chain{
                cfg:      cc,
                writer:   writer,
                listener: listener,
                stopChn:  stopChn,
        }, nil</span>
}

func (c *Chain) Start() error <span class="cov0" title="0">{
        err := c.listener.StartPollingBlocks()
        if err != nil </span><span class="cov0" title="0">{
                return err
        }</span>
        <span class="cov0" title="0">go func() </span><span class="cov0" title="0">{
                &lt;-c.stopChn
                if c.client != nil </span><span class="cov0" title="0">{
                        c.client.Close()
                }</span>
        }()
        <span class="cov0" title="0">log.Debug().Msg("Chain started!")
        return nil</span>
}

func (c *Chain) ID() msg.ChainId <span class="cov0" title="0">{
        return c.cfg.ID
}</span>

func (c *Chain) Name() string <span class="cov0" title="0">{
        return c.cfg.Name
}</span>

//func (c *Chain) LatestBlock() *metrics.LatestBlock {
//        return c.listener.LatestBlock()
//}
</pre>
		
		<pre class="file" id="file2" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package client

import (
        "context"
        "errors"
        "fmt"
        "math/big"
        "sync"
        "time"

        "github.com/ChainSafe/chainbridge-utils/crypto/secp256k1"
        eth "github.com/ethereum/go-ethereum"
        "github.com/ethereum/go-ethereum/accounts/abi/bind"
        ethcommon "github.com/ethereum/go-ethereum/common"
        "github.com/ethereum/go-ethereum/core/types"
        ethcrypto "github.com/ethereum/go-ethereum/crypto"
        "github.com/ethereum/go-ethereum/ethclient"
        "github.com/ethereum/go-ethereum/rpc"
        "github.com/rs/zerolog/log"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000

var BlockRetryInterval = time.Second * 5

type Client struct {
        *ethclient.Client
        endpoint    string
        http        bool
        kp          *secp256k1.Keypair
        gasLimit    *big.Int
        maxGasPrice *big.Int
        opts        *bind.TransactOpts
        callOpts    *bind.CallOpts
        nonce       uint64
        nonceLock   sync.Mutex
        optsLock    sync.Mutex
        stop        chan int // All routines should exit when this channel is closed
}

type LogFilterWithLatestBlock interface {
        FilterLogs(ctx context.Context, q eth.FilterQuery) ([]types.Log, error)
        LatestBlock() (*big.Int, error)
}

// NewConnection returns an uninitialized connection, must call Client.Connect() before using.
func NewClient(endpoint string, http bool, kp *secp256k1.Keypair, gasLimit *big.Int, gasPrice *big.Int) (*Client, error) <span class="cov0" title="0">{
        c := &amp;Client{
                endpoint:    endpoint,
                http:        http,
                kp:          kp,
                maxGasPrice: gasPrice,
                gasLimit:    gasLimit,
                stop:        make(chan int),
        }
        if err := c.Connect(); err != nil </span><span class="cov0" title="0">{
                return nil, err
        }</span>
        <span class="cov0" title="0">return c, nil</span>
}

// Connect starts the ethereum WS connection
func (c *Client) Connect() error <span class="cov0" title="0">{
        log.Info().Str("url", c.endpoint).Msg("Connecting to ethereum chain...")
        var rpcClient *rpc.Client
        var err error
        // Start http or ws client
        if c.http </span><span class="cov0" title="0">{
                rpcClient, err = rpc.DialHTTP(c.endpoint)
        }</span> else<span class="cov0" title="0"> {
                rpcClient, err = rpc.DialWebsocket(context.Background(), c.endpoint, "/ws")
        }</span>
        <span class="cov0" title="0">if err != nil </span><span class="cov0" title="0">{
                return err
        }</span>
        <span class="cov0" title="0">c.Client = ethclient.NewClient(rpcClient)

        // Construct tx opts, call opts, and nonce mechanism
        opts, _, err := c.newTransactOpts(big.NewInt(0), c.gasLimit, c.maxGasPrice)
        if err != nil </span><span class="cov0" title="0">{
                return err
        }</span>
        <span class="cov0" title="0">c.opts = opts
        c.nonce = 0
        c.callOpts = &amp;bind.CallOpts{From: c.kp.CommonAddress()}
        return nil</span>
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (c *Client) newTransactOpts(value, gasLimit, gasPrice *big.Int) (*bind.TransactOpts, uint64, error) <span class="cov0" title="0">{
        privateKey := c.kp.PrivateKey()
        address := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

        nonce, err := c.PendingNonceAt(context.Background(), address)
        if err != nil </span><span class="cov0" title="0">{
                return nil, 0, err
        }</span>

        <span class="cov0" title="0">auth := bind.NewKeyedTransactor(privateKey)
        auth.Nonce = big.NewInt(int64(nonce))
        auth.Value = value
        auth.GasLimit = uint64(gasLimit.Int64())
        auth.GasPrice = gasPrice
        auth.Context = context.Background()

        return auth, nonce, nil</span>
}

func (c *Client) Keypair() *secp256k1.Keypair <span class="cov0" title="0">{
        return c.kp
}</span>

func (c *Client) Opts() *bind.TransactOpts <span class="cov0" title="0">{
        return c.opts
}</span>

func (c *Client) CallOpts() *bind.CallOpts <span class="cov0" title="0">{
        return c.callOpts
}</span>

func (c *Client) LockAndUpdateNonce() error <span class="cov0" title="0">{
        c.nonceLock.Lock()
        nonce, err := c.PendingNonceAt(context.Background(), c.opts.From)
        if err != nil </span><span class="cov0" title="0">{
                c.nonceLock.Unlock()
                return err
        }</span>
        <span class="cov0" title="0">c.opts.Nonce.SetUint64(nonce)
        return nil</span>
}

func (c *Client) UnlockNonce() <span class="cov0" title="0">{
        c.nonceLock.Unlock()
}</span>

func (c *Client) UnlockOpts() <span class="cov0" title="0">{
        c.optsLock.Unlock()
}</span>

// LatestBlock returns the latest block from the current chain
func (c *Client) LatestBlock() (*big.Int, error) <span class="cov0" title="0">{
        header, err := c.HeaderByNumber(context.Background(), nil)
        if err != nil </span><span class="cov0" title="0">{
                return nil, err
        }</span>
        <span class="cov0" title="0">return header.Number, nil</span>
}

// EnsureHasBytecode asserts if contract code exists at the specified address
func (c *Client) EnsureHasBytecode(addr ethcommon.Address) error <span class="cov0" title="0">{
        code, err := c.CodeAt(context.Background(), addr, nil)
        if err != nil </span><span class="cov0" title="0">{
                return err
        }</span>

        <span class="cov0" title="0">if len(code) == 0 </span><span class="cov0" title="0">{
                return fmt.Errorf("no bytecode found at %s", addr.Hex())
        }</span>
        <span class="cov0" title="0">return nil</span>
}

// WaitForBlock will poll for the block number until the current block is equal or greater than
func (c *Client) WaitForBlock(block *big.Int) error <span class="cov0" title="0">{
        for </span><span class="cov0" title="0">{
                select </span>{
                case &lt;-c.stop:<span class="cov0" title="0">
                        return errors.New("connection terminated")</span>
                default:<span class="cov0" title="0">
                        currBlock, err := c.LatestBlock()
                        if err != nil </span><span class="cov0" title="0">{
                                return err
                        }</span>

                        // Equal or greater than target
                        <span class="cov0" title="0">if currBlock.Cmp(block) &gt;= 0 </span><span class="cov0" title="0">{
                                return nil
                        }</span>
                        <span class="cov0" title="0">log.Trace().Interface("target", block).Interface("current", currBlock).Msg("Block not ready, waiting")
                        time.Sleep(BlockRetryInterval)
                        continue</span>
                }
        }
}

// LockAndUpdateOpts acquires a lock on the opts before updating the nonce
// and gas price.
func (c *Client) LockAndUpdateOpts() error <span class="cov0" title="0">{
        c.optsLock.Lock()

        gasPrice, err := c.SafeEstimateGas(context.TODO())
        if err != nil </span><span class="cov0" title="0">{
                return err
        }</span>
        <span class="cov0" title="0">c.opts.GasPrice = gasPrice

        nonce, err := c.PendingNonceAt(context.Background(), c.opts.From)
        if err != nil </span><span class="cov0" title="0">{
                c.optsLock.Unlock()
                return err
        }</span>
        <span class="cov0" title="0">c.opts.Nonce.SetUint64(nonce)
        return nil</span>
}

func (c *Client) SafeEstimateGas(ctx context.Context) (*big.Int, error) <span class="cov0" title="0">{
        gasPrice, err := c.SuggestGasPrice(context.TODO())
        if err != nil </span><span class="cov0" title="0">{
                return nil, err
        }</span>

        // Check we aren't exceeding our limit

        <span class="cov0" title="0">if gasPrice.Cmp(c.maxGasPrice) == 1 </span><span class="cov0" title="0">{
                return c.maxGasPrice, nil
        }</span> else<span class="cov0" title="0"> {
                return gasPrice, nil
        }</span>
}

// Close terminates the client connection and stops any running routines
func (c *Client) Close() <span class="cov0" title="0">{
        if c.Client != nil </span><span class="cov0" title="0">{
                c.Client.Close()
        }</span>
}
</pre>
		
		<pre class="file" id="file3" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package config

import (
        "math/big"
        "strconv"

        "github.com/ChainSafe/chainbridge-celo/chain/client"
        "github.com/ChainSafe/chainbridge-celo/cmd/cfg"
        "github.com/ChainSafe/chainbridge-celo/flags"
        "github.com/ChainSafe/chainbridge-celo/msg"
        utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
        "github.com/ethereum/go-ethereum/common"
        "github.com/pkg/errors"
        "github.com/urfave/cli/v2"
)

const DefaultGasLimit = 6721975
const DefaultGasPrice = 20000000000

type CeloChainConfig struct {
        ID                     msg.ChainId // ChainID
        Name                   string      // Human-readable chain name
        Endpoint               string      // url for rpc endpoint
        From                   string      // address of key to use // TODO: name should be changed
        KeystorePath           string      // Location of keyfiles
        BlockstorePath         string
        FreshStart             bool // Disables loading from blockstore at start
        BridgeContract         common.Address
        Erc20HandlerContract   common.Address
        Erc721HandlerContract  common.Address
        GenericHandlerContract common.Address
        GasLimit               *big.Int
        MaxGasPrice            *big.Int
        Http                   bool // Config for type of connection
        StartBlock             *big.Int
        LatestBlock            bool
        Insecure               bool
}

func (cfg *CeloChainConfig) EnsureContractsHaveBytecode(conn *client.Client) error <span class="cov0" title="0">{
        err := conn.EnsureHasBytecode(cfg.BridgeContract)
        if err != nil </span><span class="cov0" title="0">{
                return err
        }</span>
        <span class="cov0" title="0">err = conn.EnsureHasBytecode(cfg.Erc20HandlerContract)
        if err != nil </span><span class="cov0" title="0">{
                return err
        }</span>
        <span class="cov0" title="0">err = conn.EnsureHasBytecode(cfg.GenericHandlerContract)
        if err != nil </span><span class="cov0" title="0">{
                return err
        }</span>
        <span class="cov0" title="0">err = conn.EnsureHasBytecode(cfg.Erc721HandlerContract)
        if err != nil </span><span class="cov0" title="0">{
                return err
        }</span>
        <span class="cov0" title="0">return nil</span>
}

// parseChainConfig uses a core.ChainConfig to construct a corresponding Config
func ParseChainConfig(rawCfg *cfg.RawChainConfig, ctx *cli.Context) (*CeloChainConfig, error) <span class="cov8" title="1">{
        var ks string
        var insecure bool
        if key := ctx.String(flags.TestKeyFlag.Name); key != "" </span><span class="cov0" title="0">{
                ks = key
                insecure = true
        }</span> else<span class="cov8" title="1"> {
                if ksPath := ctx.String(flags.KeystorePathFlag.Name); ksPath != "" </span><span class="cov8" title="1">{
                        ks = ksPath
                }</span>
        }
        <span class="cov8" title="1">chainId, err := strconv.Atoi(rawCfg.Id)
        if err != nil </span><span class="cov0" title="0">{
                return nil, err
        }</span>

        <span class="cov8" title="1">config := &amp;CeloChainConfig{
                Name:                   rawCfg.Name,
                ID:                     msg.ChainId(chainId),
                Endpoint:               rawCfg.Endpoint,
                From:                   rawCfg.From,
                KeystorePath:           ks,
                BlockstorePath:         ctx.String(flags.BlockstorePathFlag.Name),
                FreshStart:             ctx.Bool(flags.FreshStartFlag.Name),
                LatestBlock:            ctx.Bool(flags.LatestBlockFlag.Name),
                BridgeContract:         utils.ZeroAddress,
                Erc20HandlerContract:   utils.ZeroAddress,
                Erc721HandlerContract:  utils.ZeroAddress,
                GenericHandlerContract: utils.ZeroAddress,
                GasLimit:               big.NewInt(DefaultGasLimit),
                MaxGasPrice:            big.NewInt(DefaultGasPrice),
                Http:                   false,
                StartBlock:             big.NewInt(0),
                Insecure:               insecure,
        }

        if contract, ok := rawCfg.Opts["bridge"]; ok &amp;&amp; contract != "" </span><span class="cov8" title="1">{
                config.BridgeContract = common.HexToAddress(contract)
        }</span> else<span class="cov0" title="0"> {
                return nil, errors.New("must provide opts.bridge field for ethereum config")
        }</span>

        <span class="cov8" title="1">config.Erc20HandlerContract = common.HexToAddress(rawCfg.Opts["erc20Handler"])

        config.Erc721HandlerContract = common.HexToAddress(rawCfg.Opts["erc721Handler"])

        config.GenericHandlerContract = common.HexToAddress(rawCfg.Opts["genericHandler"])

        if gasPrice, ok := rawCfg.Opts["maxGasPrice"]; ok </span><span class="cov8" title="1">{
                price := big.NewInt(0)
                _, pass := price.SetString(gasPrice, 10)
                if pass </span><span class="cov8" title="1">{
                        config.MaxGasPrice = price
                        delete(rawCfg.Opts, "maxGasPrice")
                }</span> else<span class="cov0" title="0"> {
                        return nil, errors.New("unable to parse max gas price")
                }</span>
        }

        <span class="cov8" title="1">if gasLimit, ok := rawCfg.Opts["gasLimit"]; ok </span><span class="cov8" title="1">{
                limit := big.NewInt(0)
                _, pass := limit.SetString(gasLimit, 10)
                if pass </span><span class="cov8" title="1">{
                        config.GasLimit = limit
                }</span> else<span class="cov0" title="0"> {
                        return nil, errors.New("unable to parse gas limit")
                }</span>
        }

        <span class="cov8" title="1">if HTTP, ok := rawCfg.Opts["http"]; ok &amp;&amp; HTTP == "true" </span><span class="cov0" title="0">{
                config.Http = true
        }</span> else<span class="cov8" title="1"> if HTTP, ok := rawCfg.Opts["http"]; ok &amp;&amp; HTTP == "false" </span><span class="cov0" title="0">{
                config.Http = false
        }</span>

        <span class="cov8" title="1">if startBlock, ok := rawCfg.Opts["startBlock"]; ok &amp;&amp; startBlock != "" </span><span class="cov8" title="1">{
                block := big.NewInt(0)
                _, pass := block.SetString(startBlock, 10)
                if pass </span><span class="cov8" title="1">{
                        config.StartBlock = block
                }</span> else<span class="cov0" title="0"> {
                        return nil, errors.New("unable to parse start block")
                }</span>
        }
        <span class="cov8" title="1">return config, nil</span>
}
</pre>
		
		<pre class="file" id="file4" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package listener

import (
        "github.com/ChainSafe/chainbridge-celo/msg"
        "github.com/ethereum/go-ethereum/accounts/abi/bind"
        "github.com/rs/zerolog/log"
)

func (l *listener) handleErc20DepositedEvent(destId msg.ChainId, nonce msg.Nonce) (*msg.Message, error) <span class="cov8" title="1">{
        log.Info().Interface("dest", destId).Interface("nonce", nonce).Msg("Handling fungible deposit event")
        //TODO no call opts. should have From in original chainbridge.
        record, err := l.erc20HandlerContract.GetDepositRecord(&amp;bind.CallOpts{}, uint64(nonce), uint8(destId))
        if err != nil </span><span class="cov8" title="1">{
                log.Error().Err(err).Msg("Error Unpacking ERC20 Deposit Record")
                return nil, err
        }</span>

        <span class="cov8" title="1">return msg.NewFungibleTransfer(
                l.cfg.ID,
                destId,
                nonce,
                record.ResourceID,
                nil,
                nil,
                record.Amount,
                record.DestinationRecipientAddress,
        ), nil</span>
}

func (l *listener) handleErc721DepositedEvent(destId msg.ChainId, nonce msg.Nonce) (*msg.Message, error) <span class="cov8" title="1">{
        log.Info().Msg("Handling nonfungible deposit event")
        //TODO no call opts. should have From in original chainbridge.
        record, err := l.erc721HandlerContract.GetDepositRecord(&amp;bind.CallOpts{}, uint64(nonce), uint8(destId))
        if err != nil </span><span class="cov8" title="1">{
                log.Error().Err(err).Msg("Error Unpacking ERC721 Deposit Record")
                return nil, err
        }</span>

        <span class="cov8" title="1">return msg.NewNonFungibleTransfer(
                l.cfg.ID,
                destId,
                nonce,
                record.ResourceID,
                nil,
                nil,
                record.TokenID,
                record.DestinationRecipientAddress,
                record.MetaData,
        ), nil</span>
}

func (l *listener) handleGenericDepositedEvent(destId msg.ChainId, nonce msg.Nonce) (*msg.Message, error) <span class="cov8" title="1">{
        log.Info().Msg("Handling generic deposit event")
        //TODO no call opts. should have From in original chainbridge.
        record, err := l.genericHandlerContract.GetDepositRecord(&amp;bind.CallOpts{}, uint64(nonce), uint8(destId))
        if err != nil </span><span class="cov8" title="1">{
                log.Error().Err(err).Msg("Error Unpacking Generic Deposit Record")
                return nil, err
        }</span>

        <span class="cov8" title="1">return msg.NewGenericTransfer(
                l.cfg.ID,
                destId,
                nonce,
                record.ResourceID,
                nil,
                nil,
                record.MetaData[:],
        ), nil</span>
}
</pre>
		
		<pre class="file" id="file5" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package listener

import (
        "context"
        "errors"
        "fmt"
        "math/big"
        "time"

        "github.com/ChainSafe/chainbridge-celo/chain/client"
        "github.com/ChainSafe/chainbridge-celo/chain/config"
        "github.com/ChainSafe/chainbridge-celo/msg"
        utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
        eth "github.com/ethereum/go-ethereum"
        "github.com/ethereum/go-ethereum/accounts/abi/bind"
        ethcommon "github.com/ethereum/go-ethereum/common"
        "github.com/rs/zerolog/log"
)

var BlockDelay = big.NewInt(10)
var BlockRetryInterval = time.Second * 5
var ErrFatalPolling = errors.New("listener block polling failed")
var ExpectedBlockTime = time.Second
var BlockRetryLimit = 5

type listener struct {
        cfg                    *config.CeloChainConfig
        router                 IRouter
        bridgeContract         IBridge // instance of bound bridge contract
        erc20HandlerContract   IERC20Handler
        erc721HandlerContract  IERC721Handler
        genericHandlerContract IGenericHandler
        blockstore             Blockstorer
        stop                   &lt;-chan struct{}
        sysErr                 chan&lt;- error // Reports fatal error to core
        syncer                 BlockSyncer
        //latestBlock            *metrics.LatestBlock
        //metrics                *metrics.ChainMetrics
        client client.LogFilterWithLatestBlock
}

type BlockSyncer interface {
        Sync(latestBlock *big.Int) error
}

type IRouter interface {
        Send(msg *msg.Message) error
}
type Blockstorer interface {
        StoreBlock(*big.Int) error
}

func NewListener(cfg *config.CeloChainConfig, client client.LogFilterWithLatestBlock, bs Blockstorer, stop &lt;-chan struct{}, sysErr chan&lt;- error, syncer BlockSyncer, router IRouter) *listener <span class="cov8" title="1">{
        return &amp;listener{
                cfg:        cfg,
                blockstore: bs,
                stop:       stop,
                sysErr:     sysErr,
                syncer:     syncer,
                router:     router,
                client:     client,
        }
}</span>

func (l *listener) SetContracts(bridge IBridge, erc20Handler IERC20Handler, erc721Handler IERC721Handler, genericHandler IGenericHandler) <span class="cov8" title="1">{
        l.bridgeContract = bridge
        l.erc20HandlerContract = erc20Handler
        l.erc721HandlerContract = erc721Handler
        l.genericHandlerContract = genericHandler
}</span>

func (l *listener) StartPollingBlocks() error <span class="cov0" title="0">{
        log.Debug().Msg("Starting listener...")

        go func() </span><span class="cov0" title="0">{
                err := l.pollBlocks()
                if err != nil </span><span class="cov0" title="0">{
                        log.Error().Err(err).Msg("Polling blocks failed")
                }</span>
        }()

        <span class="cov0" title="0">return nil</span>
}

// TODO this is metrics latest block, naming mess
//func (l *listener) LatestBlock() *metrics.LatestBlock {
//        return l.latestBlock
//}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.cfg.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before continuing to the next block.
func (l *listener) pollBlocks() error <span class="cov8" title="1">{
        log.Info().Msg("Polling Blocks...")
        var currentBlock = l.cfg.StartBlock
        var retry = BlockRetryLimit
        for </span><span class="cov8" title="1">{
                select </span>{
                case &lt;-l.stop:<span class="cov8" title="1">
                        return errors.New("polling terminated")</span>
                default:<span class="cov8" title="1">
                        // No more retries, goto next block
                        if retry == 0 </span><span class="cov0" title="0">{
                                log.Error().Msg("Polling failed, retries exceeded")
                                l.sysErr &lt;- ErrFatalPolling
                                return nil
                        }</span>

                        <span class="cov8" title="1">latestBlock, err := l.client.LatestBlock()
                        if err != nil </span><span class="cov8" title="1">{
                                log.Error().Err(err).Str("block", currentBlock.String()).Msg("Unable to get latest block")
                                retry--
                                time.Sleep(BlockRetryInterval)
                                continue</span>
                        }

                        // Sleep if the difference is less than BlockDelay; (latest - current) &lt; BlockDelay
                        <span class="cov8" title="1">if big.NewInt(0).Sub(latestBlock, currentBlock).Cmp(BlockDelay) == -1 </span><span class="cov0" title="0">{
                                log.Debug().Str("target", currentBlock.String()).Str("latest", latestBlock.String()).Msg("Block not ready, will retry")
                                time.Sleep(BlockRetryInterval)
                                continue</span>
                        }

                        <span class="cov8" title="1">err = l.syncer.Sync(currentBlock)
                        if err != nil </span><span class="cov0" title="0">{
                                log.Error().Str("block", currentBlock.String()).Err(err).Msg("Failed to sync validators for block")
                                continue</span>
                        }

                        // Parse out events
                        <span class="cov8" title="1">err = l.getDepositEventsAndProofsForBlock(currentBlock)
                        if err != nil </span><span class="cov0" title="0">{
                                log.Error().Str("block", currentBlock.String()).Err(err).Msg("Failed to get events for block")
                                retry--
                                continue</span>
                        }

                        // Write to block store. Not a critical operation, no need to retry
                        <span class="cov8" title="1">err = l.blockstore.StoreBlock(currentBlock)
                        if err != nil </span><span class="cov0" title="0">{
                                log.Error().Str("block", currentBlock.String()).Err(err).Msg("Failed to write latest block to blockstore")
                        }</span>

                        //if l.metrics != nil {
                        //        l.metrics.BlocksProcessed.Inc()
                        //        l.metrics.LatestProcessedBlock.Set(float64(latestBlock.Int64()))
                        //}
                        //
                        //l.latestBlock.Height = big.NewInt(0).Set(latestBlock)
                        //l.latestBlock.LastUpdated = time.Now()

                        // Goto next block and reset retry counter
                        <span class="cov8" title="1">currentBlock.Add(currentBlock, big.NewInt(1))
                        retry = BlockRetryLimit</span>
                }
        }
}

// TODO: Proof construction.
func (l *listener) getDepositEventsAndProofsForBlock(latestBlock *big.Int) error <span class="cov8" title="1">{
        log.Debug().Str("block", latestBlock.String()).Msg("Querying block for deposit events")
        query := buildQuery(l.cfg.BridgeContract, utils.Deposit, latestBlock, latestBlock)

        // querying for logs
        logs, err := l.client.FilterLogs(context.Background(), query)
        if err != nil </span><span class="cov0" title="0">{
                return fmt.Errorf("unable to Filter Logs: %w", err)
        }</span>

        // read through the log events and handle their deposit event if handler is recognized
        <span class="cov8" title="1">for _, eventLog := range logs </span><span class="cov8" title="1">{
                var m *msg.Message
                destId := msg.ChainId(eventLog.Topics[1].Big().Uint64())
                rId := msg.ResourceId(eventLog.Topics[2])
                nonce := msg.Nonce(eventLog.Topics[3].Big().Uint64())

                addr, err := l.bridgeContract.ResourceIDToHandlerAddress(&amp;bind.CallOpts{}, rId)
                if err != nil </span><span class="cov0" title="0">{
                        return fmt.Errorf("failed to get handler from resource ID %x, reason: %w", rId, err)
                }</span>

                <span class="cov8" title="1">if addr == l.cfg.Erc20HandlerContract </span><span class="cov8" title="1">{
                        m, err = l.handleErc20DepositedEvent(destId, nonce)
                }</span> else<span class="cov8" title="1"> if addr == l.cfg.Erc721HandlerContract </span><span class="cov8" title="1">{
                        m, err = l.handleErc721DepositedEvent(destId, nonce)
                }</span> else<span class="cov8" title="1"> if addr == l.cfg.GenericHandlerContract </span><span class="cov8" title="1">{
                        m, err = l.handleGenericDepositedEvent(destId, nonce)
                }</span> else<span class="cov8" title="1"> {
                        log.Error().Err(err).Str("handler", addr.Hex()).Msg("event has unrecognized handler")
                        return nil
                }</span>

                <span class="cov8" title="1">if err != nil </span><span class="cov0" title="0">{
                        return err
                }</span>

                <span class="cov8" title="1">err = l.router.Send(m)
                if err != nil </span><span class="cov0" title="0">{
                        log.Error().Err(err).Msg("subscription error: failed to route message")
                }</span>
        }

        <span class="cov8" title="1">return nil</span>
}

//TODO removenolint
//nolint
//COMMENTED SINCE CURRENTLTY UNUSED. SEEMS TO BE USED FOR BLOCK PROOF BUILDING
//func (l *listener) getBlockHashFromTransactionHash(txHash ethcommon.Hash) (blockHash ethcommon.Hash, err error) {
//
//        receipt, err := l.conn.Client().TransactionReceipt(context.Background(), txHash)
//        if err != nil {
//                return txHash, fmt.Errorf("unable to get BlockHash: %w", err)
//        }
//        return receipt.BlockHash, nil
//}
//
////TODO removenolint
////nolint
//func (l *listener) getTransactionsFromBlockHash(blockHash ethcommon.Hash) (txHashes []ethcommon.Hash, txRoot ethcommon.Hash, err error) {
//        block, err := l.conn.Client().BlockByHash(context.Background(), blockHash)
//        if err != nil {
//                return nil, ethcommon.Hash{}, fmt.Errorf("unable to get BlockHash: %w", err)
//        }
//
//        var transactionHashes []ethcommon.Hash
//
//        transactions := block.Transactions()
//        for _, transaction := range transactions {
//                transactionHashes = append(transactionHashes, transaction.Hash())
//        }
//
//        return transactionHashes, block.Root(), nil
//}
//
//nolint
// buildQuery constructs a query for the bridgeContract by hashing sig to get the event topic
func buildQuery(contract ethcommon.Address, sig utils.EventSig, startBlock *big.Int, endBlock *big.Int) eth.FilterQuery <span class="cov8" title="1">{
        query := eth.FilterQuery{
                FromBlock: startBlock,
                ToBlock:   endBlock,
                Addresses: []ethcommon.Address{contract},
                Topics: [][]ethcommon.Hash{
                        {sig.GetTopic()},
                },
        }
        return query
}</span>
</pre>
		
		<pre class="file" id="file6" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
//nolint
//TODO remove nolint when start using this package
package validator

import (
        "context"
        "math/big"

        "github.com/ChainSafe/chainbridge-celo/chain/client"
        "github.com/celo-org/celo-bls-go/bls"
        "github.com/ethereum/go-ethereum/consensus/istanbul"
        "github.com/ethereum/go-ethereum/core/types"
        blscrypto "github.com/ethereum/go-ethereum/crypto/bls"
        "github.com/pkg/errors"
)

type ValidatorSyncer struct {
        client *client.Client
}

// ExtractValidators pulls the extra data from the block header and extract
// validators and returns an array of validator data
func (v *ValidatorSyncer) ExtractValidators(blockNumber uint64) ([]*istanbul.ValidatorData, error) <span class="cov0" title="0">{
        header, err := v.client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
        if err != nil </span><span class="cov0" title="0">{
                return nil, errors.Wrap(err, "getting the block header by number failed")
        }</span>

        <span class="cov0" title="0">extra, err := types.ExtractIstanbulExtra(header)
        if err != nil </span><span class="cov0" title="0">{
                return nil, errors.Wrap(err, "failed to extract istanbul extra from header")
        }</span>
        <span class="cov0" title="0">var validators []*istanbul.ValidatorData

        for i := range extra.AddedValidators </span><span class="cov0" title="0">{
                validator := &amp;istanbul.ValidatorData{
                        Address:      extra.AddedValidators[i],
                        BLSPublicKey: extra.AddedValidatorsPublicKeys[i],
                }

                validators = append(validators, validator)
        }</span>

        <span class="cov0" title="0">return validators, nil</span>
}

// AggregatePublicKeys merges all the validators public keys into one
// and returns it as an aggeragated public key
func (v *ValidatorSyncer) AggregatePublicKeys(validators []*istanbul.ValidatorData) (*bls.PublicKey, error) <span class="cov0" title="0">{
        var publicKeys []blscrypto.SerializedPublicKey
        for _, validator := range validators </span><span class="cov0" title="0">{
                publicKeys = append(publicKeys, validator.BLSPublicKey)
        }</span>

        <span class="cov0" title="0">publicKeyObjs := []*bls.PublicKey{}
        for _, publicKey := range publicKeys </span><span class="cov0" title="0">{
                publicKeyObj, err := bls.DeserializePublicKeyCached(publicKey[:])
                if err != nil </span><span class="cov0" title="0">{
                        return nil, err
                }</span>
                <span class="cov0" title="0">defer publicKeyObj.Destroy()
                publicKeyObjs = append(publicKeyObjs, publicKeyObj)</span>
        }
        <span class="cov0" title="0">apk, err := bls.AggregatePublicKeys(publicKeyObjs)
        if err != nil </span><span class="cov0" title="0">{
                return nil, err
        }</span>
        <span class="cov0" title="0">defer apk.Destroy()

        return apk, nil</span>
}

// ExtractValidatorsDiff extracts all values of the IstanbulExtra (aka diff) from the header
func (v *ValidatorSyncer) ExtractValidatorsDiff(num uint64, validators []*istanbul.ValidatorData) ([]*istanbul.ValidatorData, []*istanbul.ValidatorData, error) <span class="cov0" title="0">{
        header, err := v.client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(num))
        if err != nil </span><span class="cov0" title="0">{
                return nil, nil, errors.Wrap(err, "getting the block header by number failed")
        }</span>

        <span class="cov0" title="0">diff, err := types.ExtractIstanbulExtra(header)
        if err != nil </span><span class="cov0" title="0">{
                return nil, nil, errors.Wrap(err, "failed to extract istanbul extra from header")
        }</span>

        <span class="cov0" title="0">var addedValidators []*istanbul.ValidatorData
        for i, addr := range diff.AddedValidators </span><span class="cov0" title="0">{
                addedValidators = append(addedValidators, &amp;istanbul.ValidatorData{Address: addr, BLSPublicKey: diff.AddedValidatorsPublicKeys[i]})
        }</span>

        <span class="cov0" title="0">bitmap := diff.RemovedValidators.Bytes()
        var removedValidators []*istanbul.ValidatorData

        for _, i := range bitmap </span><span class="cov0" title="0">{
                removedValidators = append(removedValidators, validators[i])
        }</span>

        <span class="cov0" title="0">return addedValidators, removedValidators, nil</span>
}

func (v *ValidatorSyncer) start() error <span class="cov0" title="0">{
        return nil
}</span>

func (v *ValidatorSyncer) close() <span class="cov0" title="0">{
        v.client.Close()
}</span>

func (v *ValidatorSyncer) Sync(latestBlock *big.Int) error <span class="cov0" title="0">{
        return nil
}</span>
</pre>
		
		<pre class="file" id="file7" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package writer

import (
        "bytes"
        "math/big"

        "github.com/ChainSafe/chainbridge-celo/msg"
        "github.com/ethereum/go-ethereum/common"
        "github.com/ethereum/go-ethereum/crypto"
)

// constructErc20ProposalData returns the bytes to construct a proposal suitable for Erc20
func ConstructErc20ProposalData(amount []byte, recipient []byte) []byte <span class="cov8" title="1">{
        b := bytes.Buffer{}
        b.Write(common.LeftPadBytes(amount, 32)) // amount (uint256)
        recipientLen := big.NewInt(int64(len(recipient))).Bytes()
        b.Write(common.LeftPadBytes(recipientLen, 32))
        b.Write(recipient)
        return b.Bytes()
}</span>

// constructGenericProposalData returns the bytes to construct a generic proposal
func ConstructGenericProposalData(metadata []byte) []byte <span class="cov8" title="1">{
        data := bytes.Buffer{}
        metadataLen := big.NewInt(int64(len(metadata))).Bytes()
        data.Write(common.LeftPadBytes(metadataLen, 32)) // length of metadata (uint256)
        data.Write(metadata)
        return data.Bytes()
}</span>

// ConstructErc721ProposalData returns the bytes to construct a proposal suitable for Erc721
func ConstructErc721ProposalData(tokenId []byte, recipient []byte, metadata []byte) []byte <span class="cov8" title="1">{
        data := bytes.Buffer{}
        data.Write(common.LeftPadBytes(tokenId, 32))

        recipientLen := big.NewInt(int64(len(recipient))).Bytes()
        data.Write(common.LeftPadBytes(recipientLen, 32))
        data.Write(recipient)

        metadataLen := big.NewInt(int64(len(metadata))).Bytes()
        data.Write(common.LeftPadBytes(metadataLen, 32))
        data.Write(metadata)
        return data.Bytes()
}</span>

// CreateProposalDataHash constructs and returns proposal data hash
// https://github.com/ChainSafe/chainbridge-celo-solidity/blob/1fae9c66a07139c277b03a09877414024867a8d9/contracts/Bridge.sol#L452-L454
func CreateProposalDataHash(data []byte, handler common.Address, mp *msg.MerkleProof, sv *msg.SignatureVerification) common.Hash <span class="cov0" title="0">{
        b := bytes.NewBuffer(data)
        b.Write(handler.Bytes())
        b.Write(mp.RootHash[:])
        b.Write(mp.Key)
        b.Write(mp.Nodes)
        b.Write(sv.AggregatePublicKey)
        b.Write(sv.BlockHash[:])
        b.Write(sv.Signature)
        return crypto.Keccak256Hash(b.Bytes())
}</span>
</pre>
		
		<pre class="file" id="file8" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package writer

import (
        "math/big"

        "github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
        "github.com/ChainSafe/chainbridge-celo/chain/client"
        "github.com/ChainSafe/chainbridge-celo/chain/config"
        "github.com/ChainSafe/chainbridge-celo/msg"
        metrics "github.com/ChainSafe/chainbridge-utils/metrics/types"
        "github.com/ethereum/go-ethereum/accounts/abi/bind"
        "github.com/ethereum/go-ethereum/common"
        "github.com/ethereum/go-ethereum/core/types"
        "github.com/rs/zerolog/log"
)

var ProposalStatusPassed uint8 = 2
var ProposalStatusTransferred uint8 = 3
var ProposalStatusCancelled uint8 = 4
var BlockRetryLimit = 5

type writer struct {
        cfg            *config.CeloChainConfig
        client         ContractCaller
        bridgeContract Bridger
        stop           &lt;-chan struct{}
        sysErr         chan&lt;- error
        metrics        *metrics.ChainMetrics
}

type Bridger interface {
        GetProposal(opts *bind.CallOpts, originChainID uint8, depositNonce uint64, dataHash [32]byte) (Bridge.BridgeProposal, error)
        HasVotedOnProposal(opts *bind.CallOpts, arg0 *big.Int, arg1 [32]byte, arg2 common.Address) (bool, error)
        VoteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, resourceID [32]byte, dataHash [32]byte) (*types.Transaction, error)
        ExecuteProposal(opts *bind.TransactOpts, chainID uint8, depositNonce uint64, data []byte, resourceID [32]byte, signatureHeader []byte, aggregatePublicKey []byte, g1 []byte, hashedMessage [32]byte, rootHash [32]byte, key []byte, nodes []byte) (*types.Transaction, error)
}

type ContractCaller interface {
        client.LogFilterWithLatestBlock
        CallOpts() *bind.CallOpts
        Opts() *bind.TransactOpts
        LockAndUpdateOpts() error
        UnlockOpts()
        WaitForBlock(block *big.Int) error
}

// NewWriter creates and returns writer
func NewWriter(client ContractCaller, cfg *config.CeloChainConfig, stop &lt;-chan struct{}, sysErr chan&lt;- error, m *metrics.ChainMetrics) *writer <span class="cov8" title="1">{
        return &amp;writer{
                cfg:     cfg,
                client:  client,
                stop:    stop,
                sysErr:  sysErr,
                metrics: m,
        }
}</span>

// setContract adds the bound receiver bridgeContract to the writer
func (w *writer) SetBridge(bridge Bridger) <span class="cov8" title="1">{
        w.bridgeContract = bridge
}</span>

// ResolveMessage handles any given message based on type
// A bool is returned to indicate failure/success
// this should be ignored except for within tests.
func (w *writer) ResolveMessage(m *msg.Message) bool <span class="cov8" title="1">{
        log.Info().Str("type", string(m.Type)).Interface("src", m.Source).Interface("dst", m.Destination).Interface("nonce", m.DepositNonce).Str("rId", m.ResourceId.Hex()).Msg("Attempting to resolve message")
        var data []byte
        var handlerContract common.Address
        var err error
        switch m.Type </span>{
        case msg.FungibleTransfer:<span class="cov0" title="0">
                data, err = w.createERC20ProposalData(m)
                handlerContract = w.cfg.Erc20HandlerContract</span>
        case msg.NonFungibleTransfer:<span class="cov0" title="0">
                data, err = w.createErc721ProposalData(m)
                handlerContract = w.cfg.Erc721HandlerContract</span>
        case msg.GenericTransfer:<span class="cov0" title="0">
                data, err = w.createGenericDepositProposalData(m)
                handlerContract = w.cfg.GenericHandlerContract</span>
        default:<span class="cov8" title="1">
                log.Error().Str("type", string(m.Type)).Msg("Unknown message type received")
                return false</span>
        }
        <span class="cov0" title="0">if err != nil </span><span class="cov0" title="0">{
                log.Error().Err(err)
                return false
        }</span>
        <span class="cov0" title="0">dataHash := CreateProposalDataHash(data, handlerContract, m.MPParams, m.SVParams)

        if !w.shouldVote(m, dataHash) </span><span class="cov0" title="0">{
                return false
        }</span>
        // Capture latest block so when know where to watch from
        <span class="cov0" title="0">latestBlock, err := w.client.LatestBlock()
        if err != nil </span><span class="cov0" title="0">{
                log.Error().Err(err).Msg("unable to fetch latest block")
                return false
        }</span>

        // watch for execution event
        <span class="cov0" title="0">go w.watchThenExecute(m, data, dataHash, latestBlock)

        w.voteProposal(m, dataHash)

        return true</span>
}
</pre>
		
		<pre class="file" id="file9" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package writer

import (
        "context"
        "errors"
        "math/big"
        "time"

        "github.com/ChainSafe/chainbridge-celo/msg"
        utils "github.com/ChainSafe/chainbridge-celo/shared/ethereum"
        eth "github.com/ethereum/go-ethereum"
        ethcommon "github.com/ethereum/go-ethereum/common"
        "github.com/rs/zerolog/log"
)

// Number of blocks to wait for an finalization event
const ExecuteBlockWatchLimit = 100

// Time between retrying a failed tx
const TxRetryInterval = time.Second * 2

// Time between retrying a failed tx
const TxRetryLimit = 10

var ErrNonceTooLow = errors.New("nonce too low")
var ErrTxUnderpriced = errors.New("replacement transaction underpriced")
var ErrFatalTx = errors.New("submission of transaction failed")
var ErrFatalQuery = errors.New("query of chain state failed")

// proposalIsComplete returns true if the proposal state is either Passed, Transferred or Cancelled
func (w *writer) proposalIsComplete(srcId msg.ChainId, nonce msg.Nonce, dataHash ethcommon.Hash) bool <span class="cov8" title="1">{
        prop, err := w.bridgeContract.GetProposal(w.client.CallOpts(), uint8(srcId), uint64(nonce), dataHash)
        if err != nil </span><span class="cov8" title="1">{
                log.Error().Err(err).Msg("Failed to check proposal existence")
                return false
        }</span>
        <span class="cov8" title="1">return prop.Status == ProposalStatusPassed || prop.Status == ProposalStatusTransferred || prop.Status == ProposalStatusCancelled</span>
}

// proposalIsFinalized returns true if the proposal state is Transferred or Cancelled
func (w *writer) proposalIsFinalized(srcId msg.ChainId, nonce msg.Nonce, dataHash ethcommon.Hash) bool <span class="cov8" title="1">{
        prop, err := w.bridgeContract.GetProposal(w.client.CallOpts(), uint8(srcId), uint64(nonce), dataHash)

        if err != nil </span><span class="cov8" title="1">{
                log.Error().Err(err).Msg("Failed to check proposal existence")
                return false
        }</span>
        <span class="cov8" title="1">return prop.Status == ProposalStatusTransferred || prop.Status == ProposalStatusCancelled</span>
}

// hasVoted checks if this relayer has already voted
func (w *writer) hasVoted(srcId msg.ChainId, nonce msg.Nonce, dataHash ethcommon.Hash) bool <span class="cov8" title="1">{
        hasVoted, err := w.bridgeContract.HasVotedOnProposal(w.client.CallOpts(), utils.IDAndNonce(srcId, nonce), dataHash, w.client.Opts().From)

        if err != nil </span><span class="cov8" title="1">{
                log.Error().Err(err).Msg("Failed to check proposal existence")
                return false
        }</span>
        <span class="cov8" title="1">return hasVoted</span>
}

func (w *writer) shouldVote(m *msg.Message, dataHash ethcommon.Hash) bool <span class="cov8" title="1">{
        // Check if proposal has passed and skip if Passed or Transferred
        if w.proposalIsComplete(m.Source, m.DepositNonce, dataHash) </span><span class="cov8" title="1">{
                log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Proposal complete, not voting")
                return false
        }</span>

        // Check if relayer has previously voted
        <span class="cov8" title="1">if w.hasVoted(m.Source, m.DepositNonce, dataHash) </span><span class="cov8" title="1">{
                log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Relayer has already voted, not voting")
                return false
        }</span>
        <span class="cov8" title="1">return true</span>
}

func (w *writer) createERC20ProposalData(m *msg.Message) ([]byte, error) <span class="cov8" title="1">{
        log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Creating erc20 proposal")
        if len(m.Payload) != 2 </span><span class="cov8" title="1">{
                return nil, errors.New("malformed payload. Len  of payload should be 2")
        }</span>
        <span class="cov8" title="1">amount, ok := m.Payload[0].([]byte)
        if !ok </span><span class="cov8" title="1">{
                return nil, errors.New("wrong payloads amount format")
        }</span>

        <span class="cov8" title="1">recipient, ok := m.Payload[1].([]byte)
        if !ok </span><span class="cov8" title="1">{
                return nil, errors.New("wrong payloads recipient format")
        }</span>
        <span class="cov8" title="1">data := ConstructErc20ProposalData(amount, recipient)
        return data, nil</span>
}

func (w *writer) createErc721ProposalData(m *msg.Message) ([]byte, error) <span class="cov8" title="1">{
        log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Creating erc721 proposal")
        if len(m.Payload) != 3 </span><span class="cov8" title="1">{
                return nil, errors.New("malformed payload. Len  of payload should be 3")
        }</span>
        <span class="cov8" title="1">tokenID, ok := m.Payload[0].([]byte)
        if !ok </span><span class="cov8" title="1">{
                return nil, errors.New("wrong payloads tokenID format")
        }</span>
        <span class="cov8" title="1">recipient, ok := m.Payload[1].([]byte)
        if !ok </span><span class="cov8" title="1">{
                return nil, errors.New("wrong payloads recipient format")
        }</span>
        <span class="cov8" title="1">metadata, ok := m.Payload[2].([]byte)
        if !ok </span><span class="cov8" title="1">{
                return nil, errors.New("wrong payloads metadata format")
        }</span>
        <span class="cov8" title="1">return ConstructErc721ProposalData(tokenID, recipient, metadata), nil</span>
}

func (w *writer) createGenericDepositProposalData(m *msg.Message) ([]byte, error) <span class="cov8" title="1">{
        log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Creating generic proposal")
        if len(m.Payload) != 1 </span><span class="cov8" title="1">{
                return nil, errors.New("malformed payload. Len  of payload should be 1")
        }</span>
        <span class="cov8" title="1">metadata, ok := m.Payload[0].([]byte)
        if !ok </span><span class="cov8" title="1">{
                return nil, errors.New("unable to convert metadata to []byte")
        }</span>
        <span class="cov8" title="1">return ConstructGenericProposalData(metadata), nil</span>
}

// watchThenExecute watches for the latest block and executes once the matching finalized event is found
func (w *writer) watchThenExecute(m *msg.Message, data []byte, dataHash ethcommon.Hash, latestBlock *big.Int) <span class="cov8" title="1">{
        log.Info().Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("Watching for finalization event")

        // watching for the latest block, querying and matching the finalized event will be retried up to ExecuteBlockWatchLimit times
        for i := 0; i &lt; ExecuteBlockWatchLimit; i++ </span><span class="cov8" title="1">{
                select </span>{
                case &lt;-w.stop:<span class="cov0" title="0">
                        return</span>
                default:<span class="cov8" title="1">
                        // watch for the lastest block, retry up to BlockRetryLimit times
                        for waitRetrys := 0; waitRetrys &lt;= BlockRetryLimit; waitRetrys++ </span><span class="cov8" title="1">{
                                err := w.client.WaitForBlock(latestBlock)
                                if err != nil </span><span class="cov8" title="1">{
                                        log.Error().Err(err).Msg("Waiting for block failed")
                                        // Exit if retries exceeded
                                        if waitRetrys == BlockRetryLimit </span><span class="cov8" title="1">{
                                                log.Error().Err(err).Msg("Waiting for block retries exceeded, shutting down")
                                                w.sysErr &lt;- ErrFatalQuery
                                                return
                                        }</span>
                                } else<span class="cov8" title="1"> {
                                        break</span>
                                }
                        }

                        // query for logs
                        <span class="cov8" title="1">query := buildQuery(w.cfg.BridgeContract, utils.ProposalEvent, latestBlock, latestBlock)
                        evts, err := w.client.FilterLogs(context.Background(), query)
                        if err != nil </span><span class="cov8" title="1">{
                                log.Error().Err(err).Msg("Failed to fetch logs")
                                return
                        }</span>

                        // execute the proposal once we find the matching finalized event
                        <span class="cov8" title="1">for _, evt := range evts </span><span class="cov8" title="1">{
                                sourceId := evt.Topics[1].Big().Uint64()
                                depositNonce := evt.Topics[2].Big().Uint64()
                                status := evt.Topics[3].Big().Uint64()

                                if m.Source == msg.ChainId(sourceId) &amp;&amp;
                                        m.DepositNonce.Big().Uint64() == depositNonce &amp;&amp;
                                        utils.IsFinalized(uint8(status)) </span><span class="cov0" title="0">{
                                        w.executeProposal(m, data, dataHash)
                                        return
                                }</span> else<span class="cov8" title="1"> {
                                        log.Trace().Interface("src", sourceId).Interface("nonce", depositNonce).Msg("Ignoring event")
                                }</span>
                        }
                        <span class="cov8" title="1">log.Trace().Interface("block", latestBlock).Interface("src", m.Source).Interface("nonce", m.DepositNonce).Msg("No finalization event found in current block")
                        latestBlock = latestBlock.Add(latestBlock, big.NewInt(1))</span>
                }
        }
        <span class="cov8" title="1">log.Warn().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Block watch limit exceeded, skipping execution")</span>
}

// voteProposal submits a vote proposal
// a vote proposal will try to be submitted up to the TxRetryLimit times
func (w *writer) voteProposal(m *msg.Message, dataHash ethcommon.Hash) <span class="cov8" title="1">{
        for i := 0; i &lt; TxRetryLimit; i++ </span><span class="cov8" title="1">{
                select </span>{
                case &lt;-w.stop:<span class="cov0" title="0">
                        return</span>
                default:<span class="cov8" title="1">
                        // Checking first does proposal complete? If so, we do not need to vote for it
                        if w.proposalIsComplete(m.Source, m.DepositNonce, dataHash) </span><span class="cov8" title="1">{
                                log.Info().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Proposal voting complete on chain")
                                return
                        }</span>
                        <span class="cov8" title="1">err := w.client.LockAndUpdateOpts()
                        if err != nil </span><span class="cov8" title="1">{
                                log.Error().Err(err).Msg("Failed to update tx opts")
                                continue</span>
                        }

                        <span class="cov8" title="1">tx, err := w.bridgeContract.VoteProposal(
                                w.client.Opts(),
                                uint8(m.Source),
                                uint64(m.DepositNonce),
                                m.ResourceId,
                                dataHash,
                        )
                        w.client.UnlockOpts()
                        if err != nil </span><span class="cov8" title="1">{
                                if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() </span><span class="cov0" title="0">{
                                        log.Debug().Msg("Nonce too low, will retry")
                                        time.Sleep(TxRetryInterval)
                                        continue</span>
                                } else<span class="cov8" title="1"> {
                                        log.Warn().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Voting failed")
                                        time.Sleep(TxRetryInterval)
                                        continue</span>
                                }
                        }
                        <span class="cov8" title="1">log.Info().Str("tx", tx.Hash().Hex()).Interface("src", m.Source).Interface("depositNonce", m.DepositNonce).Msg("Submitted proposal vote")
                        return</span>
                }
        }
        <span class="cov8" title="1">log.Error().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Submission of Vote transaction failed")
        w.sysErr &lt;- ErrFatalTx</span>
}

// executeProposal executes the proposal
func (w *writer) executeProposal(m *msg.Message, data []byte, dataHash ethcommon.Hash) <span class="cov0" title="0">{
        for i := 0; i &lt; TxRetryLimit; i++ </span><span class="cov0" title="0">{
                select </span>{
                case &lt;-w.stop:<span class="cov0" title="0">
                        return</span>
                default:<span class="cov0" title="0">
                        err := w.client.LockAndUpdateOpts()
                        if err != nil </span><span class="cov0" title="0">{
                                log.Error().Err(err).Msg("Failed to update nonce")
                                return
                        }</span>

                        <span class="cov0" title="0">tx, err := w.bridgeContract.ExecuteProposal(
                                w.client.Opts(),
                                uint8(m.Source),
                                uint64(m.DepositNonce),
                                data,
                                m.ResourceId,
                                //
                                m.SVParams.Signature,
                                m.SVParams.AggregatePublicKey,
                                // TODO: Remove once G1 has been removed from contracts
                                []byte{},
                                m.SVParams.BlockHash,
                                m.MPParams.RootHash,
                                m.MPParams.Key,
                                m.MPParams.Nodes,
                        )
                        w.client.UnlockOpts()

                        if err == nil </span><span class="cov0" title="0">{
                                log.Info().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Str("tx", tx.Hash().Hex()).Msg("Submitted proposal execution")
                                return
                        }</span> else<span class="cov0" title="0"> if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() </span><span class="cov0" title="0">{
                                log.Error().Err(err).Msg("Nonce too low, will retry")
                                time.Sleep(TxRetryInterval)
                        }</span> else<span class="cov0" title="0"> {
                                log.Error().Err(err).Msg("Execution failed, proposal may already be complete")
                                time.Sleep(TxRetryInterval)
                        }</span>

                        // Verify proposal is still open for execution, tx will fail if we aren't the first to execute,
                        // but there is no need to retry
                        <span class="cov0" title="0">if w.proposalIsFinalized(m.Source, m.DepositNonce, dataHash) </span><span class="cov0" title="0">{
                                log.Info().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Proposal finalized on chain")
                                return
                        }</span>
                }
        }
        <span class="cov0" title="0">log.Error().Interface("source", m.Source).Interface("dest", m.Destination).Interface("nonce", m.DepositNonce).Msg("Submission of Execute transaction failed")
        w.sysErr &lt;- ErrFatalTx</span>
}

// buildQuery constructs a query for the bridgeContract by hashing sig to get the event topic
func buildQuery(contract ethcommon.Address, sig utils.EventSig, startBlock *big.Int, endBlock *big.Int) eth.FilterQuery <span class="cov8" title="1">{
        query := eth.FilterQuery{
                FromBlock: startBlock,
                ToBlock:   endBlock,
                Addresses: []ethcommon.Address{contract},
                Topics: [][]ethcommon.Hash{
                        {sig.GetTopic()},
                },
        }
        return query
}</span>
</pre>
		
		<pre class="file" id="file10" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package cfg

import (
        "encoding/json"
        "fmt"
        "os"
        "path/filepath"

        "github.com/ChainSafe/chainbridge-celo/flags"
        "github.com/ethereum/go-ethereum/log"
        "github.com/urfave/cli/v2"
)

const DefaultConfigPath = "./config.json"
const DefaultBlockTimeout = int64(180) // 3 minutes

type Config struct {
        Chains []RawChainConfig `json:"chains"`
}

// RawChainConfig is parsed directly from the config file and should be using to construct the core.ChainConfig
type RawChainConfig struct {
        Name     string            `json:"name"`
        Type     string            `json:"type"`
        Id       string            `json:"id"`       // ChainID
        Endpoint string            `json:"endpoint"` // url for rpc endpoint
        From     string            `json:"from"`     // address of key to use
        Opts     map[string]string `json:"opts"`
}

func NewConfig() *Config <span class="cov8" title="1">{
        return &amp;Config{
                Chains: []RawChainConfig{},
        }
}</span>

func (c *Config) ToJSON(file string) *os.File <span class="cov8" title="1">{
        var (
                newFile *os.File
                err     error
        )

        var raw []byte
        if raw, err = json.Marshal(*c); err != nil </span><span class="cov0" title="0">{
                log.Warn("error marshalling json", "err", err)
                os.Exit(1)
        }</span>

        <span class="cov8" title="1">newFile, err = os.Create(file)
        if err != nil </span><span class="cov0" title="0">{
                log.Warn("error creating config file", "err", err)
        }</span>
        <span class="cov8" title="1">_, err = newFile.Write(raw)
        if err != nil </span><span class="cov0" title="0">{
                log.Warn("error writing to config file", "err", err)
        }</span>

        <span class="cov8" title="1">if err := newFile.Close(); err != nil </span><span class="cov0" title="0">{
                log.Warn("error closing file", "err", err)
        }</span>
        <span class="cov8" title="1">return newFile</span>
}

func (c *Config) validate() error <span class="cov8" title="1">{
        for _, chain := range c.Chains </span><span class="cov8" title="1">{
                if chain.Type == "" </span><span class="cov8" title="1">{
                        return fmt.Errorf("required field chain.Type empty for chain %s", chain.Id)
                }</span>
                <span class="cov8" title="1">if chain.Endpoint == "" </span><span class="cov8" title="1">{
                        return fmt.Errorf("required field chain.Endpoint empty for chain %s", chain.Id)
                }</span>
                <span class="cov8" title="1">if chain.Name == "" </span><span class="cov8" title="1">{
                        return fmt.Errorf("required field chain.Name empty for chain %s", chain.Id)
                }</span>
                <span class="cov8" title="1">if chain.Id == "" </span><span class="cov0" title="0">{
                        return fmt.Errorf("required field chain.Id empty for chain %s", chain.Id)
                }</span>
                <span class="cov8" title="1">if chain.From == "" </span><span class="cov0" title="0">{
                        return fmt.Errorf("required field chain.From empty for chain %s", chain.Id)
                }</span>
        }
        <span class="cov8" title="1">return nil</span>
}

func GetConfig(ctx *cli.Context) (*Config, error) <span class="cov8" title="1">{
        var fig Config
        path := DefaultConfigPath
        if file := ctx.String(flags.ConfigFileFlag.Name); file != "" </span><span class="cov8" title="1">{
                path = file
        }</span>
        <span class="cov8" title="1">err := loadConfig(path, &amp;fig)
        if err != nil </span><span class="cov0" title="0">{
                log.Warn("err loading json file", "err", err.Error())
                return &amp;fig, err
        }</span>
        <span class="cov8" title="1">log.Debug("Loaded config", "path", path)
        err = fig.validate()
        if err != nil </span><span class="cov0" title="0">{
                return nil, err
        }</span>
        <span class="cov8" title="1">return &amp;fig, nil</span>
}

func loadConfig(file string, config *Config) error <span class="cov8" title="1">{
        ext := filepath.Ext(file)
        fp, err := filepath.Abs(file)
        if err != nil </span><span class="cov0" title="0">{
                return err
        }</span>

        <span class="cov8" title="1">log.Debug("Loading configuration", "path", filepath.Clean(fp))

        f, err := os.Open(filepath.Clean(fp))
        if err != nil </span><span class="cov0" title="0">{
                return err
        }</span>

        <span class="cov8" title="1">if ext == ".json" </span><span class="cov8" title="1">{
                if err = json.NewDecoder(f).Decode(&amp;config); err != nil </span><span class="cov0" title="0">{
                        return err
                }</span>
        } else<span class="cov0" title="0"> {
                return fmt.Errorf("unrecognized extention: %s", ext)
        }</span>

        <span class="cov8" title="1">return nil</span>
}
</pre>
		
		<pre class="file" id="file11" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package main

import (
        "os"

        "github.com/ChainSafe/chainbridge-celo/cmd"
        "github.com/ChainSafe/chainbridge-celo/flags"
        "github.com/rs/zerolog/log"
        "github.com/urfave/cli/v2"
)

var app = cli.NewApp()

var cliFlags = []cli.Flag{
        flags.ConfigFileFlag, // path to config file
        flags.VerbosityFlag,  // logger flag
        flags.KeystorePathFlag,
        flags.BlockstorePathFlag, // seems to be used only in tests
        flags.FreshStartFlag,     // start blocks from scratch. Used on chain initialization
        flags.LatestBlockFlag,    // latest block to start listen from. Used on chain initialization
        flags.MetricsFlag,
        flags.MetricsPort,
}

//
var generateFlags = []cli.Flag{
        flags.PasswordFlag,
        flags.Secp256k1Flag,
}

//
var devFlags = []cli.Flag{
        flags.TestKeyFlag,
}

var importFlags = []cli.Flag{
        flags.EthereumImportFlag,
        flags.PrivateKeyFlag,
        flags.Secp256k1Flag,
        flags.PasswordFlag,
}

var accountCommand = cli.Command{
        Name:  "accounts",
        Usage: "manage bridge keystore",
        Description: "The accounts command is used to manage the bridge keystore. \n" +
                "\tTo generate a new account (key type generated is determined on the flag passed in): chainbridge accounts generate\n" +
                "\tTo import a keystore file: chainbridge accounts import path/to/file\n" +
                "\tTo import a geth keystore file: chainbridge accounts import --ethereum path/to/file\n" +
                "\tTo import a private key file: chainbridge accounts import --privateKey private_key\n" +
                "\tTo list keys: chainbridge accounts list",
        Subcommands: []*cli.Command{
                {
                        Action: wrapHandler(handleGenerateCmd),
                        Name:   "generate",
                        Usage:  "generate bridge keystore, key type determined by flag",
                        Flags:  generateFlags,
                        Description: "The generate subcommand is used to generate the bridge keystore.\n" +
                                "\tIf no options are specified, a secp256k1 key will be made.",
                },
                {
                        Action: wrapHandler(handleImportCmd),
                        Name:   "import",
                        Usage:  "import bridge keystore",
                        Flags:  importFlags,
                        Description: "The import subcommand is used to import a keystore for the bridge.\n" +
                                "\tA path to the keystore must be provided\n" +
                                "\tUse --ethereum to import an ethereum keystore from external sources such as geth\n" +
                                "\tUse --privateKey to create a keystore from a provided private key.",
                },
                {
                        Action:      wrapHandler(handleListCmd),
                        Name:        "list",
                        Usage:       "list bridge keystore",
                        Description: "The list subcommand is used to list all of the bridge keystores.\n",
                },
        },
}

// init initializes CLI
func init() <span class="cov8" title="1">{
        app.Action = cmd.Run
        app.Copyright = "Copyright 2019 ChainSafe Systems Authors"
        app.Name = "chainbridge-celo"
        app.Usage = "ChainBridge-celo"
        app.Authors = []*cli.Author{{Name: "ChainSafe Systems 2020"}}
        app.Version = "0.0.1"
        app.EnableBashCompletion = true
        app.Commands = []*cli.Command{
                &amp;accountCommand,
        }

        app.Flags = append(app.Flags, cliFlags...)
        app.Flags = append(app.Flags, devFlags...)
}</span>

func main() <span class="cov0" title="0">{
        if err := app.Run(os.Args); err != nil </span><span class="cov0" title="0">{
                log.Error().Err(err).Msg("Start failed")
                os.Exit(1)
        }</span>
}
</pre>
		
		<pre class="file" id="file12" style="display: none">// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package router

import (
        "fmt"
        "sync"

        "github.com/ChainSafe/chainbridge-celo/msg"
        "github.com/rs/zerolog/log"
)

// Writer consumes a message and makes the requried on-chain interactions.
type MessageResolver interface {
        ResolveMessage(message *msg.Message) bool
}

// BaseRouter forwards messages from their source to their destination
type BaseRouter struct {
        registry map[msg.ChainId]MessageResolver
        lock     *sync.RWMutex
}

func NewRouter() *BaseRouter <span class="cov8" title="1">{
        return &amp;BaseRouter{
                registry: make(map[msg.ChainId]MessageResolver),
                lock:     &amp;sync.RWMutex{},
        }
}</span>

// Send passes a message to the destination Writer if it exists
func (r *BaseRouter) Send(msg *msg.Message) error <span class="cov8" title="1">{
        r.lock.Lock()
        defer r.lock.Unlock()

        log.Trace().Interface("src", msg.Source).Interface("dest", msg.Destination).Interface("nonce", msg.DepositNonce).Interface("rId", msg.ResourceId.Hex()).Msg("Routing message")
        w := r.registry[msg.Destination]
        if w == nil </span><span class="cov0" title="0">{
                return fmt.Errorf("unknown destination chainId: %d", msg.Destination)
        }</span>

        <span class="cov8" title="1">go w.ResolveMessage(msg)
        return nil</span>
}

// Register registers a Writer with a ChainId which BaseRouter.Send can then use to propagate messages
func (r *BaseRouter) Register(id msg.ChainId, w MessageResolver) <span class="cov8" title="1">{
        r.lock.Lock()
        defer r.lock.Unlock()
        log.Debug().Interface("id", id).Msg("Registering new chain in router")
        r.registry[id] = w
}</span>
</pre>
		
		</div>
	</body>
	<script>
	(function() {
		var files = document.getElementById('files');
		var visible;
		files.addEventListener('change', onChange, false);
		function select(part) {
			if (visible)
				visible.style.display = 'none';
			visible = document.getElementById(part);
			if (!visible)
				return;
			files.value = part;
			visible.style.display = 'block';
			location.hash = part;
		}
		function onChange() {
			select(files.value);
			window.scrollTo(0, 0);
		}
		if (location.hash != "") {
			select(location.hash.substr(1));
		}
		if (!visible) {
			select("file0");
		}
	})();
	</script>
</html>
