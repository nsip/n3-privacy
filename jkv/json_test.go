package jkv

import (
	"errors"
	"io/ioutil"
	"testing"
	"time"

	cmn "github.com/nsip/n3-privacy/common"
	pp "github.com/nsip/n3-privacy/preprocess"
)

func TestSplitJSONArr(t *testing.T) {
	defer cmn.TmTrack(time.Now())
	jarrstr := pp.FmtJSONFile("../../JSON-Mask/data/xapi.json", "../preprocess/utils/")
	// jarrstr := pp.FmtJSONFile("../../Server/config/meta.json", "../preprocess/utils/")
	if jarrstr == "" {
		fPln("Read JSON file error")
		return
	}
	jarrstr = sReplaceAll(jarrstr, "\r\n", "\n")
	if arr := SplitJSONArr(jarrstr); arr != nil {
		jarrmstr := MergeJSON(arr...)
		fPln(jarrmstr)
		if jarrstr != jarrmstr {
			panic("abc")
		}
	} else {
		cmn.FailOnErr("%v", errors.New("non-formatted json array"))
	}
}

func TestScan(t *testing.T) {
	defer cmn.TmTrack(time.Now())
	json := pp.FmtJSONFile("../../JSON-Mask/data/NAPCodeFrame.json", "../preprocess/utils/")
	jkv := NewJKV(json, "")
	LVL, mLvlFParr, mFPosLvl, _ := jkv.scan()
	fPln("levels:", LVL)
	for k, v := range mLvlFParr {
		fPln(k, v)
	}
	for k, v := range mFPosLvl {
		fPln(k, v)
	}
}

func TestFieldByPos(t *testing.T) {
	defer cmn.TmTrack(time.Now())
	if jsonbytes, e := ioutil.ReadFile("../JSON-Mask/data/NAPCodeFrame.json"); e == nil {
		jkv := NewJKV(sReplaceAll(string(jsonbytes), "\r\n", "\n"), "")
		LVL, mLvlFParr, _, _ := jkv.scan()
		// for k, v := range mLvlFParr {
		// 	fPln(k, v)
		// }
		mFPosFNameList := jkv.fields(mLvlFParr)
		for i := 1; i <= LVL; i++ {
			fPln("---------------->", i)
			mFPosFName := mFPosFNameList[i]
			for k, v := range mFPosFName {
				_, t := jkv.fValueType(k)
				fPf("%-8d%-20s%-10s\n", k, v, t.Str())
				// if t.IsPrimitive() {
				// 	fPf("%-8d%-20s%-10s\n", k, v, t.Str())
				// } else {
				// 	fPf("%-8d%-20s\n", k, v)
				// }
			}
		}
	}
}

func TestFType(t *testing.T) {
	defer cmn.TmTrack(time.Now())
	if jsonbytes, e := ioutil.ReadFile("../JSON-Mask/data/NAPCodeFrame.json"); e == nil {
		jkv := NewJKV(sReplaceAll(string(jsonbytes), "\r\n", "\n"), "")
		value, typ := jkv.fValueType(1617)
		fPln(typ.Str())
		if typ == ARR|OBJ {
			objs := fValuesOnObjList(value)
			fPln(objs[1])
		}
	}
}

func TestInit(t *testing.T) {
	defer cmn.TmTrack(time.Now())
	if jsonbytes, e := ioutil.ReadFile("../JSON-Mask/data/NAPCodeFrame.json"); e == nil {
		NewJKV(sReplaceAll(string(jsonbytes), "\r\n", "\n"), "")
	}
	fPln("break")
}

func TestWrap(t *testing.T) {
	defer cmn.TmTrack(time.Now())
	if jsonbytes, e := ioutil.ReadFile("../JSON-Mask/data/xapi1.json"); e == nil {
		json := pp.FmtJSONStr(string(jsonbytes), "../preprocess/utils/")
		jkv := NewJKV(sReplaceAll(json, "\r\n", "\n"), "root")
		fPln("--- Init ---")
		fPln(jkv.JSON)

	}
}

func TestUnfold(t *testing.T) {
	defer cmn.TmTrack(time.Now())
	if jsonbytes, e := ioutil.ReadFile("../JSON-Mask/data/xapi1.json"); e == nil {
		// fPln(string(jsonbytes))

		json := pp.FmtJSONStr(string(jsonbytes), "../preprocess/utils/")
		jkv := NewJKV(sReplaceAll(json, "\r\n", "\n"), "root")
		fPln("--- Init ---")
		fPln(jkv.Wrapped)
		fPln(jkv.Unfold(0, nil))

		// fPln(jkv.mOIDLvl["fe7262a928bbe05f8a42bab98ebec56a8e1e9379"])
		// fPln(jkv.mOIDLvl["887450b46a52ccad78f6a74f34c2699c649b17cd"]).

		fPln(" -------------------------------------- ")

		jkv = jkv.UnwrapDefault()
		// fPln(jkv.Unfold(0, nil))
		fPln(jkv.JSON)
	}
}
