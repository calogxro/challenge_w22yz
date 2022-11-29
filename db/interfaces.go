package db

import "github.com/calogxro/qaservice/domain"

type IEventStore interface {
	GetEvents() ([]*domain.Event, error)
	AddEvent(event *domain.Event) error
	GetHistory(key string) ([]*domain.Event, error)
	Subscribe(onEvent func(*domain.Event)) error
}

type IReadRepository interface {
	GetAnswer(key string) (*domain.Answer, error)
	CreateAnswer(answer domain.Answer) error
	UpdateAnswer(answer domain.Answer) error
	DeleteAnswer(answer domain.Answer) error
}
