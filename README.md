# llm-mock

A mock server for Large Language Model (LLM) APIs, supporting OpenAI, Ollama, Gemini and AWS Bedrock endpoints. Useful for testing and development without relying on actual LLM services.

## Installation

1. Clone this repository:
	```sh
	git clone rancher-sandbox/rancher-ai-llm-mock.git
	cd rancher-ai-llm-mock
	```

2. Build and run with Go:
	```sh
	go run ./cmd/main.go
	```
	Or use Docker:
	```sh
	docker build -t llm-mock .
	docker run -p 8083:8083 llm-mock
	```
    Or use Helm:
    ```sh
    helm install llm-mock chart/llm-mock \
        --namespace your-namespace \
        --create-namespace
    ```

The server will start on http://localhost:8083

## Controlling Responses
You can control the mock responses using the `/v1/control` endpoints:

- `POST /v1/control/push`: Push mock responses to the queue. Example body:
	```json
	{
		"agentResponses": [
			{
				"agent": "rancher",
				"mcpTool": {
					"name": "mcp_tool_name",
					"args": {
						"key": "value"
					}
				}
			}
		],
		"text": {
			"chunks": ["Hello", " world!"]
		},
		"uiTools": [
			{
				"name": "explore",
				"args": {
					"arg1": ["Option 1", "Option 2"],
					"arg2": "Additional info",
					"arg3": 123
				}
			}
		]
	}
	```
	- `agentResponses`: Array of agent responses with optional MCP tool calls
	- `text`: Optional text response containing message chunks
	- `uiTools`: Optional array of UI tools to display. Valid tools can be found in rancher-ai-ui repository in `ui-tools.json`
	- The `args` field in tools accepts either a single object or an array of objects (e.g. confirmation request payload for multiple resources)
	- The next model API call will stream text chunks as response, use mcpTool for MCP invocation and uiTools for returning the ui tools calls.
	- The MCP tool must be one of the supported tools of the agent in request.
	- If there are less than two agents configured in Rancher, the agent must not be provided.

- `POST /v1/control/clear`: Clear the mock response queue.

If the response queue is empty, default mock responses will be used.

See the [OpenAPI spec](openapi.yaml) for full API details.
