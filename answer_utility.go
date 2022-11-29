package main

import (
	"encoding/json"
)

func RecreateAnswerState(store IEventStore, key string) (*Answer, error) {
	var answer *Answer
	events, _ := store.GetHistory(key)
	if len(events) > 0 {
		lastEvent := events[len(events)-1]
		if lastEvent != nil && lastEvent.Type != ANSWER_DELETED_EVENT {
			err := json.Unmarshal([]byte(lastEvent.Data), &answer)
			if err != nil {
				return nil, err
			}
		}
	}
	return answer, nil
}
