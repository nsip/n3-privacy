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
	sSplit      = strings.Split

	failOnErr     = cmn.FailOnErr
	failOnErrWhen = cmn.FailOnErrWhen
	localIP       = cmn.LocalIP
	cfgRepl       = cmn.CfgRepl
	struct2Env    = cmn.Struct2Env
)
