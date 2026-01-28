package modelHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"llm-mock/internal/response"
	"llm-mock/internal/types"
)

type GeminiHandler struct {
	response *response.Handler
}

func NewGeminiHandler(response *response.Handler) *GeminiHandler {
	return &GeminiHandler{
		response: response,
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
		s.handleStreamGenerateContent(c)
	default:
		c.Status(404)
	}
}

func (s *GeminiHandler) handleStreamGenerateContent(c *gin.Context) {
	w := c.Writer
	alt := c.Query("alt")
	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	response := s.response.Pop()

	if alt == "sse" {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		if response.Tool.Name != "" {
			toolResp := s.buildToolResponse(response.Tool)
			b, _ := json.Marshal(toolResp)
			w.Write([]byte("data: "))
			w.Write(b)
			w.Write([]byte("\n\n"))
			flusher.Flush()
		} else {
			for _, chunk := range response.Text.Chunks {
				resp := s.buildTextResponse(chunk)
				b, _ := json.Marshal(resp)
				w.Write([]byte("data: "))
				w.Write(b)
				w.Write([]byte("\n\n"))
				flusher.Flush()
				time.Sleep(100 * time.Millisecond)
			}
		}
		w.Write([]byte("data: {\"done\": true}\n\n"))
		flusher.Flush()
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Transfer-Encoding", "chunked")
		fmt.Fprint(w, "[")
		flusher.Flush()

		if response.Tool.Name != "" {
			toolResp := s.buildToolResponse(response.Tool)
			enc := json.NewEncoder(w)
			if err := enc.Encode(toolResp); err != nil {
				return
			}
		} else {
			for i, chunk := range response.Text.Chunks {
				resp := s.buildTextResponse(chunk)
				enc := json.NewEncoder(w)
				if err := enc.Encode(resp); err != nil {
					return
				}
				if i < len(response.Text.Chunks)-1 {
					fmt.Fprint(w, ",")
				}
				flusher.Flush()
				time.Sleep(100 * time.Millisecond)
			}
		}

		// Close the array
		fmt.Fprint(w, "]")
		flusher.Flush()
	}
}

func (s *GeminiHandler) buildTextResponse(chunk string) map[string]interface{} {
	return map[string]interface{}{
		"candidates": []map[string]interface{}{
			{
				"content": map[string]interface{}{
					"parts": []map[string]interface{}{
						{"text": chunk},
					},
				},
				"finishReason": "STOP",
				"index":        0,
			},
		},
	}
}

func (s *GeminiHandler) buildToolResponse(tool types.Tool) map[string]interface{} {
	return map[string]interface{}{
		"candidates": []map[string]interface{}{
			{
				"content": map[string]interface{}{
					"parts": []map[string]interface{}{
						{
							"function_call": map[string]interface{}{
								"name": tool.Name,
								"args": tool.Args,
							},
						},
					},
				},
				"finishReason": "STOP",
				"index":        0,
			},
		},
	}
}
