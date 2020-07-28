package process

import (
	"fmt"

	"github.com/cdutwhu/debog/fn"
	"github.com/cdutwhu/gotil/io"
	"github.com/cdutwhu/gotil/misc"
	"github.com/cdutwhu/n3-util/jkv"
	"github.com/cdutwhu/n3-util/n3json"
)

var (
	fPf  = fmt.Printf
	fPln = fmt.Println

	failOnErrWhen = fn.FailOnErrWhen
	trackTime     = misc.TrackTime
	mustWriteFile = io.MustWriteFile
	fmtJSON       = n3json.Fmt
	fmtJSONFile   = n3json.FmtFile
	maybeJSONArr  = n3json.MaybeArr
	splitJSONArr  = n3json.SplitArr
	makeJSONArr   = n3json.MakeArr
	newJKV        = jkv.NewJKV
)
