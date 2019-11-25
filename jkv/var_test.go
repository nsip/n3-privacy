package jkv

import "testing"

func TestProjectV(t *testing.T) {
	ss := []string{
		"12@a@1~b",
		"b~c",
		"c~a",
		"a~b~c",
		"b~c~d~e",
		"c~a",
		"d~e~b~a~f",
		"a",
		"***@b@***",
	}
	fPln(ProjectV(ss, "~", "@", "@"))
}
