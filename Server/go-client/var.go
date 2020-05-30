package client

import (
	"fmt"
	"strings"

	cmn "github.com/cdutwhu/n3-util/common"
)

var (
	fSf           = fmt.Sprintf
	fPln          = fmt.Println
	sReplaceAll   = strings.ReplaceAll
	sReplace      = strings.Replace
	sJoin         = strings.Join
	sTrimRight    = strings.TrimRight
	cfgRepl       = cmn.CfgRepl
	failOnErr     = cmn.FailOnErr
	failOnErrWhen = cmn.FailOnErrWhen
	warnOnErr     = cmn.WarnOnErr
	env2Struct    = cmn.Env2Struct
	setLog        = cmn.SetLog
	logWhen       = cmn.LogWhen
	warnOnErrWhen = cmn.WarnOnErrWhen
	isJSON        = cmn.IsJSON
	struct2Env    = cmn.Struct2Env
	mapFromStruct = cmn.MapFromStruct
	mapKeys       = cmn.MapKeys
)

// Args is arguments for "Route"
type Args struct {
	ID     string
	User   string
	Ctx    string
	Object string
	RW     string
	Policy []byte
	Data   []byte
}

func initMapFnURL(protocol, ip string, port int, route interface{}) (map[string]string, []string) {
	mFnURL := make(map[string]string)
	for k, v := range mapFromStruct(route) {
		mFnURL[k] = fSf("%s://%s:%d%s", protocol, ip, port, v)
	}
	return mFnURL, mapKeys(mFnURL).([]string)
}
