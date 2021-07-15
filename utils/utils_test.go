//nolint
package utils

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/common/hexutil"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/crypto"
	blscrypto "github.com/celo-org/celo-blockchain/crypto/bls"
	"github.com/celo-org/celo-blockchain/rlp"
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
		t.Fatalf("error encoding header to bytes: %v", err)
	}

	sampleRlpEncodedHeader := "f902bfa00000000000000000000000000000000000000000000000000000000000000000940000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000001a00000000000000000000000000000000000000000000000000000000000000002a00000000000000000000000000000000000000000000000000000000000000003b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000007b8080b9011d0000000000000000000000000000000000000000000000000000000000000000f8fbea9444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212f8c4b86035b46d6f783958831a85deb6ad15d4daf428cbfafd7c59845e17941681dbc09bcf4fc516ffc3c89979c678075033ea006b749d100a0200ff53527570a45239388a20336ea19278f334ee2c3d2a383be4f6769e4be31ff9ae3406d7a642c44700b8601c5303fab0c8bd479422c6199b6a3a6e05f129fc815481732da91f408cc0e5229abf6395410244d7caac41a9b074d80015ffc2ce2ce49c678942460d0af93911e0caaf300921c6ec49937ba6e4aff16a0441457ff6b4795e157c828344ab3d018080c380807bc3808080"

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
		RemovedValidators: big.NewInt(0),
		Seal:              []byte{},
		AggregatedSeal: types.IstanbulAggregatedSeal{
			Round: big.NewInt(123),
		},
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
func TestPrepareAPKForContract(t *testing.T) {
	apk := []byte("aab506de1ef9b0df75f202b0813904e08d99ba0dbbf2084c3a983d9190c41f5f773489a6ee530da67d517d3151805101860dfdb8d7d72d768643af1b07a468f93ba1e08edb2a7f22c85bad3c2c02545f036647f11ce63eed3bd44e2cc080c480")

	// encode APK
	preparedApk, err := PrepareAPKForContract(apk)

	if err != nil {
		t.Fatalf("could not prepare APK for contract: %v", err)
	}

	encodedApk := []byte("0000000000000000000000000000000001518051317d517da60d53eea68934775f1fc490913d983a4c08f2bb0dba998de0043981b002f275dfb0f91ede06b5aa0000000000000000000000000000000000c480c02c4ed43bed3ee61cf14766035f54022c3cad5bc8227f2adb8ee0a13bf968a4071baf4386762dd7d7b8fd0d8600000000000000000000000000000000002f95dc47ed1188c0cc5762cf3fe3de3cbd0584eedcda9581d78dbc44dbbbc8f773003f3d024c94a04b099158d677bf0000000000000000000000000000000000eca490fa6b49e659f7e143bf06e230b6417cd7caaa65ae2e9b6a610b84559e9244879348c1a9ac06259f6fdba605e5")

	result := bytes.Compare(preparedApk, encodedApk)

	if result != 0 {
		t.Fatal("preparedAPK != encodedAPK; bytes do not match")
	}
}

// TestPrepareSignatureForContract tests PrepareSignatureForContract to ensure
// that it properly encodes a SignatureVerification.Signature for use within
// a contract
func TestPrepareSignatureForContract(t *testing.T) {
	// sample signature
	signature := []byte("063c39b0d4f7fa61cef3b97fe6705f02ffc209a4f3c91442588f15a884e25e68fb63cc60215ac6912f67b9fee9276501")

	preparedSignature, err := PrepareSignatureForContract(signature)
	if err != nil {
		t.Fatalf("could not prepare signature for contract: %v", err)
	}

	encodedSignature := []byte("00000000000000000000000000000000016527e9feb9672f91c65a2160cc63fb685ee284a8158f584214c9f3a409c2ff025f70e67fb9f3ce61faf7d4b0393c060000000000000000000000000000000000a3a8fc6a0258b5bba7bb306d8e7abc147999201f61df1e65f1adbd6ba7b6037e08e966178165ba8f28101a94574853")

	result := bytes.Compare(preparedSignature, encodedSignature)

	if result != 0 {
		t.Fatal("preparedSignature != encodedSignature; bytes do not match")
	}
}

// TestCommitedSealSuffix tests CommitedSealSuffix to ensure that it
// properly creates a CommitedSealSuffix
func TestCommitedSealSuffix(t *testing.T) {
	// generate block header with dummy data
	header, err := generateBlockHeader()
	if err != nil {
		t.Fatalf("could not generate block header: %v", err)
	}

	// init new block with custom header
	block := types.NewBlockWithHeader(header)

	// extra Istanbul extra data from block header
	extra, err := types.ExtractIstanbulExtra(block.Header())
	if err != nil {
		t.Fatalf("could not extract Istanbul extra data: %v", err)
	}

	// generate commited seal suffix from round
	commitedSealSuffix := CommitedSealSuffix(extra.AggregatedSeal.Round)

	// init new byte slice to hold dummy data
	preparedSuffixSlice := make([]byte, 0)

	// append first element (round)
	preparedSuffixSlice = append(preparedSuffixSlice, []byte{byte(123)}...)

	// append second element (istanbul msg commit)
	preparedSuffixSlice = append(preparedSuffixSlice, []byte{byte(2)}...)

	result := bytes.Compare(preparedSuffixSlice, commitedSealSuffix)

	if result != 0 {
		t.Fatal("preparedSuffix != CommitedSealSuffix; bytes do not match")
	}
}

