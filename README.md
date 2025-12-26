# llm-mock

A mock server for Large Language Model (LLM) APIs, supporting OpenAI, Ollama, and Gemini endpoints. Useful for testing and development without relying on actual LLM services.

## Installation

1. Clone this repository:
	```sh
	git clone https://github.com/your-org/llm-mock.git
	cd llm-mock
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

- `POST /v1/control/push`: Push a mock response to the queue. Example body:
	```json
	{
		"text": {
			"chunks": ["Hello", " world!"]
		},
		"tool": {
			"name": "example_tool",
			"args": {
				"key": "value"
			}
		}
	}
	```
	The next model API call will stream text chunks as response and use tool for MCP invocation.

- `POST /v1/control/clear`: Clear the mock response queue.

If the queue is empty, default mock responses will be used.

See the [OpenAPI spec](openapi.yaml) for full API details.
