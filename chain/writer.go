package chain

import (
	"github.com/ChainSafe/chainbridge-utils/core"
	metrics "github.com/ChainSafe/chainbridge-utils/metrics/types"
	"github.com/ChainSafe/chainbridge-utils/msg"
	"github.com/ChainSafe/log15"
	"github.com/ChainSafe/chainbridge-celo/bindings/Bridge"
)

var _ core.Writer = &writer{}

type writer struct {
	cfg           Config
	conn          Connection
	bridgeContract *Bridge.Bridge
	log           log15.Logger
	stop          <-chan int
	sysErr        chan<- error
	metrics       *metrics.ChainMetrics
}

// NewWriter creates and returns writer
func NewWriter(conn Connection, cfg *Config, log log15.Logger, stop <-chan int, sysErr chan<- error, m *metrics.ChainMetrics) *writer {
	return &writer{
		cfg:     *cfg,
		conn:    conn,
		log:     log,
		stop:    stop,
		sysErr:  sysErr,
		metrics: m,
	}
}

func (w *writer) start() error {
	w.log.Debug("Starting celo writer...")
	return nil
}

// setContract adds the bound receiver bridgeContract to the writer
func (w *writer) setContract(bridge *Bridge.Bridge) {
	w.bridgeContact = bridge
}

// ResolveMessage handles any given message based on type
// A bool is returned to indicate failure/success
// this should be ignored except for within tests.
func (w *writer) ResolveMessage(m msg.Message) bool {
	w.log.Info("Attempting to resolve message", "type", m.Type, "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce, "rId", m.ResourceId.Hex())

	switch m.Type {
	case msg.FungibleTransfer:
		return w.createErc20Proposal(m)
	case msg.NonFungibleTransfer:
		return w.createErc721Proposal(m)
	case msg.GenericTransfer:
		return w.createGenericDepositProposal(m)
	default:
		w.log.Error("Unknown message type received", "type", m.Type)
		return false
	}
}
