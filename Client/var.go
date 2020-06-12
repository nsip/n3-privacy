package main

import (
	"fmt"
	"strings"

	cmn "github.com/cdutwhu/n3-util/common"
)

var (
	fPt  = fmt.Print
	fPf  = fmt.Printf
	fPln = fmt.Println
	fSf  = fmt.Sprintf

	sJoin       = strings.Join
	sReplaceAll = strings.ReplaceAll

	xin           = cmn.XIn
	setLog        = cmn.SetLog
	resetLog      = cmn.ResetLog
	failOnErr     = cmn.FailOnErr
	failOnErrWhen = cmn.FailOnErrWhen
	warnOnErrWhen = cmn.WarnOnErrWhen
	isFLog        = cmn.IsFLog
	isJSON        = cmn.IsJSON
	cfgRepl       = cmn.CfgRepl
	struct2Env    = cmn.Struct2Env
	struct2Map    = cmn.Struct2Map
	structFields  = cmn.StructFields
)
