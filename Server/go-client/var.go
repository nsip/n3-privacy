package client

import (
	"fmt"
	"strings"

	cmn "github.com/cdutwhu/n3-util/common"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
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
	for k, v := range struct2Map(route) {
		mFnURL[k] = fSf("%s://%s:%d%s", protocol, ip, port, v)
	}
	return mFnURL, mapKeys(mFnURL).([]string)
}

func initTracer(serviceName string) opentracing.Tracer {
	cfg, err := config.FromEnv()
	failOnErr("%v: ", err)
	cfg.ServiceName = serviceName
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1

	tracer, _, err := cfg.NewTracer()
	failOnErr("%v: ", err)
	return tracer
}
