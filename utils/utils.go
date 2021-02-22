package utils

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"math/big"
)

type EventSig string

func (es EventSig) GetTopic() common.Hash {
	return crypto.Keccak256Hash([]byte(es))
}

func IsActive(status uint8) bool {
	return ProposalStatus(status) == Active
}

func IsPassed(status uint8) bool {
	return ProposalStatus(status) == Passed
}

func IsExecuted(status uint8) bool {
	return ProposalStatus(status) == Executed
}

func SliceTo32Bytes(in []byte) [32]byte {
	var res [32]byte
	copy(res[:], in)
	return res
}

func SliceTo4Bytes(in []byte) [4]byte {
	var res [4]byte
	copy(res[:], in)
	return res
}

func GetSolidityFunctionSig(in string) [4]byte {
	res := crypto.Keccak256(bytes.NewBufferString(in).Bytes())
	return SliceTo4Bytes(res)
}

func UserAmountToReal(amount string, decimal *big.Int) (*big.Int, error) {
	amountFloat, ok := big.NewFloat(0).SetString(amount)
	if !ok {
		return nil, errors.New("wrong amount format")
	}
	//deimcalFloat, ok := big.NewFloat(0).SetString(decimal.String())
	//if !ok {
	//	return nil, errors.New("wrong decimal format")
	//}
	powerTo := big.NewInt(0).Exp(big.NewInt(10), decimal, nil)

	powerToFloat, ok := big.NewFloat(0).SetString(powerTo.String())
	if !ok {
		return nil, errors.New("wrong decimal format")
	}
	res := big.NewFloat(0)
	res.Mul(amountFloat, powerToFloat)
	resInt, _ := big.NewInt(0).SetString(res.String(), 10)

	return resInt, nil
}

//
//func RealAmountToUser(amount *big.Int, decimals *big.Int, result *big.Float) {
//	one := big.NewFloat(1)
//	oneDivisor := big.NewInt(0)
//	r := big.NewInt(0)
//	oneDivisor.Exp(bigTen, decimals, nil).String()
//	log.Debug().Msgf("%s", oneDivisor.String())
//	log.Debug().Msgf("%s", r.Div(one, oneDivisor).())
//	result.Mul(amount, one.Div(one, bigTen.Exp(bigTen, decimals, nil)))
//}
