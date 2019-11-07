package storage

import cmn "github.com/nsip/n3-privacy/common"

// UpdatePolicy :
func UpdatePolicy(uid, ctx, rw, mask string) {
	mid := cmn.SHA1Str(mask)
	mMIDRWMask[ssLink(mid, rw)] = mask
	mUIDlsCtx[uid] = append(mUIDlsCtx[uid], ctx)
	mUIDlsMID[uid] = append(mUIDlsMID[uid], mid)
	mCtxlsMID[ctx] = append(mCtxlsMID[ctx], mid)
}
