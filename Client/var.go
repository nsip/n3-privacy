package main

import (
	"fmt"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/judge"
	"github.com/cdutwhu/gotil/rflx"
	"github.com/cdutwhu/n3-util/cfg"
)

var (
	fPt  = fmt.Print
	fPf  = fmt.Printf
	fPln = fmt.Println
	fSf  = fmt.Sprintf

	sJoin       = strings.Join
	sReplaceAll = strings.ReplaceAll

	setLog        = fn.SetLog
	resetLog      = fn.ResetLog
	failOnErr     = fn.FailOnErr
	failOnErrWhen = fn.FailOnErrWhen
	warnOnErrWhen = fn.WarnOnErrWhen
	isJSON        = judge.IsJSON
	exist         = judge.Exist
	cfgRepl       = cfg.Modify
	struct2Env    = rflx.Struct2Env
	struct2Map    = rflx.Struct2Map
	structFields  = rflx.StructFields
)
