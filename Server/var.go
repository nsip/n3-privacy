package main

import (
	"fmt"

	cmn "github.com/cdutwhu/n3-util/common"
)

var (
	fPln = fmt.Println

	failOnErrWhen = cmn.FailOnErrWhen
	localIP       = cmn.LocalIP
	setLog        = cmn.SetLog
	logWhen       = cmn.LogWhen
)
