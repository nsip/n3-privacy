package main

import (
	"fmt"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/net"
	"github.com/cdutwhu/gotil/rflx"
	"github.com/cdutwhu/n3-util/n3log"
)

var (
	fPt  = fmt.Print
	fPln = fmt.Println
	fSf  = fmt.Sprintf

	failOnErr     = fn.FailOnErr
	failOnErrWhen = fn.FailOnErrWhen
	localIP       = net.LocalIP
	enableLog2F   = fn.EnableLog2F
	logWhen       = fn.LoggerWhen
	logger        = fn.Logger
	env2Struct    = rflx.Env2Struct
	loggly        = n3log.Loggly
	logBind       = n3log.Bind
	setLoggly     = n3log.SetLoggly
)
