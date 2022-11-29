package main

import "encoding/json"

type Projector struct {
	repository IReadRepository
}

func NewProjector(repository IReadRepository) *Projector {
	return &Projector{
		repository: repository,
	}
}

func (p *Projector) Project(event *Event) {
	var answer Answer
	json.Unmarshal([]byte(event.Data), &answer)

	if event.Type == ANSWER_CREATED_EVENT {
		p.repository.CreateAnswer(answer)
	}
	if event.Type == ANSWER_UPDATED_EVENT {
		p.repository.UpdateAnswer(answer)
	}
	if event.Type == ANSWER_DELETED_EVENT {
		p.repository.DeleteAnswer(answer)
	}
}
