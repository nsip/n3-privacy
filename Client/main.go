package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"time"

	glb "github.com/nsip/n3-privacy/Client/global"
	cmn "github.com/nsip/n3-privacy/common"
	pp "github.com/nsip/n3-privacy/preprocess"
)

func main() {
	protocol, ip, port := "", "", 0
	cfgOK := glb.Init()

	protocolPtr := flag.String("protocol", "", "e.g. http/https/...")
	ipPtr := flag.String("ip", "", "Server IP address")
	portPtr := flag.Int("port", 0, "Server Port Number")
	fnPtr := flag.String("fn", "", "Select From: ["+sJoin(getCfgRouteFields(), " ")+"]")
	argsPtr := flag.String("args", "", "e.g. id=value1&user=value2&...")
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
			fPln("flag [-fn] is missing or invalid")
			fPln("Select From: [" + sJoin(getCfgRouteFields(), " ") + "]")
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
	case "GetID", "GetHash", "Get", "ListPolicyID", "ListUser", "ListContext", "ListObject": // GET
		go func() {
			resp, err := http.Get(url)
			cmn.FailOnErr("%v", err)
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			cmn.FailOnErr("%v", err)
			if resp.StatusCode == 200 {
				fPln(pp.FmtJSONStr(string(data)))
			} else {
				fPln(string(data))
			}
			done <- true
		}()
	case "Update": // POST

	case "Delete": // DELETE
	}

	select {
	case <-timeout:
		cmn.FailOnErr("%v", errors.New(fSf("Didn't Get Policy API Response in time. %d(s)", glb.Cfg.Access.Timeout)))
	case <-done:
	}
}
