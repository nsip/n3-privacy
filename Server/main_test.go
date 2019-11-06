package main

import (
	"testing"

	pp "github.com/nsip/n3-privacy/preprocess"
)

func TestMain(t *testing.T) {
	main()
}

func TestParsePolicy(t *testing.T) {
	mask := sReplaceAll(pp.FmtJSONFile("../../Server/config/mask.json", "../preprocess/utils"), "\r\n", "\n")
	fPln(mask)
	recPolicy(mask)
}
