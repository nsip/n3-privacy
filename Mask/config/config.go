package config

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/burntsushi/toml"
)

// config is toml, type name MUST be identical to file name "config.toml"
type config struct {
	Path    string
	LogFile string
}

var (
	// toml file name must be identical to config struct definition name
	lsCfg = []interface{}{&config{}}
)

// NewCfg :
func NewCfg(cfgPaths ...string) interface{} {
	for _, f := range cfgPaths {
		if _, e := os.Stat(f); e == nil {
			if abs, e := filepath.Abs(f); e == nil {
				return set(f, abs)
			}
		}
	}
	return nil
}

func set(f, abs string) interface{} {
	for _, cfg := range lsCfg {
		name := reflect.TypeOf(cfg).Elem().Name()
		if sHasSuffix(f, "/"+name+".toml") {
			if _, e := toml.DecodeFile(f, cfg); e == nil {
				reflect.ValueOf(cfg).Elem().FieldByName("Path").SetString(abs)
				return save(f, cfg)
			}
		}
	}
	return nil
}

func save(path string, cfg interface{}) interface{} {
	if f, e := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, os.ModePerm); e == nil {
		defer f.Close()
		toml.NewEncoder(f).Encode(cfg)
	}
	return cfg
}
