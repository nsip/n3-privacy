package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"

	cmn "github.com/nsip/n3-privacy/common"
	"github.com/nsip/n3-privacy/jkv"
	pp "github.com/nsip/n3-privacy/preprocess"
)

// when a mask is coming, parse and record it to meta.json
func recPolicyMeta(mask, metafile string) (updated bool) {
	jkvM := jkv.NewJKV(mask, cmn.SHA1Str(mask))
	object := jkvM.LsL12Fields[1][0]
	md := &MetaData{Object: object, Fields: jkvM.LsL12Fields[2]}
	if b, e := json.Marshal(md); e == nil {
		newPolicy := pp.FmtJSONStr(string(b), "../../preprocess/utils")
		// first meta.
		if _, err := os.Stat(metafile); err != nil && os.IsNotExist(err) {
			newPolicy, _ = jkv.Indent(newPolicy, 2, false)
			newPolicy = fSf("[\n%s]", newPolicy)
			ioutil.WriteFile(metafile, []byte(newPolicy), 0666) // make sure meta.json is formated

		} else {
			b, _ = ioutil.ReadFile(metafile)
			policies := jkv.SplitJSONArr(string(b))
			// update
			for i, policy := range policies {
				md := MetaData{}
				json.Unmarshal([]byte(policy), &md)
				if md.Object == object {
					policies[i] = newPolicy
					updated = true
					break
				}
			}
			if !updated {
				policies = append(policies, newPolicy)
			}
			ioutil.WriteFile(metafile, []byte(jkv.MergeJSONs(policies...)), 0666)
		}
	}
	return
}
