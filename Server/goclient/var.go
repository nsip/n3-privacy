package goclient

import (
	"fmt"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/judge"
	"github.com/cdutwhu/gotil/rflx"
	"github.com/cdutwhu/n3-util/n3cfg"
)

var (
	fSf           = fmt.Sprintf
	fPln          = fmt.Println
	sReplaceAll   = strings.ReplaceAll
	sReplace      = strings.Replace
	sJoin         = strings.Join
	sTrimRight    = strings.TrimRight
	cfgRepl       = n3cfg.Modify
	failOnErr     = fn.FailOnErr
	failOnErrWhen = fn.FailOnErrWhen
	warnOnErr     = fn.WarnOnErr
	enableLog2F   = fn.EnableLog2F
	logWhen       = fn.LoggerWhen
	logger        = fn.Logger
	warnOnErrWhen = fn.WarnOnErrWhen
	isJSON        = judge.IsJSON
	struct2Env    = rflx.Struct2Env
	struct2Map    = rflx.Struct2Map
	mapKeys       = rflx.MapKeys
	env2Struct    = rflx.Env2Struct
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
	for k, v := range struct2Map(route) {
		mFnURL[k] = fSf("%s://%s:%d%s", protocol, ip, port, v)
	}
	return mFnURL, mapKeys(mFnURL).([]string)
}
