package main

import "testing"

func TestMain(t *testing.T) {
	doMask("./samples/xapi.json", "./samples/xapiMask.json", "./out.json")
}
