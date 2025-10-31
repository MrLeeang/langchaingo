package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

type Tool struct {
	Conn       ConnSpec
	RemoteName string
	RemoteDesc string
	Args       string
	Output     string
}

func (t *Tool) Name() string {
	return t.RemoteName
}

func (t *Tool) Description() string {
	desc := fmt.Sprintf("\n工具名称: %s，描述: %s, 参数: %s", t.RemoteName, t.RemoteDesc, t.Args)
	return desc
}

func (t *Tool) Call(ctx context.Context, input string) (string, error) {
	transport, err := newTransportFromSpec(t.Conn)
	if err != nil {
		return "", fmt.Errorf("创建传输失败: %w", err)
	}

	c := mcpclient.NewClient(transport)
	if err := c.Start(ctx); err != nil {
		return "", fmt.Errorf("启动客户端失败: %w", err)
	}
	defer c.Close()

	initReq := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo:      mcp.Implementation{Name: "remote-mcp-tool", Version: "0.1.0"},
		},
	}
	if _, err := c.Initialize(ctx, initReq); err != nil {
		return "", fmt.Errorf("初始化失败: %w", err)
	}

	// 尝试把 input 当作 JSON 参数解析，否则放到 "input" 字段
	var args map[string]any
	if err := json.Unmarshal([]byte(input), &args); err != nil || args == nil {
		args = map[string]any{"input": input}
	}

	result, err := c.CallTool(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Name: t.RemoteName, Arguments: args}})
	if err != nil {
		return "", fmt.Errorf("调用远端工具 %s 失败: %w", t.RemoteName, err)
	}
	// 返回首个文本片段（如无文本则返回摘要）
	for _, part := range result.Content {
		switch v := part.(type) {
		case mcp.TextContent:
			return v.Text, nil
		case *mcp.TextContent:
			return v.Text, nil
		}
	}
	return fmt.Sprintf("工具 %s 已返回 %d 段内容（非纯文本）", t.RemoteName, len(result.Content)), nil
}
