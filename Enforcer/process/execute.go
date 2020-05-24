package process

import (
	"sync"
	"time"

	eg "github.com/cdutwhu/n3-util/n3errs"
)

// Execute :
func Execute(data, policy string) (ret string) {
	data = fmtJSON(data, 2)
	policy = fmtJSON(policy, 2)
	jkvP := newJKV(policy, "root", false)

	if maybeJSONArr(data) {
		jsonArr := splitJSONArr(data, 2)
		wg := sync.WaitGroup{}
		wg.Add(len(jsonArr))
		jsonList := make([]string, len(jsonArr))
		for i, json := range jsonArr {
			go func(i int, json string) {
				defer wg.Done()
				jkvD := newJKV(json, "root", false)
				all, _ := jkvD.Unfold(0, jkvP)
				jkvEnforced := newJKV(all, "", false)
				jkvEnforced.Wrapped = jkvD.Wrapped
				jsonList[i] = jkvEnforced.UnwrapDefault().JSON
			}(i, json)
		}
		wg.Wait()
		ret = makeJSONArr(jsonList...)

	} else {
		jkvD := newJKV(data, "root", false)
		all, _ := jkvD.Unfold(0, jkvP)
		jkvEnforced := newJKV(all, "", false)
		jkvEnforced.Wrapped = jkvD.Wrapped
		json := jkvEnforced.UnwrapDefault().JSON
		ret = json
	}

	return ret
}

// FileExe :
func FileExe(inFilePath, policyFilePath, output string) {
	defer trackTime(time.Now())

	data := fmtJSONFile(inFilePath, 2)
	policy := fmtJSONFile(policyFilePath, 2)
	failOnErrWhen(data == "", "%v: check input file path", eg.FILE_EMPTY)
	failOnErrWhen(policy == "", "%v: check policy file path", eg.FILE_EMPTY)

	if output != "" {
		mustWriteFile(output, []byte(Execute(data, policy)))
	}
}
