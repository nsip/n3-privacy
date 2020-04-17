package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	glb "github.com/nsip/n3-privacy/Client/global"
)

func v2(cfgOK bool) {
	failOnErrWhen(!cfgOK, "%v", fEf("Config File Init Failed"))
	failOnErrWhen(len(os.Args) < 2, "%v", fEf("need subcommands: ["+sJoin(getCfgRouteFields(), " ")+"]"))
	failOnErrWhen(!initMapFnURL(glb.Cfg.Server.Protocol, glb.Cfg.Server.IP, glb.Cfg.WebService.Port), "%v", fEf("initMapFnURL Failed"))
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
		cmd.Parse(os.Args[2:])

		mngMode := false

		switch arg1 {
		case "GetID":
			failOnErrWhen(*user == "", "%v", fEf("[-u] user must be provided"))
			failOnErrWhen(*ctx == "", "%v", fEf("[-c] context must be provided"))
			failOnErrWhen(*object == "", "%v", fEf("[-o] object must be provided"))
			failOnErrWhen(*rw == "", "%v", fEf("[-rw] read/write must be provided"))
			url += fSf("?user=%s&ctx=%s&object=%s&rw=%s", *user, *ctx, *object, *rw)
			resp, err = http.Get(url)

		case "GetHash", "Get":
			failOnErrWhen(*id == "", "%v", fEf("[-id] ID must be provided"))
			url += fSf("?id=%s", *id)
			resp, err = http.Get(url)

		case "Update":
			failOnErrWhen(*user == "", "%v", fEf("[-u] user must be provided"))
			failOnErrWhen(*ctx == "", "%v", fEf("[-c] context must be provided"))
			failOnErrWhen(*rw == "", "%v", fEf("[-rw] read/write must be provided"))
			warnOnErrWhen(*object == "", "%v", fEf("if [-o] object is ignored, an auto-name will be assigned"))
			url += fSf("?name=%s&user=%s&ctx=%s&rw=%s", *object, *user, *ctx, *rw)
			failOnErrWhen(*policyPtr == "", "%v", fEf("[-p] policy must be provided"))
			policy, err := ioutil.ReadFile(*policyPtr)
			failOnErr("%v: %v", err, "Is [-p] policy provided correctly?")
			failOnErrWhen(!isJSON(string(policy)), "%v", fEf("policy is invalid JSON, abort"))
			resp, err = http.Post(url, "application/json", bytes.NewBuffer(policy))

		case "Delete":
			failOnErrWhen(*id == "", "%v", fEf("[-id] ID must be provided"))
			url += fSf("?id=%s", *id)
			req, err := http.NewRequest("DELETE", url, nil)
			failOnErr("%v", err)
			resp, err = (&http.Client{}).Do(req)

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
			failOnErr("%v", fEf("unknown subcommand: %v", arg1))
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

		done <- true
	}()

	select {
	case <-time.After(time.Duration(glb.Cfg.Access.Timeout) * time.Second):
		failOnErr("%v", fEf(fSf("Didn't Get Server Response in %d(s)", glb.Cfg.Access.Timeout)))
	case <-done:
	}
}
