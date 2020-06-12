package client

import (
	"os"
	"path/filepath"
	"time"

	"github.com/burntsushi/toml"
)

// Config is toml
type Config struct {
	Path        string
	LogFile     string
	ServiceName string
	Route       struct {
		HELP      string
		GetID     string
		GetHash   string
		Get       string
		Update    string
		Delete    string
		LsID      string
		LsUser    string
		LsContext string
		LsObject  string
		Enforce   string
	}
	Server struct {
		Protocol string
		IP       string
		Port     int
	}
	Access struct {
		Timeout int64
	}
}

// newCfg :
func newCfg(configs ...string) *Config {
	for _, f := range configs {
		if _, e := os.Stat(f); e == nil {
			return (&Config{Path: f}).set()
		}
	}
	return nil
}

// set is
func (cfg *Config) set() *Config {
	f := cfg.Path /* make a copy of original for restoring */
	if _, e := toml.DecodeFile(f, cfg); e == nil {
		// modify some to save
		cfg.Path = f
		if abs, e := filepath.Abs(f); e == nil {
			cfg.Path = abs
		}

		// save
		cfg.save()

		ICfg, e := cfgRepl(cfg, map[string]interface{}{
			"[DATE]": time.Now().Format("2006-01-02"),
		})
		failOnErr("%v", e)
		return ICfg.(*Config)
	}
	return nil
}

func (cfg *Config) save() {
	if f, e := os.OpenFile(cfg.Path, os.O_WRONLY|os.O_TRUNC, os.ModePerm); e == nil {
		defer f.Close()
		toml.NewEncoder(f).Encode(cfg)
	}
}

// initEnvVarFromTOML : initialize the global variables
func initEnvVarFromTOML(key string, configs ...string) bool {
	configs = append(configs, "./config.toml")
	Cfg := newCfg(configs...)
	if Cfg == nil {
		return false
	}
	struct2Env(key, Cfg)
	return true
}
