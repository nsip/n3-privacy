package config

import (
	"fmt"
	"strings"

	cmn "github.com/cdutwhu/n3-util/common"
)

var (
	fSf  = fmt.Sprintf
	fPln = fmt.Println

	sReplaceAll = strings.ReplaceAll
	sHasSuffix  = strings.HasSuffix
	sIndex      = strings.Index
	sContains   = strings.Contains

	failOnErr = cmn.FailOnErr
	localIP   = cmn.LocalIP
	cfgRepl   = cmn.CfgRepl
)
