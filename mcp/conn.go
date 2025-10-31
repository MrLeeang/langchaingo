package mcp

import (
	"fmt"

	mcpxport "github.com/mark3labs/mcp-go/client/transport"
)

type ConnSpec struct {
	Name      string
	Transport string
	Endpoint  string
	Command   string
	Args      []string
}

func newTransportFromSpec(spec ConnSpec) (mcpxport.Interface, error) {
	switch spec.Transport {
	case "sse":
		if spec.Endpoint == "" {
			return nil, fmt.Errorf("sse 需要 endpoint")
		}
		return mcpxport.NewSSE(spec.Endpoint)
	case "streamable_http":
		if spec.Endpoint == "" {
			return nil, fmt.Errorf("streamable_http 需要 endpoint")
		}
		return mcpxport.NewStreamableHTTP(spec.Endpoint)
	case "stdio":
		if spec.Command == "" {
			return nil, fmt.Errorf("stdio 需要 command")
		}
		tr := mcpxport.NewStdio(spec.Command, spec.Args, []string{}...)
		return tr, nil
	default:
		return nil, fmt.Errorf("不支持的传输: %s", spec.Transport)
	}
}
