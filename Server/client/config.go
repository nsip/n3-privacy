package client

import "github.com/cdutwhu/n3-util/n3cfg"

// Config : AUTO Created From /home/qmiao/Desktop/temp/n3-privacy/Server/client/config.toml
type Config struct {
	Service string
	Route struct {
		LsContext string
		Get string
		Update string
		FetchEnfWin string
		Enforce string
		FetchEnfLinux string
		LsObject string
		LsUser string
		Delete string
		GetID string
		GetHash string
		Help string
		FetchEnfMac string
		LsID string
	}
	Server struct {
		IP string
		Port int
		Protocol string
	}
	Access struct {
		Timeout int
	}
}

// NewCfg :
func NewCfg(cfgStruName string, mReplExpr map[string]string, cfgPaths ...string) interface{} {
	var cfg interface{}
	switch cfgStruName {
	case "Config":
		cfg = &Config{}
	default:
		return nil
	}
	return n3cfg.InitEnvVar(cfg, mReplExpr, cfgStruName, cfgPaths...)
}
