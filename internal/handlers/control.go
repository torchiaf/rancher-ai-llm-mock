package controlHandler

import (
	"llm-mock/internal/queue"
	"llm-mock/internal/types"

	"github.com/gin-gonic/gin"
)

type ControlHandler struct {
	queue *queue.Queue
}

func NewControlHandler(queue *queue.Queue) *ControlHandler {
	return &ControlHandler{
		queue: queue,
	}
}

func (s *ControlHandler) HandlePushRequest(c *gin.Context) {
	var req types.MockResponse

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(400)
		return
	}

	if req.Tool.Name != "" && req.Tool.Args == nil {
		c.JSON(400, gin.H{"error": "Invalid payload: Tool.Args must be provided when Tool is set"})
		return
	}

	if req.Tool.Name == "" && req.Tool.Args != nil {
		c.JSON(400, gin.H{"error": "Invalid payload: Tool.Name must be provided when Tool is set"})
		return
	}

	if (req.Text.Chunks == nil || len(req.Text.Chunks) == 0) && (req.Tool.Name == "" || req.Tool.Args == nil) {
		c.JSON(400, gin.H{"error": "Invalid payload: one of Text or Tool fields must be provided"})
		return
	}

	// If both Text and Tool are provided, push Tool first, then Text so that the Agent simulate mcp call behavior
	if len(req.Text.Chunks) > 0 && req.Tool.Name != "" {
		s.queue.Push(types.MockResponse{Tool: req.Tool})
		s.queue.Push(types.MockResponse{Text: req.Text})
	} else {
		s.queue.Push(req)
	}

	c.Status(204)
}

func (s *ControlHandler) HandleClearRequest(c *gin.Context) {
	s.queue.Clear()
	c.Status(204)
}
