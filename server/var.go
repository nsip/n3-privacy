package main

import (
	"fmt"
	"strings"
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
	sHasPrefix  = strings.HasPrefix
	sHasSuffix  = strings.HasSuffix
)

const (
	linker = "#"
)

var (
	mCtxUNum  = make(map[string]int)    // contextID       - User Count
	mCIdxUID  = make(map[string]string) // contextID#Index - UUID
	mUOrwMask = make(map[string]string) // UUID#Object#RW  - mask.json
)

// SILink :
func SILink(s string, i int) string {
	return fSf("%s%s%d", s, linker, i)
}
