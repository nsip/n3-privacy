package storage

import (
	"fmt"
	"strings"

	u "github.com/cdutwhu/go-util"
	cmn "github.com/nsip/n3-privacy/common"
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

	xin = u.XIn
)

const (
	linker = "#"
)

var (
	mMIDMask  = make(map[string]string)
	mMIDHash  = make(map[string]string)
	lsMID     []string // = u.MapKeys(mMIDMask).([]string)
	hash      = cmn.SHA256Str
	lenOfHash = len(hash("1"))
	lenOfOID  = lenOfHash / 4 // length of Object-Hash-ID Occupied
	lenOfFID  = lenOfHash / 4 // length of Fields-Hash-ID Occupied
	lenOfSID  = lenOfHash / 2 // length of Suffix-Hash-ID Occupied ( Suffix: UserID+ContextID+RW )

	mUIDlsCTX = make(map[string][]string)
	mCTXlsUID = make(map[string][]string)
)

// MetaData :
type MetaData struct {
	Object string   `json:"object"`
	Fields []string `json:"fields"`
}

// siLink :
func siLink(s string, i int) string {
	return fSf("%s%s%d", s, linker, i)
}

// ssLink :
func ssLink(s1, s2 string) string {
	return fSf("%s%s%s", s1, linker, s2)
}
