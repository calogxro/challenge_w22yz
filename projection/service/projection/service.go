package projection

import (
	"github.com/calogxro/qaservice/domain"
)

type ReadRepository interface {
	GetAnswer(key string) (*domain.Answer, error)
	// CreateAnswer(answer domain.Answer) error
	// UpdateAnswer(answer domain.Answer) error
	// DeleteAnswer(answer domain.Answer) error
}

type Projection struct {
	repository ReadRepository
}

func New(repository ReadRepository) *Projection {
	return &Projection{
		repository: repository,
	}
}

func (p *Projection) GetAnswer(key string) (*domain.Answer, error) {
	answer, err := p.repository.GetAnswer(key)
	return answer, err
}
