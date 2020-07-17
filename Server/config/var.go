package config

import (
	"fmt"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/net"
	"github.com/cdutwhu/gotil/rflx"
	"github.com/cdutwhu/n3-util/cfg"
)

var (
	fSf  = fmt.Sprintf
	fPln = fmt.Println

	sReplaceAll = strings.ReplaceAll
	sHasSuffix  = strings.HasSuffix
	sIndex      = strings.Index
	sContains   = strings.Contains
	sSplit      = strings.Split

	failOnErr     = fn.FailOnErr
	failOnErrWhen = fn.FailOnErrWhen
	localIP       = net.LocalIP
	cfgRepl       = cfg.Modify
	gitver        = cfg.GitVer
	struct2Env    = rflx.Struct2Env
)
