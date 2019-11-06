package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	cmn "github.com/nsip/n3-privacy/common"
	"github.com/nsip/n3-privacy/jkv"
	pp "github.com/nsip/n3-privacy/preprocess"
)

// when a mask is coming, parse and record it to meta.json
func recPolicy(mask string) (updated bool) {

	jkvM := jkv.NewJKV(mask, "ToBeNamed")
	if jkvM.Wrapped {
		fPln("wrapped")
	} else {
		fPln("Not Wrapped")
	}

	object := jkvM.LsLvlFields[1][0]
	// object = sSpl(object, "@")[0]

	fields := []string{}
	for _, field := range jkvM.LsLvlFields[2] {
		//ifield := sSpl(ipath, jkv.PathLinker)[1]
		// field := sSpl(ifield, "@")[0]
		fields = append(fields, field)
	}

	md := &MetaData{Object: object, Fields: fields}
	if b, e := json.Marshal(md); e == nil {
		newPolicy := pp.FmtJSONStr(string(b), "../preprocess/utils")

		// first meta.
		if _, err := os.Stat("./config/meta.json"); err != nil && os.IsNotExist(err) {
			newPolicy, _ = jkv.Indent(newPolicy, 2, false)
			newPolicy = fSf("[\n%s]", newPolicy)
			ioutil.WriteFile("./config/meta.json", []byte(newPolicy), 0666) // make sure meta.json is filled with formated json
			return
		}

		b, _ = ioutil.ReadFile("./config/meta.json")
		policies := jkv.SplitJSONArr(string(b))

		// update
		for i, policy := range policies {
			md := MetaData{}
			json.Unmarshal([]byte(policy), &md)
			if md.Object == object {
				policies[i] = newPolicy
				updated = true
			}
		}
		if !updated {
			policies = append(policies, newPolicy)
		}
		ioutil.WriteFile("./config/meta.json", []byte(jkv.MergeJSONs(policies...)), 0666)
	}

	return
}

// policyObject :
func policyObject(mask string) string {
	return "xapi"
}

// UpdatePolicy :
func UpdatePolicy(uid, ctx, rw, mask string) {
	mid := cmn.SHA1Str(mask)
	mMIDRWMask[ssLink(mid, rw)] = mask
	mUIDlsCtx[uid] = append(mUIDlsCtx[uid], ctx)
	mUIDlsMID[uid] = append(mUIDlsMID[uid], mid)
	mCtxlsMID[ctx] = append(mCtxlsMID[ctx], mid)
}

// GetPolicy :
func GetPolicy(uid, ctx, object, rw string) (string, bool) {
	if xin(ctx, mUIDlsCtx[uid]) {
		lsMIDu, lsMIDc := mUIDlsMID[uid], mCtxlsMID[ctx]
		lsMIDuc := []string{}
		for _, midu := range lsMIDu {
			for _, midc := range lsMIDc {
				if midu == midc {
					lsMIDuc = append(lsMIDuc, midu)
				}
			}
		}
		for _, mid := range lsMIDuc {
			mask := mMIDRWMask[ssLink(mid, rw)]
			if object == policyObject(mask) {
				return mask, true
			}
		}
	}
	return "", false
}

func main() {
	UpdatePolicy("qm", "ctx1", "r", "policy.json")
	policy, ok := GetPolicy("qm", "ctx1", "xapi", "r")
	fPln(policy, ok)

	md := &MetaData{Object: "A", Fields: []string{"B", "C"}}
	if b, e := json.Marshal(md); e == nil {
		fPln(string(b))
	}
}
