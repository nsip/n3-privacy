package db

import (
	"fmt"
	"sort"
	"strings"

	cmn "github.com/cdutwhu/json-util/common"
	"github.com/cdutwhu/json-util/jkv"
)

var (
	fP          = fmt.Print
	fPf         = fmt.Printf
	fPln        = fmt.Println
	fSf         = fmt.Sprintf
	fEf         = fmt.Errorf
	sSpl        = strings.Split
	sJoin       = strings.Join
	sCount      = strings.Count
	sReplace    = strings.Replace
	sReplaceAll = strings.ReplaceAll
	sIndex      = strings.Index
	sLastIndex  = strings.LastIndex
	sTrim       = strings.Trim
	sTrimLeft   = strings.TrimLeft
	sHasPrefix  = strings.HasPrefix
	sHasSuffix  = strings.HasSuffix
	sToLower    = strings.ToLower
	sToUpper    = strings.ToUpper
	xin         = cmn.XIn
)

var (
	hash      = cmn.SHA1Str // 32 [40] 64
	lenOfHash = len(hash("1"))
	lenOfOID  = lenOfHash / 4 // length of Object-Hash-ID Occupied
	lenOfFID  = lenOfHash / 4 // length of Fields-Hash-ID Occupied
	lenOfUID  = lenOfHash / 4 // length of UserID-Hash-ID Occupied
	lenOfCTX  = lenOfHash / 4 // length of Context-Hash-ID Occupied
	listID    = []string{}    // Policy ID List in running time
)

const (
	linker = "#"
)

func siLink(s string, i int) string {
	return fSf("%s%s%d", s, linker, i)
}

func ssLink(s1, s2 string) string {
	return fSf("%s%s%s", s1, linker, s2)
}

func genPolicyID(policy, name, user, ctx, rw string) (string, string) {

	genPolicyCode := func(policy, name string) (string, string) {
		autoname := false
		if name == "" {
			jkvTmp := jkv.NewJKV(policy, "", false)
			attris := jkvTmp.LsL12Fields[1]
			if len(attris) == 1 {
				name = attris[0]
			} else {
				sort.Slice(attris, func(i, j int) bool {
					return attris[i][0] < attris[j][0]
				})
				for _, a := range attris {
					name += string(a[0])
				}
			}
			autoname = true
		}
		// fPln(name)
		jkvM := jkv.NewJKV(policy, name, false)
		object := jkvM.LsL12Fields[1][0]
		if !autoname {
			object = name
		}

		fields := jkvM.LsL12Fields[2]
		sort.Strings(fields)
		oCode := hash(object)[:lenOfOID]
		fCode := hash(sJoin(fields, ""))[:lenOfFID]
		return oCode + fCode, object
	}
	uCode := hash(user)[:lenOfUID]
	cCode := hash(ctx)[:lenOfCTX]
	pCode, object := genPolicyCode(policy, name)
	return pCode + uCode + cCode + rw[:1], object
}

func oCodeByPID(pid string) string {
	return pid[:lenOfOID]
}

func fCodeByPID(pid string) string {
	return pid[lenOfOID : lenOfOID+lenOfFID]
}

func uCodeByPID(pid string) string {
	return pid[lenOfOID+lenOfFID : lenOfOID+lenOfFID+lenOfUID]
}

func cCodeByPID(pid string) string {
	return pid[lenOfOID+lenOfFID+lenOfUID : lenOfOID+lenOfFID+lenOfUID+lenOfCTX]
}

func validate(policy string) (string, error) {
	if !cmn.IsJSON(policy) {
		return "", fEf("Not a valid JSON")
	}
	// return pp.FmtJSONStr(policy), nil
	return jkv.FmtJSON(policy, 2), nil
}

// [listID] has already been loaded
func getPolicyID(user, ctx, rw string, objects ...string) (lsID []string) {
	suffix := hash(user)[:lenOfUID] + hash(ctx)[:lenOfCTX] + rw[:1]
	if len(objects) > 0 {
		for _, object := range objects {
			oid := hash(object)[:lenOfOID]
			for _, id := range listID {
				if sHasPrefix(id, oid) && sHasSuffix(id, suffix) {
					lsID = append(lsID, id)
				}
			}
		}
	} else {
		for _, id := range listID {
			if sHasSuffix(id, suffix) {
				lsID = append(lsID, id)
			}
		}
	}
	return lsID
}
