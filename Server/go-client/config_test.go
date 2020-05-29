package client

import "testing"

func TestLoad(t *testing.T) {
	cfg := newCfg("./config.toml")
	fPln(cfg.Path)
	fPln(cfg.LogFile)
	fPln(cfg.Route)
	fPln(cfg.Server)
	fPln(cfg.Access)
}
