package configs

import (
	"Messaggio/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Server *ServerConfig `mapstructure:"http_server"`
	DB     *DBConfig     `mapstructure:"db"`
	Kafka  *KafkaConfig  `mapstructure:"kafka"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DBConfig struct {
	Postgres *db.PostgresConfig `mapstructure:"postgres"`
}

type KafkaConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func Load() *Config {
	viper.SetDefault("port", "8080")
	viper.SetDefault("host", "localhost")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./producer")
	viper.AddConfigPath("./producer/configs")
	viper.AddConfigPath("./configs") //to docker container

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Msgf("failed to read config file: %s", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatal().Msgf("failed to unmarshal config: %s", err)
	}

	if err != nil {
		log.Fatal().Msgf("Error loading .env file")
	}

	config.DB.Postgres.Password = os.Getenv("DB_PASSWORD")
	return config
}
