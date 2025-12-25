package modelHandlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"llm-mock/internal/queue"
	"llm-mock/internal/types"
)

type OllamaHandler struct {
	queue *queue.Queue
}

func NewOllamaHandler(queue *queue.Queue) *OllamaHandler {
	return &OllamaHandler{
		queue: queue,
	}
}

func (s *OllamaHandler) HandleRequest(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	response := s.queue.Pop()

	if response.Tool.Name != "" {
		resp := s.buildToolResponse(response.Tool)
		enc := json.NewEncoder(w)
		if err := enc.Encode(resp); err != nil {
			return
		}
		flusher.Flush()
	} else {
		for i, text := range response.Text.Chunks {
			resp := s.buildTextResponse(text, i == len(response.Text.Chunks)-1)
			enc := json.NewEncoder(w)
			if err := enc.Encode(resp); err != nil {
				return
			}
			flusher.Flush()
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *OllamaHandler) buildTextResponse(chunk string, done bool) map[string]interface{} {
	return map[string]interface{}{
		"message": map[string]interface{}{
			"role":    "assistant",
			"content": chunk,
		},
		"done": done,
	}
}

func (s *OllamaHandler) buildToolResponse(tool types.Tool) map[string]interface{} {
	return map[string]interface{}{
		"message": map[string]interface{}{
			"role":    "assistant",
			"content": "",
			"tool_calls": []map[string]interface{}{
				{
					"function": map[string]interface{}{
						"name":      tool.Name,
						"arguments": tool.Args,
					},
				},
			},
		},
		"done_reason": "stop",
		"done":        true,
	}
}
