package client

import "github.com/cdutwhu/n3-util/n3cfg"

// Config : AUTO Created From /home/qmiao/Desktop/temp/n3-privacy/Server/client/config.toml
type Config struct {
	Service string
	Route struct {
		GetHash string
		LsContext string
		LsUser string
		FetchEnfMac string
		Get string
		LsID string
		FetchEnfLinux string
		Delete string
		Help string
		GetID string
		LsObject string
		Enforce string
		FetchEnfWin string
		Update string
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
