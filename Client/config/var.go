package config

import (
	"fmt"
	"strings"

	cmn "github.com/cdutwhu/json-util/common"
)

var (
	fSf         = fmt.Sprintf
	fPln        = fmt.Println
	sReplaceAll = strings.ReplaceAll
	cfgRepl     = cmn.CfgRepl
	failOnErr   = cmn.FailOnErr
)
