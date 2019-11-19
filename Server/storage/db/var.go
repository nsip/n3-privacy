package db

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	u "github.com/cdutwhu/go-util"
	cmn "github.com/nsip/n3-privacy/common"
	"github.com/nsip/n3-privacy/jkv"
	pp "github.com/nsip/n3-privacy/preprocess"
)

var (
	fP          = fmt.Print
	fPf         = fmt.Printf
	fPln        = fmt.Println
	fSf         = fmt.Sprintf
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
	xin         = u.XIn
)

var (
	hash      = cmn.SHA1Str
	lenOfHash = len(hash("1"))
	lenOfOID  = lenOfHash / 4 // length of Object-Hash-ID Occupied
	lenOfFID  = lenOfHash / 4 // length of Fields-Hash-ID Occupied
	lenOfSID  = lenOfHash / 2 // length of Suffix-Hash-ID Occupied ( Suffix: UserID+ContextID+RW )
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

func genPolicyID(policy, uid, ctx, rw string) string {
	genPolicyCode := func(policy string) string {
		jkvM := jkv.NewJKV(policy, hash(policy))
		object := jkvM.LsL12Fields[1][0]
		fields := jkvM.LsL12Fields[2]
		sort.Strings(fields)
		oCode := hash(object)[:lenOfOID]
		fCode := hash(sJoin(fields, ""))[:lenOfFID]
		return oCode + fCode
	}
	suffix := hash(uid + ctx + rw)[:lenOfSID]
	return genPolicyCode(policy) + suffix
}

func validate(policy string) (string, error) {
	if !jkv.IsJSON(policy) {
		return "", errors.New("Not a valid JSON")
	}
	return pp.FmtJSONStr(policy), nil
}

func getPolicyID(uid, ctx, rw string, objects ...string) (lsID []string) {
	suffix := hash(uid + ctx + rw)[:lenOfSID]
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
