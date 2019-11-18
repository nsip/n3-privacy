package storage

import (
	"fmt"
	"testing"

	glb "github.com/nsip/n3-privacy/Server/global"
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
	ctx := "ctx123"
	// policy := `{
	// 	"Test": {
	// 		"F1": "-----",
	// 		"F2": "*****"
	// 	}
	// }`
	// policy = pp.FmtJSONStr(policy)
	// db.UpdatePolicy(policy, uid, ctx, "r")

	// fPln(db.PolicyCount())

	policy := db.PolicyID(uid, ctx, "r", "Test")[0]
	fPln(policy)

	for _, id := range db.PolicyIDs(uid, ctx, "r") {
		fPln(id)
		policy, _ := db.Policy(id)
		fPln(policy)
	}
}
