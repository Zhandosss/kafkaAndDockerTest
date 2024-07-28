package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"time"
)

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Name     string `mapstructure:"name"`
	Password string
}

func NewPostgres(config *PostgresConfig) *sqlx.DB {
	dataSource := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", config.User, config.Password, config.Host, config.Port, config.Name)
	log.Info().Msgf("Connecting to db: %s", dataSource)

	timeout := time.NewTicker(30 * time.Second)

	var err error
	var conn *sqlx.DB
	fl := false

	for {
		select {
		case <-timeout.C:
			log.Fatal().Msgf("sqlx connect: %s", err)
		default:
			conn, err = sqlx.Connect("postgres", dataSource)
			if err == nil {
				fl = true
				break
			}
		}
		if fl {
			break
		}
	}

	err = conn.Ping()
	if err != nil {
		log.Fatal().Msgf("ping db: %s", err)
	}

	return conn
}
