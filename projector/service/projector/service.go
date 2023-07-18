package projector

import (
	"encoding/json"

	"github.com/calogxro/qaservice/domain"
)

type Repository interface {
	CreateAnswer(answer domain.Answer) error
	UpdateAnswer(answer domain.Answer) error
	DeleteAnswer(answer domain.Answer) error
}

type Projector struct {
	repository Repository
}

func New(repository Repository) *Projector {
	return &Projector{
		repository: repository,
	}
}

func (p *Projector) Project(event *domain.Event) {
	var answer domain.Answer
	json.Unmarshal([]byte(event.Data), &answer)

	if event.Type == domain.ANSWER_CREATED_EVENT {
		p.repository.CreateAnswer(answer)
	}
	if event.Type == domain.ANSWER_UPDATED_EVENT {
		p.repository.UpdateAnswer(answer)
	}
	if event.Type == domain.ANSWER_DELETED_EVENT {
		p.repository.DeleteAnswer(answer)
	}
}
