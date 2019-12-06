package main

import (
	"bytes"
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	glb "github.com/nsip/n3-privacy/Client/global"
	cmn "github.com/nsip/n3-privacy/common"
	"github.com/nsip/n3-privacy/jkv"
)

func main() {

	protocol, ip, port := "", "", 0
	cfgOK := glb.Init()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-u", "--usage", "usage", "-h", "--help", "help":
			fPf(`Usage:
[-f]                 mandatory, one of [%[1]v];
[-protocol]          optional, config file [protocol] applies when missing;
[-ip]                optional, config file [ip] applies when missing;
[-port]              optional, config file [port] applies when missing;
[-args]              mandatory for [GetID/GetHash/Get/Update/Delete]. 
[-policy]            conditional, only available for [-f=Update]. It is the path of policy which is to be uploaded

Examples:
Get a Policy ID:     "%[2]v -f=GetID -args='user=(nsip)&ctx=(ctx123)&object=(StaffPersonal)&rw=(r)'"
Get a Policy Hash:   "%[2]v -f=GetHash -args='id=(policy id)'"
Get a Policy:        "%[2]v -f=Get -args='id=(policy id)'"
Update a Policy:     "%[2]v -f=Update -policy='./policy.json' -args='user=(nsip)&ctx=(ctx123)&rw=(r)'"
Delete a Policy:     "%[2]v -f=Delete -args='id=(policy id)'"
List some Policy ID: "%[2]v -f=ListID [-args='']"
List some Context:   "%[2]v -f=ListContext [-args='']"
List some Users:     "%[2]v -f=ListUser [-args='']"
List some Objects:   "%[2]v -f=ListObject [-args='']"

[GetID/GetHash/Get/Update/Delete] are BASIC functions, return {"data": "***", "empty": true/false, "error": "***" }
If "error" is NOT empty string (""), other fields' outcome are useless!
"empty" is only apply to [GetID/GetHash/Get]; If "error" is ("") and "empty" is true, which means found nothing without errors.
For [Update/Delete], if successful, "data" returns processed (Policy-ID), otherwise it returns ("") and check "error" for details.

[ListID/ListContext/ListUser/ListObject] are MANAGEMENT functions, only return { "condition": array of [Policy-ID/Context/User/Object] }
No "error" or "empty" fields.
`, sJoin(getCfgRouteFields(), " "), "privacy-client")
			return
		}
	}

	protocolPtr := flag.String("protocol", "", "e.g. http/https/...")
	ipPtr := flag.String("ip", "", "Server IP address")
	portPtr := flag.Int("port", 0, "Server Port Number")
	fnPtr := flag.String("f", "", "Select From: ["+sJoin(getCfgRouteFields(), " ")+"]")
	argsPtr := flag.String("args", "", "e.g. id=value1&user=value2&...")
	policyPtr := flag.String("policy", "", "the path of policy which is to be uploaded")
	flag.Parse()

	if protocol = *protocolPtr; protocol == "" && cfgOK {
		protocol = glb.Cfg.Server.Protocol
	}
	if ip = *ipPtr; ip == "" && cfgOK {
		ip = glb.Cfg.Server.IP
	}
	if port = *portPtr; port == 0 && cfgOK {
		port = glb.Cfg.WebService.Port
	}

	if ok := initMapFnURL(protocol, ip, port); ok {
		if _, ok := mFnURL[*fnPtr]; !ok {
			fPln("flag [-f] is missing or invalid. use [-h] for help")
			return
		}
	} else {
		cmn.FailOnErr("%v", errors.New("initMapFnURL fatal"))
	}

	if *argsPtr != "" {
		*argsPtr = "?" + *argsPtr
	}
	url := mFnURL[*fnPtr] + *argsPtr
	fPln("accessing ... " + url)

	timeout := time.After(time.Duration(glb.Cfg.Access.Timeout) * time.Second)
	done := make(chan bool)

	switch *fnPtr {
	case "GetID", "GetHash", "Get", "ListID", "ListUser", "ListContext", "ListObject": // GET
		go func() {
			resp, err := http.Get(url)
			cmn.FailOnErr("%v", err)
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			cmn.FailOnErr("%v", err)

			// var objmap map[string]interface{}
			// json.Unmarshal(data, &objmap)
			// fPln(objmap["data"])
			// if !jkv.IsJSON(string(objmap["data"].(string))) {
			// 	panic("return error")
			// }

			if data != nil {
				fPln(string(data))
			}
			done <- true
		}()
	case "Update": // POST
		go func() {
			policy, err := ioutil.ReadFile(*policyPtr)
			cmn.FailOnErr("%v: %v", err, "Is [-policy] provided correctly?")
			if !jkv.IsJSON(string(policy)) {
				cmn.FailOnErr("%v", errors.New("policy is not valid JSON file, failed to upload"))
			}
			if resp, err := http.Post(url, "application/json", bytes.NewBuffer(policy)); err == nil {
				defer resp.Body.Close()
				data, err := ioutil.ReadAll(resp.Body)
				cmn.FailOnErr("%v", err)
				if data != nil {
					fPln(string(data))
				}
			}
			done <- true
		}()
	case "Delete": // DELETE
		go func() {
			client := &http.Client{}
			req, err := http.NewRequest("DELETE", url, nil)
			cmn.FailOnErr("%v", err)
			resp, err := client.Do(req)
			cmn.FailOnErr("%v", err)
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			cmn.FailOnErr("%v", err)
			if data != nil {
				fPln(string(data))
			}
			done <- true
		}()
	}

	select {
	case <-timeout:
		cmn.FailOnErr("%v", errors.New(fSf("Didn't Get Policy Server Response in time. %d(s)", glb.Cfg.Access.Timeout)))
	case <-done:
	}
}
