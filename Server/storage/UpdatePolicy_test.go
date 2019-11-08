package storage

import (
	"testing"

	pp "github.com/nsip/n3-privacy/preprocess"
)

func TestUpdatePolicy(t *testing.T) {

	uid := "u123456"
	ctx := "c123fff"

	mask := `{
		"testobj": {				
			"t1": "-----",
			"f12":    "*333*333***"			
		}
	}`
	mask = pp.FmtJSONStr(mask, "../../preprocess/utils")
	UpdatePolicy(uid, ctx, "r", mask)

	mask = `{
		"testobj": {	
			"f12T": "*444*444***",				
			"t1": "-----"					
		}
	}`
	mask = pp.FmtJSONStr(mask, "../../preprocess/utils")
	UpdatePolicy(uid, ctx, "r", mask)

	lsIDs := GetPolicyID(uid, ctx, "testobj", "r")
	for _, code := range lsIDs {
		fPln(code)
		policy, _ := GetPolicy(code)
		fPln(policy)
	}

	fPln(mCTXlsUID)
	fPln(mUIDlsCTX)
}
