package storage

import (
	"testing"

	eg "github.com/cdutwhu/n3-util/n3errs"
	cfg "github.com/nsip/n3-privacy/Server/config"
)

const Config = "../config/config.toml"

func TestAdapter(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDB("badger")
	fPln(db)

	// db = NewDB("map")
	// fmt.Println(db)
}

func TestUpdatePolicy(t *testing.T) {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg", Config), "%v: Config Init Error", eg.CFG_INIT_ERR)
	db := NewDB("badger")
	fPln(db)

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
