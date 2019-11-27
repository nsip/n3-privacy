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

	uid := "qmiao1"
	ctx := "ctx123"
	policy := `{
		"Test1": {
			"F1": "-----",
			"F2": "*****",
			"F3": "~~~~~"
		}
	}`
	policy = pp.FmtJSONStr(policy)
	db.UpdatePolicy(policy, uid, ctx, "r")
	fPln(db.PolicyCount())

	if id := db.PolicyID(uid, ctx, "r", "Test"); len(id) > 0 {
		fPln(id)
	}

	for _, id := range db.PolicyIDs(uid, ctx, "r", "Test") {
		fPln(id)
		policy, _ := db.Policy(id)
		fPln(policy)
	}

	fPln(db.ListAllUser())
	fPln(db.ListUserOfOneCtx("ctx123"))
	fPln(db.ListAllCtx())
	fPln(db.ListCtxOfOneUser("qmiao1"))
	fPln(db.ListAllObject())
	fPln(db.ListObjectOfOneUser("qmiao1"))
	fPln(db.ListObjectOfOneCtx("ctx123"))
}
