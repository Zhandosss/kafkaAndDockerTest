package services

type iRepository interface {
	EmptyProcess(id string) error
}

type Services struct {
	repository iRepository
}

func New(repository iRepository) *Services {
	return &Services{
		repository: repository,
	}
}
