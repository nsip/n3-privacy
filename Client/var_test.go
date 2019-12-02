package main

import (
	"testing"

	glb "github.com/nsip/n3-privacy/Client/global"
)

func TestInitURL(t *testing.T) {
	glb.Init()
	protocol := glb.Cfg.Server.Protocol
	ip := glb.Cfg.Server.IP
	port := glb.Cfg.WebService.Port
	initMapFnURL(protocol, ip, port)
	fPln(mFnURL)
}
