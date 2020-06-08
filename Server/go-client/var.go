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
	logger        = cmn.Log
	warnOnErrWhen = cmn.WarnOnErrWhen
	isJSON        = cmn.IsJSON
	struct2Env    = cmn.Struct2Env
	struct2Map    = cmn.Struct2Map
	mapKeys       = cmn.MapKeys
)

const (
	envVarName = "CfgClt-PRI"
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
	m, err := struct2Map(route)
	failOnErr("%v", err)
	for k, v := range m {
		mFnURL[k] = fSf("%s://%s:%d%s", protocol, ip, port, v)
	}
	Ikeys, err := mapKeys(mFnURL)
	failOnErr("%v", err)
	return mFnURL, Ikeys.([]string)
}
