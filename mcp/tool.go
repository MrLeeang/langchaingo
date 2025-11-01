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
	desc := fmt.Sprintf("\nname: %s，desc: %s, args_schema: %s", t.RemoteName, t.RemoteDesc, t.Args)
	return desc
}

func (t *Tool) Call(ctx context.Context, input string) (string, error) {
	transport, err := newTransportFromSpec(t.Conn)
	if err != nil {
		return "", err
	}

	c := mcpclient.NewClient(transport)
	if err := c.Start(ctx); err != nil {
		return "", err
	}
	defer c.Close()

	initReq := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo:      mcp.Implementation{Name: "remote-mcp-tool", Version: "0.1.0"},
		},
	}
	if _, err := c.Initialize(ctx, initReq); err != nil {
		return "", err
	}

	// prepare args, try to parse input as json or put it in "input" field
	var args map[string]any
	if err := json.Unmarshal([]byte(input), &args); err != nil || args == nil {
		args = map[string]any{"input": input}
	}

	result, err := c.CallTool(ctx, mcp.CallToolRequest{Params: mcp.CallToolParams{Name: t.RemoteName, Arguments: args}})
	if err != nil {
		return "", err
	}

	// return first text content found
	for _, part := range result.Content {
		// Can be TextContent, ImageContent, AudioContent, or EmbeddedResource
		switch v := part.(type) {
		case mcp.TextContent:
			return v.Text, nil
		case *mcp.TextContent:
			return v.Text, nil
		case mcp.ImageContent:
			return v.Data, nil
		case *mcp.ImageContent:
			return v.Data, nil
		case mcp.AudioContent:
			return v.Data, nil
		case *mcp.AudioContent:
			return v.Data, nil
			// case mcp.EmbeddedResource:
			// 	return v.Resource, nil
			// case *mcp.EmbeddedResource:
			// 	return v.Resource, nil
		}
	}

	// return "" if no text content found
	return "", nil
}
