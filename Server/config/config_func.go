package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	mux sync.Mutex
)

// newCfg :
func newCfg(configs ...string) *Config {
	defer func() { mux.Unlock() }()
	mux.Lock()
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
	_, e := toml.DecodeFile(f, cfg)
	failOnErr("%v", e)

	// modify path
	cfg.Path = f
	if abs, e := filepath.Abs(f); e == nil {
		cfg.Path = abs
	}
	if ver, e := gitver(); e == nil && ver != "" { /* successfully got git ver */
		cfg.Version = ver
	}

	// save
	cfg.save()

	home, e := os.UserHomeDir()
	failOnErr("%v", e)
	return cfgRepl(cfg, map[string]interface{}{
		"~":      home,
		"[DATE]": time.Now().Format("2006-01-02"),
		"[IP]":   localIP(),
		"[port]": cfg.WebService.Port,
		"[db]":   cfg.Storage.DB,
		"[s]":    cfg.Service,
		"[v]":    cfg.Version,
	}).(*Config)
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
	newCfg(filename).save()
}

// InitEnvVarFromTOML : initialize global variables
func InitEnvVarFromTOML(key string, configs ...string) bool {
	Cfg := newCfg(append(configs, "./config.toml", "./config/config.toml")...)
	if Cfg == nil {
		return false
	}
	struct2Env(key, Cfg)
	return true
}
