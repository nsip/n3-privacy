package webapi

import (
	"fmt"
	"net/url"
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

// ---------------------------------------------- //

// urlValues :
func urlValues(values url.Values, params ...string) (ok bool, lsValues [][]string) {
	for _, param := range params {
		if pv, ok := values[param]; ok {
			lsValues = append(lsValues, pv)
		}
	}
	if len(lsValues) == len(params) {
		return true, lsValues
	}
	return false, nil
}

func url1stValues(values url.Values, params ...string) (ok bool, ls1stValue []string) {
	if ok, lsValues := urlValues(values, params...); ok {
		for _, vs := range lsValues {
			ls1stValue = append(ls1stValue, vs[0])
		}
	}
	if len(params) == len(ls1stValue) {
		return true, ls1stValue
	}
	return false, nil
}

func url1stValuesOf1(values url.Values, params ...string) (bool, string) {
	if ok, ls1stValue := url1stValues(values, params...); ok {
		return true, ls1stValue[0]
	}
	return false, ""
}

func url1stValuesOf2(values url.Values, params ...string) (bool, string, string) {
	if ok, ls1stValue := url1stValues(values, params...); ok {
		return true, ls1stValue[0], ls1stValue[1]
	}
	return false, "", ""
}

func url1stValuesOf3(values url.Values, params ...string) (bool, string, string, string) {
	if ok, ls1stValue := url1stValues(values, params...); ok {
		return true, ls1stValue[0], ls1stValue[1], ls1stValue[2]
	}
	return false, "", "", ""
}

func url1stValuesOf4(values url.Values, params ...string) (bool, string, string, string, string) {
	if ok, ls1stValue := url1stValues(values, params...); ok {
		return true, ls1stValue[0], ls1stValue[1], ls1stValue[2], ls1stValue[3]
	}
	return false, "", "", "", ""
}

func url1stValuesOf5(values url.Values, params ...string) (bool, string, string, string, string, string) {
	if ok, ls1stValue := url1stValues(values, params...); ok {
		return true, ls1stValue[0], ls1stValue[1], ls1stValue[2], ls1stValue[3], ls1stValue[4]
	}
	return false, "", "", "", "", ""
}
