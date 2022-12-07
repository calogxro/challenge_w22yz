package service

import (
	"encoding/json"

	"github.com/calogxro/qaservice/db"
	"github.com/calogxro/qaservice/domain"
)

type Projector struct {
	repository db.ReadRepository
}

func NewProjector(repository db.ReadRepository) *Projector {
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
