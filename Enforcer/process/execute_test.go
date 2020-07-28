package process

import (
	"testing"
)

func TestFileExe(t *testing.T) {
	save := ""
	FileExe("../samples/xapi1.json", "../samples/xapiPolicyP.json", save)
}
