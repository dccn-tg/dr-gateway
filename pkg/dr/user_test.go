package dr

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	log "github.com/Donders-Institute/tg-toolset-golang/pkg/logger"
	"github.com/spf13/viper"
)

func init() {
	cfg := log.Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      log.Debug,
		EnableFile:        true,
		FileJSONFormat:    true,
		FileLocation:      "log/dr-gateway_test.log",
		FileLevel:         log.Debug,
	}
	// initialize logger
	log.NewLogger(cfg, log.InstanceZapLogger)
}

// load configuration to DRConfig structure
func loadDRConfig(cpath string) (Config, error) {

	var conf Config

	// load configuration
	cfg, err := filepath.Abs(cpath)
	if err != nil {
		return conf, err
	}

	if _, err := os.Stat(cfg); err != nil {
		return conf, fmt.Errorf("cannot load config: %s", cfg)
	}

	viper.SetConfigFile(cfg)
	if err := viper.ReadInConfig(); err != nil {
		return conf, fmt.Errorf("cannot read config file, %s", err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		return conf, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return conf, nil
}

func TestGetAllUsers(t *testing.T) {

	pathCfg := os.Getenv("DR_GATEWAY_CONFIG")

	cfg, err := loadDRConfig(pathCfg)
	if err != nil {
		t.Fatal(err)
	}

	users, err := GetAllUsers(cfg)
	if err != nil {
		t.Fatal(err)
	}

	for u := range users {
		t.Logf("user: %+v\n", u)
	}
}
