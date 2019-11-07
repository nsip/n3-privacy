package storage

import (
	cmn "github.com/nsip/n3-privacy/common"
	"github.com/nsip/n3-privacy/jkv"
)

// policyObject :
func policyObject(mask string) string {
	jkvM := jkv.NewJKV(mask, cmn.SHA1Str(mask))
	object := jkvM.LsL12Fields[1][0]
	return object
}
