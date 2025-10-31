package mcp

import (
	"context"
	"log"

	"encoding/json"

	"github.com/tmc/langchaingo/tools"

	client "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func InitializeMCP(ctx context.Context, configs []*Config) ([]tools.Tool, error) {
	var tools []tools.Tool

	for _, cfg := range configs {
		if cfg.Disabled {
			continue
		}
		spec := ConnSpec{Name: cfg.Name, Transport: cfg.Transport, Endpoint: cfg.URL, Command: cfg.Command, Args: cfg.Args}

		// 尝试列出远端工具并逐个注册 MCPTool
		transport, err := newTransportFromSpec(spec)
		if err != nil {
			log.Fatalln(err)
		}
		c := client.NewClient(transport)
		if err := c.Start(context.Background()); err != nil {
			log.Fatalln(err)
		}

		defer c.Close()
		// 初始化客户端
		if _, err := c.Initialize(context.Background(), mcp.InitializeRequest{
			Params: mcp.InitializeParams{
				ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
				ClientInfo:      mcp.Implementation{Name: "tool-enumerator", Version: "0.1.0"},
			},
		}); err != nil {
			log.Printf("初始化 MCP 客户端失败 (%s): %v", cfg.Name, err)
			continue
		}
		toolsList, err := c.ListTools(context.Background(), mcp.ListToolsRequest{})
		if err != nil {
			log.Printf("列出远端工具失败 (%s): %v", cfg.Name, err)
			continue
		}

		for _, rt := range toolsList.Tools {
			if rt.Name == "" {
				continue
			}

			inputSchema, err := json.Marshal(rt.InputSchema)
			if err != nil {
				log.Printf("序列化输入模式失败 (%s - %s): %v", cfg.Name, rt.Name, err)
				continue
			}

			outputSchema, err := json.Marshal(rt.OutputSchema)
			if err != nil {
				log.Printf("序列化输入模式失败 (%s - %s): %v", cfg.Name, rt.Name, err)
				continue
			}

			tools = append(tools, &Tool{
				Conn:       spec,
				RemoteName: rt.Name,
				RemoteDesc: rt.Description,
				Args:       string(inputSchema),
				Output:     string(outputSchema),
			})
		}
	}

	return tools, nil
}
