package chain

import (
	"fmt"
	"math/big"

	"github.com/rs/zerolog/log"

	"github.com/pkg/errors"

	bridgeHandler "github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
	erc20Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC20Handler"
	erc721Handler "github.com/ChainSafe/chainbridge-celo/bindings/ERC721Handler"
	"github.com/ChainSafe/chainbridge-celo/bindings/GenericHandler"
	"github.com/ChainSafe/chainbridge-celo/chain/connection"
	"github.com/ChainSafe/chainbridge-celo/core"
	"github.com/ChainSafe/chainbridge-utils/blockstore"
	metrics "github.com/ChainSafe/chainbridge-utils/metrics/types"
	"github.com/ChainSafe/chainbridge-utils/msg"
)

type BlockDB interface {
	blockstore.Blockstorer
	TryLoadLatestBlock() (*big.Int, error)
}

// checkBlockstore queries the blockstore for the latest known block. If the latest block is
// greater than cfg.startBlock, then cfg.startBlock is replaced with the latest known block.
type Listener interface {
	SetRouter(*core.Router)
	Start() error
	SetContracts(bridge *bridgeHandler.Bridge, erc20Handler *erc20Handler.ERC20Handler, erc721Handler *erc721Handler.ERC721Handler, genericHandler *GenericHandler.GenericHandler)
	LatestBlock() *metrics.LatestBlock
}

type Writer interface {
	Start() error
	SetBridge(bc *bridgeHandler.Bridge)
	core.Writer
}

func InitializeChain(cc *CeloChainConfig, sysErr chan<- error, conn *connection.Connection, listener Listener, writer Writer, blockDB BlockDB) (*Chain, error) {

	stop := make(chan int)

	err := conn.EnsureHasBytecode(cc.BridgeContract)
	if err != nil {
		return nil, err
	}
	err = conn.EnsureHasBytecode(cc.Erc20HandlerContract)
	if err != nil {
		return nil, err
	}
	err = conn.EnsureHasBytecode(cc.GenericHandlerContract)
	if err != nil {
		return nil, err
	}

	bridgeContract, err := bridgeHandler.NewBridge(cc.BridgeContract, conn.Client())
	if err != nil {
		return nil, err
	}

	chainId, err := bridgeContract.ChainID(conn.CallOpts())
	if err != nil {
		return nil, err
	}

	if chainId != uint8(cc.ID) {
		return nil, errors.New(fmt.Sprintf("chainId (%d) and configuration chainId (%d) do not match", chainId, cc.ID))
	}

	erc20HandlerContract, err := erc20Handler.NewERC20Handler(cc.Erc20HandlerContract, conn.Client())
	if err != nil {
		return nil, err
	}

	erc721HandlerContract, err := erc721Handler.NewERC721Handler(cc.Erc721HandlerContract, conn.Client())
	if err != nil {
		return nil, err
	}

	genericHandlerContract, err := GenericHandler.NewGenericHandler(cc.GenericHandlerContract, conn.Client())
	if err != nil {
		return nil, err
	}

	if cc.LatestBlock {
		curr, err := conn.LatestBlock()
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
		stop:     stop,
	}, nil
}

type Chain struct {
	cfg      *CeloChainConfig // The config of the chain
	listener Listener         // The listener of this chain
	writer   Writer           // The writer of the chain
	stop     chan<- int
}

func (c *Chain) SetRouter(r *core.Router) {
	r.Listen(c.cfg.ID, c.writer)
	c.listener.SetRouter(r)
}

func (c *Chain) Start() error {
	err := c.listener.Start()
	if err != nil {
		return err
	}

	err = c.writer.Start()
	if err != nil {
		return err
	}

	log.Debug().Msg("Chain started!")
	return nil
}

func (c *Chain) ID() msg.ChainId {
	return c.cfg.ID
}

func (c *Chain) Name() string {
	return c.cfg.Name
}

func (c *Chain) LatestBlock() *metrics.LatestBlock {
	return c.listener.LatestBlock()
}

// Stop signals to any running routines to exit
func (c *Chain) Stop() {
	close(c.stop)
	// TODO not forget add conn close on end conn users
	//if c.conn != nil {
	//	c.conn.Close()
	//}
}
