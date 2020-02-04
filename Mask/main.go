package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	cmn "github.com/cdutwhu/json-util/common"
	"github.com/cdutwhu/json-util/jkv"
	pp "github.com/cdutwhu/json-util/preprocess"
	cfg "github.com/nsip/n3-privacy/Mask/config"
)

func main() {
	exe := filepath.Base(os.Args[0])
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s [-o='output'] <inputdata.json> <mask.json>\n", exe)
		return
	}

	inFilePath, maskFilePath := os.Args[1], os.Args[2]
	if len(os.Args) == 4 {
		inFilePath, maskFilePath = os.Args[2], os.Args[3]
	}

	outputPtr := flag.String("o", "result.json", "a string")
	flag.Parse()
	output := *outputPtr
	if !strings.HasSuffix(output, ".json") {
		output = output + ".json"
	}

	config := cfg.NewCfg("./Config.toml").(*cfg.Config) // Config.toml is hard-coded

	data := pp.FmtJSONFile(inFilePath, config.JQDir)
	mask := pp.FmtJSONFile(maskFilePath, config.JQDir)
	cmn.FailOnErrWhen(data == "", "%v", fmt.Errorf("input data is empty, check path"))
	cmn.FailOnErrWhen(mask == "", "%v", fmt.Errorf("input mask is empty, check path"))

	jkvM := jkv.NewJKV(mask, "root", false)

	if jkv.IsJSONArr(data) {
		jsonArr := jkv.SplitJSONArr(data)
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
		ioutil.WriteFile(output, []byte(jkv.MergeJSON(jsonList...)), 0666)

	} else {
		jkvD := jkv.NewJKV(data, "root", false)
		maskroot, _ := jkvD.Unfold(0, jkvM)
		jkvMR := jkv.NewJKV(maskroot, "", false)
		jkvMR.Wrapped = jkvD.Wrapped
		json := jkvMR.UnwrapDefault().JSON
		ioutil.WriteFile(output, []byte(json), 0666)
	}
}
