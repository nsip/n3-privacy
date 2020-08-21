package main

import (
	"os"
	"os/signal"

	"github.com/cdutwhu/n3-util/n3cfg"
	"github.com/cdutwhu/n3-util/n3err"
)

func main() {
	Cfg := n3cfg.ToEnvN3privacyServer(map[string]string{
		"[s]": "Service",
		"[v]": "Version",
	}, envKey)
	failOnErrWhen(Cfg == nil, "%v: Config Init Error", n3err.CFG_INIT_ERR)

	ws, logfile, service := Cfg.WebService, Cfg.Log, Cfg.Service.(string)
	os.Setenv("JAEGER_SERVICE_NAME", service)
	os.Setenv("JAEGER_SAMPLER_TYPE", "const")
	os.Setenv("JAEGER_SAMPLER_PARAM", "1")

	// --- LOGGLY --- //
	setLoggly(true, Cfg.Loggly.Token, service)
	syncBindLog(true)

	enableLog2F(true, logfile)
	enableWarnDetail(false)
	logGrp.Do(fSf("local log file @ [%s]", logfile))
	logGrp.Do(fSf("[%s] Hosting on: [%v:%d], version [%v]", service, localIP(), ws.Port, Cfg.Version))
	logGrp.Do(fSf("Storage Database: [%s], @ [%s]", Cfg.Storage.DB, Cfg.Storage.DBPath))

	done := make(chan string)
	c := make(chan os.Signal)
	signal.Notify(c, os.Kill, os.Interrupt)
	go HostHTTPAsync(c, done)
	logGrp.Do(<-done)
}
