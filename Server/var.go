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

	failOnErr      = fn.FailOnErr
	failOnErrWhen  = fn.FailOnErrWhen
	localIP        = net.LocalIP
	setLog         = fn.SetLog
	logWhen        = fn.LoggerWhen
	logger         = fn.Logger
	env2Struct     = rflx.Env2Struct
	lrInit         = n3log.LrInit
	lrOut          = n3log.LrOut
	enableLoggly   = n3log.EnableLoggly
	setLogglyToken = n3log.SetLogglyToken
)
