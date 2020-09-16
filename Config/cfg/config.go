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
		LsID string
		LsObject string
		Update string
		FetchEnfWin string
		Enforce string
		GetHash string
		Help string
		LsContext string
		LsUser string
		FetchEnfLinux string
		FetchEnfMac string
		Delete string
		Get string
		GetID string
	}
	File struct {
		EnforcerMac string
		EnforcerWin64 string
		EnforcerLinux64 string
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
