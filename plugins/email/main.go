package main

import (
	"strings"

	"github.com/caleb-llh/data-enrichment-pipeline/shared" // replace with actual path
	"github.com/hashicorp/go-plugin"
)

// EmailPlugin validates the format of an email.
type EmailPlugin struct{}

func (EmailPlugin) Enrich(data map[string]string) map[string]string {
	if email, exists := data["email"]; exists {
		if strings.Contains(email, "@") {
			data["email_valid"] = "true"
		} else {
			data["email_valid"] = "false"
		}
	}
	return data
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "ENRICHMENT_PLUGIN",
			MagicCookieValue: "email",
		},
		Plugins: map[string]plugin.Plugin{
			"email": &shared.EnrichmentPluginRPC{Impl: EmailPlugin{}},
		},
	})
}
