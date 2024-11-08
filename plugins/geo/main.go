package main

import (
	"github.com/caleb-llh/data-enrichment-pipeline/shared" // replace with actual path
	"github.com/hashicorp/go-plugin"
)

// GeoPlugin provides geolocation information based on an IP address.
type GeoPlugin struct{}

func (GeoPlugin) Enrich(data map[string]string) map[string]string {
	if ip, exists := data["ip"]; exists {
		data["location"] = "MockLocationForIP_" + ip // Mock geolocation
	}
	return data
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "ENRICHMENT_PLUGIN",
			MagicCookieValue: "geo",
		},
		Plugins: map[string]plugin.Plugin{
			"geo": &shared.EnrichmentPluginRPC{Impl: GeoPlugin{}},
		},
	})
}
