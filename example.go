package main

import (
	"context"
	"fmt"

	"github.com/MrLeeang/langchaingo/mcp"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
)

var (
	configs = []*mcp.Config{
		{
			Name:      "my-mcp-server",
			Transport: "sse",
			URL:       "http://localhost:8080/sse",
			Disabled:  false,
		},
	}
)

func InitializeAgent(ctx context.Context, conversationId, systemPrompt string) (*agents.Executor, error) {

	// Initialize MCP tools
	tools, err := mcp.InitializeMCP(ctx, configs)
	if err != nil {
		return nil, err
	}

	// Initialize LLM with DeepSeek
	llm, err := openai.New(
		openai.WithBaseURL("https://api.deepseek.com/v1"),
		openai.WithToken("sk-deepseek-your-token"),
		openai.WithModel("deepseek-chat"),
	)

	if err != nil {
		return nil, err
	}

	// Initialize conversation memory
	mem := memory.NewConversationBuffer()

	// add system prompt
	mem.ChatHistory.AddUserMessage(ctx, systemPrompt)

	// add history messages
	// user input
	mem.ChatHistory.AddUserMessage(ctx, "hello world")
	// AI response
	mem.ChatHistory.AddAIMessage(ctx, "Hello! How can I assist you today?")

	agent := agents.NewExecutor(
		agents.NewConversationalAgent(llm, tools),
		agents.WithMemory(mem),
		// agents.WithCallbacksHandler(callbacks.LogHandler{}), // debug
	)
	return agent, nil
}

func main() {

	ctx := context.Background()

	systemPrompt := `
		You are an intelligent agent skilled in leveraging external tools.
	Always prioritize using available MCP tools to obtain accurate answers; 
	resort to direct reasoning only when no suitable tool is available. 
	Keep responses concise and chain multiple tools when necessary.
	`

	conversationId := "conversation-1234"

	agent, err := InitializeAgent(ctx, conversationId, systemPrompt)
	if err != nil {
		panic(err)
	}

	prompt := "I just greeted you, do you remember what I said?"

	result, err := chains.Run(
		ctx,
		agent,
		prompt,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
