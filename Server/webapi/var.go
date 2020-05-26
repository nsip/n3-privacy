package webapi

import (
	"fmt"
	"sync"

	cmn "github.com/cdutwhu/n3-util/common"
	ext "github.com/cdutwhu/n3-util/external"
	glb "github.com/nsip/n3-privacy/Server/global"
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
	mapFromStruct = cmn.MapFromStruct

	prepare = ext.Prepare
)

func initMutex() map[string]*sync.Mutex {
	mMtx := make(map[string]*sync.Mutex)
	for _, v := range mapFromStruct(glb.Cfg.Route) {
		mMtx[v.(string)] = &sync.Mutex{}
	}
	return mMtx
}

type result struct {
	Data  *string `json:"data"`
	Empty *bool   `json:"empty"`
	Error string  `json:"error"`
}

// Const
var (
	fFalse    = false
	fTrue     = true
	fEmptyStr = ""
	False     = &fFalse
	True      = &fTrue
	EmptyStr  = &fEmptyStr
)
