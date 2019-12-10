package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	glb "github.com/nsip/n3-privacy/Client/global"
	cmn "github.com/nsip/n3-privacy/common"
	"github.com/nsip/n3-privacy/jkv"
)

func v2(cfgOK bool) {
	cmn.FailOnCondition(!cfgOK, "%v", fEf("Config File Init Failed"))
	cmn.FailOnCondition(len(os.Args) < 2, "%v", fEf("need subcommands: ["+sJoin(getCfgRouteFields(), " ")+"]"))
	cmn.FailOnCondition(!initMapFnURL(glb.Cfg.Server.Protocol, glb.Cfg.Server.IP, glb.Cfg.WebService.Port), "%v", fEf("initMapFnURL fatal"))

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
		cmd.Parse(os.Args[2:])

		switch os.Args[1] {
		case "GetID":
			cmn.FailOnCondition(*user == "", "%v", fEf("[-user] must be provided"))
			cmn.FailOnCondition(*ctx == "", "%v", fEf("[-ctx] must be provided"))
			cmn.FailOnCondition(*object == "", "%v", fEf("[-object] must be provided"))
			cmn.FailOnCondition(*rw == "", "%v", fEf("[-rw] must be provided"))
			url += fSf("?user=%s&ctx=%s&object=%s&rw=%s", *user, *ctx, *object, *rw)
			fPln("accessing ... " + url)
			resp, err = http.Get(url)

		case "GetHash", "Get":
			cmn.FailOnCondition(*id == "", "%v", fEf("[-id] must be provided"))
			url += fSf("?id=%s", *id)
			fPln("accessing ... " + url)
			resp, err = http.Get(url)

		case "Update":
			cmn.FailOnCondition(*user == "", "%v", fEf("[-user] must be provided"))
			cmn.FailOnCondition(*ctx == "", "%v", fEf("[-ctx] must be provided"))
			cmn.FailOnCondition(*rw == "", "%v", fEf("[-rw] must be provided"))
			url += fSf("?user=%s&ctx=%s&rw=%s", *user, *ctx, *rw)
			fPln("accessing ... " + url)
			cmn.FailOnCondition(*policyPtr == "", "%v", fEf("[-policy] must be provided"))
			policy, err := ioutil.ReadFile(*policyPtr)
			cmn.FailOnErr("%v: %v", err, "Is [-policy] provided correctly?")
			cmn.FailOnCondition(!jkv.IsJSON(string(policy)), "%v", fEf("policy is not valid JSON file, abort"))
			resp, err = http.Post(url, "application/json", bytes.NewBuffer(policy))

		case "Delete":
			cmn.FailOnCondition(*id == "", "%v", fEf("[-id] must be provided"))
			url += fSf("?id=%s", *id)
			fPln("accessing ... " + url)
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
			fPln("accessing ... " + url)
			resp, err = http.Get(url)

		default:
			cmn.FailOnErr("%v", fEf("unknown subcommand: %v", os.Args[1]))
		}

		cmn.FailOnErr("%v", err)
		defer resp.Body.Close()
		data, err = ioutil.ReadAll(resp.Body)
		cmn.FailOnErr("%v", err)
		if data != nil {
			fPln(string(data))
		}

		done <- true
	}()

	select {
	case <-timeout:
		cmn.FailOnErr("%v", fEf(fSf("Didn't Get Server Response in time. %d(s)", glb.Cfg.Access.Timeout)))
	case <-done:
	}
}
