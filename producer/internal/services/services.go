package services

import (
	"Messaggio/producer/internal/model"
)

type iRepository interface {
	SaveMessage(message *model.Message) (string, error)
	GetMessage(id string) (*model.Message, error)
	GetMessages() ([]*model.Message, error)
	DeleteMessages() error
	GetStatsByHour() ([]*model.ByHours, error)
}

type Services struct {
	repository iRepository
}

func New(repository iRepository) *Services {
	return &Services{
		repository: repository,
	}
}
