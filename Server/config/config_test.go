package config

import (
	"flag"
	"os"
	"testing"

	"github.com/cdutwhu/n3-util/n3cfg"
	"github.com/cdutwhu/n3-util/n3errs"
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
func TestGenCltCfg(t *testing.T) {
	cfg := newCfg("./config.toml")
	failOnErrWhen(cfg == nil, "%v", n3errs.CFG_INIT_ERR)
	temp := "./temp.toml"
	cfg.SaveAs(temp)
	if !flag.Parsed() {
		flag.Parse()
	}
	n3cfg.RmFileAttrL1(temp, "../goclient/config.toml", flag.Args()...)
	os.Remove(temp)
	fPln(flag.Args(), "were removed")
}
