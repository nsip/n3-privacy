package webapi

import (
	"fmt"
	"reflect"
	"sync"

	glb "github.com/nsip/n3-privacy/Server/global"
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
	v := reflect.ValueOf(glb.Cfg.Route)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface().(string)
		mMtx[field] = &sync.Mutex{}
	}
}

func initDB() {
	db = storage.NewDB(glb.Cfg.Storage.DataBase)
}
