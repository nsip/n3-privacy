package db

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/nsip/n3-privacy/jkv"

	pp "github.com/nsip/n3-privacy/preprocess"
)

// MetaData :
type MetaData struct {
	Object string   `json:"object"`
	Fields []string `json:"fields"`
}

// RecordMeta :
func RecordMeta(policy, metafile string) (updated bool) {
	jkvM := jkv.NewJKV(policy, hash(policy))
	object := jkvM.LsL12Fields[1][0]
	md := &MetaData{Object: object, Fields: jkvM.LsL12Fields[2]}
	if b, e := json.Marshal(md); e == nil {
		newPolicy := pp.FmtJSONStr(string(b))
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
			ioutil.WriteFile(metafile, []byte(jkv.MergeJSON(policies...)), 0666)
		}
	}
	return
}