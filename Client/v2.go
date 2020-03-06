package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	cmn "github.com/cdutwhu/json-util/common"
	glb "github.com/nsip/n3-privacy/Client/global"
)

func v2(cfgOK bool) {
	cmn.FailOnErrWhen(!cfgOK, "%v", fEf("Config File Init Failed"))
	cmn.FailOnErrWhen(len(os.Args) < 2, "%v", fEf("need subcommands: ["+sJoin(getCfgRouteFields(), " ")+"]"))
	cmn.FailOnErrWhen(!initMapFnURL(glb.Cfg.Server.Protocol, glb.Cfg.Server.IP, glb.Cfg.WebService.Port), "%v", fEf("initMapFnURL fatal"))

	timeout := time.After(time.Duration(glb.Cfg.Access.Timeout) * time.Second)
	done := make(chan bool)

	go func() {
		var resp *http.Response = nil
		var err error = nil
		var data []byte = nil
		url := mFnURL[os.Args[1]]
		cmd := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
		id := cmd.String("id", "", "policy ID")
		user := cmd.String("user", "", "user")
		ctx := cmd.String("ctx", "", "context")
		object := cmd.String("object", "", "object")
		rw := cmd.String("rw", "", "read/write")
		policyPtr := cmd.String("policy", "", "the path of policy which is to be uploaded")
		fullFlag := cmd.Bool("f", false, "output all attributes content from response")
		cmd.Parse(os.Args[2:])

		switch os.Args[1] {
		case "GetID":
			cmn.FailOnErrWhen(*user == "", "%v", fEf("[-user] must be provided"))
			cmn.FailOnErrWhen(*ctx == "", "%v", fEf("[-ctx] must be provided"))
			cmn.FailOnErrWhen(*object == "", "%v", fEf("[-object] must be provided"))
			cmn.FailOnErrWhen(*rw == "", "%v", fEf("[-rw] must be provided"))
			url += fSf("?user=%s&ctx=%s&object=%s&rw=%s", *user, *ctx, *object, *rw)
			resp, err = http.Get(url)

		case "GetHash", "Get":
			cmn.FailOnErrWhen(*id == "", "%v", fEf("[-id] must be provided"))
			url += fSf("?id=%s", *id)
			resp, err = http.Get(url)

		case "Update":
			cmn.FailOnErrWhen(*user == "", "%v", fEf("[-user] must be provided"))
			cmn.FailOnErrWhen(*ctx == "", "%v", fEf("[-ctx] must be provided"))
			cmn.FailOnErrWhen(*rw == "", "%v", fEf("[-rw] must be provided"))
			url += fSf("?user=%s&ctx=%s&rw=%s", *user, *ctx, *rw)
			cmn.FailOnErrWhen(*policyPtr == "", "%v", fEf("[-policy] must be provided"))
			policy, err := ioutil.ReadFile(*policyPtr)
			cmn.FailOnErr("%v: %v", err, "Is [-policy] provided correctly?")
			cmn.FailOnErrWhen(!cmn.IsJSON(string(policy)), "%v", fEf("policy is not valid JSON file, abort"))
			resp, err = http.Post(url, "application/json", bytes.NewBuffer(policy))

		case "Delete":
			cmn.FailOnErrWhen(*id == "", "%v", fEf("[-id] must be provided"))
			url += fSf("?id=%s", *id)
			req, err := http.NewRequest("DELETE", url, nil)
			cmn.FailOnErr("%v", err)
			resp, err = (&http.Client{}).Do(req)

		case "ListID", "ListContext", "ListUser", "ListObject":
			switch {
			case *user != "" && *ctx != "":
				url += fSf("?user=%s&ctx=%s", *user, *ctx)
			case *user != "":
				url += fSf("?user=%s", *user)
			case *ctx != "":
				url += fSf("?ctx=%s", *ctx)
			}
			resp, err = http.Get(url)

		default:
			cmn.FailOnErr("%v", fEf("unknown subcommand: %v", os.Args[1]))
		}

		cmn.FailOnErr("%v", err)
		defer resp.Body.Close()

		data, err = ioutil.ReadAll(resp.Body)
		cmn.FailOnErr("%v", err)

		// fPln(string(data))
		// Manage output TODO

		if data != nil {
			m := make(map[string]interface{})
			cmn.FailOnErr("json.Unmarshal ... %v", json.Unmarshal(data, &m))
			if *fullFlag {
				fPln("accessing ... " + url)
				fPln("-----------------------------")

				if m["empty"] != nil && m["empty"] != "" {
					fPf("Is Empty? %v\n", m["empty"])
				}
				fPln("-----------------------------")

				if m["error"] != nil && m["error"] != "" {
					fPf("ERROR: %v\n", m["error"])
				}
				fPln("-----------------------------")
			}
			if m["data"] != nil && m["data"] != "" {
				fPf("%s\n", m["data"])
			}
		}

		done <- true
	}()

	select {
	case <-timeout:
		cmn.FailOnErr("%v", fEf(fSf("Didn't Get Server Response in time. %d(s)", glb.Cfg.Access.Timeout)))
	case <-done:
	}
}
