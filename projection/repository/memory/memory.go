package readrepository

import "github.com/calogxro/qaservice/domain"

type ReadRepository struct {
	answersByKey map[string]*domain.Answer
}

func New() *ReadRepository {
	return &ReadRepository{
		answersByKey: make(map[string]*domain.Answer),
	}
}

func (rr *ReadRepository) GetAnswer(key string) (*domain.Answer, error) {
	answer, exists := rr.answersByKey[key]
	if !exists {
		return nil, domain.ErrKeyNotFound
	}
	return answer, nil
}

func (rr *ReadRepository) CreateAnswer(answer domain.Answer) error {
	rr.answersByKey[answer.Key] = &answer
	return nil
}

func (rr *ReadRepository) UpdateAnswer(answer domain.Answer) error {
	rr.answersByKey[answer.Key] = &answer
	return nil
}

func (rr *ReadRepository) DeleteAnswer(answer domain.Answer) error {
	delete(rr.answersByKey, answer.Key)
	return nil
}
