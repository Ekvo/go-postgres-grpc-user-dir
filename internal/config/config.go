package config

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/Ekvo/go-postgres-grpc-user-dir/pkg/utils"
)

type Config struct {
	DBURL string `mapstructure:"DB_URL"`

	SRVPort    string `mapstructure:"SRV_PORT_USER"`
	SRVNetwork string `mapstructure:"SRV_NETWORK"`

	JWTSecretKey string `mapstructure:"JWT_SECRET"`
}

func NewConfig(pathToEnv string) (*Config, error) {
	if err := godotenv.Load(pathToEnv); err != nil {
		// work with ENV
		log.Printf("config: .env file error - %v", err)
	}
	viper.AutomaticEnv()
	for _, env := range getNameEnv() {
		if err := viper.BindEnv(env); err != nil {
			return nil, fmt.Errorf("config: ENV error - %w", err)
		}
	}
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config: Unmarshal config error - %w", err)
	}
	if err := cfg.validConfig(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func getNameEnv() []string {
	return []string{
		"DB_URL",
		"SRV_PORT_USER",
		"SRV_NETWORK",
		"JWT_SECRET",
	}
}

// reURLDB - pattern for check 'DB_URL'
// ^postgresql:\/\/[a-zA-Z0-9]+:[a-zA-Z0-9]+@[a-zA-Z0-9\.:]+\/[a-zA-Z0-9]+$
var reURLDB = regexp.MustCompile(`^postgresql:\/\/[a-zA-Z0-9_-]+:[a-zA-Z0-9]+@(?:[a-zA-Z0-9]+|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):[0-9]{1,5}\/[a-zA-Z0-9_-]+$`)

func (cfg *Config) validConfig() error {
	msgErr := utils.Message{}
	if cfg.DBURL == "" {
		msgErr["url-DB"] = "empty"
	} else if !reURLDB.MatchString(cfg.DBURL) {
		msgErr["url-DB"] = "invalid"
	}
	if cfg.SRVPort == "" {
		msgErr["server-port"] = "empty"
	} else if _, err := strconv.ParseUint(cfg.SRVPort, 10, 16); err != nil {
		msgErr["server-port"] = "no numeric"
	}
	if cfg.SRVNetwork == "" {
		msgErr["server-network"] = "empty"
	}
	if cfg.JWTSecretKey == "" {
		msgErr["jwt-secret-key"] = "empty"
	}
	if len(msgErr) > 0 {
		return fmt.Errorf("config: invalid config - %s", msgErr.String())
	}
	return nil
}
