package jkv

import (
	cmn "github.com/nsip/n3-privacy/common"
)

// QueryPV : value ("*.*") for no path checking
func (jkv *JKV) QueryPV(path string, value interface{}) (mLvlOIDs map[int][]string, maxLvl int) {
	mLvlOIDs = make(map[int][]string)
	nGen, valstr := sCount(path, pLinker), ""
	switch value.(type) {
	case string:
		valstr = fSf("\"%v\"", value)
	default:
		valstr = fSf("%v", value)
	}

	for i := 0; i < jkv.mPathMAXIdx[path]; i++ {
		iPath := fSf("%s@%d", path, i)
		if v, ok := jkv.MapIPathValue[iPath]; ok && v == valstr {
			pos, PIPath := jkv.mIPathPos[iPath], ""
			for upGen := 1; upGen <= nGen; upGen++ {
				pPath := S(iPath).RmTailFromLastN(pLinker, upGen).V()
				for j := 0; j < jkv.mPathMAXIdx[pPath]; j++ {
					piPath := fSf("%s@%d", pPath, j)
					pPos := jkv.mIPathPos[piPath]
					if pPos > pos {
						break
					}
					PIPath = piPath
				}
				if pid, ok := jkv.MapIPathValue[PIPath]; ok {
					if _, ok := jkv.mOIDObj[pid]; ok {
						iLvl := nGen - upGen + 1
						if !cmn.XIn(pid, mLvlOIDs[iLvl]) {
							mLvlOIDs[iLvl] = append(mLvlOIDs[iLvl], pid)
							if iLvl > maxLvl {
								maxLvl = iLvl
							}
						}
					}
					// break // if search only the first one, break here !
				}
			}
		}
	}
	return mLvlOIDs, maxLvl
}

// Query : unfinished ...
// func Query(paths []string, values []interface{}) map[int][]string {
// 	lPaths, lValues := len(paths), len(values)
// 	if lPaths != lValues {
// 		panic("paths' count & values' count are not same!")
// 	}

// 	mLvlOIDs, pathShared, maxLvl := make(map[int][]string), "", 0
// 	for i := 0; i < lPaths; i++ {
// 		path, value := paths[i], values[i]
// 		mlvloids, maxl := QueryPV(path, value)
// 		if len(mlvloids) == 0 {
// 			return nil
// 		}
// 		if i == 0 {
// 			mLvlOIDs, pathShared, maxLvl = mlvloids, path, maxl
// 			continue
// 		}

// 		pathShared = func(s1, s2 string) string {
// 			minl := int(math.Min(float64(len(s1)), float64(len(s2))))
// 			for i := 0; i < minl; i++ {
// 				if s1[i] != s2[i] {
// 					return s1[:i]
// 				}
// 			}
// 			return ""
// 		}(pathShared, path)

// 		if maxl > maxLvl {
// 			maxLvl = maxl
// 		}

// 		lvl := sCount(pathShared, pLinker)
// 		IDs1, IDs2 := mLvlOIDs[lvl], mlvloids[lvl]
// 	NEXT:
// 		for j, id1 := range IDs1 {
// 			for _, id2 := range IDs2 {
// 				if id1 == id2 {
// 					continue NEXT
// 				}
// 			}
// 			// remove id1 from IDs1
// 			IDs1[j] = IDs1[len(IDs1)-1]
// 			mLvlOIDs[lvl] = mLvlOIDs[lvl][:len(mLvlOIDs[lvl])-1]
// 		}
// 		if len(mLvlOIDs[lvl]) == 0 {
// 			return nil
// 		}

// 		// refresh mLvlIDs
// 		// if i > 0 {
// 		// 	for l := 1; l <= maxLvl; l++ {
// 		// 		IDs1, IDs2 = mLvlIDs[l], mlvlids[l]

// 		// 	}
// 		// }
// 	}
// 	return mLvlOIDs
// }
