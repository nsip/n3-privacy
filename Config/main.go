package main

import (
	"os"

	"github.com/cdutwhu/n3-util/n3cfg/strugen"
)

func main() {
	cfgsrc := "./cfg/config.go"
	os.Remove(cfgsrc)
	pkgname := "cfg"
	strugen.GenStruct("./config.toml", "Config", pkgname, cfgsrc)
	strugen.GenNewCfg(cfgsrc)
}
