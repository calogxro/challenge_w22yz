package db

import "github.com/calogxro/qaservice/domain"

type ReadRepositoryStub struct {
	answersByKey map[string]*domain.Answer
}

func NewReadRepositoryStub() *ReadRepositoryStub {
	return &ReadRepositoryStub{
		answersByKey: make(map[string]*domain.Answer),
	}
}

func (rr *ReadRepositoryStub) GetAnswer(key string) (*domain.Answer, error) {
	answer, exists := rr.answersByKey[key]
	if !exists {
		return nil, &domain.KeyNotFound{}
	}
	return answer, nil
}

func (rr *ReadRepositoryStub) CreateAnswer(answer domain.Answer) error {
	rr.answersByKey[answer.Key] = &answer
	return nil
}

func (rr *ReadRepositoryStub) UpdateAnswer(answer domain.Answer) error {
	rr.answersByKey[answer.Key] = &answer
	return nil
}

func (rr *ReadRepositoryStub) DeleteAnswer(answer domain.Answer) error {
	delete(rr.answersByKey, answer.Key)
	return nil
}
