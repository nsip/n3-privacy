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
	p, e := ioutil.ReadFile("../../../Enforcer/samples/xapiPolicy1.json")
	failOnErr("%v", e)
	id, obj, e := db.UpdatePolicy(string(p), "object", "user", "n3ctx", "r")
	fPln(id, obj, e)
}

func TestDeletePolicy(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	// db.DeletePolicy("92c8797efc18b369ed0a12dea96fec4024700fd9r")
	for _, pid := range db.PolicyIDs("user", "n3ctx", "r", "object", "root") {
		db.DeletePolicy(pid)
	}
}

func TestPolicyCount(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.PolicyCount())
}

func TestPolicyID(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.PolicyID("user", "n3ctx", "r", "root"))
}

func TestPolicyIDs(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.PolicyIDs("user", "n3ctx", "r"))
}

func TestPolicyHash(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.PolicyHash("dc76e9f0c018b369ed0a12dea96fec4024700fd9r"))
}

func TestPolicy(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.Policy("dc76e9f0c018b369ed0a12dea96fec4024700fd9r"))
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
	fPln(db.MapRW2lsPID("user", "n3ctx"))
}

func TestMapCtx2lsUser(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.MapCtx2lsUser())
	fPln(db.MapCtx2lsUser("n3ctx"))
}

func TestMapUser2lsCtx(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.MapUser2lsCtx())
	fPln(db.MapUser2lsCtx("user"))
}

func TestMapUC2lsObject(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDBByBadger().(*badgerDB)
	fPln(db.MapUC2lsObject("user", "n3ctx"))
}
