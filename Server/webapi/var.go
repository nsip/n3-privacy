package webapi

import (
	"fmt"
	"sync"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/judge"
	"github.com/cdutwhu/gotil/net"
	"github.com/cdutwhu/gotil/rflx"
	ext "github.com/cdutwhu/n3-util/external"
	"github.com/cdutwhu/n3-util/n3json"
	"github.com/cdutwhu/n3-util/n3log"
	"github.com/cdutwhu/n3-util/rest"
)

var (
	fSf  = fmt.Sprintf
	fPf  = fmt.Printf
	fPln = fmt.Println
	fPt  = fmt.Print

	localIP          = net.LocalIP
	isJSON           = judge.IsJSON
	jsonRoot         = n3json.JSONRoot
	failOnErr        = fn.FailOnErr
	failOnErrWhen    = fn.FailOnErrWhen
	warnOnErr        = fn.WarnOnErr
	warnOnErrWhen    = fn.WarnOnErrWhen
	setLog           = fn.SetLog
	logger           = fn.Logger
	url1Value        = rest.URL1Value
	url2Values       = rest.URL2Values
	url3Values       = rest.URL3Values
	url4Values       = rest.URL4Values
	urlValues        = rest.URLValues
	struct2Map       = rflx.Struct2Map
	env2Struct       = rflx.Env2Struct
	mustInvokeWithMW = rflx.MustInvokeWithMW
	toGeneralSlc     = rflx.ToGeneralSlc
	prepare          = ext.Prepare
	lrOut            = n3log.LrOut
)

func initMutex(route interface{}) map[string]*sync.Mutex {
	mMtx := make(map[string]*sync.Mutex)
	for _, v := range struct2Map(route) {
		mMtx[v.(string)] = &sync.Mutex{}
	}
	return mMtx
}

type result struct {
	Data  string `json:"data"`
	Empty bool   `json:"empty"`
	Error string `json:"error"`
}
