package config

import "testing"

func TestLoad(t *testing.T) {
	cfg := NewCfg("../config.toml")
	fPln(cfg.Path)
	fPln(cfg.LogFile)
	fPln(cfg.WebService)
	fPln(cfg.Route)
	fPln(cfg.Storage)
	fPln(cfg.File)
	fPln(cfg.Server)
	fPln(cfg.Access)

	cfg.SaveAs("./copy")
}
