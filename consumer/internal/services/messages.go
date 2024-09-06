package services

func (s *Services) TestProcessMessage(id string) error {
	return s.repository.EmptyProcess(id)
}
