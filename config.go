package topspin

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	configFileKey = "config.file"
)

var (
	defaultFile    string = filepath.FromSlash("configs/config.yml")
	defaultK8SFile string = "config.yml"
)

type (
	Config struct {
		*SimpleWorker
		appName        string
		v              *viper.Viper
		file           string
		k8s            bool
		onConfigChange func(e fsnotify.Event)
	}
)

func NewConfig(appName string, log Logger) *Config {
	return &Config{
		SimpleWorker: NewWorker(configName(appName), log),
		appName:      appName,
		v:            viper.New(),
		k8s:          false,
	}
}

func (cfg *Config) Load() (updated *Config, err error) {
	flag.StringVar(&cfg.file, "configfile", defaultFile, "path to configuration file")
	flag.BoolVar(&cfg.k8s, "configk8s", false, "is app running in kubernetes")
	flag.Parse()

	if cfg.k8s {
		return cfg.loadK8sConfig()
	}

	return cfg.loadConfig()
}

func (cfg *Config) loadConfig() (updated *Config, err error) {
	cfg.v.SetDefault(configFileKey, cfg.file)
	cfg.v.AutomaticEnv()
	cfg.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.v.SetTypeByDefaultValue(true)
	cfg.v.SetConfigFile(cfg.file)

	cfg.Log().Debugf("Reading configuration from file: %s", cfg.file)
	err = cfg.v.ReadInConfig()
	if _, ok := err.(*os.PathError); ok {
		cfg.Log().Infof("No config file at '%s', using default values")

	} else if err != nil {
		return cfg, fmt.Errorf("error reading config: %w", err)
	}

	cfg.v.WatchConfig()

	cfg.v.OnConfigChange(cfg.defaultOnConfigChangeFunc())

	return cfg, nil
}

func (cfg *Config) loadK8sConfig() (updated *Config, err error) {
	cfg.v.SetDefault(configFileKey, defaultK8SFile)
	cfg.v.AutomaticEnv()
	cfg.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.v.SetTypeByDefaultValue(true)

	path, err := cfg.pathToK8SConfig()
	if err != nil {
		err = fmt.Errorf("cannot get config file path: %e", err)
		return cfg, err
	}

	cfg.Log().Debugf("Reading configuration from k8s config map", cfg.file)
	cfg.v.SetConfigFile(path)

	err = cfg.v.ReadInConfig()
	if _, ok := err.(*os.PathError); ok {
		cfg.Log().Infof("No config file at '%s', using default values")

	} else if err != nil {
		return cfg, fmt.Errorf("error reading config: %w", err)
	}

	cfg.v.WatchConfig()

	cfg.v.OnConfigChange(cfg.defaultOnConfigChangeFunc())

	return cfg, nil
}

func (cfg *Config) setDefaultLookupPaths() {
	cfg.v.SetConfigType("yaml")
	cfg.v.AddConfigPath("/etc/" + cfg.appName)
	cfg.v.AddConfigPath("$HOME/." + cfg.appName)
	cfg.v.AddConfigPath(cfg.file)
	cfg.v.AddConfigPath(".")
}

func (cfg *Config) pathToK8SConfig() (filePath string, err error) {
	// WIP: There is no associated logic at the moment,
	// returning the default path
	return defaultK8SFile, nil
}

func (cfg *Config) SetOnConfigChange(onConfigChangeFunc func(e fsnotify.Event)) {
	cfg.onConfigChange = onConfigChangeFunc
}

func (cfg *Config) List() {
	for k, v := range cfg.v.AllSettings() {
		cfg.Log().Debugf("%s: %v\n", k, v)
	}
}

func configName(appName string) string {
	return appName + "-config"
}

func (cfg *Config) defaultOnConfigChangeFunc() func(e fsnotify.Event) {
	return func(e fsnotify.Event) {
		cfg.Log().Infof("Config file updated: %s", e.Name)
	}
}
