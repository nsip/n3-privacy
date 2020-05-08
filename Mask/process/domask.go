package process

import (
	"sync"
	"time"

	eg "github.com/cdutwhu/json-util/n3errs"
)

// DoMask :
func DoMask(data, mask string) (ret string) {
	data = fmtJSON(data, 2)
	mask = fmtJSON(mask, 2)
	jkvM := newJKV(mask, "root", false)

	if maybeJSONArr(data) {
		jsonArr := splitJSONArr(data, 2)
		wg := sync.WaitGroup{}
		wg.Add(len(jsonArr))
		jsonList := make([]string, len(jsonArr))
		for i, json := range jsonArr {
			go func(i int, json string) {
				defer wg.Done()
				jkvD := newJKV(json, "root", false)
				maskroot, _ := jkvD.Unfold(0, jkvM)
				jkvMR := newJKV(maskroot, "", false)
				jkvMR.Wrapped = jkvD.Wrapped
				jsonList[i] = jkvMR.UnwrapDefault().JSON
			}(i, json)
		}
		wg.Wait()
		ret = makeJSONArr(jsonList...)

	} else {
		jkvD := newJKV(data, "root", false)
		maskroot, _ := jkvD.Unfold(0, jkvM)
		jkvMR := newJKV(maskroot, "", false)
		jkvMR.Wrapped = jkvD.Wrapped
		json := jkvMR.UnwrapDefault().JSON
		ret = json
	}

	return ret
}

// FileMask :
func FileMask(inFilePath, maskFilePath, output string) {
	defer trackTime(time.Now())

	data := fmtJSONFile(inFilePath, 2)
	mask := fmtJSONFile(maskFilePath, 2)
	failOnErrWhen(data == "", "%v: check input file path", eg.FILE_EMPTY)
	failOnErrWhen(mask == "", "%v: check mask file path", eg.FILE_EMPTY)

	if output != "" {
		mustWriteFile(output, []byte(DoMask(data, mask)))
	}
}
