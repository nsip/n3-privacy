package db

import (
	"io/ioutil"
	"testing"

	"github.com/cdutwhu/n3-util/n3cfg"
)

const (
	Config = "../../config.toml"
	user   = "foo"
	ctx    = "bar"
)

func TestPolicyCount(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "PolicyCount")
	fPln(rets, ok)
}

func TestUpdatePolicy(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	p, e := ioutil.ReadFile("../../../Enforcer/samples/xapiPolicy.json")
	failOnErr("%v", e)
	rets, ok := tryInvokeWithMW(db, "UpdatePolicy", string(p), "object", user, ctx, "r")
	fPln(rets, ok)
}

func TestDeletePolicy(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	// db.DeletePolicy("92c8797efc18b369ed0a12dea96fec4024700fd9r")

	if rets, ok := tryInvokeWithMW(db, "PolicyIDs", user, ctx, "r", "object"); ok {
		for _, pid := range rets {
			if ids, ok := pid.([]string); ok && len(ids) > 0 {
				mustInvokeWithMW(db, "DeletePolicy", ids[0])
			}
		}
	}
}

func TestPolicyID(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "PolicyID", user, ctx, "r", "object")
	fPln(rets, ok)
}

func TestPolicyIDs(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "PolicyIDs", user, ctx, "r")
	fPln(rets, ok)
}

func TestPolicyHash(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "PolicyHash", "1615307cc418b369ed0a12dea96fecf9f90f0abbr")
	fPln(rets, ok)
}

func TestPolicy(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "Policy", "1615307cc418b369ed0a12dea96fecf9f90f0abbr")
	fPln(rets, ok)
}

// --------------------- //

func TestListPolicyID(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listPolicyID(user, ctx, "r"))
}

func TestListUser(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listUser(ctx))
}

func TestListCtx(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listCtx(user))
}

func TestListObject(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listObject(user, ctx))
}

// --------------------- //

func TestMapRW2lsPID(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "MapRW2lsPID", "", "")
	fPln(rets, ok)
	rets, ok = tryInvokeWithMW(db, "MapRW2lsPID", user, ctx)
	fPln(rets, ok)
}

func TestMapCtx2lsUser(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "MapCtx2lsUser")
	fPln(rets, ok)
	rets, ok = tryInvokeWithMW(db, "MapCtx2lsUser", ctx)
	fPln(rets, ok)
}

func TestMapUser2lsCtx(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "MapUser2lsCtx")
	fPln(rets, ok)
	rets, ok = tryInvokeWithMW(db, "MapUser2lsCtx", user)
	fPln(rets, ok)
}

func TestMapUC2lsObject(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "MapUC2lsObject", "", "")
	fPln(rets, ok)
	rets, ok = tryInvokeWithMW(db, "MapUC2lsObject", user, ctx)
	fPln(rets, ok)
}
