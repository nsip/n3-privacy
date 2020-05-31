package webapi

import (
	"fmt"
	"sync"

	cmn "github.com/cdutwhu/n3-util/common"
	ext "github.com/cdutwhu/n3-util/external"
)

var (
	fSf  = fmt.Sprintf
	fPf  = fmt.Printf
	fPln = fmt.Println

	localIP       = cmn.LocalIP
	isJSON        = cmn.IsJSON
	jsonRoot      = cmn.JSONRoot
	failOnErrWhen = cmn.FailOnErrWhen
	warnOnErr     = cmn.WarnOnErr
	warnOnErrWhen = cmn.WarnOnErrWhen
	setLog        = cmn.SetLog
	url1Value     = cmn.URL1Value
	url2Values    = cmn.URL2Values
	url3Values    = cmn.URL3Values
	url4Values    = cmn.URL4Values
	urlValues     = cmn.URLValues
	struct2Map    = cmn.Struct2Map
	env2Struct    = cmn.Env2Struct

	prepare = ext.Prepare
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
