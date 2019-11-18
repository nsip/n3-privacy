package main

import (
	g "github.com/nsip/n3-privacy/Server/global"
	api "github.com/nsip/n3-privacy/Server/webapi"
)

func main() {
	if !g.Init() {
		panic("Global Config Init Error")
	}
	done := make(chan string)
	go api.HostHTTPAsync()
	<-done
}
