package topspin

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
		appName        string
		v              *viper.Viper
		file           string
		k8s            bool
		onConfigChange func(e fsnotify.Event)
		status         map[time.Time]string
	}
)

func NewConfig(appName string) *Config {
	return &Config{
		appName: appName,
		v:       viper.New(),
		k8s:     false,
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

	err = cfg.v.ReadInConfig()
	if _, ok := err.(*os.PathError); ok {
		cfg.addStatus(fmt.Sprintf("mo config file at '%s', using default values", cfg.file))

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

	cfg.v.SetConfigFile(path)

	err = cfg.v.ReadInConfig()
	if _, ok := err.(*os.PathError); ok {
		cfg.addStatus(fmt.Sprintf("mo config file at '%s', using default values", cfg.file))

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

func (cfg *Config) List() string {
	var sb strings.Builder
	for k, v := range cfg.v.AllSettings() {
		entry := fmt.Sprintf("%s: %v\n", k, v)
		sb.WriteString(entry)
	}
	return sb.String()
}

func (cfg *Config) addStatus(message string) {
	cfg.status[time.Now()] = message
}

func configName(appName string) string {
	return appName + "-config"
}

func (cfg *Config) defaultOnConfigChangeFunc() func(e fsnotify.Event) {
	return func(e fsnotify.Event) {
		cfg.addStatus(fmt.Sprintf("Config file updated: %s", e.Name))
	}
}
