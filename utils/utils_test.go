package utils

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	blscrypto "github.com/ethereum/go-ethereum/crypto/bls"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestGetFunctionBytes(t *testing.T) {
	resb := GetSolidityFunctionSig("submitAsk(uint256,uint256)")
	ress := common.Bytes2Hex(resb[:])
	if ress != "7288a28a" {
		t.Fatal(fmt.Sprintf("Result is %s", ress))
	}

}

func TestUserAmountToReal(t *testing.T) {
	amount := big.NewInt(123)
	decimal := big.NewInt(10)
	res, err := UserAmountToWei(amount.String(), decimal)
	if err != nil {
		t.Fatal(err.Error())
	}
	if res.Cmp(big.NewInt(1230000000000)) != 0 {
		t.Fatal(res.String())
	}
}

func TestUserAmountToRealFloat(t *testing.T) {
	amount := big.NewFloat(1.2345)
	decimal := big.NewInt(5)
	res, err := UserAmountToWei(amount.String(), decimal)
	if err != nil {
		t.Fatal(err.Error())
	}
	if res.Cmp(big.NewInt(123450)) != 0 {
		t.Fatal(res.String())

	}
}

func TestWeiAmountToUser(t *testing.T) {
	amount := big.NewInt(100000001)
	decimal := big.NewInt(5)
	res, err := WeiAmountToUser(amount, decimal)
	if err != nil {
		t.Fatal(err.Error())
	}
	if res.Text('f', int(decimal.Int64())) != "1000.00001" {
		t.Fatal(res.Text('f', int(decimal.Int64())))
	}
}

func TestRlpEncodeHeader(t *testing.T) {
	// init new header with sample data
	header, err := generateBlockHeader()
	if err != nil {
		t.Fatalf("error generating test block header: %v", err)
	}

	// init new block with custom header
	block := types.NewBlockWithHeader(header)

	// encode copied header into local byte slice variable
	rlpEncodedHeader, err := rlp.EncodeToBytes(block.Header())
	if err != nil {
		t.Fatal(err.Error())
	}

	sampleRlpEncodedHeader := "f902bfa00000000000000000000000000000000000000000000000000000000000000000940000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000001a00000000000000000000000000000000000000000000000000000000000000002a00000000000000000000000000000000000000000000000000000000000000003b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007b8080b9011d0000000000000000000000000000000000000000000000000000000000000000f8fbea9444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212f8c4b86035b46d6f783958831a85deb6ad15d4daf428cbfafd7c59845e17941681dbc09bcf4fc516ffc3c89979c678075033ea006b749d100a0200ff53527570a45239388a20336ea19278f334ee2c3d2a383be4f6769e4be31ff9ae3406d7a642c44700b8601c5303fab0c8bd479422c6199b6a3a6e05f129fc815481732da91f408cc0e5229abf6395410244d7caac41a9b074d80015ffc2ce2ce49c678942460d0af93911e0caaf300921c6ec49937ba6e4aff16a0441457ff6b4795e157c828344ab3d018080c3808080c3808080"

	if strings.TrimPrefix(hexutil.Encode(rlpEncodedHeader), "0x") != sampleRlpEncodedHeader {
		t.Fatal("rlp encoded headers do not match")
	}
}

// borrowed and modified from validatorsync/sync_test.go
func generateBlockHeader() (*types.Header, error) {
	testKey2, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292")
	testKey3, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f293")
	blsPK2, _ := blscrypto.ECDSAToBLS(testKey2)
	blsPK3, _ := blscrypto.ECDSAToBLS(testKey3)
	pubKey2, _ := blscrypto.PrivateToPublic(blsPK2)
	pubKey3, _ := blscrypto.PrivateToPublic(blsPK3)
	extra, err := rlp.EncodeToBytes(&types.IstanbulExtra{
		AddedValidators: []common.Address{
			common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a")),
			common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212")),
		},
		AddedValidatorsPublicKeys: []blscrypto.SerializedPublicKey{
			pubKey2,
			pubKey3,
		},
		RemovedValidators:    big.NewInt(0),
		Seal:                 []byte{},
		AggregatedSeal:       types.IstanbulAggregatedSeal{},
		ParentAggregatedSeal: types.IstanbulAggregatedSeal{},
	})
	if err != nil {
		return nil, err
	}
	h := &types.Header{
		Number:      big.NewInt(123),
		Root:        common.HexToHash("0x1"),
		TxHash:      common.HexToHash("0x2"),
		ReceiptHash: common.HexToHash("0x3"),
		Extra:       append(make([]byte, types.IstanbulExtraVanity), extra...),
	}
	return h, nil
}

// TestPrepareAPKForContract tests PrepareAPKForContract to ensure
// that it properly encodes APK for use within a contract
func TestPrepareAPKForContract(t *testing.T) {}

// TestPrepareSignatureForContract tests PrepareSignatureForContract to ensure
// that it properly encodes a SignatureVerification.Signature for use within
// a contract
func TestPrepareSignatureForContract(t *testing.T) {}
