package main

import (
	"github.com/caleb-llh/data-enrichment-pipeline/shared" // replace with actual path
	"github.com/hashicorp/go-plugin"
)

// CompanyPlugin provides mock information for a company.
type CompanyPlugin struct{}

func (CompanyPlugin) Enrich(data map[string]string) map[string]string {
	if company, exists := data["company"]; exists {
		data["company_info"] = "MockInfoFor_" + company
	}
	return data
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "ENRICHMENT_PLUGIN",
			MagicCookieValue: "company",
		},
		Plugins: map[string]plugin.Plugin{
			"company": &shared.EnrichmentPluginRPC{Impl: CompanyPlugin{}},
		},
	})
}
