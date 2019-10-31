package main

import (
	"encoding/json"

	cmn "../common"
)

func parsePolicy(mask string) string {
	return "xapi"
}

// PolicyObject :
func PolicyObject(mask string) string {
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
			if object == PolicyObject(mask) {
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
