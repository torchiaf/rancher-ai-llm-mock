package modelHandlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"llm-mock/internal/response"
)

type BedrockHandler struct {
	response *response.Handler
}

func NewBedrockHandler(response *response.Handler) *BedrockHandler {
	return &BedrockHandler{
		response: response,
	}
}

/**
 * TODO: Currently, the entire response is sent in a single stream.
 * In the real-world scenario, the response should be sent in smaller chunks as they become available.
 */
func (s *BedrockHandler) HandleRequest(c *gin.Context) {
	response := s.response.Pop()

	c.Header("Content-Type", "application/json")
	c.Header("Transfer-Encoding", "chunked")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return
	}

	c.Stream(func(w io.Writer) bool {
		w.Write([]byte(`{"output":{"message":{"role":"assistant","content":[`))

		totalChunks := len(response.Text.Chunks)
		first := true

		if response.Tool.Name != "" {
			jsonTool, _ := json.Marshal(buildTool(response.Tool.Name, response.Tool.Args))
			w.Write(jsonTool)
			first = false
		}

		for _, chunkText := range response.Text.Chunks {
			if !first {
				w.Write([]byte(","))
			}
			jsonBlock, _ := json.Marshal(buildText(chunkText))
			w.Write(jsonBlock)
			first = false
			flusher.Flush()

			time.Sleep(100 * time.Millisecond)
		}

		usageJSON := map[string]interface{}{
			"inputTokens":  15,
			"outputTokens": totalChunks,
			"totalTokens":  15 + totalChunks,
		}
		usageBytes, _ := json.Marshal(usageJSON)

		w.Write([]byte(`]}},"stopReason":"end_turn","usage":`))
		w.Write(usageBytes)
		w.Write([]byte(`}`))

		return false
	})
}

func buildTool(name string, args interface{}) map[string]interface{} {
	return map[string]interface{}{
		"toolUse": map[string]interface{}{
			"toolUseId": "call_123",
			"name":      name,
			"input":     args,
		},
	}
}

func buildText(text string) map[string]interface{} {
	return map[string]interface{}{"text": text}
}
