package main

// ObjFromPolicy :
func ObjFromPolicy(policy string) string {
	return "xapi"
}

// AddPolicy :
func AddPolicy(ctx, uid, rw, policy string) {
	mCtxUNum[ctx]++
	idx := mCtxUNum[ctx]
	mCIdxUID[SILink(ctx, idx)] = uid
	object := ObjFromPolicy(policy)
	mUOrwMask[sJoin([]string{uid, object, rw}, linker)] = policy
}

func main() {
	fPln(mCtxUNum["ass"])
	AddPolicy("ctx", "qm", "r", "policy.json")
	fPln(mUOrwMask["qm#xapi#r"])
}
