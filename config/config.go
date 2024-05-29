package config

import (
	"time"

	"github.com/zer0day88/tinder/pkg/environment"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type config struct {
	ServerPort  string            `mapstructure:"server_port"`
	Environment environment.Level `mapstructure:"environment"`
	LogLevel    zerolog.Level     `mapstructure:"log_level"`
	Database    struct {
		Postgres struct {
			Host              string `mapstructure:"host"`
			User              string `mapstructure:"user"`
			Password          string `mapstructure:"password"`
			DbName            string `mapstructure:"dbname"`
			Port              int    `mapstructure:"port"`
			SSLMode           string `mapstructure:"sslmode"`
			ConMaxLifetime    int    `mapstructure:"conMaxLifetime"`
			ConMaxOpen        int32  `mapstructure:"conMaxOpen"`
			HealthCheckPeriod int    `mapstructure:"healthCheckPeriod"`
			ConMaxIdleTime    int    `mapstructure:"conMaxIdleTime"`
		} `mapstructure:"postgres"`
		Redis struct {
			Host     string `mapstructure:"host"`
			Password string `mapstructure:"password"`
			Port     int    `mapstructure:"port"`
			UseTLS   bool   `mapstructure:"use_tls"`
		} `mapstructure:"redis"`
	} `mapstructure:"database"`

	AccessTokenPrivateKey string        `mapstructure:"access_token_private_key"`
	AccessTokenPublicKey  string        `mapstructure:"access_token_public_key"`
	AccessTokenExpiredIn  time.Duration `mapstructure:"access_token_expired_in"`

	RefreshTokenPrivateKey string        `mapstructure:"refresh_token_private_key"`
	RefreshTokenPublicKey  string        `mapstructure:"refresh_token_public_key"`
	RefreshTokenExpiredIn  time.Duration `mapstructure:"refresh_token_expired_in"`

	JwtRedisKey string `mapstructure:"jwt_redis_key"`
}

var (
	Key config
)

func Load() {

	readConfig()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info().Msgf("Config file changed: %s", e.Name)
		readConfig()
	})

}
func LoadWithPath(path string) {

	readConfigPath(path)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info().Msgf("Config file changed: %s", e.Name)
		readConfigPath(path)
	})

}

func readConfig() {
	viper.SetConfigName("app-config") // name of config file (without extension)
	viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("config")

	viper.SetDefault("environment", "development")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	err := viper.Unmarshal(&Key)
	if err != nil {
		panic(err)
	}

}

func readConfigPath(path string) {
	viper.SetConfigName("app-config") // name of config file (without extension)
	viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)

	viper.SetDefault("environment", "development")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	err := viper.Unmarshal(&Key)
	if err != nil {
		panic(err)
	}

}
