package jkv

import (
	"io/ioutil"
	"sync"
	"testing"
	"time"

	cmn "../common"
	pp "../preprocess"
)

func TestJSONPolicy(t *testing.T) {
	defer cmn.TmTrack(time.Now())
	data := pp.FmtJSONFile("../../JSON-Mask/data/xapi.json", "../preprocess/utils/")
	data = sReplaceAll(data, "\r\n", "\n")
	mask := pp.FmtJSONFile("../../JSON-Mask/data/xapiMask.json", "../preprocess/utils/")
	mask = sReplaceAll(mask, "\r\n", "\n")

	jkvM := NewJKV(mask, "root")

	if IsJSONArr(data) {
		jsonarr := SplitJSONArr(data)
		wg := sync.WaitGroup{}
		wg.Add(len(jsonarr))
		jsons := make([]string, len(jsonarr))
		for i, json := range jsonarr {
			// jkvD := NewJKV(json, "root")
			// maskroot, _ := jkvD.Unfold(0, jkvM.MapIPathValue)
			// jkvMR := NewJKV(maskroot, "")
			// jkvMR.Wrapped = jkvD.Wrapped
			// jsons[i] = jkvMR.UnwrapDefault().JSON

			go func(i int, json string) {
				defer wg.Done()
				jkvD := NewJKV(json, "root")
				maskroot, _ := jkvD.Unfold(0, jkvM.MapIPathValue)
				jkvMR := NewJKV(maskroot, "")
				jkvMR.Wrapped = jkvD.Wrapped
				jsons[i] = jkvMR.UnwrapDefault().JSON
			}(i, json)
		}
		wg.Wait()
		ioutil.WriteFile("array.json", []byte(MergeJSONs(jsons...)), 0666)

	} else {
		jkvD := NewJKV(data, "root")
		maskroot, _ := jkvD.Unfold(0, jkvM.MapIPathValue)
		jkvMR := NewJKV(maskroot, "")
		jkvMR.Wrapped = jkvD.Wrapped
		json := jkvMR.UnwrapDefault().JSON
		ioutil.WriteFile("single.json", []byte(json), 0666)
	}
}