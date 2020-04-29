package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	eg "github.com/cdutwhu/json-util/n3errs"
)

func doMask(inFilePath, maskFilePath, output string) {
	defer trackTime(time.Now())

	data := fmtJSONFile(inFilePath, 2)
	mask := fmtJSONFile(maskFilePath, 2)

	failOnErrWhen(data == "", "%v: check input file path", eg.FILE_EMPTY)
	failOnErrWhen(mask == "", "%v: check mask file path", eg.FILE_EMPTY)

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
		ioutil.WriteFile(output, []byte(makeJSONArr(jsonList...)), 0666)

	} else {
		jkvD := newJKV(data, "root", false)
		maskroot, _ := jkvD.Unfold(0, jkvM)
		jkvMR := newJKV(maskroot, "", false)
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
