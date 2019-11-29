package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/burntsushi/toml"
)

var (
	fSf         = fmt.Sprintf
	fPln        = fmt.Println
	sReplaceAll = strings.ReplaceAll
)

// Config is toml
type Config struct {
	Path       string
	ErrLog     string
	WebService struct {
		Port    int
		Version string
		Service string
	}
	Route struct {
		GetID        string
		GetHash      string
		Get          string
		Update       string
		Delete       string
		ListOfPID    string
		ListOfUser   string
		ListOfCtx    string
		ListOfObject string
	}
	// Client
	Server struct {
		Protocol string
		IP       string
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
	ver := fSf("%s", cfg.WebService.Version)
	svr := fSf("%s", cfg.WebService.Service)
	return cfg.modCfg(map[string]string{"#v": ver, "#s": svr}) // *** replace version & service-name *** //
}

func (cfg *Config) modCfg(mRepl map[string]string) *Config {
	if mRepl == nil || len(mRepl) == 0 {
		return cfg
	}
	nField := reflect.ValueOf(cfg.Route).NumField()
	for i := 0; i < nField; i++ {
		for key, value := range mRepl {
			replaced := sReplaceAll(reflect.ValueOf(cfg.Route).Field(i).Interface().(string), key, value)
			reflect.ValueOf(&cfg.Route).Elem().Field(i).SetString(replaced)
		}
	}
	return cfg
}
