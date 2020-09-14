package db

import (
	"io/ioutil"
	"testing"

	"github.com/cdutwhu/n3-util/n3err"
	"github.com/nsip/n3-privacy/Config/cfg"
)

const (
	Config = "../../../Config/config_test.toml"
	user   = "foo"
	ctx    = "bar"
)

var mReplExpr = map[string]string{
	"[db]": "Storage.DB",
}

func TestPolicyCount(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "PolicyCount")
	fPln(rets, ok)
}

func TestUpdatePolicy(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	p, e := ioutil.ReadFile("../../../Enforcer/samples/xapiPolicy.json")
	failOnErr("%v", e)
	rets, ok := tryInvokeWithMW(db, "UpdatePolicy", string(p), "object", user, ctx, "r")
	fPln(rets, ok)
}

func TestDeletePolicy(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
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
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "PolicyID", user, ctx, "r", "object")
	fPln(rets, ok)
}

func TestPolicyIDs(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "PolicyIDs", user, ctx, "r")
	fPln(rets, ok)
}

func TestPolicyHash(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "PolicyHash", "1615307cc418b369ed0a0beec7b5ea62cdb7020fr")
	fPln(rets, ok)
}

func TestPolicy(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "Policy", "1615307cc418b369ed0a0beec7b5ea62cdb7020fr")
	fPln(rets, ok)
}

// --------------------- //

func TestListPolicyID(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listPolicyID(user, ctx, "r"))
}

func TestListUser(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listUser(ctx))
}

func TestListCtx(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listCtx(user))
}

func TestListObject(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listObject(user, ctx))
}

// --------------------- //

func TestMapRW2lsPID(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "MapRW2lsPID", "", "")
	fPln(rets, ok)
	rets, ok = tryInvokeWithMW(db, "MapRW2lsPID", user, ctx)
	fPln(rets, ok)
}

func TestMapCtx2lsUser(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "MapCtx2lsUser")
	fPln(rets, ok)
	rets, ok = tryInvokeWithMW(db, "MapCtx2lsUser", ctx)
	fPln(rets, ok)
}

func TestMapUser2lsCtx(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "MapUser2lsCtx")
	fPln(rets, ok)
	rets, ok = tryInvokeWithMW(db, "MapUser2lsCtx", user)
	fPln(rets, ok)
}

func TestMapUC2lsObject(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	rets, ok := tryInvokeWithMW(db, "MapUC2lsObject", "", "")
	fPln(rets, ok)
	rets, ok = tryInvokeWithMW(db, "MapUC2lsObject", user, ctx)
	fPln(rets, ok)
}
