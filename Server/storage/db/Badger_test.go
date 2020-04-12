package db

import (
	"io/ioutil"
	"testing"

	cmn "github.com/cdutwhu/json-util/common"
	glb "github.com/nsip/n3-privacy/Server/global"
)

func TestPolicyCount(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.PolicyCount())
}

func TestPolicyID(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.PolicyID("user", "ctx", "r", "6375f102e4d2562317fc841a9be2a56b8dfda4ad"))
}

func TestPolicyIDs(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.PolicyIDs("user", "ctx", "r"))
}

func TestUpdatePolicy(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	p, e := ioutil.ReadFile("../../../Mask/samples/xapiMask.json")
	cmn.FailOnErr("%v", e)
	id, obj, e := db.UpdatePolicy(string(p), "user", "ctx", "r")
	fPln(id, obj, e)
}

func TestDeletePolicy(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	db.DeletePolicy("92c8797efc18b369ed0aa9993e3647589c22335ar")
}

func TestPolicyHash(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.PolicyHash("92c8797efc18b369ed0aa9993e3647589c22335ar"))
}

func TestPolicy(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.Policy("92c8797efc18b369ed0aa9993e3647589c22335ar"))
}

// --------------------- //

func TestListPolicyID(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listPolicyID("user", "ctx", "r", "w"))
}

func TestListUser(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listUser("ctx"))
}

func TestListCtx(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listCtx("user"))
}

func TestListObject(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listObject("user", "ctx"))
}

// --------------------- //

func TestMapRW2lsPID(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.MapRW2lsPID("user", "ctx"))
}

func TestMapCtx2lsUser(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.MapCtx2lsUser())
	fPln(db.MapCtx2lsUser("ctx"))
	fPln(db.MapCtx2lsUser("def"))
}

func TestMapUser2lsCtx(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.MapUser2lsCtx())
	fPln(db.MapUser2lsCtx("user"))
	fPln(db.MapUser2lsCtx("abc"))
}

func TestMapUC2lsObject(t *testing.T) {
	glb.Init()
	db := NewDBByBadger().(*badgerDB)
	fPln(db.MapUC2lsObject("user", "ctx"))
	fPln(db.MapUC2lsObject("abc", "def"))
}
