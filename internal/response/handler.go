package response

import (
	types "llm-mock/internal/types"
	"runtime"
	"strings"
)

type Handler struct {
	queue *Queue
}

func NewHandler() *Handler {
	return &Handler{
		queue: NewQueue(),
	}
}

func getModelNameFromCaller(index int) string {
	// Send default response with file name (model name)
	_, file, _, _ := runtime.Caller(index)

	// Remove directory path, keep only file name and remove .go extension
	model := file[strings.LastIndex(file, "/")+1 : strings.LastIndex(file, ".")]

	return model
}

func generateResponse() types.MockResponse {
	model := getModelNameFromCaller(3)

	chunks := []string{
		"Mock service: messages queue is empty. ",
		"This is ",
		"a default mock response ",
		"from the ",
		model,
		" model."}

	return types.MockResponse{Text: types.Text{Chunks: chunks}}
}

func (h *Handler) Push(req types.MockResponse) {
	// If Agent is provided, push it first so the we can simulate llm's agent selection behavior
	if req.Agent != "" {
		h.queue.Push(types.MockResponse{Text: types.Text{Chunks: []string{req.Agent}}})
	}

	// If both Text and Tool are provided, push Tool first, then Text so that the Agent simulate mcp call behavior
	if len(req.Text.Chunks) > 0 && req.Tool.Name != "" {
		h.queue.Push(types.MockResponse{Tool: req.Tool})
		h.queue.Push(types.MockResponse{Text: req.Text})
	} else {
		h.queue.Push(req)
	}
}

func (h *Handler) Pop() types.MockResponse {
	if len(h.queue.messages) == 0 {
		return generateResponse()
	}

	return h.queue.Pop()
}

func (h *Handler) Clear() {
	h.queue.Clear()
}
