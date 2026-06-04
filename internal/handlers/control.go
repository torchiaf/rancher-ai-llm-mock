package controlHandler

import (
	"llm-mock/internal/response"
	"llm-mock/internal/types"

	"github.com/gin-gonic/gin"
)

type ControlHandler struct {
	response *response.Handler
}

func NewControlHandler(response *response.Handler) *ControlHandler {
	return &ControlHandler{
		response: response,
	}
}

func (s *ControlHandler) HandlePushRequest(c *gin.Context) {
	var req types.PushRequestArgs

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(400)
		return
	}

	for i := 0; i < len(req.AgentResponses); i++ {
		var resp = req.AgentResponses[i]

		if len(req.AgentResponses) > 1 && resp.Agent == "" {
			c.JSON(400, gin.H{"error": "Invalid payload: Agent must be provided for all AgentResponses when there are multiple AgentResponses"})
			return
		}

		if resp.MCPTool.Name != "" && resp.MCPTool.Args == nil {
			c.JSON(400, gin.H{"error": "Invalid payload: MCPTool.Args must be provided when MCPTool is set"})
			return
		}

		if resp.MCPTool.Name == "" && resp.MCPTool.Args != nil {
			c.JSON(400, gin.H{"error": "Invalid payload: MCPTool.Name must be provided when MCPTool is set"})
			return
		}
	}

	if req.Text.Chunks == nil || len(req.Text.Chunks) == 0 {
		c.JSON(400, gin.H{"error": "Invalid payload: Text.Chunks must be provided"})
		return
	}

	s.response.Push(req)

	c.Status(204)
}

func (s *ControlHandler) HandleClearRequest(c *gin.Context) {
	s.response.Clear()
	c.Status(204)
}
