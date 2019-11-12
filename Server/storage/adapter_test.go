package storage

import (
	"fmt"
	"testing"

	pp "github.com/nsip/n3-privacy/preprocess"
)

func TestAdapter(t *testing.T) {
	db := NewDB()
	fmt.Println(db)
}

func TestUpdatePolicy(t *testing.T) {

	db := NewDB()

	uid := "u123456"
	ctx := "c123fff"

	policy := `{
		"testobj": {
			"t1": "-----",
			"f12":    "*333*333***"
		}
	}`
	policy = pp.FmtJSONStr(policy)
	db.UpdatePolicy(policy, uid, ctx, "r")

	policy = `{
		"testobj": {
			"f12T": "*444*444***",
			"t1": "-----"
		}
	}`
	policy = pp.FmtJSONStr(policy)
	db.UpdatePolicy(policy, uid, ctx, "r")

	lsIDs := db.GetPolicyID(uid, ctx, "testobj", "r")
	for _, id := range lsIDs {
		fPln(id)
		policy, _ := db.GetPolicy(id)
		fPln(policy)
	}

	// fPln(mCTXlsUID)
	// fPln(mUIDlsCTX)
}

func TestRecPolicyMeta(t *testing.T) {
	policy := pp.FmtJSONFile("../meta/mask.json")

	db := NewDB()
	db.RecordMeta(policy, "../meta/meta1.json")
}
