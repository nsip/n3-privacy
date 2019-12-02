package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	glb "github.com/nsip/n3-privacy/Client/global"
	pp "github.com/nsip/n3-privacy/preprocess"
)

func main() {
	glb.Init()
	fmt.Println(glb.Cfg.Path)

	protocol := glb.Cfg.Server.Protocol
	ip := glb.Cfg.Server.IP
	port := glb.Cfg.WebService.Port
	fn := "ListOfObject"

	if ok := initMapFnURL(protocol, ip, port); ok {
		if resp, err := http.Get(mFnURL[fn]); err == nil {
			data, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(pp.FmtJSONStr(string(data)))
		}
	}
}