// TestCommitedSealPrefix tests CommitedSealPrefix to ensure that it
// properly creates a CommitedSealPrefix
func TestCommitedSealPrefix(t *testing.T) {
	// generate concatenated blockHashAndSuffix from dummy header + block data
	blockHashAndSuffix, err := generateBlockHashAndSuffix()
	if err != nil {
		t.Fatalf("could not generate blockHashAndSuffix: %v", err)
	}

	// resulting prefix from sample
	preparedPrefix := []byte("03")

	// derive prefix
	commitedSealPrefix, err := CommitedSealPrefix(blockHashAndSuffix)
	if err != nil {
		t.Fatalf("could not derive CommitedSealPrefix: %v", err)
	}

	// match bytes of results
	result := bytes.Compare(preparedPrefix, commitedSealPrefix[:])

	if result != 0 {
		t.Fatal("preparedPrefix != CommitedSealPrefix; bytes do not match")
	}
}

// TestCommitedSealHints tests CommitedSealHints to ensure that it
// properly creates a CommitedSealHints
func TestCommitedSealHints(t *testing.T) {
	// sample data
	preparedBlockHashAndSuffix := []byte("8972a43b239781a90cc8cd56e6ea0b0c62a02ede612d58d71d1fd5211ed65c2402")

	// generate concatenated blockHashAndSuffix from dummy header + block data
	blockHashAndSuffix, err := generateBlockHashAndSuffix()
	if err != nil {
		t.Fatalf("could not generate blockHashAndSuffix: %v", err)
	}

	result := bytes.Compare(preparedBlockHashAndSuffix, blockHashAndSuffix)

	if result != 0 {
		t.Fatal("preparedBlockHashAndSuffix != blockHashAndSuffix; bytes do not match")
	}

	// data generated by referencing Celo example
	// https://github.com/ChainSafe/chainbridge-celo/celo-bls-go/examples/prepare_for_contract/prepare_for_contract.go#L53-L70
	preparedHints := []byte("000000000000000000000000000000000057a618735d7522cbe64234384ecc00fb53a9425455ec896abd45636c88e249413b261915a5cb6aff21589dff42a1cc000000000000000000000000000000000156942da4679bc7fa54c38c34527d3a1ecf30b0ac9f2705b4361ccc4d8065b6d5d0372b1a5a349585e7676200bd5e35")

	commitedSealHints, err := CommitedSealHints(blockHashAndSuffix)
	if err != nil {
		t.Fatalf("could not generate commitedSealHints: %v", err)
	}

	result = bytes.Compare(preparedHints, commitedSealHints)

	if result != 0 {
		t.Fatal("preparedPrefix != CommitedSealPrefix; bytes do not match")
	}
}

// TestConcatBlockHashAndCommitedSealSuffix tests
// ConcatBlockHashAndCommitedSealSuffix to ensure that it properly
// concatenates the blockhash with the commited seal suffix
func TestConcatBlockHashAndCommitedSealSuffix(t *testing.T) {
	// init new header with dummy data
	header, err := generateBlockHeader()
	if err != nil {
		t.Fatalf("error generating test block header: %v", err)
	}

	// init new block with custom header
	block := types.NewBlockWithHeader(header)

	// init new byte slice to hold resulting Commited Seal Hints
	blockHashAndSuffix := make([]byte, 0)

	// append block hash bytes to commited seal hints
	blockHashAndSuffix = append(blockHashAndSuffix, block.Hash().Bytes()...)

	// extract Istanbul Extra data from block header
	extra, err := types.ExtractIstanbulExtra(block.Header())
	if err != nil {
		t.Fatalf("error extracting Istanbul extra data: %v", err)
	}

	// derive commited seal suffix data
	commitedSealSuffix := CommitedSealSuffix(extra.AggregatedSeal.Round)

	// append commited seal suffix to commited seal hints
	blockHashAndSuffix = append(blockHashAndSuffix, commitedSealSuffix...)

}

// generateBlockHashAndSuffix is a private function to concatenate the
// commited seal suffix with the block hash
func generateBlockHashAndSuffix() ([]byte, error) {
	// init new header with dummy data
	// dummy data: round = 123
	header, err := generateBlockHeader()
	if err != nil {
		return []byte{}, err
	}

	// init new block with custom header
	block := types.NewBlockWithHeader(header)

	// extract Istanbul Extra data from block header
	extra, err := types.ExtractIstanbulExtra(block.Header())
	if err != nil {
		return []byte{}, err
	}

	// derive commited seal suffix data
	commitedSealSuffix := CommitedSealSuffix(extra.AggregatedSeal.Round)

	// init new byte slice to hold resulting Commited Seal Hints
	blockHashAndSuffix := make([]byte, 0)

	// append block hash bytes to commited seal hints
	blockHashAndSuffix = append(blockHashAndSuffix, block.Hash().Bytes()...)

	// append commited seal suffix to commited seal hints
	blockHashAndSuffix = append(blockHashAndSuffix, commitedSealSuffix...)

	return blockHashAndSuffix, nil
}
