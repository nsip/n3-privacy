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
	ctx := "ctx1234"
	// policy := `{
	// 	"Test": {
	// 		"F1": "-----",
	// 		"F2": "*****"
	// 	}
	// }`
	// policy = pp.FmtJSONStr(policy)
	// db.UpdatePolicy(policy, uid, ctx, "r")

	// fPln(db.PolicyCount())

	if id := db.PolicyID(uid, ctx, "r", "Test"); len(id) > 0 {
		fPln(id)
	}

	for _, id := range db.PolicyIDs(uid, ctx, "r", "testobj") {
		fPln(id)
		policy, _ := db.Policy(id)
		fPln(policy)
	}
}
