package db

import (
	"fmt"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/endec"
	"github.com/cdutwhu/gotil/judge"
	"github.com/cdutwhu/gotil/rflx"
	"github.com/cdutwhu/gotil/str"
	"github.com/cdutwhu/n3-util/jkv"
	"github.com/cdutwhu/n3-util/n3json"
)

var (
	fP               = fmt.Print
	fPf              = fmt.Printf
	fPln             = fmt.Println
	fSf              = fmt.Sprintf
	sSpl             = strings.Split
	sJoin            = strings.Join
	sCount           = strings.Count
	sReplace         = strings.Replace
	sReplaceAll      = strings.ReplaceAll
	sIndex           = strings.Index
	sLastIndex       = strings.LastIndex
	sTrim            = strings.Trim
	sTrimLeft        = strings.TrimLeft
	sHasPrefix       = strings.HasPrefix
	sHasSuffix       = strings.HasSuffix
	sToLower         = strings.ToLower
	sToUpper         = strings.ToUpper
	isJSON           = judge.IsJSON
	exist            = judge.Exist
	indent           = str.IndentTxt
	encrypt          = endec.Encrypt
	decrypt          = endec.Decrypt
	hash             = endec.SHA1Str // 32 [40] 64
	failOnErrWhen    = fn.FailOnErrWhen
	logger           = fn.Logger
	enableLog2F      = fn.EnableLog2F
	failOnErr        = fn.FailOnErr
	tryInvoke        = rflx.TryInvoke
	tryInvokeWithMW  = rflx.TryInvokeWithMW
	mustInvokeWithMW = rflx.MustInvokeWithMW
	invokeRst        = rflx.InvokeRst
	gslc             = rflx.ToGeneralSlc
	toSet            = rflx.ToSet
	env2Struct       = rflx.Env2Struct
	fmtJSON          = n3json.Fmt
	fmtJSONFile      = n3json.FmtFile
	maybeJSONArr     = n3json.MaybeArr
	splitJSONArr     = n3json.SplitArr
	makeJSONArr      = n3json.MakeArr
	newJKV           = jkv.NewJKV
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
	envKey = "PRISvr"
	linker = "#"
)

func siLink(s string, i int) string {
	return fSf("%s%s%d", s, linker, i)
}

func ssLink(s1, s2 string) string {
	return fSf("%s%s%s", s1, linker, s2)
}
