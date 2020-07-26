package preprocess

import (
	"fmt"
	"testing"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/n3-util/n3cfg"
)

var (
	fPln          = fmt.Println
	failOnErrWhen = fn.FailOnErrWhen
)

// ****************************** //
// Auto generate config go struct file.
// Once 'config.toml' modified, run this func to update 'config.go' for new 'config.toml'
func TestGenSvrCfgStruct(t *testing.T) {
	n3cfg.GenStruct("../config/config.toml", "Config", "config", "../config/config_auto.go")
}

// ****************************** //
// Auto generate go struct file.
// Once 'config.toml' modified, run this func to update 'config.go' for new 'config.toml'
func TestGenCltCfgStruct(t *testing.T) {
	n3cfg.GenStruct("../goclient/config.toml", "Config", "goclient", "../goclient/config_auto.go")
}
