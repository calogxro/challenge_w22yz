package eventstore

import (
	"log"

	"github.com/calogxro/qaservice/domain"
)

type EventStore interface {
	GetEvents() ([]*domain.Event, error)
	AddEvent(event *domain.Event) error
	GetHistory(key string) ([]*domain.Event, error)
}

type Service struct {
	eventStore EventStore
}

func New(es EventStore) *Service {
	return &Service{
		eventStore: es,
	}
}

func (s *Service) CreateAnswer(answer domain.Answer) (*domain.Event, error) {
	if answer.Key == "" || answer.Value == "" {
		return nil, domain.ErrInputNotValid
	}

	if AnswerExists(s.eventStore, answer.Key) {
		return nil, domain.ErrKeyExists
	}

	event, err := domain.NewAnswerCreatedEvent(answer)

	if err != nil {
		log.Println("domain.NewAnswerCreatedEvent", err)
		return nil, err
	}

	err = s.eventStore.AddEvent(event)

	if err != nil {
		log.Println("(db.EventStore).AddEvent", err)
		return nil, err
	}

	return event, nil
}

func (s *Service) UpdateAnswer(answer domain.Answer) (*domain.Event, error) {
	if answer.Key == "" || answer.Value == "" {
		return nil, domain.ErrInputNotValid
	}

	if !AnswerExists(s.eventStore, answer.Key) {
		return nil, domain.ErrKeyNotFound
	}

	event, _ := domain.NewAnswerUpdatedEvent(answer)
	s.eventStore.AddEvent(event)
	return event, nil
}

func (s *Service) DeleteAnswer(key string) (*domain.Event, error) {
	if !AnswerExists(s.eventStore, key) {
		return nil, domain.ErrKeyNotFound
	}
	event, _ := domain.NewAnswerDeletedEvent(domain.Answer{Key: key})
	s.eventStore.AddEvent(event)
	return event, nil
}

func (s *Service) GetHistory(key string) ([]*domain.Event, error) {
	events, err := s.eventStore.GetHistory(key)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, domain.ErrKeyNotFound
	}

	return events, nil
}
