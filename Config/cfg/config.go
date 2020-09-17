package cfg

import "github.com/cdutwhu/n3-util/n3cfg"

// Config : AUTO Created From /home/qmiao/Desktop/temp/n3-privacy/Config/config.toml
type Config struct {
	Log string
	Service interface{}
	Version interface{}
	Storage struct {
		DB string
		DBPath string
	}
	Loggly struct {
		Token string
	}
	WebService struct {
		Port int
	}
	Route struct {
		LsObject string
		Help string
		LsContext string
		Update string
		Delete string
		GetHash string
		GetID string
		LsUser string
		FetchEnfMac string
		Enforce string
		Get string
		LsID string
		FetchEnfLinux string
		FetchEnfWin string
	}
	File struct {
		EnforcerLinux64 string
		EnforcerMac string
		EnforcerWin64 string
	}
	Server struct {
		IP interface{}
		Port interface{}
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
