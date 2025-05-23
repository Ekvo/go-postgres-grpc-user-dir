// get and parse initial data from files or ENV variables
package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"github.com/Ekvo/go-postgres-grpc-user-dir/pkg/utils"
)

var ErrConfigEmpty = errors.New("empty")

// Config - contains url for database, server port with server network, secret key for jwt
type Config struct {
	DB         DataBaseConfig  `envPrefix:"DB_"`
	Migrations MigrationConfig `envPrefix:"MIGRATION_"`
	Server     ServerConfig    `envPrefix:"SRV_"`

	JWTSecretKey string `env:"JWT_SECRET"`

	msgErr utils.Message `env:"-"`
}

// NewConfig - load data from ENV (file or ENV variables)
func NewConfig(pathToEnv string) (*Config, error) {
	log.Print("config: config start")

	cfg := &Config{msgErr: utils.Message{}}
	if err := cfg.parse(pathToEnv); err != nil {
		return nil, fmt.Errorf("config: env.Parse error - {%w};", err)
	}

	if !cfg.validConfig() {
		return nil, fmt.Errorf("config: invalid config - {%s}", cfg.msgErr.String())
	}

	log.Print("config: config created")

	return cfg, nil
}

func (cfg *Config) parse(pathToEnv string) error {
	if err := godotenv.Load(pathToEnv); err != nil {
		// work with ENV
		log.Printf("config: .env file error - {%v};", err)
	}
	if err := env.Parse(cfg); err != nil {
		return err
	}
	cfg.DB.URL = cfg.DB.url()
	cfg.Migrations.DBURL = cfg.DB.URL + "?sslmode=disable"

	log.Print("config: parse end")

	return nil
}

func (cfg *Config) validConfig() bool {
	cfg.DB.validConfig(cfg.msgErr)
	cfg.Migrations.validConfig(cfg.msgErr)
	cfg.Server.validConfig(cfg.msgErr)

	if cfg.JWTSecretKey == "" {
		cfg.msgErr["jwt-secret-key"] = ErrConfigEmpty
	}
	return len(cfg.msgErr) == 0
}

type DataBaseConfig struct {
	Host     string `env:"HOST"`
	Port     uint16 `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`

	URL string `env:"-"`

	MaxConn           uint16        `env:"MAX_CONN"`
	MinConn           uint16        `env:"MIN_CONN"`
	ConnMaxLifeTime   time.Duration `env:"CONN_MAX_LIFE_TIME"`
	ConnMaxIdleTime   time.Duration `env:"CONN_MAX_IDLE_TIME"`
	ConnTime          time.Duration `env:"CONN_TIMEOUT"`
	HealthCheckPeriod time.Duration `env:"HEALTH_CHECK_PERIOD"`
}

func (cfgDB *DataBaseConfig) url() string {
	return fmt.Sprintf(`postgresql://%s:%s@%s:%d/%s`,
		cfgDB.User,
		cfgDB.Password,
		cfgDB.Host,
		cfgDB.Port,
		cfgDB.Name,
	)
}

func (cfgDB *DataBaseConfig) validConfig(msgErr utils.Message) {
	if cfgDB.Host == "" {
		msgErr["db-host"] = ErrConfigEmpty
	}
	if cfgDB.Port == 0 {
		msgErr["db-port"] = ErrConfigEmpty
	}
	if cfgDB.User == "" {
		msgErr["db-user"] = ErrConfigEmpty
	}
	if cfgDB.Password == "" {
		msgErr["db-password"] = ErrConfigEmpty
	}
	if cfgDB.Name == "" {
		msgErr["db-name"] = ErrConfigEmpty
	}
	if cfgDB.MaxConn == 0 {
		msgErr["db-max-conn"] = ErrConfigEmpty
	}
	if cfgDB.MinConn == 0 {
		msgErr["db-min-conn"] = ErrConfigEmpty
	}
	if cfgDB.ConnMaxLifeTime == 0 {
		msgErr["db-conn-max-life"] = ErrConfigEmpty
	}
	if cfgDB.ConnMaxIdleTime == 0 {
		msgErr["db-conn-idle-tile"] = ErrConfigEmpty
	}
	if cfgDB.ConnTime == 0 {
		msgErr["db-conn-time"] = ErrConfigEmpty
	}
	if cfgDB.HealthCheckPeriod == 0 {
		msgErr["db-health-check-period"] = ErrConfigEmpty
	}
}

type MigrationConfig struct {
	PathToMigrations string `env:"PATH"`
	DBURL            string `env:"-"`
}

func (cfgMig *MigrationConfig) validConfig(msgErr utils.Message) {
	if cfgMig.PathToMigrations == "" {
		msgErr["migration-path"] = ErrConfigEmpty
	}
}

type ServerConfig struct {
	Port    uint16 `env:"PORT"`
	Network string `env:"NETWORK"`
}

func (cfgSrv *ServerConfig) validConfig(msgErr utils.Message) {
	if cfgSrv.Port == 0 {
		msgErr["srv-port"] = ErrConfigEmpty
	}
	if cfgSrv.Network == "" {
		msgErr["srv-network"] = ErrConfigEmpty
	}
}
