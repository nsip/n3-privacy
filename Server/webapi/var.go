package webapi

import (
	"fmt"
	"reflect"
	"sync"

	g "github.com/nsip/n3-privacy/Server/global"
	"github.com/nsip/n3-privacy/Server/storage"
)

var (
	fSf  = fmt.Sprintf
	fPln = fmt.Println
)

var (
	mMtx = map[string]*sync.Mutex{}
	db   = storage.NewDB("map")
)

func initMutex() {
	v := reflect.ValueOf(g.Cfg.Route)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface().(string)
		mMtx[field] = &sync.Mutex{}
	}
}
