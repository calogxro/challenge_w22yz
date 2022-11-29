package main

type QAProjection struct {
	repository IReadRepository
}

func NewQAProjection(repository IReadRepository) *QAProjection {
	return &QAProjection{
		repository: repository,
	}
}

func (p *QAProjection) GetAnswer(key string) (*Answer, error) {
	answer, err := p.repository.GetAnswer(key)
	return answer, err
}
