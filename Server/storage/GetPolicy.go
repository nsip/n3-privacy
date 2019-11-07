package storage

// GetPolicy :
func GetPolicy(uid, ctx, object, rw string) (string, bool) {
	if xin(ctx, mUIDlsCtx[uid]) {
		lsMIDuc := []string{}
		for _, midu := range mUIDlsMID[uid] {
			for _, midc := range mCtxlsMID[ctx] {
				if midu == midc {
					lsMIDuc = append(lsMIDuc, midu)
				}
			}
		}
		for _, mid := range lsMIDuc {
			mask := mMIDRWMask[ssLink(mid, rw)]
			if sToLower(object) == sToLower(policyObject(mask)) {
				return mask, true
			}
		}
	}
	return "", false
}
