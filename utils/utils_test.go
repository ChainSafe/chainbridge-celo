package utils

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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

// TestCommitedSealSuffix is ...
func TestCommitedSealSuffix(t *testing.T) {
	// init sample pointer to big int
	istAggSealRound := big.NewInt(123)
	// generate commited seal suffix from round
	commitedSealSuffix := CommitedSealSuffix(istAggSealRound)
	// fail if length of commited seal suffix less than 1
	if len(commitedSealSuffix) < 1 {
		t.Fatal("could not generate CommitedSealSuffix")
	}
}

// TestCommitedSealPrefix is ...
func TestCommitedSealPrefix(t *testing.T) {
	t.Log()
}

// TestCommitedSealHints is ...
func TestCommitedSealHints(t *testing.T) {
	t.Log()
}
