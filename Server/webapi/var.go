package webapi

import (
	"fmt"
	"reflect"
	"sync"

	g "github.com/nsip/n3-privacy/Server/global"
)

var (
	fSf  = fmt.Sprintf
	fPln = fmt.Println
)

var (
	mMtx = map[string]*sync.Mutex{}
)

func initMutex() {
	v := reflect.ValueOf(g.Cfg.Route)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface().(string)
		mMtx[field] = &sync.Mutex{}
	}
}
