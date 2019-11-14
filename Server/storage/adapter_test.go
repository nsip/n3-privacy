package storage

import (
	"fmt"
	"testing"

	glb "github.com/nsip/n3-privacy/Server/global"
	dbfn "github.com/nsip/n3-privacy/Server/storage/db"
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
	db := NewDB("map")

	uid := "u123456"
	ctx := "c123fff"

	policy := `{
		"test": {
			"T1": "-----",
			"f1":    "*333*333***"
		}
	}`
	policy = pp.FmtJSONStr(policy)
	db.UpdatePolicy(policy, uid, ctx, "r")

	lsID := dbfn.GetPolicyID(uid, ctx, "test", "r")
	for _, id := range lsID {
		fPln(id)
		policy, _ := db.GetPolicy(id)
		fPln(policy)
	}
}

func TestRecPolicyMeta(t *testing.T) {
	policy := pp.FmtJSONFile("../meta/mask.json")
	dbfn.RecordMeta(policy, "../meta/meta1.json")
}
