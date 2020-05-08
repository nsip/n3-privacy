package process

import "testing"

func TestFileMask(t *testing.T) {
	FileMask("../samples/xapi.json", "../samples/xapiMask.json", "./out.json")
}
