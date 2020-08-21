package cfggen

import (
	"flag"
	"os"
	"testing"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/n3-util/n3cfg"
	"github.com/cdutwhu/n3-util/n3cfg/attrim"
	"github.com/cdutwhu/n3-util/n3err"
)

func TestMkCltCfg(t *testing.T) {
	cfg := n3cfg.ToEnvN3privacyServer(
		map[string]string{
			"[v]":    "Version",
			"[s]":    "Service",
			"[port]": "WebService.Port",
			"[db]":   "Storage.DB",
		},
		"FakeKey",
		"../../Server/config.toml",
	)
	fn.FailOnErrWhen(cfg == nil, "%v", n3err.CFG_INIT_ERR)

	temp := "./temp.toml"
	defer func() { os.Remove(temp) }()
	n3cfg.Save(temp, cfg)

	// A Copy for Server Executable
	attrim.RmCfgAttrL1(temp, "../../Server/config_rel.toml")

	// A Copy for Server's goclient
	if !flag.Parsed() {
		flag.Parse()
	}
	attrim.SelCfgAttrL1(temp, "../../Server/goclient/config.toml", flag.Args()...)
}
