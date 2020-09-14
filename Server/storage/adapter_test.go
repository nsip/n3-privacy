package storage

import (
	"testing"

	"github.com/cdutwhu/n3-util/n3err"
	"github.com/nsip/n3-privacy/Config/cfg"
)

const Config = "../../Config/config_test.toml"

var mReplExpr = map[string]string{
	"[db]": "Storage.DB",
}

func TestAdapter(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDB("badger")
	fPln(db)
}

func TestUpdatePolicy(t *testing.T) {
	failOnErrWhen(cfg.NewCfg("Config", mReplExpr, Config) == nil, "%v", n3err.CFG_INIT_ERR)
	db := NewDB("badger").(DB)
	fPln(db)

	fPln(db.MapRW2lsPID("foo", "bar"))
	fPln(db.MapCtx2lsUser("bar"))
	fPln(db.MapUser2lsCtx("foo"))
	fPln(db.MapUC2lsObject("foo", "bar"))
}
