package storage

// Dump :
type Dump interface {
	ListCTXByUID(string) []string
	ListUIDByCTX(string) []string
}

// func dumpMIDByUID(uid string) {
// 	lsMID := mUIDlsMID[uid]
// 	fPln(uid, lsMID)
// }

// func dumpMIDByCtx(ctx string) {
// 	lsMID := mCTXlsMID[ctx]
// 	fPln(ctx, lsMID)
// }

// func dumpMIDByOBJ(obj string) {
// 	lsMID := mOBJlsMID[obj]
// 	fPln(obj, lsMID)
// }

// func dumpMaskByUID(uid string) {
// 	lsMID := mUIDlsMID[uid]
// 	for _, mid := range lsMID {
// 		mask := mMIDMask[mid]
// 		fPln(mid, mask)
// 	}
// }
