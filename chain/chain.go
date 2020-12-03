package chain

import (
	"context"
	"fmt"
	"math/big"

	bridgeHandler "github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	erc20Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
	erc721Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
	"github.com/ChainSafe/chainbridge-celo/chain/client"
	"github.com/ChainSafe/chainbridge-utils/blockstore"
	"github.com/ChainSafe/chainbridge-utils/msg"
	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type BlockDB interface {
	blockstore.Blockstorer
	TryLoadLatestBlock() (*big.Int, error)
}

// checkBlockstore queries the blockstore for the latest known block. If the latest block is
// greater than cfg.startBlock, then cfg.startBlock is replaced with the latest known block.
type Listener interface {
	StartPollingBlocks() error
	SetContracts(bridge *bridgeHandler.Bridge, erc20Handler *erc20Handler.ERC20Handler, erc721Handler *erc721Handler.ERC721Handler, genericHandler *GenericHandler.GenericHandler)
	//LatestBlock() *metrics.LatestBlock
}

type Writer interface {
	SetBridge(bc *bridgeHandler.Bridge)
}

type ContractBackendWithBlockFinder interface {
	bind.ContractBackend
	LatestBlock() (*big.Int, error)
}

type LogFilterWithLatestBlock interface {
	FilterLogs(ctx context.Context, q eth.FilterQuery) ([]types.Log, error)
	LatestBlock() (*big.Int, error)
}

type Chain struct {
	cfg      *CeloChainConfig // The config of the chain
	listener Listener         // The listener of this chain
	writer   Writer           // The writer of the chain
	client   *client.Client
	stopChn  <-chan struct{}
}

func InitializeChain(cc *CeloChainConfig, c *client.Client, listener Listener, writer Writer, stopChn <-chan struct{}) (*Chain, error) {

	bridgeContract, err := bridgeHandler.NewBridge(cc.BridgeContract, c)
	if err != nil {
		return nil, err
	}

	//TODO
	chainId, err := bridgeContract.ChainID(c.CallOpts())
	if err != nil {
		return nil, err
	}
	if chainId != uint8(cc.ID) {
		return nil, errors.New(fmt.Sprintf("chainId (%d) and configuration chainId (%d) do not match", chainId, cc.ID))
	}

	erc20HandlerContract, err := erc20Handler.NewERC20Handler(cc.Erc20HandlerContract, c)
	if err != nil {
		return nil, err
	}
	erc721HandlerContract, err := erc721Handler.NewERC721Handler(cc.Erc721HandlerContract, c)
	if err != nil {
		return nil, err
	}
	genericHandlerContract, err := GenericHandler.NewGenericHandler(cc.GenericHandlerContract, c)
	if err != nil {
		return nil, err
	}
	if cc.LatestBlock {
		curr, err := c.LatestBlock()
		if err != nil {
			return nil, err
		}
		cc.StartBlock = curr
	}
	listener.SetContracts(bridgeContract, erc20HandlerContract, erc721HandlerContract, genericHandlerContract)
	writer.SetBridge(bridgeContract)

	return &Chain{
		cfg:      cc,
		writer:   writer,
		listener: listener,
		stopChn:  stopChn,
	}, nil
}

func (c *Chain) Start() error {
	err := c.listener.StartPollingBlocks()
	if err != nil {
		return err
	}
	go func() {
		<-c.stopChn
		if c.client != nil {
			c.client.Close()
		}
	}()
	log.Debug().Msg("Chain started!")
	return nil
}

func (c *Chain) ID() msg.ChainId {
	return c.cfg.ID
}

func (c *Chain) Name() string {
	return c.cfg.Name
}

//func (c *Chain) LatestBlock() *metrics.LatestBlock {
//	return c.listener.LatestBlock()
//}