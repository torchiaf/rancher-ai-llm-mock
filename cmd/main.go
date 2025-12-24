package main

import (
	"fmt"
	"llm-mock/internal"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	port := ":8083"

	log.Println(fmt.Sprintf("Starting LLM Mock Service on %s...", port))

	r := gin.Default()
	svc := internal.NewLLMService()

	/**
	 * Control endpoints
	 * These endpoints are used to control the behavior of the mock service
	 */
	r.POST("/v1/control/push", svc.Control.HandlePushRequest)
	r.POST("/v1/control/clear", svc.Control.HandleClearRequest)

	/**
	 * Ollama endpoint
	 * https://docs.ollama.com/api/chat
	 */
	r.POST("/api/chat", svc.Models.Ollama.HandleRequest)

	/**
	 * Gemini endpoint
	 * https://ai.google.dev/api/generate-content#method:-models.streamgeneratecontent
	 */
	r.POST("/v1beta/models/:path", svc.Models.Gemini.HandleRequest)

	/**
	 * OpenAI endpoint
	 * https://platform.openai.com/docs/api-reference/chat/completions/create
	 */
	r.POST("/chat/completions", svc.Models.OpenAI.HandleRequest)

	r.Run(port)
}
