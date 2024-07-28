package handlers

import (
	"Messaggio/producer/internal/model"
	"github.com/IBM/sarama"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type iServices interface {
	SaveMessage(message *model.Message) (string, error)
	GetMessage(id string) (*model.Message, error)
	GetMessages() ([]*model.Message, error)
	DeleteMessages() error
}

type Handlers struct {
	services iServices
	producer sarama.SyncProducer
}

func New(e *echo.Echo, services iServices, producer sarama.SyncProducer) *Handlers {
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

	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "ok")
	})

	api := e.Group("/api")
	{
		messages := api.Group("/messages")
		{
			messages.POST("", h.createMessage)
			messages.GET("/:id", h.getMessage)
			messages.GET("", h.getMessages)
			messages.DELETE("", h.deleteMessages)
		}
	}

	return h
}
