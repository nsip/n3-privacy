package goclient

import (
	"testing"

	"github.com/cdutwhu/n3-util/n3cfg"
	"github.com/davecgh/go-spew/spew"
)

func TestLoad(t *testing.T) {
	cfg := newCfg("./config.toml")
	spew.Dump(cfg)
}

// ****************************** //
// Auto generate go struct file.
// Once 'config.toml' modified, run this func to update 'config.go' for new 'config.toml'
func TestAttrTypes(t *testing.T) {
	n3cfg.GenStruct("./config.toml", "", "", "./config_auto.go")
}
