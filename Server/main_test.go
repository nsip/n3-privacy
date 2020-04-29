package main

import (
	"testing"

	toml "github.com/cdutwhu/json-util/n3toml"
)

func TestMain(t *testing.T) {
	main()
}

// Create a copy of config for Client. Excluding [Storage], [File] attributes.
// Once building, move it to Client Directory.
// DO NOT delete or modify this Test!
func TestCreateClientCfg(t *testing.T) {
	toml.RmFileAttrL1("./config.toml", "config_client", "Storage", "File")
}
