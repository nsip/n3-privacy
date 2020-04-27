package config

import "testing"

func TestLoad(t *testing.T) {
	cfg := NewCfg("../config.toml")
	fPln(cfg.Route.ROOT)
	fPln(cfg.Path)
	fPln(cfg.LogFile)
	fPln(cfg.WebService)
	fPln(cfg.Route)
	fPln(cfg.Storage)
	fPln(cfg.File)
}
