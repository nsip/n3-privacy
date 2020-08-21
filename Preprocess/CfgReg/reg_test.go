package cfgreg

import (
	"flag"
	"os"
	"os/user"
	"strings"
	"testing"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/n3-util/n3cfg"
)

// Under n3-privacy/Preprocess/CfgReg/
// echo password | sudo -S env "PATH=$PATH" go test -v -count=1 ./ -run TestRegCfg -args `whoami` server ...
func TestRegCfg(t *testing.T) {
	if !flag.Parsed() {
		flag.Parse()
	}
	user, _ := user.Current()
	osuser := user.Name
	if len(flag.Args()) > 0 {
		osuser = flag.Args()[0]
	}
	project := n3cfg.PrjName()
	// fmt.Println(osuser, project)

	mPkgConfig := map[string]string{
		"server":   "../../Server/config.toml",
		"goclient": "../../Server/goclient/config.toml",
	}

	for _, pkg := range flag.Args()[1:] {
		config := mPkgConfig[strings.ToLower(pkg)]
		_, err := os.Stat(config)
		if fn.WarnOnErr("%v @ %s %s", err, pkg, config) != nil {
			continue
		}
		ok, file := n3cfg.Register(osuser, config, project, pkg)
		fn.Logger("%v %v", ok, file)
	}
}
