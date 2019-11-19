package main

import (
	"log"

	g "github.com/nsip/n3-privacy/Server/global"
	api "github.com/nsip/n3-privacy/Server/webapi"
)

func main() {
	if !g.Init() {
		panic("Global Config Init Error")
	}
	log.Printf("Working on Database: [%s]", g.Cfg.Storage.DataBase)
	done := make(chan string)
	go api.HostHTTPAsync()
	<-done
}
