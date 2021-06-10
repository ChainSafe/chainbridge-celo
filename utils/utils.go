package utils

import (
	"bytes"
	"fmt"
	gomath "math"
	"math/big"
	"reflect"
	"strings"

	"github.com/celo-org/celo-bls-go/bls"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
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

// UserAmountToWei converts decimal user friendly representation of token amount to 'Wei' representation with provided amount of decimal places
// eg UserAmountToWei(1, 5) => 100000
func UserAmountToWei(amount string, decimal *big.Int) (*big.Int, error) {
	amountFloat, ok := big.NewFloat(0).SetString(amount)
	if !ok {
		return nil, errors.New("wrong amount format")
	}
	ethValueFloat := new(big.Float).Mul(amountFloat, big.NewFloat(gomath.Pow10(int(decimal.Int64()))))
	ethValueFloatString := strings.Split(ethValueFloat.Text('f', int(decimal.Int64())), ".")

	i, ok := big.NewInt(0).SetString(ethValueFloatString[0], 10)
	if !ok {
		return nil, errors.New(ethValueFloat.Text('f', int(decimal.Int64())))
	}

	return i, nil
}

func WeiAmountToUser(amount *big.Int, decimals *big.Int) (*big.Float, error) {
	amountFloat, ok := big.NewFloat(0).SetString(amount.String())
	if !ok {
		return nil, errors.New("wrong amount format")
	}
	return new(big.Float).Quo(amountFloat, big.NewFloat(gomath.Pow10(int(decimals.Int64())))), nil
}

func ConstructErc20DepositData(destRecipient []byte, amount *big.Int) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(amount, 32)...)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...)
	data = append(data, destRecipient...)
	return data
}

// constructErc20Data constructs the data field to be passed into an erc721 deposit call
func ConstructErc721DepositData(tokenId *big.Int, destRecipient []byte) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(tokenId, 32)...)                               // Resource Id + Token Id
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...) // Length of recipient
	data = append(data, destRecipient...)                                                  // Recipient

	return data
}

func ConstructGenericDepositData(metadata []byte) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 32)...)
	data = append(data, metadata...)
	return data
}

// TODO:
// move all below to new package

// borrowed from Celo
// https://github.com/celo-org/celo-bls-go/blob/kobigurk/arkworks/examples/utils/utils.go#L8-L13
func ReverseAnyAndPad(s []byte) []byte {
	s = ReverseAny(s)
	padding := make([]byte, FIELD_SIZE_IN_CONTRACT-(len(s)%FIELD_SIZE_IN_CONTRACT))
	z := append(padding, s...)
	return z
}

// borrowed from Celo
// https://github.com/celo-org/celo-bls-go/blob/kobigurk/arkworks/examples/utils/utils.go#L15-L24
func ReverseAny(s []byte) []byte {
	z := make([]byte, len(s))
	copy(z, s)
	n := reflect.ValueOf(z).Len()
	swap := reflect.Swapper(z)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
	return z
}

// borrowed from Celo
// https://github.com/celo-org/celo-bls-go/blob/kobigurk/arkworks/examples/utils/utils.go#L15-L24
const FIELD_SIZE = 48
const FIELD_SIZE_IN_CONTRACT = 32

// RlpEncodeHeader is method to RLP encode data stored in a block header
func RlpEncodeHeader(header *types.Header) ([]byte, error) {
	// deep copy of header
	newHeader := types.CopyHeader(header)

	// encode copied header into byte slice
	rlpEncodedHeader, err := rlp.EncodeToBytes(newHeader)
	if err != nil {
		// return empty byte slice, error
		return []byte{}, fmt.Errorf("error encoding header: %w", err)
	}

	return rlpEncodedHeader, nil
}

// PrepareAPKForContract properly encodes APK for use within a contract
// NOTE: uses new functionality from celo-bls-go PR #23
// https://github.com/celo-org/celo-bls-go/examples/utils
func PrepareAPKForContract(apk []byte) ([]byte, error) {
	// registration required for celo-bls-go package
	bls.InitBLSCrypto()

	// init new byte slice to hold newly encoded APK
	encodedAPK := make([]byte, 0)

	// deserialize public key
	key, err := bls.DeserializePublicKey(apk)
	if err != nil {
		return encodedAPK, fmt.Errorf("could not deserialize public key: %w", err)
	}

	// serialize uncompressed data
	// new functionality from celo-bls-go PR #23
	// https://github.com/celo-org/celo-bls-go/pull/23
	encodedData, err := key.SerializeUncompressed()
	if err != nil {
		return encodedAPK, fmt.Errorf("could not serialize data: %w", err)
	}

	// new functionality from celo-bls-go PR #23
	// https://github.com/celo-org/celo-bls-go/examples/utils
	// https://github.com/celo-org/celo-bls-go/examples/prepare_for_contract/prepare_for_contract.go#L23-L35
	encodedDataPart1 := encodedData[0:FIELD_SIZE]
	encodedDataPart1 = ReverseAnyAndPad(encodedDataPart1)
	encodedDataPart2 := encodedData[FIELD_SIZE : 2*FIELD_SIZE]
	encodedDataPart2 = ReverseAnyAndPad(encodedDataPart2)
	encodedDataPart3 := encodedData[2*FIELD_SIZE : 3*FIELD_SIZE]
	encodedDataPart3 = ReverseAnyAndPad(encodedDataPart3)
	encodedDataPart4 := encodedData[3*FIELD_SIZE : 4*FIELD_SIZE]
	encodedDataPart4 = ReverseAnyAndPad(encodedDataPart4)

	// append encoded data to APK byte slice
	encodedAPK = append(encodedAPK, encodedDataPart1...)
	encodedAPK = append(encodedAPK, encodedDataPart2...)
	encodedAPK = append(encodedAPK, encodedDataPart3...)
	encodedAPK = append(encodedAPK, encodedDataPart4...)

	return encodedAPK, nil
}

// PrepareSignatureForContract properly encodes Signature field within
// the SignatureVerification struct to be used within a contract
// NOTE: uses new functionality from celo-bls-go PR #23
// https://github.com/celo-org/celo-bls-go/examples/utils
func PrepareSignatureForContract(signature []byte) ([]byte, error) {
	// registration required for celo-bls-go package
	bls.InitBLSCrypto()

	// init new byte slice to hold newly encoded signature
	encodedSignature := make([]byte, 0)

	// deserialize signature
	key, err := bls.DeserializeSignature(signature)
	if err != nil {
		return encodedSignature, fmt.Errorf("could not deserialize public key: %w", err)
	}

	// serialize uncompressed data
	// new functionality from celo-bls-go PR #23
	// https://github.com/celo-org/celo-bls-go/pull/23
	encodedData, err := key.SerializeUncompressed()
	if err != nil {
		return encodedSignature, fmt.Errorf("could not serialize data: %w", err)
	}

	// new functionality from celo-bls-go PR #23
	// https://github.com/celo-org/celo-bls-go/examples/utils
	// https://github.com/celo-org/celo-bls-go/examples/prepare_for_contract/prepare_for_contract.go#L23-L3
	encodedDataPart1 := encodedData[0:FIELD_SIZE]
	encodedDataPart1 = ReverseAnyAndPad(encodedDataPart1)
	encodedDataPart2 := encodedData[FIELD_SIZE : 2*FIELD_SIZE]
	encodedDataPart2 = ReverseAnyAndPad(encodedDataPart2)

	// append encoded data to encoded signature byte slice
	encodedSignature = append(encodedSignature, encodedDataPart1...)
	encodedSignature = append(encodedSignature, encodedDataPart2...)

	return encodedSignature, nil
}
