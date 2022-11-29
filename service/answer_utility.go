package service

import (
	"encoding/json"

	"github.com/calogxro/qaservice/db"
	"github.com/calogxro/qaservice/domain"
)

func RecreateAnswerState(store db.IEventStore, key string) (*domain.Answer, error) {
	var answer *domain.Answer
	events, _ := store.GetHistory(key)
	if len(events) > 0 {
		lastEvent := events[len(events)-1]
		if lastEvent != nil && lastEvent.Type != domain.ANSWER_DELETED_EVENT {
			err := json.Unmarshal([]byte(lastEvent.Data), &answer)
			if err != nil {
				return nil, err
			}
		}
	}
	return answer, nil
}
