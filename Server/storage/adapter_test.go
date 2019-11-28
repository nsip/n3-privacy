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
	ctx := "ctx124"
	policy := `{
		"Test": {
			"F1": "-----",
			"F2": "*****",
			"F3": "~~~~~",
			"F4": "     "
		}
	}`
	policy = pp.FmtJSONStr(policy)
	db.UpdatePolicy(policy, uid, ctx, "r")
	fPln(db.PolicyCount())

	fPln(db.MapRWListOfPID("qmiao", "ctx122"))
	fPln(db.MapCtxListOfUser("ctx122", "ctx124"))
	fPln(db.MapUserListOfCtx())
	fPln(db.MapUCListOfObject("qmiao", "ctx122"))

	if id := db.PolicyID(uid, ctx, "w", "Test"); len(id) > 0 {
		fPln(id)
	}
	for _, id := range db.PolicyIDs(uid, ctx, "w", "Test") {
		fPln(id)
		policy, _ := db.Policy(id)
		fPln(policy)
	}
}
