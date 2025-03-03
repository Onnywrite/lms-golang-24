package orchestrator

type Service struct {
}

type Dependencies struct {
}

type Config struct {
	Dependencies
}

func New(deps Dependencies) *Service {
	return NewWithConfig(Config{
		Dependencies: deps,
	})
}

func NewWithConfig(Config) *Service {
	return nil
}
