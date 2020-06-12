package db

import (
	"io/ioutil"
	"testing"

	eg "github.com/cdutwhu/n3-util/n3errs"
	cfg "github.com/nsip/n3-privacy/Server/config"
)

const Config = "../../config/config.toml"

func TestUpdatePolicy(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	p, e := ioutil.ReadFile("../../../Enforcer/samples/xapiPolicy.json")
	failOnErr("%v", e)
	{
		fPln("Normal Calling")
		id, obj, e := db.UpdatePolicy(string(p), "object", "user", "n3ctx", "r")
		fPln(id, obj, e)
	}
	{
		fPln("WM Calling")
		rets, ok, e := tryInvokeWithMW(db, "UpdatePolicy", string(p), "object", "user", "n3ctx", "r")
		fPln(rets, ok, e)
	}
}

func TestDeletePolicy(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	// db.DeletePolicy("92c8797efc18b369ed0a12dea96fec4024700fd9r")
	{
		fPln("Normal Calling")
		for _, pid := range db.PolicyIDs("user", "n3ctx", "r", "object") {
			db.DeletePolicy(pid)
		}
	}
	{
		fPln("WM Calling")
		if rets, ok, e := tryInvokeWithMW(db, "PolicyIDs", "user", "n3ctx", "r", "object"); e == nil && ok {
			for _, pid := range rets {
				mustInvokeWithMW(db, "DeletePolicy", pid)
			}
		}
	}
}

func TestPolicyCount(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	{
		fPln("Normal Calling")
		fPln(db.PolicyCount())
	}
	{
		fPln("WM Calling")
		rets, ok, e := tryInvokeWithMW(db, "PolicyCount")
		if failOnErr("%v", e); ok {
			fPln(rets)
		}
	}
}

func TestPolicyID(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	{
		fPln("Normal Calling")
		fPln(db.PolicyID("user", "n3ctx", "r", "object"))
	}
	{
		fPln("WM Calling")
		rets, ok, e := tryInvokeWithMW(db, "PolicyID", "user", "n3ctx", "r", "object")
		if failOnErr("%v", e); ok {
			fPln(rets)
		}
	}
}

func TestPolicyIDs(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	{
		fPln("Normal Calling")
		fPln(db.PolicyIDs("user", "n3ctx", "r"))
	}
	{
		fPln("WM Calling")
		rets, ok, e := tryInvokeWithMW(db, "PolicyIDs", "user", "n3ctx", "r")
		if failOnErr("%v", e); ok {
			fPln(rets)
		}
	}
}

func TestPolicyHash(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	{
		fPln("Normal Calling")
		fPln(db.PolicyHash("1615307cc418b369ed0a12dea96fecf9f90f0abbr"))
	}
	{
		fPln("WM Calling")
		rets, ok, e := tryInvokeWithMW(db, "PolicyHash", "1615307cc418b369ed0a12dea96fecf9f90f0abbr")
		if failOnErr("%v", e); ok {
			fPln(rets)
		}
	}
}

func TestPolicy(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	{
		fPln("Normal Calling")
		fPln(db.Policy("1615307cc418b369ed0a12dea96fecf9f90f0abbr"))
	}
	{
		fPln("WM Calling")
		rets, ok, e := tryInvokeWithMW(db, "Policy", "1615307cc418b369ed0a12dea96fecf9f90f0abbr")
		if failOnErr("%v", e); ok {
			fPln(rets)
		}
	}
}

// --------------------- //

func TestListPolicyID(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listPolicyID("user", "n3ctx", "r"))
}

func TestListUser(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listUser("n3ctx"))
}

func TestListCtx(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listCtx("user"))
}

func TestListObject(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.listObject("user", "n3ctx"))
}

// --------------------- //

func TestMapRW2lsPID(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	{
		fPln("Normal Calling")
		fPln(db.MapRW2lsPID("user", "n3ctx"))
	}
	{
		fPln("WM Calling")
		rets, ok, e := tryInvokeWithMW(db, "MapRW2lsPID", "user", "n3ctx")
		if failOnErr("%v", e); ok {
			fPln(rets)
		}
	}
}

func TestMapCtx2lsUser(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	{
		fPln("Normal Calling")
		fPln(db.MapCtx2lsUser())
		fPln(db.MapCtx2lsUser("n3ctx"))
	}
	{
		fPln("WM Calling")
		rets, ok, e := tryInvokeWithMW(db, "MapCtx2lsUser")
		if failOnErr("%v", e); ok {
			fPln(rets)
		}
		rets, ok, e = tryInvokeWithMW(db, "MapCtx2lsUser", "n3ctx")
		if failOnErr("%v", e); ok {
			fPln(rets)
		}
	}
}

func TestMapUser2lsCtx(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	{
		fPln("Normal Calling")
		fPln(db.MapUser2lsCtx())
		fPln(db.MapUser2lsCtx("user"))
	}
	{
		fPln("WM Calling")
		rets, ok, e := tryInvokeWithMW(db, "MapUser2lsCtx")
		if failOnErr("%v", e); ok {
			fPln(rets)
		}
		rets, ok, e = tryInvokeWithMW(db, "MapUser2lsCtx", "user")
		if failOnErr("%v", e); ok {
			fPln(rets)
		}
	}
}

func TestMapUC2lsObject(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	{
		fPln("Normal Calling")
		fPln(db.MapUC2lsObject("user", "n3ctx"))
	}
	{
		fPln("WM Calling")
		rets, ok, e := tryInvokeWithMW(db, "MapUC2lsObject", "user", "n3ctx")
		if failOnErr("%v", e); ok {
			fPln(rets)
		}
	}
}
