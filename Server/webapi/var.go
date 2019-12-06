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
	db   storage.DB
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

type result struct {
	Data  *string `json:"data"`
	Empty *bool   `json:"empty"`
	Error string  `json:"error"`
}

// Const
var (
	fFalse    = false
	fTrue     = true
	fEmptyStr = ""
	False     = &fFalse
	True      = &fTrue
	EmptyStr  = &fEmptyStr
)

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

// pick up one index-fixed value item from each values array, then combine them into one array
func urlOneValueList(values url.Values, idx int, params ...string) (ok bool, lsOneValue []string) {
	if ok, lsValues := urlValues(values, params...); ok {
		for _, vs := range lsValues {
			lsOneValue = append(lsOneValue, vs[idx])
		}
	}
	if len(params) == len(lsOneValue) {
		return true, lsOneValue
	}
	return false, nil
}

func url1Value(values url.Values, idx int, params ...string) (bool, string) {
	if ok, ls1Value := urlOneValueList(values, idx, params...); ok {
		return true, ls1Value[0]
	}
	return false, ""
}

func url2Values(values url.Values, idx int, params ...string) (bool, string, string) {
	if ok, ls2Values := urlOneValueList(values, idx, params...); ok {
		return true, ls2Values[0], ls2Values[1]
	}
	return false, "", ""
}

func url3Values(values url.Values, idx int, params ...string) (bool, string, string, string) {
	if ok, ls3Values := urlOneValueList(values, idx, params...); ok {
		return true, ls3Values[0], ls3Values[1], ls3Values[2]
	}
	return false, "", "", ""
}

func url4Values(values url.Values, idx int, params ...string) (bool, string, string, string, string) {
	if ok, ls4Values := urlOneValueList(values, idx, params...); ok {
		return true, ls4Values[0], ls4Values[1], ls4Values[2], ls4Values[3]
	}
	return false, "", "", "", ""
}

func url5Values(values url.Values, idx int, params ...string) (bool, string, string, string, string, string) {
	if ok, ls5Values := urlOneValueList(values, idx, params...); ok {
		return true, ls5Values[0], ls5Values[1], ls5Values[2], ls5Values[3], ls5Values[4]
	}
	return false, "", "", "", "", ""
}

func url6Values(values url.Values, idx int, params ...string) (bool, string, string, string, string, string, string) {
	if ok, ls6Values := urlOneValueList(values, idx, params...); ok {
		return true, ls6Values[0], ls6Values[1], ls6Values[2], ls6Values[3], ls6Values[4], ls6Values[5]
	}
	return false, "", "", "", "", "", ""
}
