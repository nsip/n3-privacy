package storage

import (
	"testing"

	pp "github.com/nsip/n3-privacy/preprocess"
)

func TestRecPolicyMeta(t *testing.T) {
	mask := pp.FmtJSONFile("../../Server/config/mask.json", "../../preprocess/utils")
	mask = sReplaceAll(mask, "\r\n", "\n")
	recPolicyMeta(mask, "../config/meta.json")
}
