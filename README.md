# langchaingo

This is an extension of [langchaingo](https://github.com/tmc/langchaingo) that adds Model Context Protocol (MCP) support. While the original langchaingo project provides excellent functionality for working with language models, it lacks built-in support for MCP. This extension bridges that gap by implementing MCP functionality that can be seamlessly integrated into langchaingo.

## Features

- Full MCP (Model Context Protocol) support for langchaingo
- Easy integration with existing langchaingo applications
- Support for remote tool execution via MCP
- Configurable transport options (HTTP, WebSocket, etc.)

## Installation

```bash
go get github.com/MrLeeang/langchaingo
```

## Usage

Here's a simple example of how to use the MCP extension with langchaingo:

```go
package main

import (
    "context"
    "log"
    "fmt"

    "github.com/MrLeeang/langchaingo/mcp"
    "github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
    "github.com/tmc/langchaingo/chains"
)

func main() {
    // Configure MCP
    configs := []*mcp.Config{
        {
            Name:      "my-mcp-server",
            Transport: "sse",
            URL:      "http://localhost:8080/sse",
            Disabled: false,
        },
    }

    ctx := context.Background()

    // Initialize MCP tools
    tools, err := mcp.InitializeMCP(ctx, configs)
    if err != nil {
        panic(err)
    }

    // Initialize LLM with DeepSeek
	llm, err := openai.New(
		openai.WithBaseURL("https://api.deepseek.com/v1"),
		openai.WithToken("sk-deepseek-your-token"),
		openai.WithModel("deepseek-chat"),
	)

    agent := agents.NewExecutor(
		agents.NewConversationalAgent(llm, tools),
	)

    prompt := "Get the latest Beijing time."

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
```

## Configuration

The MCP extension supports various configuration options:

```go
type Config struct {
    Name      string   // Name of the MCP server
    Transport string   // Transport protocol (http, ws, etc.)
    URL       string   // Endpoint URL
    Command   string   // Command to execute (optional)
    Args      []string // Command arguments (optional)
    Disabled  bool     // Whether the tool is disabled
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the same terms as langchaingo.
