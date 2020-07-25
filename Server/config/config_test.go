package config

import (
	"flag"
	"os"
	"testing"

	"github.com/cdutwhu/n3-util/n3cfg"
	eg "github.com/cdutwhu/n3-util/n3errs"
	"github.com/davecgh/go-spew/spew"
)

func TestLoad(t *testing.T) {
	cfg := newCfg("./config.toml")
	spew.Dump(cfg)
	cfg.SaveAs("./copy")
}

// ****************************** //
// Auto generate go struct file.
// Once 'config.toml' modified, run this func to update 'config.go' for new 'config.toml'
func TestAttrTypes(t *testing.T) {
	n3cfg.GenStruct("./config.toml", "", "", "./config_auto.go")
}

// ****************************** //
// Create a copy of config for Client. Excluding some attributes.
// Once building, move it to Client config Directory.
func TestGenClientCfg(t *testing.T) {
	cfg := newCfg("./config.toml")
	failOnErrWhen(cfg == nil, "%v", eg.CFG_INIT_ERR)
	temp := "./temp.toml"
	cfg.SaveAs(temp)
	if !flag.Parsed() {
		flag.Parse()
	}
	fPln(flag.Args(), "were removed")
	n3cfg.RmFileAttrL1(temp, "../goclient/config.toml", flag.Args()...)
	os.Remove(temp)
}
