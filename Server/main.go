package main

import (
	"log"

	cmn "github.com/cdutwhu/json-util/common"
	eg "github.com/cdutwhu/json-util/n3errs"
	g "github.com/nsip/n3-privacy/Server/global"
	api "github.com/nsip/n3-privacy/Server/webapi"
)

func main() {
	cmn.FailOnErrWhen(!g.Init(), "%v: Global Config Init Error", eg.CFG_INIT_ERR)
	log.Printf("Working on Database: [%s]", g.Cfg.Storage.DataBase)
	done := make(chan string)
	go api.HostHTTPAsync()
	<-done
}
