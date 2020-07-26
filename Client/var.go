package main

import (
	"fmt"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/judge"
	"github.com/cdutwhu/gotil/rflx"
	"github.com/cdutwhu/n3-util/n3cfg"
)

var (
	fPt         = fmt.Print
	fPf         = fmt.Printf
	fPln        = fmt.Println
	fSf         = fmt.Sprintf
	sJoin       = strings.Join
	sReplaceAll = strings.ReplaceAll

	enableLog2F   = fn.EnableLog2F
	failOnErr     = fn.FailOnErr
	failOnErrWhen = fn.FailOnErrWhen
	warnOnErrWhen = fn.WarnOnErrWhen
	isJSON        = judge.IsJSON
	exist         = judge.Exist
	cfgRepl       = n3cfg.Modify
	struct2Env    = rflx.Struct2Env
	struct2Map    = rflx.Struct2Map
	structFields  = rflx.StructFields
)
