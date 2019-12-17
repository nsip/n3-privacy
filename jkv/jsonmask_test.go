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
	data := pp.FmtJSONFile("../../JSON-Mask/data/NAPCodeFrame.json", "../preprocess/utils/")
	mask1 := pp.FmtJSONFile("../../JSON-Mask/data/NAPCodeFrameMaskP.json", "../preprocess/utils/")
	mask2 := pp.FmtJSONFile("../../JSON-Mask/data/NAPCodeFrameMaskPcopy.json", "../preprocess/utils/")

	cmn.FailOnCondition(data == "", "%v", fEf("input data is empty, check its path"))
	cmn.FailOnCondition(mask1 == "", "%v", fEf("input mask1 is empty, check its path"))
	cmn.FailOnCondition(mask2 == "", "%v", fEf("input mask2 is empty, check its path"))

	jkvM1 := NewJKV(mask1, "root")
	jkvM2 := NewJKV(mask2, "root")

	if IsJSONArr(data) {
		jsonArr := SplitJSONArr(data)
		wg := sync.WaitGroup{}
		wg.Add(len(jsonArr))
		jsonList := make([]string, len(jsonArr))
		for i, json := range jsonArr {
			go func(i int, json string) {
				defer wg.Done()
				jkvD := NewJKV(json, "root")
				maskroot, _ := jkvD.Unfold(0, jkvM1)
				jkvMR := NewJKV(maskroot, "")
				jkvMR.Wrapped = jkvD.Wrapped
				jsonList[i] = jkvMR.UnwrapDefault().JSON
			}(i, json)
		}
		wg.Wait()
		ioutil.WriteFile("array.json", []byte(MergeJSON(jsonList...)), 0666)

	} else {

		jkvD := NewJKV(data, "root")
		maskroot, _ := jkvD.Unfold(0, jkvM1)
		jkvMR := NewJKV(maskroot, "")
		jkvMR.Wrapped = jkvD.Wrapped
		json := jkvMR.UnwrapDefault().JSON
		json = pp.FmtJSONStr(json, "../preprocess/utils/")

		jkvD = NewJKV(json, "root")
		maskroot, _ = jkvD.Unfold(0, jkvM2)
		jkvMR = NewJKV(maskroot, "")
		jkvMR.Wrapped = jkvD.Wrapped
		json = jkvMR.UnwrapDefault().JSON
		json = pp.FmtJSONStr(json, "../preprocess/utils/")

		ioutil.WriteFile("single.json", []byte(json), 0666)
	}
}
