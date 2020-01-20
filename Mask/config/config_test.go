package config

import "testing"

func TestLoad(t *testing.T) {
	cfg := NewCfg("../Config.toml")
	fPln(cfg)
}
