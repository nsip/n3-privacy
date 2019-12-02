package main

import (
	"fmt"
	"reflect"

	glb "github.com/nsip/n3-privacy/Client/global"
)

var (
	fPf  = fmt.Printf
	fPln = fmt.Println
	fSf  = fmt.Sprintf
)

var (
	mFnURL = map[string]string{}
)

func initMapFnURL(protocol, ip string, port int) bool {
	v := reflect.ValueOf(glb.Cfg.Route)
	typeOfT := reflect.ValueOf(&glb.Cfg.Route).Elem().Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface().(string)
		mFnURL[typeOfT.Field(i).Name] = fSf("%s://%s:%d%s", protocol, ip, port, field)
	}
	return len(mFnURL) > 0
}
