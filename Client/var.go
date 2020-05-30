package main

import (
	"fmt"
	"strings"

	cmn "github.com/cdutwhu/n3-util/common"
)

var (
	fPt  = fmt.Print
	fPf  = fmt.Printf
	fPln = fmt.Println
	fSf  = fmt.Sprintf

	sJoin       = strings.Join
	sReplaceAll = strings.ReplaceAll

	setLog        = cmn.SetLog
	resetLog      = cmn.ResetLog
	failOnErr     = cmn.FailOnErr
	failOnErrWhen = cmn.FailOnErrWhen
	warnOnErrWhen = cmn.WarnOnErrWhen
	isFLog        = cmn.IsFLog
	isJSON        = cmn.IsJSON
	cfgRepl       = cmn.CfgRepl
	struct2Env    = cmn.Struct2Env
)

// Args is arguments for "Route"
type Args struct {
	ID     string
	User   string
	Ctx    string
	Object string
	RW     string
	Policy []byte
	Data   []byte
}
