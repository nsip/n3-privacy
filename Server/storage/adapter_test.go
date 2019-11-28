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

	uid := "qmiao"
	ctx := "ctx122"
	policy := `{
		"Test": {
			"F1": "-----",
			"F2": "*****",
			"F3": "~~~~~"
		}
	}`
	policy = pp.FmtJSONStr(policy)
	db.UpdatePolicy(policy, uid, ctx, "w")
	fPln(db.PolicyCount())

	if id := db.PolicyID(uid, ctx, "w", "Test"); len(id) > 0 {
		fPln(id)
	}

	for _, id := range db.PolicyIDs(uid, ctx, "w", "Test") {
		fPln(id)
		policy, _ := db.Policy(id)
		fPln(policy)
	}

	fPln(db.ListPolicyID("qmiao1", "", "r", "w"))

	fPln(db.ListUser("ctx122", "ctx123"))
	fPln(db.ListCtx("qmiao1", "qmiao"))
	fPln(db.ListObject("qmiao", "ctx123"))
}
