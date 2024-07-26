package services

import "Messaggio/consumer/internal/model"

func (s *Services) TestProcessMessage(message *model.Message) error {
	return s.repository.EmptyProcess(message)
}
