package main

import (
	"fmt"
	"os/exec"

	"github.com/caleb-llh/data-enrichment-pipeline/shared" // replace with actual path
	"github.com/hashicorp/go-plugin"
)

func main() {
	handshake := plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "ENRICHMENT_PLUGIN",
		MagicCookieValue: "geo",
	}

	plugins := map[string]string{
		"geo":     "../plugins/geo/geo",
		"email":   "../plugins/email/email",
		"company": "../plugins/company/company",
	}

	data := map[string]string{
		"ip":      "192.168.1.1",
		"email":   "user@example.com",
		"company": "OpenAI",
	}

	for name, path := range plugins {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: handshake,
			Plugins: map[string]plugin.Plugin{
				name: &shared.EnrichmentPluginRPC{},
			},
			Cmd: exec.Command(path),
		})

		rpcClient, err := client.Client()
		if err != nil {
			fmt.Println("Error starting plugin:", name, err)
			continue
		}

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
