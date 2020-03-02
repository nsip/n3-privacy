package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	cmn "github.com/cdutwhu/json-util/common"
	jkv "github.com/cdutwhu/json-util/jkv"
)

func doMask(inFilePath, maskFilePath, output string) {
	defer cmn.TrackTime(time.Now())

	data := jkv.FmtJSONFile(inFilePath, 2)
	mask := jkv.FmtJSONFile(maskFilePath, 2)

	cmn.FailOnErrWhen(data == "", "%v", fEf("input data is empty, check path"))
	cmn.FailOnErrWhen(mask == "", "%v", fEf("input mask is empty, check path"))

	jkvM := jkv.NewJKV(mask, "root", false)

	if jkv.MaybeJSONArr(data) {
		jsonArr := jkv.SplitJSONArr(data, 2)
		wg := sync.WaitGroup{}
		wg.Add(len(jsonArr))
		jsonList := make([]string, len(jsonArr))
		for i, json := range jsonArr {
			go func(i int, json string) {
				defer wg.Done()
				jkvD := jkv.NewJKV(json, "root", false)
				maskroot, _ := jkvD.Unfold(0, jkvM)
				jkvMR := jkv.NewJKV(maskroot, "", false)
				jkvMR.Wrapped = jkvD.Wrapped
				jsonList[i] = jkvMR.UnwrapDefault().JSON
			}(i, json)
		}
		wg.Wait()
		ioutil.WriteFile(output, []byte(jkv.MakeJSONArray(jsonList...)), 0666)

	} else {
		jkvD := jkv.NewJKV(data, "root", false)
		maskroot, _ := jkvD.Unfold(0, jkvM)
		jkvMR := jkv.NewJKV(maskroot, "", false)
		jkvMR.Wrapped = jkvD.Wrapped
		json := jkvMR.UnwrapDefault().JSON
		ioutil.WriteFile(output, []byte(json), 0666)
	}
}

func main() {
	exe := filepath.Base(os.Args[0])
	if len(os.Args) < 3 {
		fPf("Usage: %s [-o='output'] <inputdata.json> <mask.json>\n", exe)
		return
	}

	inFilePath, maskFilePath := os.Args[1], os.Args[2]
	if len(os.Args) == 4 {
		inFilePath, maskFilePath = os.Args[2], os.Args[3]
	}

	outputPtr := flag.String("o", "out.json", "a string")
	flag.Parse()
	output := *outputPtr
	if !strings.HasSuffix(output, ".json") {
		output = output + ".json"
	}

	doMask(inFilePath, maskFilePath, output)
}
