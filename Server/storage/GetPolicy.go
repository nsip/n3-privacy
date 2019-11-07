package storage

// GetPolicyCode :
func GetPolicyCode(uid, ctx, object, rw string) (code string, ok bool) {
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
			code := ssLink(mid, rw)
			mask := mMIDRWMask[code]
			if sToLower(object) == sToLower(policyObject(mask)) {
				return code, true
			}
		}
	}
	return "", false
}

// GetPolicy :
func GetPolicy(code string) (string, bool) {
	if mask, ok := mMIDRWMask[code]; ok {
		return mask, ok
	}
	return "", false
}
