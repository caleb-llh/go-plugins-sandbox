package main

import (
	"fmt"
	"os/exec"

	"github.com/caleb-llh/data-enrichment-pipeline/shared" // replace with actual path
	"github.com/hashicorp/go-plugin"
)

func main() {

	plugins := map[string]string{
		"geo":     "plugins/geo/geo",
		"email":   "plugins/email/email",
		"company": "plugins/company/company",
	}

	// input
	data := map[string]string{
		"ip":      "192.168.1.1",
		"email":   "user@example.com",
		"company": "OpenAI",
	}

	for name, path := range plugins {
		// Discover and Initialise
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: plugin.HandshakeConfig{
				ProtocolVersion:  1,
				MagicCookieKey:   "ENRICHMENT_PLUGIN",
				MagicCookieValue: name,
			},
			Plugins: map[string]plugin.Plugin{
				name: &shared.EnrichmentPluginRPC{},
			},
			Cmd: exec.Command(path),
		})

		// Connect via RPC
		rpcClient, err := client.Client()
		if err != nil {
			fmt.Println("Error starting plugin:", name, err)
			continue
		}

		// Plugin hook
		raw, err := rpcClient.Dispense(name)
		if err != nil {
			fmt.Println("Error dispensing plugin:", name, err)
			continue
		}

		enricher := raw.(shared.EnrichmentPlugin)
		data = enricher.Enrich(data)

		client.Kill()
	}

	fmt.Println("Enriched Data:", data)
}
