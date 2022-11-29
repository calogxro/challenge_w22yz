package main

type ReadRepositoryStub struct {
	answersByKey map[string]*Answer
}

func NewReadRepositoryStub() *ReadRepositoryStub {
	return &ReadRepositoryStub{
		answersByKey: make(map[string]*Answer),
	}
}

func (rr *ReadRepositoryStub) GetAnswer(key string) (*Answer, error) {
	answer, exists := rr.answersByKey[key]
	if !exists {
		return nil, &KeyNotFound{}
	}
	return answer, nil
}

func (rr *ReadRepositoryStub) CreateAnswer(answer Answer) error {
	rr.answersByKey[answer.Key] = &answer
	return nil
}

func (rr *ReadRepositoryStub) UpdateAnswer(answer Answer) error {
	rr.answersByKey[answer.Key] = &answer
	return nil
}

func (rr *ReadRepositoryStub) DeleteAnswer(answer Answer) error {
	delete(rr.answersByKey, answer.Key)
	return nil
}
