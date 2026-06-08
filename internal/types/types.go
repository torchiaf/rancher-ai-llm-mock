package types

type Text struct {
	Chunks []string `json:"chunks"`
}

type Args interface{}

type Tool struct {
	Name string `json:"name"`
	Args Args   `json:"args"`
}

type MockResponse struct {
	Text    Text   `json:"text"`
	MCPTool Tool   `json:"mcpTool,omitempty"`
	UITools []Tool `json:"uiTools,omitempty"`
}

type AgentResponse struct {
	Agent   string `json:"agent,omitempty"`
	MCPTool Tool   `json:"mcpTool,omitempty"`
}

type PushRequestArgs struct {
	AgentResponses []AgentResponse `json:"agentResponses,omitempty"`
	Text           Text            `json:"text,omitempty"`
	UITools        []Tool          `json:"uiTools,omitempty"`
}
