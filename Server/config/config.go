package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/burntsushi/toml"
)

// Config is toml
type Config struct {
	Path    string
	LogFile string

	Storage struct {
		DataBase     string
		MetaPath     string
		BadgerDBPath string
	}

	WebService struct {
		Port    int
		Service string
		Version string
	}

	Route struct {
		HELP        string
		GetID       string
		GetHash     string
		Get         string
		Update      string
		Delete      string
		LsID        string
		LsUser      string
		LsContext   string
		LsObject    string
		GetEnforced string
	}

	File struct {
		ClientLinux64   string
		ClientMac       string
		ClientWin64     string
		ClientConfig    string
		EnforcerLinux64 string
		EnforcerMac     string
		EnforcerWin64   string
		EnforcerConfig  string
	}

	Server struct {
		Protocol string
		IP       string
		Port     interface{}
	}

	Access struct {
		Timeout int
	}
}

// NewCfg :
func NewCfg(configs ...string) *Config {
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
		// modify path
		cfg.Path = f
		if abs, e := filepath.Abs(f); e == nil {
			cfg.Path = abs
		}

		// DO NOT SAVE

		ICfg, e := cfgRepl(cfg, map[string]interface{}{
			"[DATE]": time.Now().Format("2006-01-02"),
			"[IP]":   localIP(),
			"[PORT]": cfg.WebService.Port,
			"[s]":    cfg.WebService.Service,
			"[v]":    cfg.WebService.Version,
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

// SaveAs :
func (cfg *Config) SaveAs(filename string) {
	bytes, err := ioutil.ReadFile(cfg.Path)
	failOnErr("%v", err)
	if !sHasSuffix(filename, ".toml") {
		filename += ".toml"
	}
	failOnErr("%v", ioutil.WriteFile(filename, bytes, 0666))
	NewCfg(filename).save()
}
