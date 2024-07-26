package services

import "Messaggio/consumer/internal/model"

type iRepository interface {
	EmptyProcess(message *model.Message) error
}

type Services struct {
	repository iRepository
}

func New(repository iRepository) *Services {
	return &Services{
		repository: repository,
	}
}
