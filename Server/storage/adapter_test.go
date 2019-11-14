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

	uid := "u123"
	ctx := "c123"
	policy := `{
		"Test": {
			"F1": "-----",
			"F3": "*****"
		}
	}`
	policy = pp.FmtJSONStr(policy)
	db.UpdatePolicy(policy, uid, ctx, "r")

	fPln(db.PolicyCount())

	for _, id := range db.PolicyIDs(uid, ctx, "r") {
		fPln(id)
		policy, _ := db.Policy(id)
		fPln(policy)
	}
}
