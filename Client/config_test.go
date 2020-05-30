package main

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestLoad(t *testing.T) {
	cfg := NewCfg("./cfg-clt-privacy.toml")
	spew.Dump(cfg)
}
