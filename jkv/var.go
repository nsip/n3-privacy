// ********** ALL Based On JQ Formatted JSON ********** //

package jkv

import (
	"fmt"
	"strings"

	u "github.com/cdutwhu/go-util"
	w "github.com/cdutwhu/go-wrappers"
	cmn "github.com/nsip/n3-privacy/common"
)

type (
	b        = byte
	S        = w.Str
	JSONTYPE int
)

var (
	StartTrait = []byte{
		b('"'), // [array : string] OR [object : field]
		// b('{'), // [array : object]
		// b('n'),         // [array : null]
		// b('t'), b('f'), // [array : bool]
		// b('1'), b('2'), b('3'), b('4'), b('5'), b('6'), b('7'), b('8'), b('9'), b('-'), b('0'), // [array : number]
	}

	LF, SP, DQ = byte('\n'), byte(' '), byte('"')
)

var (
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
	sTrimRight  = strings.TrimRight
	sHasPrefix  = strings.HasPrefix
	sHasSuffix  = strings.HasSuffix
	IF          = u.IF
	MapKeys     = u.MapKeys
	MapKVs      = u.MapKVs
	MapsJoin    = u.MapsJoin
	MapsMerge   = u.MapsMerge
	MapPrint    = u.MapPrint
	SliceCover  = u.SliceCover
	MatchAssign = u.MatchAssign
	XIn         = u.XIn
)

var (
	BLANK = " \t\n\r"
	hash  = func(str string) string {
		return "\"" + cmn.SHA1Str(str) + "\""
	}
	// hash     = cmn.SHA1Str
	hashRExp = cmn.RExpSHA1 // compiled with ""
)

const (
	TraitScan = "\n                                                                " // 64 spaces

	AOS0  = "[\n  {\n    "                                                 // 2, 4
	AOE0  = "\n  }\n]"                                                     // 2, 0
	AOS1  = "[\n    {\n      "                                             // 4, 6
	AOE1  = "\n    }\n  ]"                                                 // 4, 2
	AOS2  = "[\n      {\n        "                                         // 6, 8
	AOE2  = "\n      }\n    ]"                                             // 6, 4
	AOS3  = "[\n        {\n          "                                     // 8, 10
	AOE3  = "\n        }\n      ]"                                         // 8, 6
	AOS4  = "[\n          {\n            "                                 // 10, 12
	AOE4  = "\n          }\n        ]"                                     // 10, 8
	AOS5  = "[\n            {\n              "                             // 12, 14
	AOE5  = "\n            }\n          ]"                                 // 12, 10
	AOS6  = "[\n              {\n                "                         // 14, 16
	AOE6  = "\n              }\n            ]"                             // 14, 12
	AOS7  = "[\n                {\n                  "                     // 16, 18
	AOE7  = "\n                }\n              ]"                         // 16, 14
	AOS8  = "[\n                  {\n                    "                 // 18, 20
	AOE8  = "\n                  }\n                ]"                     // 18, 16
	AOS9  = "[\n                    {\n                      "             // 20, 22
	AOE9  = "\n                    }\n                  ]"                 // 20, 18
	AOS10 = "[\n                      {\n                        "         // 22, 24
	AOE10 = "\n                      }\n                    ]"             // 22, 20
	AOS11 = "[\n                        {\n                          "     // 24, 26
	AOE11 = "\n                        }\n                      ]"         // 24, 22
	AOS12 = "[\n                          {\n                            " // 26, 28
	AOE12 = "\n                          }\n                        ]"     // 26, 24

	TraitFV    = "\": "
	Trait1EndV = ",\n"
	Trait2EndV = "\n"
	Trait3EndV = "],\n"
	Trait4EndV = "]\n"

	PathLinker = "~~"
	LvlMax     = 20 // init 20 max level in advances
)

// readonly var
var (
	sTAOStart = []string{AOS0, AOS1, AOS2, AOS3, AOS4, AOS5, AOS6, AOS7, AOS8, AOS9, AOS10, AOS11, AOS12}
	sTAOEnd   = []string{AOE0, AOE1, AOE2, AOE3, AOE4, AOE5, AOE6, AOE7, AOE8, AOE9, AOE10, AOE11, AOE12}
	pLinker   = PathLinker
)

// JKV :
type JKV struct {
	JSON          string
	LsL12Fields   [][]string          // 2D slice for each Level's each ifield
	lsLvlIPaths   [][]string          // 2D slice for each Level's each ipath
	mPathMAXIdx   map[string]int      //
	mIPathPos     map[string]int      //
	MapIPathValue map[string]string   //
	mIPathOID     map[string]string   //
	mOIDiPath     map[string]string   //
	mOIDObj       map[string]string   //
	mOIDLvl       map[string]int      // from 1 ...
	mOIDType      map[string]JSONTYPE // OID-type is OBJ or ARR|OBJ
	Wrapped       bool                //
}

// ********************************************************** //

// T : JSON line Search Feature.
func T(lvl int) string {
	return TraitScan[0 : 2*lvl+1]
}

// PT :
func PT(T string) string {
	return T[0 : len(T)-2]
}

// NT :
func NT(T string) string {
	return T[0 : len(T)+2]
}

// TL : get T & L by nchar
func TL(nChar int) (string, int) {
	lvl := (nChar - 1) / 2
	return T(lvl), lvl
}

// IndentFmt :
func IndentFmt(str string) (string, bool) {
	str = sTrim(str, BLANK)
	i := len(str) - 1
	N := 0
	if str[i] == '}' {
		for i = i - 1; i >= 0; i-- {
			if str[i] == ' ' {
				N++
				continue
			}
			break
		}
	}
	return Indent(str, -N, true)
}

// Indent :
func Indent(str string, n int, ignoreFirstLine bool) (string, bool) {
	if n == 0 {
		return str, false
	}
	S := 0
	if ignoreFirstLine {
		S = 1
	}
	lines := sSpl(str, "\n")
	if n > 0 {
		space := ""
		for i := 0; i < n; i++ {
			space += " "
		}
		for i := S; i < len(lines); i++ {
			if sTrim(lines[i], " \n\t") != "" {
				lines[i] = space + lines[i]
			}
		}
	} else {
		for i := S; i < len(lines); i++ {
			if len(lines[i]) == 0 { //                                         ignore empty string line
				continue
			}
			if len(lines[i]) <= -n || sTrimLeft(lines[i][:-n], " ") != "" { // cannot be indented as <n>, give up indent
				return str, false
			}
			lines[i] = lines[i][-n:]
		}
	}
	return sJoin(lines, "\n"), true
}

// ProjectV :
func ProjectV(strlist []string, sep, trimToL, trimFromR string) [][]string {
	nSep := 0
	for _, str := range strlist {
		if n := sCount(str, sep); n > nSep {
			nSep = n
		}
	}
	rtStrList := make([][]string, nSep+1)
	for _, str := range strlist {
		for i, s := range sSpl(str, sep) {
			if trimToL != "" {
				if fd := sIndex(s, trimToL); fd >= 0 {
					s = s[fd+1:]
				}
			}
			if trimFromR != "" {
				if fd := sLastIndex(s, trimFromR); fd >= 0 {
					s = s[:fd]
				}
			}
			rtStrList[i] = append(rtStrList[i], s)
		}
	}
	for i := 0; i < len(rtStrList); i++ {
		rtStrList[i] = cmn.ToSet(rtStrList[i]).([]string)
	}
	return rtStrList
}
