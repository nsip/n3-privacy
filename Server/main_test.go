package main

import (
	"testing"

	"github.com/nsip/n3-privacy/jkv"
	pp "github.com/nsip/n3-privacy/preprocess"	
)

func TestMain(t *testing.T) {
	main()
}

func TestParsePolicy(t *testing.T) {
	mask := sReplaceAll(pp.FmtJSONFile("../../Server/config/mask.json", "../preprocess/utils"), "\r\n", "\n")
	jkvM := jkv.NewJKV(mask, "ToBeNamed")
	if jkvM.Wrapped {
		fPln("wrapped")
	} else {
		fPln("Not Wrapped")
	}

	object := jkvM.LsLvlIPaths[1][0]
	object = sSpl(object, "@")[0]
	fPln(object)

	for _, ipath := range jkvM.LsLvlIPaths[2] {
		ifield := sSpl(ipath, jkv.PathLinker)[1]
		field := sSpl(ifield, "@")[0]
		fPln(field)
	}

	// fPln(jkvM.LsLvlIPaths[2])
	// fPln(mask)
}
