package main

import (
	"fmt"

	cmn "github.com/cdutwhu/json-util/common"
)

var (
	fEf = fmt.Errorf
	fPf = fmt.Printf

	failOnErrWhen = cmn.FailOnErrWhen
	trackTime     = cmn.TrackTime
)
