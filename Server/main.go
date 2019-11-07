package main

import (
	g "github.com/nsip/n3-privacy/Server/global"
	"github.com/nsip/n3-privacy/Server/webapi"
)

func main() {
	if !g.Init() {
		panic("global Init Error")
	}
	done := make(chan string)
	go webapi.HostHTTPAsync()
	<-done
}
