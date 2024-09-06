package main

import (
	"Messaggio/consumer/configs"
	"Messaggio/consumer/internal/repository"
	"Messaggio/consumer/internal/services"
	"Messaggio/db"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	conf := configs.Load()
	timeout := time.NewTicker(60 * time.Second)

	var consumer sarama.Consumer
	var err error
	fl := false

	for {
		select {
		case <-timeout.C:
			log.Fatal().Msgf("sarama producer: %s", err)

		default:
			consumer, err = sarama.NewConsumer([]string{conf.Kafka.Host + ":" + conf.Kafka.Port}, nil)
			if err == nil {
				fl = true
				break
			}
		}

		if fl {
			break
		}
	}
	defer consumer.Close()

	partConsumer, err := consumer.ConsumePartition("messages", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal().Err(err).Msg("Error consuming partition")
	}
	defer partConsumer.Close()

	conn := db.NewPostgres(conf.DB.Postgres)

	rep := repository.New(conn)

	serv := services.New(rep)

	for {
		select {
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				log.Error().Msg("Consumer closed")
				return
			}

			var id string

			id = string(msg.Key)
			if err != nil {
				log.Error().Err(err).Msg("Error unmarshalling message")
				continue
			}

			err = serv.TestProcessMessage(id)
			if err != nil {
				log.Error().Err(err).Msg("Error processing message")
				continue
			}

			log.Info().Msgf("message with ID: %s processed", id)
		}
	}
}
