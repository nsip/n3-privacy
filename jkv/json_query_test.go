package jkv

import (
	"io/ioutil"
	"testing"
	"time"

	cmn "github.com/nsip/n3-privacy/common"
)

func TestQuery(t *testing.T) {
	defer cmn.TmTrack(time.Now())
	param := "NAPTestItemLocalId"
	value := "x00101935"
	if jsonbytes, e := ioutil.ReadFile("../JSON-Mask/data/NAPCodeFrame.json"); e == nil {
		// jstr := jStr(string(jsonbytes))
		jkv := NewJKV(sReplaceAll(string(jsonbytes), "\r\n", "\n"), "")
		fPln("--- Init ---")

		path := func(string) string {
			return "NAPCodeFrame~~TestletList~~Testlet~~TestItemList~~TestItem~~TestItemContent~~NAPTestItemLocalId"
		}(param)

		//path1 := "NAPCodeFrame~~TestletList~~Testlet~~TestItemList~~TestItem~~TestItemContent~~NAPTestItemLocalId"
		//value1 := "\"x00101923-00-AIA\""
		// path2 := "NAPCodeFrame~~TestletList~~Testlet~~NAPTestletRefId"
		// value2 := "\"2b7c9606-09b9-43c2-a935-6a2db78bf2c9\""

		if mLvlOIDs, maxL := jkv.QueryPV(path, value); mLvlOIDs != nil && len(mLvlOIDs) > 0 {

			for _, oid := range mLvlOIDs[maxL] {
				fPln(oid, jkv.mOIDObj[oid])
			}

			// for _, lvl := range MapKeys(mLvlOIDs).([]int) {
			// 	for _, oid := range mLvlOIDs[lvl] {
			// 		fPf("[%s] %s\n", oid, mOIDObj[oid])
			// 		if mOIDType[oid].IsObjArr() {
			// 			fPf("ex: array object\n")
			// 			for _, oid := range AOIDStrToOIDs(mOIDObj[oid]) {
			// 				fPf("[%s] %s\n", oid, mOIDObj[oid])
			// 			}
			// 		}
			// 	}
			// 	fPln(" ----------------------------------------------------------------- ")
			// }

			// fPln(mOIDLvl["fe7262a928bbe05f8a42bab98ebec56a8e1e9379"])
			// fPln(mOIDLvl["887450b46a52ccad78f6a74f34c2699c649b17cd"])
		}
	}
}
