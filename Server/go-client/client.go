package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	eg "github.com/cdutwhu/n3-util/n3errs"
)

// DO : fn ["HELP", ...]
func DO(configfile, fn string, args Args) (string, error) {
	failOnErrWhen(!initEnvVarFromTOML("Cfg-Clt-PRI", configfile), "%v", eg.CFG_INIT_ERR)

	Cfg := env2Struct("Cfg-Clt-PRI", &config{}).(*config)
	server := Cfg.Server
	protocol, ip, port := server.Protocol, server.IP, server.Port
	timeout := Cfg.Access.Timeout
	setLog(Cfg.LogFile)

	mFnURL, fields := initMapFnURL(protocol, ip, port, Cfg.Route)
	url, ok := mFnURL[fn]
	if err := warnOnErrWhen(!ok, "%v: Need ["+sJoin(fields, " ")+"]", eg.PARAM_NOT_SUPPORTED); err != nil {
		return "", err
	}

	chStr, chErr := make(chan string), make(chan error)
	go func() {
		rest(fn, url, args, chStr, chErr)
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		return "", warnOnErr("%v: Didn't get response in %d(s)", eg.NET_TIMEOUT, timeout)
	case str := <-chStr:
		err := <-chErr
		if err == eg.NO_ERROR {
			return str, nil
		}
		return str, err
	}
}

// rest :
func rest(fn, url string, args Args, chStr chan string, chErr chan error) {

	var (
		Resp *http.Response
		Err  error
		Data []byte
	)

	id, user, ctx, object, rw, policy, file := args.ID, args.User, args.Ctx, args.Object, args.RW, args.Policy, args.File

	switch fn {
	case "HELP":
		if Resp, Err = http.Get(url); Err != nil {
			goto ERR_RET
		}

	case "GetID":
		if user == "" || ctx == "" || object == "" || rw == "" {
			Err = warnOnErr("%v: [User] [Ctx] [Object] [RW] all are required", eg.PARAM_INVALID)
			goto ERR_RET
		}
		url += fSf("?user=%s&ctx=%s&object=%s&rw=%s", user, ctx, object, rw)
		if Resp, Err = http.Get(url); Err != nil {
			goto ERR_RET
		}

	case "GetHash", "Get":
		if id == "" {
			Err = warnOnErr("%v: [ID] is required", eg.PARAM_INVALID)
			goto ERR_RET
		}
		url += fSf("?id=%s", id)
		if Resp, Err = http.Get(url); Err != nil {
			goto ERR_RET
		}

	case "Update":
		if user == "" || ctx == "" || rw == "" || policy == "" {
			Err = warnOnErr("%v: [User] [Ctx] [RW] [Policy] all are required", eg.PARAM_INVALID)
			goto ERR_RET
		}
		if Data, Err = ioutil.ReadFile(policy); Err != nil {
			goto ERR_RET
		}
		if Err = warnOnErrWhen(!isJSON(string(Data)), "%v: policy", eg.PARAM_INVALID_JSON); Err != nil {
			goto ERR_RET
		}
		url += fSf("?name=%s&user=%s&ctx=%s&rw=%s", object, user, ctx, rw)
		if Resp, Err = http.Post(url, "application/json", bytes.NewBuffer(Data)); Err != nil {
			goto ERR_RET
		}

	case "Delete":
		if id == "" {
			Err = warnOnErr("%v: [ID] is required", eg.PARAM_INVALID)
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
		if user == "" || ctx == "" || rw == "" || file == "" {
			Err = warnOnErr("%v: [User] [Ctx] [RW] [File] all are required", eg.PARAM_INVALID)
			goto ERR_RET
		}
		if Data, Err = ioutil.ReadFile(file); Err != nil {
			goto ERR_RET
		}
		if Err = warnOnErrWhen(!isJSON(string(Data)), "%v: file", eg.PARAM_INVALID_JSON); Err != nil {
			goto ERR_RET
		}
		url += fSf("?name=%s&user=%s&ctx=%s&rw=%s", object, user, ctx, rw)
		if Resp, Err = http.Post(url, "application/json", bytes.NewBuffer(Data)); Err != nil {
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

	if Data, Err = ioutil.ReadAll(Resp.Body); Err != nil {
		goto ERR_RET
	}

ERR_RET:
	if Err != nil {
		chStr <- ""
		chErr <- Err
		return
	}

	chStr <- string(Data)
	chErr <- eg.NO_ERROR
	return
}
