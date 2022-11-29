package main

type QAService struct {
	eventStore IEventStore
}

func NewQAService(es IEventStore) *QAService {
	return &QAService{
		eventStore: es,
	}
}

func (s *QAService) answerExists(key string) bool {
	answer, _ := RecreateAnswerState(s.eventStore, key)
	return answer != nil
}

func (s *QAService) CreateAnswer(answer Answer) (*Event, error) {
	if s.answerExists(answer.Key) {
		return nil, &KeyExists{}
	}
	event, _ := NewAnswerCreatedEvent(answer)
	s.eventStore.AddEvent(event)
	return event, nil
}

func (s *QAService) UpdateAnswer(answer Answer) (*Event, error) {
	if !s.answerExists(answer.Key) {
		return nil, &KeyNotFound{}
	}
	event, _ := NewAnswerUpdatedEvent(answer)
	s.eventStore.AddEvent(event)
	return event, nil
}

func (s *QAService) DeleteAnswer(key string) (*Event, error) {
	if !s.answerExists(key) {
		return nil, &KeyNotFound{}
	}
	event, _ := NewAnswerDeletedEvent(Answer{Key: key})
	s.eventStore.AddEvent(event)
	return event, nil
}

func (s *QAService) GetHistory(key string) ([]*Event, error) {
	events, _ := s.eventStore.GetHistory(key)
	return events, nil
}
