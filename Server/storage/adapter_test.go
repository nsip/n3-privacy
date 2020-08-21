package storage

import (
	"testing"

	"github.com/cdutwhu/n3-util/n3cfg"
)

const Config = "../config.toml"

func TestAdapter(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDB("badger")
	fPln(db)
}

func TestUpdatePolicy(t *testing.T) {
	n3cfg.ToEnvN3privacyServer(nil, envKey, Config)
	db := NewDB("badger").(DB)
	fPln(db)

	fPln(db.MapRW2lsPID("foo", "bar"))
	fPln(db.MapCtx2lsUser("bar"))
	fPln(db.MapUser2lsCtx("foo"))
	fPln(db.MapUC2lsObject("foo", "bar"))
}
