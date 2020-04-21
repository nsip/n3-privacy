package main

import (
	"fmt"

	cmn "github.com/cdutwhu/json-util/common"
	jkv "github.com/cdutwhu/json-util/jkv"
	n3json "github.com/cdutwhu/json-util/n3json"
)

var (
	fPf = fmt.Printf

	failOnErrWhen = cmn.FailOnErrWhen
	trackTime     = cmn.TrackTime

	fmtJSONFile  = n3json.FmtFile
	maybeJSONArr = n3json.MaybeArr
	splitJSONArr = n3json.SplitArr
	makeJSONArr  = n3json.MakeArr
	newJKV       = jkv.NewJKV
)
