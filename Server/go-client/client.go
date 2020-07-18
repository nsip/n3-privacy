package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"

	eg "github.com/cdutwhu/n3-util/n3errs"
	"github.com/opentracing/opentracing-go"
	tags "github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go/config"
)

// DOwithTrace :
func DOwithTrace(ctx context.Context, configfile, fn string, args *Args) (string, error) {
	failOnErrWhen(!initEnvVarFromTOML(envVarName, configfile), "%v", eg.CFG_INIT_ERR)
	Cfg := env2Struct(envVarName, &Config{}).(*Config)
	serviceName := Cfg.ServiceName

	if ctx != nil {
		if span := opentracing.SpanFromContext(ctx); span != nil {
			initTracer := func(serviceName string) opentracing.Tracer {
				cfg, err := config.FromEnv()
				failOnErr("%v: ", err)
				cfg.ServiceName = serviceName
				cfg.Sampler.Type = "const"
				cfg.Sampler.Param = 1
				tracer, _, err := cfg.NewTracer()
				failOnErr("%v: ", err)
				return tracer
			}

			tracer := initTracer(serviceName)
			span := tracer.StartSpan(fn, opentracing.ChildOf(span.Context()))
			defer span.Finish()
			tags.SpanKindRPCClient.Set(span)
			tags.PeerService.Set(span, serviceName)
			if args != nil {
				span.SetTag(fn, *args)
			}
			span.LogEvent("Test log")
			span.SetBaggageItem("Test baggage", "baggage")
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
	}
	return DO(configfile, fn, args)
}

// DO : fn ["HELP", ...]
func DO(configfile, fn string, args *Args) (string, error) {
	failOnErrWhen(!initEnvVarFromTOML(envVarName, configfile), "%v", eg.CFG_INIT_ERR)
	Cfg := env2Struct(envVarName, &Config{}).(*Config)

	server := Cfg.Server
	protocol, ip, port := server.Protocol, server.IP, server.Port
	timeout := Cfg.Access.Timeout
	setLog(Cfg.LogFile)

	mFnURL, fields := initMapFnURL(protocol, ip, port, &Cfg.Route)
	url, ok := mFnURL[fn]
	if _, err := warnOnErrWhen(!ok, "%v: Need %v", eg.PARAM_NOT_SUPPORTED, fields); err != nil {
		return "", err
	}

	chStr, chErr := make(chan string), make(chan error)
	go func() {
		rest(fn, url, args, chStr, chErr)
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		_, err := warnOnErr("%v: Didn't get response in %d(s)", eg.NET_TIMEOUT, timeout)
		return "", err
	case str := <-chStr:
		err := <-chErr
		if err == eg.NO_ERROR {
			return str, nil
		}
		return str, err
	}
}

// rest :
func rest(fn, url string, args *Args, chStr chan string, chErr chan error) {

	var (
		Resp    *http.Response
		Err     error
		RetData []byte

		id     string
		user   string
		ctx    string
		object string
		rw     string
		policy []byte
		data   []byte
	)

	if args != nil {
		id, user, ctx, object, rw, policy, data = args.ID, args.User, args.Ctx, args.Object, args.RW, args.Policy, args.Data
	}

	switch fn {
	case "HELP":
		if Resp, Err = http.Get(url); Err != nil {
			goto ERR_RET
		}

	case "GetID":
		if user == "" || ctx == "" || object == "" || rw == "" {
			_, Err = warnOnErr("%v: [User] [Ctx] [Object] [RW] all are required", eg.PARAM_INVALID)
			goto ERR_RET
		}
		url += fSf("?user=%s&ctx=%s&object=%s&rw=%s", user, ctx, object, rw)
		if Resp, Err = http.Get(url); Err != nil {
			goto ERR_RET
		}

	case "GetHash", "Get":
		if id == "" {
			_, Err = warnOnErr("%v: [ID] is required", eg.PARAM_INVALID)
			goto ERR_RET
		}
		url += fSf("?id=%s", id)
		if Resp, Err = http.Get(url); Err != nil {
			goto ERR_RET
		}

	case "Update":
		if user == "" || ctx == "" || rw == "" {
			_, Err = warnOnErr("%v: [User] [Ctx] [RW] all are required", eg.PARAM_INVALID)
			goto ERR_RET
		}
		if _, Err = warnOnErrWhen(!isJSON(string(policy)), "%v: policy", eg.PARAM_INVALID_JSON); Err != nil {
			goto ERR_RET
		}
		url += fSf("?name=%s&user=%s&ctx=%s&rw=%s", object, user, ctx, rw)
		if Resp, Err = http.Post(url, "application/json", bytes.NewBuffer(policy)); Err != nil {
			goto ERR_RET
		}

	case "Delete":
		if id == "" {
			_, Err = warnOnErr("%v: [ID] is required", eg.PARAM_INVALID)
			goto ERR_RET
		}
		url += fSf("?id=%s", id)
		if req, err := http.NewRequest("DELETE", url, nil); err != nil {
			Err = err
			goto ERR_RET
		} else {
			if Resp, Err = (&http.Client{}).Do(req); Err != nil {
				goto ERR_RET
			}
		}

	case "Enforce":
		if user == "" || ctx == "" || rw == "" {
			_, Err = warnOnErr("%v: [User] [Ctx] [RW] all are required", eg.PARAM_INVALID)
			goto ERR_RET
		}
		if _, Err = warnOnErrWhen(!isJSON(string(data)), "%v: input data", eg.PARAM_INVALID_JSON); Err != nil {
			goto ERR_RET
		}
		url += fSf("?name=%s&user=%s&ctx=%s&rw=%s", object, user, ctx, rw)
		if Resp, Err = http.Post(url, "application/json", bytes.NewBuffer(data)); Err != nil {
			goto ERR_RET
		}

	case "LsID", "LsContext", "LsUser", "LsObject":
		switch fn {
		case "LsContext":
			ctx = ""
		case "LsUser":
			user = ""
		}
		switch {
		case user != "" && ctx != "":
			url += fSf("?user=%s&ctx=%s", user, ctx)
		case user != "":
			url += fSf("?user=%s", user)
		case ctx != "":
			url += fSf("?ctx=%s", ctx)
		}
		if Resp, Err = http.Get(url); Err != nil {
			goto ERR_RET
		}

	default:
		Err = eg.PARAM_NOT_SUPPORTED
		goto ERR_RET
	}

	if Resp == nil {
		Err = eg.NET_NO_RESPONSE
		goto ERR_RET
	}
	defer Resp.Body.Close()

	if RetData, Err = ioutil.ReadAll(Resp.Body); Err != nil {
		goto ERR_RET
	}

ERR_RET:
	if Err != nil {
		chStr <- ""
		chErr <- Err
		return
	}

	chStr <- string(RetData)
	chErr <- eg.NO_ERROR
	return
}
