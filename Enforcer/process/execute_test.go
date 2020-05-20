package process

import "testing"

func TestFileExe(t *testing.T) {
	FileExe("../samples/xapi.json", "../samples/xapiPolicy.json", "./out.json")
}
