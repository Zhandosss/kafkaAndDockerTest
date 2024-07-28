package main

import (
	"Messaggio/db"
	"Messaggio/producer/configs"
	"Messaggio/producer/internal/handlers"
	"Messaggio/producer/internal/repository"
	"Messaggio/producer/internal/services"
	"context"
	"embed"
	"errors"
	"github.com/IBM/sarama"
	"github.com/labstack/echo/v4"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	conf := configs.Load()

	conn := db.NewPostgres(conf.DB.Postgres)
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatal().Msgf("Error closing connection: %s", err)
		}
		log.Info().Msg("DB connection closed")
	}()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal().Msgf("Error setting dialect: %s", err)
	}

	if err := goose.Up(conn.DB, "migrations"); err != nil {
		log.Fatal().Msgf("Error applying migrations: %s", err)
	}

	repo := repository.New(conn)

	serv := services.New(repo)

	e := echo.New()

	timeout := time.NewTicker(30 * time.Second)

	var producer sarama.SyncProducer
	var err error
	fl := false
	//TODO: add kafka host and port to config
	for {
		select {
		case <-timeout.C:
			log.Fatal().Msgf("Error creating producer: %s", err)
		default:
			producer, err = sarama.NewSyncProducer([]string{conf.Kafka.Host + ":" + conf.Kafka.Port}, nil)
			if err == nil {
				fl = true
				break
			}
		}
		if fl {
			break
		}
	}
	defer producer.Close()

	handlers.New(e, serv, producer)

	server := &http.Server{
		Addr:         conf.Server.Host + ":" + conf.Server.Port,
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	log.Info().Msgf("Server started on %s", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Msgf("Server shut down: %s", err)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err)
	}

}
