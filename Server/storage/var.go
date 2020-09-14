package storage

import (
	"fmt"
	"strings"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/judge"
)

var (
	fP              = fmt.Print
	fPf             = fmt.Printf
	fPln            = fmt.Println
	fSf             = fmt.Sprintf
	sSpl            = strings.Split
	sJoin           = strings.Join
	sCount          = strings.Count
	sReplace        = strings.Replace
	sReplaceAll     = strings.ReplaceAll
	sIndex          = strings.Index
	sLastIndex      = strings.LastIndex
	sTrim           = strings.Trim
	sTrimLeft       = strings.TrimLeft
	sHasPrefix      = strings.HasPrefix
	sHasSuffix      = strings.HasSuffix
	sToLower        = strings.ToLower
	sToUpper        = strings.ToUpper
	failOnErr       = fn.FailOnErr
	failOnErrWhen   = fn.FailOnErrWhen
	failP1OnErrWhen = fn.FailP1OnErrWhen
	exist           = judge.Exist
)
