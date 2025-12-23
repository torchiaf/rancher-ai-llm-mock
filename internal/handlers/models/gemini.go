package modelHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"rancher-ai-llm-mock/internal/queue"
)

type GeminiHandler struct {
	queue *queue.Queue
}

func NewGeminiHandler(queue *queue.Queue) *GeminiHandler {
	return &GeminiHandler{
		queue: queue,
	}
}

func (s *GeminiHandler) HandleRequest(c *gin.Context) {
	/**
	 * Expected path format: {model}:{api-name}
	 * example: gemini-flash-2.0:streamGenerateContent
	 */
	path := c.Param("path")
	parts := strings.Split(path, ":")

	switch parts[1] {
	case "streamGenerateContent":
		s.HandleStreamGenerateContent(c)
	default:
		c.Status(404)
	}
}

func (s *GeminiHandler) HandleStreamGenerateContent(c *gin.Context) {
	w := c.Writer
	alt := c.Query("alt")
	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}
	response := s.queue.Pop()
	if alt == "sse" {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		for _, text := range response.Chunks {
			resp := map[string]interface{}{
				"candidates": []map[string]interface{}{
					{
						"content": map[string]interface{}{
							"parts": []map[string]interface{}{
								{"text": text},
							},
						},
						"finishReason": "length",
						"index":        0,
					},
				},
				"modelVersion": "gemini-mock-v1",
				"responseId":   "resp-mock-12345",
			}
			b, _ := json.Marshal(resp)
			w.Write([]byte("data: "))
			w.Write(b)
			w.Write([]byte("\n\n"))
			flusher.Flush()
			time.Sleep(200 * time.Millisecond)
		}
		w.Write([]byte("data: [DONE]\n\n"))
		flusher.Flush()
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "chunked")
		fmt.Fprint(w, "[")
		flusher.Flush()
		for i, text := range response.Chunks {
			resp := map[string]interface{}{
				"candidates": []map[string]interface{}{
					{
						"content": map[string]interface{}{
							"parts": []map[string]interface{}{
								{"text": text},
							},
						},
						"finishReason": "length",
						"index":        0,
					},
				},
				"modelVersion": "gemini-mock-v1",
				"responseId":   "resp-mock-12345",
			}
			enc := json.NewEncoder(w)
			if err := enc.Encode(resp); err != nil {
				return
			}
			if i < len(response.Chunks)-1 {
				fmt.Fprint(w, ",")
			}
			flusher.Flush()
			time.Sleep(200 * time.Millisecond)
		}
		// Close the array
		fmt.Fprint(w, "]")
		flusher.Flush()
	}
}
