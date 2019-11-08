package config

import (
	"os"
	"reflect"

	"github.com/burntsushi/toml"
)

// Config is toml
type Config struct {
	Path       string
	ErrLog     string
	JMPath     string
	JQPath     string
	WebService struct {
		Port    int
		Version string
	}
	Route struct {
		Peek   string
		Get    string
		Update string
		GetJM  string
		GetJQ  string
	}
}

// NewCfg :
func NewCfg(configs ...string) *Config {
	for _, f := range configs {
		if _, e := os.Stat(f); e == nil {
			cfg := &Config{Path: f}
			return cfg.set()
		}
	}
	return nil
}

// set is
func (cfg *Config) set() *Config {
	path := cfg.Path /* make a copy of original path for restoring */
	toml.DecodeFile(cfg.Path, cfg)
	cfg.Path = path
	return cfg.modCfg()
}

func (cfg *Config) modCfg() *Config {
	// *** replace version *** //
	ver := fSf("%s", cfg.WebService.Version)
	v := reflect.ValueOf(cfg.Route)
	for i := 0; i < v.NumField(); i++ {
		vv := sReplaceAll(v.Field(i).Interface().(string), "#", ver)
		reflect.ValueOf(&cfg.Route).Elem().Field(i).SetString(vv)
	}
	return cfg
}
