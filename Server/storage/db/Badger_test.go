package db

import (
	"testing"

	glb "github.com/nsip/n3-privacy/Server/global"
)

func TestLoadIDList(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	// for _, id := range listID {
	// 	fPln(id)
	// }
	fPln(db.PolicyCount())
}

func TestAllPolicyIDs(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	// for _, id := range listID {
	// 	fPln(id)
	// }
	IDs := db.PolicyIDListOfOneCtx("ctx123", "r")
	//IDs[0] = "sssss"
	fPln(IDs)
	// for _, id := range listID {
	// 	fPln(id)
	// }

	IDs = db.PolicyIDs("qmiao", "ctx123", "r")
	fPln(IDs)
}
