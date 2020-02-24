package main

import "testing"

func TestMain(t *testing.T) {
	// main()
	doMask("./samples/inquiry_skills.json", "./samples/inquiry_skillsMask.json", "./out.json")
}
