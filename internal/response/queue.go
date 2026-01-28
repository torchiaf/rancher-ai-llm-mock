package response

import (
	types "llm-mock/internal/types"
	"sync"
)

type Queue struct {
	mu       sync.RWMutex
	messages []types.MockResponse
}

func NewQueue() *Queue {
	return &Queue{
		messages: []types.MockResponse{},
	}
}

func (q *Queue) Push(response types.MockResponse) {
	q.mu.Lock()
	q.messages = append(q.messages, response)
	q.mu.Unlock()
}

func (q *Queue) Pop() types.MockResponse {
	q.mu.Lock()
	defer q.mu.Unlock()

	resp := types.MockResponse{}

	if len(q.messages) > 0 {
		resp = q.messages[0]
		q.messages = q.messages[1:]
	}

	return resp
}

func (q *Queue) Clear() {
	q.mu.Lock()
	q.messages = []types.MockResponse{}
	q.mu.Unlock()
}
