package main

import (
	"fmt"
	"sync"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/judge"
	"github.com/cdutwhu/gotil/net"
	"github.com/cdutwhu/gotil/rflx"
	"github.com/cdutwhu/n3-util/n3json"
	"github.com/cdutwhu/n3-util/n3log"
	"github.com/cdutwhu/n3-util/rest"
)

var (
	fPt              = fmt.Print
	fPln             = fmt.Println
	fSf              = fmt.Sprintf
	fPf              = fmt.Printf
	localIP          = net.LocalIP
	enableLog2F      = fn.EnableLog2F
	enableWarnDetail = fn.EnableWarnDetail
	failOnErr        = fn.FailOnErr
	failOnErrWhen    = fn.FailOnErrWhen
	logWhen          = fn.LoggerWhen
	logger           = fn.Logger
	warnOnErr        = fn.WarnOnErr
	warnOnErrWhen    = fn.WarnOnErrWhen
	warner           = fn.Warner
	debug            = fn.Debug
	isJSON           = judge.IsJSON
	url1Value        = rest.URL1Value
	url2Values       = rest.URL2Values
	url3Values       = rest.URL3Values
	url4Values       = rest.URL4Values
	urlValues        = rest.URLValues
	struct2Map       = rflx.Struct2Map
	env2Struct       = rflx.Env2Struct
	mustInvokeWithMW = rflx.MustInvokeWithMW
	toGeneralSlc     = rflx.ToGeneralSlc
	loggly           = n3log.Loggly
	logBind          = n3log.Bind
	setLoggly        = n3log.SetLoggly
	syncBindLog      = n3log.SyncBindLog
	jsonRoot         = n3json.JSONRoot
)

const (
	envKey = "PRISvr"
)

var (
	logGrp  = logBind(logger) // logBind(logger, loggly("info"))
	warnGrp = logBind(warner) // logBind(warner, loggly("warn"))
)

func initMutex(route interface{}) map[string]*sync.Mutex {
	mMtx := make(map[string]*sync.Mutex)
	for _, v := range struct2Map(route) {
		mMtx[v.(string)] = &sync.Mutex{}
	}
	return mMtx
}
