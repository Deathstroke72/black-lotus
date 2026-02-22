package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Deathstroke72/black-lotus/lotus-agents/config"
	"github.com/Deathstroke72/black-lotus/lotus-agents/orchestrator"
)

func main() {
	cfg := config.Load()

	if cfg.AnthropicAPIKey == "" {
		log.Fatal("ANTHROPIC_API_KEY environment variable is required")
	}

	// ---------------------------------------------------------------
	// Define your microservice here.
	// Swap this out for config.PaymentsService(), config.NotificationsService(),
	// or define your own ServiceDefinition from scratch:
	//
	//	svc := &config.ServiceDefinition{
	//	    Name:        "orders",
	//	    Description: "Manages the order lifecycle for an e-commerce platform",
	//	    Language:    "Go",
	//	    Entities:    []string{"Order", "OrderItem", "Address"},
	//	    Operations:  []string{"Place order", "Cancel order", "Track shipment"},
	//	    Integrations: []string{"Inventory Service (Kafka)", "Payment Service (REST)", "PostgreSQL"},
	//	}
	//
	// ---------------------------------------------------------------
	svc := config.InventoryService()

	outputDir := "./generated"
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}

	ctx := context.Background()

	fmt.Printf("🚀 Microservice Agent Pipeline\n")
	fmt.Printf("   Service:  %s\n", svc.Name)
	fmt.Printf("   Language: %s\n", svc.Language)
	fmt.Printf("   Output:   %s\n\n", outputDir)

	pipeline := orchestrator.NewPipeline(cfg)

	result, err := pipeline.Run(ctx, svc)
	if err != nil {
		log.Fatalf("Pipeline failed: %v", err)
	}

	fmt.Printf("✅ Pipeline completed in %s\n", result.Duration.Round(1e9))
	fmt.Printf("💾 Saving artifacts to %s/%s/...\n", outputDir, svc.Name)

	if err := orchestrator.SaveArtifacts(result, outputDir); err != nil {
		log.Fatalf("Failed to save artifacts: %v", err)
	}

	fmt.Printf("\n📁 Generated files:\n")
	for _, r := range result.Results {
		fmt.Printf("  %-30s %d artifact(s)\n", r.AgentName, len(r.Artifacts))
		for _, a := range r.Artifacts {
			if a.Filename != "" {
				fmt.Printf("    └─ %s\n", a.Filename)
			}
		}
	}

	fmt.Printf("\n✨ Done! See %s/%s/README.md for a summary.\n", outputDir, svc.Name)

}
