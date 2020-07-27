package db

import (
	"sort"

	"github.com/cdutwhu/n3-util/n3err"
)

func genPolicyID(policy, name, user, ctx, rw string) (string, string) {
	genPolicyCode := func(policy, name string) (string, string) {
		autoname := false
		if name == "" {
			jkvTmp := newJKV(policy, "", false)
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
		jkvM := newJKV(policy, name, false)
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
	if !isJSON(policy) {
		return "", n3err.JSON_INVALID
	}
	return fmtJSON(policy, 2), nil
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
