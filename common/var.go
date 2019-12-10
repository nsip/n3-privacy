package common

import (
	"fmt"
	"regexp"
	"strings"
)

var (
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

	RExpMD5, _    = regexp.Compile("\"[a-f0-9]{32}\"")
	RExpSHA1, _   = regexp.Compile("\"[a-f0-9]{40}\"")
	RExpSHA256, _ = regexp.Compile("\"[a-f0-9]{64}\"")
)
