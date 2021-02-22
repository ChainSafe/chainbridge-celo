package bridge

import (
	"fmt"
	"testing"

	"github.com/status-im/keycard-go/hexutils"
)

func TestGetFunctionBytes(t *testing.T) {
	resb := getFunctionBytes("submitAsk(uint256,uint256)")
	ress := hexutils.BytesToHex(resb[:])
	if ress != "7288A28A" {
		t.Fatal(fmt.Sprintf("Result is %s", ress))
	}

}
