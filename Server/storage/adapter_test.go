package storage

import (
	"fmt"
	"testing"

	glb "github.com/nsip/n3-privacy/Server/global"
	pp "github.com/nsip/n3-privacy/preprocess"
)

func TestAdapter(t *testing.T) {
	glb.Init()
	db := NewDB("badger")
	fmt.Println(db)
	db = NewDB("map")
	fmt.Println(db)
}

func TestUpdatePolicy(t *testing.T) {
	glb.Init()
	db := NewDB("badger")

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

	return

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

	db := NewDB("badger")
	db.RecordMeta(policy, "../meta/meta1.json")
}
