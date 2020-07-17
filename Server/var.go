package main

import (
	"fmt"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/net"
	"github.com/cdutwhu/gotil/rflx"
)

var (
	fPln = fmt.Println

	failOnErr     = fn.FailOnErr
	failOnErrWhen = fn.FailOnErrWhen
	localIP       = net.LocalIP
	setLog        = fn.SetLog
	logWhen       = fn.LoggerWhen
	logger        = fn.Logger
	env2Struct    = rflx.Env2Struct
)
