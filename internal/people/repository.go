package people

import "go_test_task_2/config"

type Repository struct {
}

type RepositoryDeps struct {
	Config *config.Config
	DB     string
}

func NewRepository(config *config.Config) *Repository {
	return &Repository{}
}
