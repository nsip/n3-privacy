package db

import (
	"fmt"
	"strings"

	cmn "github.com/cdutwhu/n3-util/common"
	"github.com/cdutwhu/n3-util/jkv"
	"github.com/cdutwhu/n3-util/n3json"
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

	isJSON           = cmn.IsJSON
	xin              = cmn.XIn
	failOnErr        = cmn.FailOnErr
	toSet            = cmn.ToSet
	indent           = cmn.Indent
	encrypt          = cmn.Encrypt
	decrypt          = cmn.Decrypt
	env2Struct       = cmn.Env2Struct
	failOnErrWhen    = cmn.FailOnErrWhen
	logger           = cmn.Log
	setLog           = cmn.SetLog
	tryInvoke        = cmn.TryInvoke
	tryInvokeWithMW  = cmn.TryInvokeWithMW
	mustInvokeWithMW = cmn.MustInvokeWithMW
	invokeRst        = cmn.InvokeRst
	hash             = cmn.SHA1Str // 32 [40] 64

	fmtJSON      = n3json.Fmt
	fmtJSONFile  = n3json.FmtFile
	maybeJSONArr = n3json.MaybeArr
	splitJSONArr = n3json.SplitArr
	makeJSONArr  = n3json.MakeArr

	newJKV = jkv.NewJKV
)

var (
	lenOfHash = len(hash("1")) //
	lenOfOID  = lenOfHash / 4  // length of Object-Hash-ID Occupied
	lenOfFID  = lenOfHash / 4  // length of Fields-Hash-ID Occupied
	lenOfUID  = lenOfHash / 4  // length of UserID-Hash-ID Occupied
	lenOfCTX  = lenOfHash / 4  // length of Context-Hash-ID Occupied
)

var (
	listID = []string{} // Policy ID List in running time
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
