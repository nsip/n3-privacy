package main

import (
	"os"

	eg "github.com/cdutwhu/n3-util/n3errs"
	cfg "github.com/nsip/n3-privacy/Server/config"
	api "github.com/nsip/n3-privacy/Server/webapi"
)

func main() {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg"), "%v: Config Init Error", eg.CFG_INIT_ERR)

	ICfg, err := env2Struct("Cfg", &cfg.Config{})
	failOnErr("%v", err)
	Cfg := ICfg.(*cfg.Config)
	ws, logfile, servicename := Cfg.WebService, Cfg.LogFile, Cfg.ServiceName

	setLog(logfile)
	fPln(logger("[%s] Hosting on: [%v:%d], version [%v]", ws.Service, localIP(), ws.Port, ws.Version))
	fPln(logger("Working on Database: [%s]", Cfg.Storage.DataBase))

	os.Setenv("JAEGER_SERVICE_NAME", servicename)
	os.Setenv("JAEGER_SAMPLER_TYPE", "const")
	os.Setenv("JAEGER_SAMPLER_PARAM", "1")

	done := make(chan string)
	go api.HostHTTPAsync()
	<-done
}
