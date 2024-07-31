package handlers

import (
	"Messaggio/producer/internal/model"
	"github.com/IBM/sarama"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "Messaggio/producer/docs"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name IServices --output mocks/
type IServices interface {
	SaveMessage(message *model.Message) (string, error)
	GetMessage(id string) (*model.Message, error)
	GetMessages() ([]*model.Message, error)
	DeleteMessages() error
	DeleteMessage(id string) error
	GetStatsByDays() (map[string]*model.ByDays, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.43.2 --name IProducer --output mocks/
type IProducer interface {
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

type Handlers struct {
	services IServices
	producer IProducer
}

func New(e *echo.Echo, services IServices, producer IProducer) *Handlers {
	h := &Handlers{
		services: services,
		producer: producer,
	}

	e.Use(middleware.RequestID())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogRemoteIP:  true,
		LogError:     true,
		LogLatency:   true,
		LogMethod:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("latentcy", v.Latency.String()).
				Str("requestID", v.RequestID).
				Err(v.Error).
				Str("remoteIP", v.RemoteIP).
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("method", v.Method).
				Msg("request")

			return nil
		},
	}))
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api")
	{
		statistic := e.Group("/statistic")
		{
			statistic.GET("/days", h.getStatsByDays)
		}

		messages := api.Group("/messages")
		{
			messages.POST("", h.createMessage)
			messages.GET("/:id", h.getMessage)
			messages.GET("", h.getMessages)
			messages.DELETE("", h.deleteMessages)
			messages.DELETE("/:id", h.deleteMessage)
		}
	}

	return h
}
