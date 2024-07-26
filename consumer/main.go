package main

import (
	"Messaggio/consumer/configs"
	"Messaggio/consumer/internal/model"
	"Messaggio/consumer/internal/repository"
	"Messaggio/consumer/internal/services"
	"Messaggio/db"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	conf := configs.Load()
	//TODO: Create a new consumer
	consumer, err := sarama.NewConsumer([]string{conf.Kafka.Host + ":" + conf.Kafka.Port}, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating consumer")
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
				//TODO: Maybe something else should be done here?
				log.Error().Msg("Consumer closed")
				return
			}

			var message model.Message

			err = json.Unmarshal(msg.Value, &message)
			if err != nil {
				log.Error().Err(err).Msg("Error unmarshalling message")
				continue
			}

			err = serv.TestProcessMessage(&message)
			if err != nil {
				log.Error().Err(err).Msg("Error processing message")
				continue
			}

			fmt.Println("Received message:", message.Content)
		}
	}
}
