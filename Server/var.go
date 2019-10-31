package main

import (
	"fmt"
	"strings"

	u "github.com/cdutwhu/go-util"
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

	xin = u.XIn
)

const (
	linker = "#"
)

var (
	mUIDlsCtx  = make(map[string][]string) // UUID         - ContextID array
	mUIDlsMID  = make(map[string][]string) // UUID         - MaskID array
	mCtxlsMID  = make(map[string][]string) // ContextID    - MaskID array
	mMIDRWMask = make(map[string]string)   // MaskID#[R/W] - Mask
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
