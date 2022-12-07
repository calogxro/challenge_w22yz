package service

import (
	"github.com/calogxro/qaservice/db"
	"github.com/calogxro/qaservice/domain"
)

type QAProjection struct {
	repository db.ReadRepository
}

func NewQAProjection(repository db.ReadRepository) *QAProjection {
	return &QAProjection{
		repository: repository,
	}
}

func (p *QAProjection) GetAnswer(key string) (*domain.Answer, error) {
	answer, err := p.repository.GetAnswer(key)
	return answer, err
}
