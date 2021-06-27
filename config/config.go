package config

import (
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/ory/viper"
	"github.com/pkg/errors"
)

type Config struct {
	Port string `mapstructure:"PORT"`
	DB   struct {
		DRV        string
		Host       string
		Port       int
		Name       string
		User       string
		Password   string
		SSLEnabled bool
	} `mapstructure:"db"`
	CORS struct {
		AllowOrigins []string `mapstructure:"allow_origins"`
	} `mapstructure:"cors"`
	FrontEnd struct {
		BaseURL            string   `mapstructure:"base_url"`
		URLPrefixesToMatch []string `mapstructure:"url_prefixes_to_match"`
	} `mapstructure:"front_end"`
}

// Load loads config file with name in path to Config
func Load() (*Config, error) {
	vpr := viper.New()
	vpr.SetDefault("mode", "development")
	vpr.AutomaticEnv()
	vpr.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vpr.SetConfigType("yaml")

	mode := vpr.GetString("mode")
	vpr.SetConfigName(mode)
	vpr.AddConfigPath("./config")

	if err := vpr.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, errors.Wrap(err, "loading configuration file error")
		}
	}

	cfg := &Config{}
	if err := vpr.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "cannot unmarshal configuration")
	}

	return cfg, nil
}

// ConnectionStr ...
func (c *Config) ConnectionString() string {
	sslMode := "disable"
	if c.DB.SSLEnabled {
		sslMode = "require"
	}

	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		c.DB.User,
		c.DB.Password,
		c.DB.Name,
		c.DB.Host,
		c.DB.Port,
		sslMode,
	)
}
