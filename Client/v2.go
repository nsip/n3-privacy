package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	eg "github.com/cdutwhu/n3-util/n3errs"
	glb "github.com/nsip/n3-privacy/Client/global"
)

func v2(cfgOK bool) {

	failOnErrWhen(!cfgOK, "%v", eg.CFG_INIT_ERR)

	cfg := glb.Cfg
	protocol, ip, port, timeout, logfile := cfg.Server.Protocol, cfg.Server.IP, cfg.Server.Port, cfg.Access.Timeout, cfg.LogFile

	setLog(logfile)

	if e := warnOnErrWhen(len(os.Args) < 2, "%v: need ["+sJoin(getCfgRouteFields(), " ")+"]", eg.CLI_SUBCMD_ERR); e != nil {
		if isFLog() {
			fPf("*** %v   abort\n", e)
		}
		return
	}

	failOnErrWhen(!initMapFnURL(protocol, ip, port), "%v: MapFnURL", eg.INTERNAL_INIT_ERR)
	arg1 := os.Args[1]
	done := make(chan bool)

	go func() {
		var resp *http.Response = nil
		var err error = nil
		var data []byte = nil
		url := mFnURL[arg1]
		cmd := flag.NewFlagSet(arg1, flag.ExitOnError)
		id := cmd.String("id", "", "policy ID")
		user := cmd.String("u", "", "user")
		ctx := cmd.String("c", "", "context")
		object := cmd.String("o", "", "object")
		rw := cmd.String("rw", "", "read/write")
		policyPtr := cmd.String("p", "", "the path of policy which is to be uploaded")
		fullDump := cmd.Bool("f", false, "output all attributes content from response")
		dataPtr := cmd.String("d", "", "the path of json which is to be uploaded")
		cmd.Parse(os.Args[2:])

		mngMode := false

		switch arg1 {
		case "HELP":
			resp, err = http.Get(url)

		case "GetID":
			failOnErrWhen(*user == "", "%v: [-u] user is required", eg.CLI_FLAG_ERR)
			failOnErrWhen(*ctx == "", "%v: [-c] context is required", eg.CLI_FLAG_ERR)
			failOnErrWhen(*object == "", "%v: [-o] object is required", eg.CLI_FLAG_ERR)
			failOnErrWhen(*rw == "", "%v: [-rw] read/write is required", eg.CLI_FLAG_ERR)
			url += fSf("?user=%s&ctx=%s&object=%s&rw=%s", *user, *ctx, *object, *rw)
			resp, err = http.Get(url)

		case "GetHash", "Get":
			failOnErrWhen(*id == "", "%v: [-id] ID is required", eg.CLI_FLAG_ERR)
			url += fSf("?id=%s", *id)
			resp, err = http.Get(url)

		case "Update":
			failOnErrWhen(*user == "", "%v: [-u] user is required", eg.CLI_FLAG_ERR)
			failOnErrWhen(*ctx == "", "%v: [-c] context is required", eg.CLI_FLAG_ERR)
			failOnErrWhen(*rw == "", "%v: [-rw] read/write is required", eg.CLI_FLAG_ERR)
			warnOnErrWhen(*object == "", "%v: if [-o] object is ignored, an auto-name will be assigned", eg.CLI_FLAG_ERR)
			url += fSf("?name=%s&user=%s&ctx=%s&rw=%s", *object, *user, *ctx, *rw)
			failOnErrWhen(*policyPtr == "", "%v: [-p] policy is required", eg.CLI_FLAG_ERR)
			policy, err := ioutil.ReadFile(*policyPtr)
			failOnErr("%v: %v", err, "Is [-p] policy provided correctly?")
			failOnErrWhen(!isJSON(string(policy)), "%v: policy", eg.PARAM_INVALID_JSON)
			resp, err = http.Post(url, "application/json", bytes.NewBuffer(policy))

		case "Delete":
			failOnErrWhen(*id == "", "%v: [-id] ID is required", eg.CLI_FLAG_ERR)
			url += fSf("?id=%s", *id)
			req, err := http.NewRequest("DELETE", url, nil)
			failOnErr("%v", err)
			resp, err = (&http.Client{}).Do(req)

			// One Step Enforce:
		case "GetEnforced":
			failOnErrWhen(*user == "", "%v: [-u] user is required", eg.CLI_FLAG_ERR)
			failOnErrWhen(*ctx == "", "%v: [-c] context is required", eg.CLI_FLAG_ERR)
			failOnErrWhen(*rw == "", "%v: [-rw] read/write is required", eg.CLI_FLAG_ERR)
			warnOnErrWhen(*object == "", "%v: if [-o] object is ignored, an auto-name will be assigned", eg.CLI_FLAG_ERR)
			url += fSf("?name=%s&user=%s&ctx=%s&rw=%s", *object, *user, *ctx, *rw)
			failOnErrWhen(*dataPtr == "", "%v: [-d] json data is required", eg.CLI_FLAG_ERR)
			data, err := ioutil.ReadFile(*dataPtr)
			failOnErr("%v: %v", err, "Is [-d] json data provided correctly?")
			failOnErrWhen(!isJSON(string(data)), "%v: data", eg.PARAM_INVALID_JSON)
			resp, err = http.Post(url, "application/json", bytes.NewBuffer(data))

			// Management Functions:
		case "LsID", "LsContext", "LsUser", "LsObject":
			switch arg1 {
			case "LsContext":
				*ctx = ""
			case "LsUser":
				*user = ""
			}
			switch {
			case *user != "" && *ctx != "":
				url += fSf("?user=%s&ctx=%s", *user, *ctx)
			case *user != "":
				url += fSf("?user=%s", *user)
			case *ctx != "":
				url += fSf("?ctx=%s", *ctx)
			}
			resp, err = http.Get(url)
			mngMode = true

		default:
			failOnErr("%v: %v", eg.CLI_SUBCMD_UNKNOWN, arg1)
		}

		failOnErr("%v", err)
		defer resp.Body.Close()

		data, err = ioutil.ReadAll(resp.Body)
		failOnErr("%v", err)

		const SepLn = "-----------------------------"

		if *fullDump {
			fPf("accessing... %s\n%s\n", url, SepLn)
		}

		if data != nil {
			if os.Args[1] == "HELP" {
				fPt(string(data))
			} else {
				m := make(map[string]interface{})
				failOnErr("json.Unmarshal ... %v", json.Unmarshal(data, &m))
				if !mngMode {
					if *fullDump {
						if m["empty"] != nil && m["empty"] != "" {
							fPf("Empty? %v\n%s\n", m["empty"], SepLn)
						}
						if m["error"] != nil && m["error"] != "" {
							fPf("ERROR: %v\n%s\n", m["error"], SepLn)
						}
					}
					if m["data"] != nil && m["data"] != "" {
						fPf("%s\n", m["data"])
					}
				} else {
					key := ""
					switch {
					case *user != "" && *ctx != "":
						key = fSf("%s@%s", *user, *ctx)
					case *user != "":
						key = *user
					case *ctx != "":
						key = *ctx
					default:
						key = "all"
					}
					fPf("%s\n", m[key])
				}
			}
		}

		done <- true
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		failOnErr("%v: Didn't Get Response in %d(s)", eg.NET_TIMEOUT, timeout)
	case <-done:
	}
}
