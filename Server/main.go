package main

import (
	"os"
	"os/signal"

	"github.com/cdutwhu/n3-util/n3err"
	cfg "github.com/nsip/n3-privacy/Server/config"
	api "github.com/nsip/n3-privacy/Server/webapi"
	"github.com/sirupsen/logrus"
)

func main() {
	failOnErrWhen(!cfg.InitEnvVarFromTOML("Cfg"), "%v: Config Init Error", n3err.CFG_INIT_ERR)

	Cfg := env2Struct("Cfg", &cfg.Config{}).(*cfg.Config)
	ws, service := Cfg.WebService, Cfg.Service

	// --- LOGGLY
	enableLoggly(true)
	setLogglyToken(Cfg.Loggly.Token)
	lrInit()
	// --- LOGGLY

	enableLog2F(true, Cfg.Log)
	msg := fSf("[%s] Hosting on: [%v:%d], version [%v]", service, localIP(), ws.Port, Cfg.Version)
	logger(msg)
	lrOut(logrus.Infof, msg) // --> LOGGLY

	msg = fSf("Working on Database: [%s]", Cfg.Storage.DB)
	logger(msg)
	lrOut(logrus.Infof, msg) // --> LOGGLY

	os.Setenv("JAEGER_SERVICE_NAME", service)
	os.Setenv("JAEGER_SAMPLER_TYPE", "const")
	os.Setenv("JAEGER_SAMPLER_PARAM", "1")

	done := make(chan string)
	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, os.Interrupt)
	go api.HostHTTPAsync(c, done)
	msg = <-done
	logger(msg)
	lrOut(logrus.Infof, msg) // --> LOGGLY
}
