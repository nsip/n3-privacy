package jkv

import (
	"io/ioutil"
	"sync"
	"testing"
	"time"

	cmn "github.com/nsip/n3-privacy/common"
	pp "github.com/nsip/n3-privacy/preprocess"
)

func TestJSONPolicy(t *testing.T) {
	defer cmn.TmTrack(time.Now())
	data := pp.FmtJSONFile("../../JSON-Mask/data/xapi1.json", "../preprocess/utils/")
	mask := pp.FmtJSONFile("../../JSON-Mask/data/xapiMask.json", "../preprocess/utils/")

	jkvM := NewJKV(mask, "root")

	if IsJSONArr(data) {
		jsonarr := SplitJSONArr(data)
		wg := sync.WaitGroup{}
		wg.Add(len(jsonarr))
		jsons := make([]string, len(jsonarr))
		for i, json := range jsonarr {
			go func(i int, json string) {
				defer wg.Done()
				jkvD := NewJKV(json, "root")
				maskroot, _ := jkvD.Unfold(0, jkvM)
				jkvMR := NewJKV(maskroot, "")
				jkvMR.Wrapped = jkvD.Wrapped
				jsons[i] = jkvMR.UnwrapDefault().JSON
			}(i, json)
		}
		wg.Wait()
		ioutil.WriteFile("array.json", []byte(MergeJSON(jsons...)), 0666)

	} else {
		jkvD := NewJKV(data, "root")
		maskroot, _ := jkvD.Unfold(0, jkvM)
		jkvMR := NewJKV(maskroot, "")
		jkvMR.Wrapped = jkvD.Wrapped
		json := jkvMR.UnwrapDefault().JSON
		ioutil.WriteFile("single.json", []byte(json), 0666)
	}
}
