package config

import (
	"flag"
	"os"
	"testing"

	eg "github.com/cdutwhu/n3-util/n3errs"
	toml "github.com/cdutwhu/n3-util/n3toml"
)

func TestLoad(t *testing.T) {
	cfg := NewCfg("./config.toml")
	fPln(cfg.Path)
	fPln(cfg.LogFile)
	fPln(cfg.ServiceName)
	fPln(cfg.WebService)
	fPln(cfg.Route)
	fPln(cfg.Storage)
	fPln(cfg.File)
	fPln(cfg.Server)
	fPln(cfg.Access)
	cfg.SaveAs("./copy")
}

// Create a copy of config for Client. Excluding some attributes.
// Once building, move it to Client config Directory.
func TestGenClientCfg(t *testing.T) {
	cfg := NewCfg("./config.toml")
	failOnErrWhen(cfg == nil, "%v", eg.CFG_INIT_ERR)
	temp := "./temp.toml"
	cfg.SaveAs(temp)
	if !flag.Parsed() {
		flag.Parse()
	}
	fPln(flag.Args(), "to be removed")
	toml.RmFileAttrL1(temp, "../../Client/config/config.toml", flag.Args()...)
	os.Remove(temp)
}
