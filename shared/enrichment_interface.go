package shared

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

// EnrichmentPlugin defines the interface for all enrichment plugins.
type EnrichmentPlugin interface {
	Enrich(data map[string]string) map[string]string
}

// EnrichmentRPC is a client for communicating with the enrichment plugins.
type EnrichmentRPC struct {
	client *rpc.Client
}

func (e *EnrichmentRPC) Enrich(data map[string]string) map[string]string {
	var resp map[string]string
	e.client.Call("Plugin.Enrich", data, &resp)
	return resp
}

// EnrichmentPluginRPCServer serves RPC requests for the plugin.
type EnrichmentPluginRPCServer struct {
	Impl EnrichmentPlugin
}

func (s *EnrichmentPluginRPCServer) Enrich(args map[string]string, resp *map[string]string) error {
	*resp = s.Impl.Enrich(args)
	return nil
}

// EnrichmentPluginRPC wraps the server to implement the plugin.Plugin interface.
type EnrichmentPluginRPC struct {
	Impl EnrichmentPlugin
}

func (p *EnrichmentPluginRPC) Server(*plugin.MuxBroker) (interface{}, error) {
	return &EnrichmentPluginRPCServer{Impl: p.Impl}, nil
}

func (p *EnrichmentPluginRPC) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &EnrichmentRPC{client: c}, nil
}
