package services

import (
	"Messaggio/producer/internal/model"
)

func (s *Services) SaveMessage(message *model.Message) (string, error) {
	return s.repository.SaveMessage(message)
}

func (s *Services) GetMessage(id string) (*model.Message, error) {
	return s.repository.GetMessage(id)
}

func (s *Services) GetMessages() ([]*model.Message, error) {
	return s.repository.GetMessages()
}

func (s *Services) DeleteMessages() error {
	return s.repository.DeleteMessages()
}
