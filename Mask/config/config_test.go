package config

import "testing"

func TestLoad(t *testing.T) {
	cfg := NewCfg("../config.toml")
	fPln(cfg)
}
