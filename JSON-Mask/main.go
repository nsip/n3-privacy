package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nsip/n3-privacy/jkv"
	pp "github.com/nsip/n3-privacy/preprocess"
)

func main() {
	exe := filepath.Base(os.Args[0])
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s [-o='output'] <inputdata.json> <mask.json>\n", exe)
		return
	}

	inputfp, maskfp := os.Args[1], os.Args[2]
	if len(os.Args) == 4 {
		inputfp, maskfp = os.Args[2], os.Args[3]
	}

	outputPtr := flag.String("o", "result.json", "a string")
	flag.Parse()
	output := *outputPtr
	if !strings.HasSuffix(output, ".json") {
		output = output + ".json"
	}

	data := pp.FmtJSONFile(inputfp)
	mask := pp.FmtJSONFile(maskfp)
	jkvM := jkv.NewJKV(mask, "root")

	if jkv.IsJSONArr(data) {
		jsonarr := jkv.SplitJSONArr(data)
		wg := sync.WaitGroup{}
		wg.Add(len(jsonarr))
		jsons := make([]string, len(jsonarr))
		for i, json := range jsonarr {
			go func(i int, json string) {
				defer wg.Done()
				jkvD := jkv.NewJKV(json, "root")
				maskroot, _ := jkvD.Unfold(0, jkvM)
				jkvMR := jkv.NewJKV(maskroot, "")
				jkvMR.Wrapped = jkvD.Wrapped
				jsons[i] = jkvMR.UnwrapDefault().JSON
			}(i, json)
		}
		wg.Wait()
		ioutil.WriteFile(output, []byte(jkv.MergeJSON(jsons...)), 0666)

	} else {
		jkvD := jkv.NewJKV(data, "root")
		maskroot, _ := jkvD.Unfold(0, jkvM)
		jkvMR := jkv.NewJKV(maskroot, "")
		jkvMR.Wrapped = jkvD.Wrapped
		json := jkvMR.UnwrapDefault().JSON
		ioutil.WriteFile(output, []byte(json), 0666)
	}
}
