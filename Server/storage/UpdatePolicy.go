package storage

import (
	"sort"

	u "github.com/cdutwhu/go-util"
	"github.com/nsip/n3-privacy/jkv"
)

func genPolicyCode(mask string) string {
	jkvM := jkv.NewJKV(mask, hash(mask))
	object := jkvM.LsL12Fields[1][0]
	fields := jkvM.LsL12Fields[2]
	sort.Strings(fields)
	oCode := hash(object)[:lenOfOID]
	fCode := hash(sJoin(fields, ""))[:lenOfFID]
	return oCode + fCode
}

func genPolicyID(mask, uid, ctx, rw string) string {
	code := genPolicyCode(mask)
	suffix := hash(uid + ctx + rw)[:lenOfSID]
	return code + suffix
}

// UpdatePolicy :
func UpdatePolicy(uid, ctx, rw, mask string) {
	id := genPolicyID(mask, uid, ctx, rw)
	mMIDMask[id] = mask
	mMIDHash[id] = hash(mask)
	lsMID = u.MapKeys(mMIDMask).([]string)
	//
	lsCTX := mUIDlsCTX[uid]
	if !xin(ctx, lsCTX) {
		mUIDlsCTX[uid] = append(lsCTX, ctx)
	}
	lsUID := mCTXlsUID[ctx]
	if !xin(uid, lsUID) {
		mCTXlsUID[ctx] = append(lsUID, uid)
	}
}
