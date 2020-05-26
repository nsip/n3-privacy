package main

import (
	"os"

	eg "github.com/cdutwhu/n3-util/n3errs"
	g "github.com/nsip/n3-privacy/Server/global"
	api "github.com/nsip/n3-privacy/Server/webapi"
)

func main() {
	failOnErrWhen(!g.Init("./config/config.toml"), "%v: Global Config Init Error", eg.CFG_INIT_ERR)

	cfg := g.Cfg
	ws, logfile, servicename := cfg.WebService, cfg.LogFile, cfg.ServiceName

	setLog(logfile)
	fPln(logWhen(true, "[%s] Hosting on: [%v:%d], version [%v]", ws.Service, localIP(), ws.Port, ws.Version))
	fPln(logWhen(true, "Working on Database: [%s]", cfg.Storage.DataBase))

	os.Setenv("JAEGER_SERVICE_NAME", servicename)
	os.Setenv("JAEGER_SAMPLER_TYPE", "const")
	os.Setenv("JAEGER_SAMPLER_PARAM", "1")

	done := make(chan string)
	go api.HostHTTPAsync()
	<-done
}
