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

	// user := "qmiao"
	// ctx := "ctx123"
	// policy := `{
	// 	"Test": {
	// 		"F1": "-----",
	// 		"F2": "*****",
	// 		"F3": "~~~~~",
	// 		"F4": "     "
	// 	}
	// }`
	// policy = pp.FmtJSONStr(policy)
	// db.UpdatePolicy(policy, user, ctx, "r")
	// fPln(db.PolicyCount())

	fPln(db.MapRW2lsPID("user", "ctx"))
	fPln(db.MapCtx2lsUser("ctx", "def"))
	fPln(db.MapUser2lsCtx())
	fPln(db.MapUC2lsObject("user", "ctx"))

	// if id := db.PolicyID(user, ctx, "w", "Test"); len(id) > 0 {
	// 	fPln(id)
	// }
	// for _, id := range db.PolicyIDs(user, ctx, "w", "Test") {
	// 	fPln(id)
	// 	policy, _ := db.Policy(id)
	// 	fPln(policy)
	// }
}
