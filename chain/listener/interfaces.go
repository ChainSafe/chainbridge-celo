package listener

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type Bridge interface {
	ResourceIDToHandlerAddress(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error)
}
