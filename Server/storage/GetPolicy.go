package storage

// GetPolicyID :
func GetPolicyID(uid, ctx, object, rw string) (lsID []string) {
	oid := hash(object)[:lenOfOID]
	sid := hash(uid + ctx + rw)[:lenOfSID]
	for _, id := range lsMID {
		if sHasPrefix(id, oid) && sHasSuffix(id, sid) {
			lsID = append(lsID, id)
		}
	}
	return lsID
}

// GetPolicyHash :
func GetPolicyHash(code string) (string, bool) {
	if hashcode, ok := mMIDHash[code]; ok {
		return hashcode, ok
	}
	return "", false
}

// GetPolicy :
func GetPolicy(code string) (string, bool) {
	if mask, ok := mMIDMask[code]; ok {
		return mask, ok
	}
	return "", false
}
