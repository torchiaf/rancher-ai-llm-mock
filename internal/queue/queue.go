package queue

import (
	types "llm-mock/internal/types"
	"runtime"
	"strings"
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
	if len(q.messages) == 0 {
		// Send default response with file name (model name)
		_, file, _, _ := runtime.Caller(1)
		// Remove directory path, keep only file name and remove .go extension
		model := file[strings.LastIndex(file, "/")+1 : strings.LastIndex(file, ".")]

		chunks := []string{
			"Mock service queue is empty. ",
			"This is ",
			"a default mock response ",
			"from the ",
			model,
			" model."}

		return types.MockResponse{Text: types.Text{Chunks: chunks}}
	}

	resp := q.messages[0]
	q.messages = q.messages[1:]
	return resp
}

func (q *Queue) Clear() {
	q.mu.Lock()
	q.messages = []types.MockResponse{}
	q.mu.Unlock()
}
