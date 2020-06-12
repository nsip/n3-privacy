package process

import (
	"fmt"

	cmn "github.com/cdutwhu/n3-util/common"
	jkv "github.com/cdutwhu/n3-util/jkv"
	n3json "github.com/cdutwhu/n3-util/n3json"
)

var (
	fPf  = fmt.Printf
	fPln = fmt.Println

	failOnErrWhen = cmn.FailOnErrWhen
	trackTime     = cmn.TrackTime
	mustWriteFile = cmn.MustWriteFile

	fmtJSON      = n3json.Fmt
	fmtJSONFile  = n3json.FmtFile
	maybeJSONArr = n3json.MaybeArr
	splitJSONArr = n3json.SplitArr
	makeJSONArr  = n3json.MakeArr
	newJKV       = jkv.NewJKV
)
