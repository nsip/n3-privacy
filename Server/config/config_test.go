package config

import (
	"flag"
	"os"
	"testing"

	eg "github.com/cdutwhu/n3-util/n3errs"
	toml "github.com/cdutwhu/n3-util/n3toml"
	"github.com/davecgh/go-spew/spew"
)

func TestLoad(t *testing.T) {
	cfg := newCfg("./config.toml")
	spew.Dump(cfg)
	cfg.SaveAs("./copy")
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
	toml.RmFileAttrL1(temp, "../go-client/config.toml", flag.Args()...)
	os.Remove(temp)
}
