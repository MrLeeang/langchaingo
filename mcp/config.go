package mcp

type Config struct {
	Name        string
	URL         string
	Transport   string
	Description string
	TimeoutSec  int
	Disabled    bool
	Command     string
	Args        []string
}
