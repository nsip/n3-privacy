package main

import (
	"testing"

	toml "github.com/cdutwhu/json-util/n3toml"
	cfg "github.com/nsip/n3-privacy/Server/config"
)

func TestMain(t *testing.T) {
	main()
}

// Create a copy of config for Client. Excluding [WebService], [Storage], [File] attributes.
// Once building, move it to Client Directory. i.e. "config-client.toml", "config_client" is used in build.sh
// DO NOT delete or modify this Test!
func TestCreateClientCfg(t *testing.T) {
	temp := "./config-client.toml"
	cfg.NewCfg("./config.toml").SaveAs(temp)
	toml.RmFileAttrL1(temp, "config_client", "WebService", "Storage", "File")
}
