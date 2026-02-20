package main

import (
“context”
“fmt”
“log”
“os”

```
"inventory-agents/config"
"inventory-agents/orchestrator"
```

)

const defaultTask = `
Build a production-ready inventory microservice for an e-commerce platform with:

- Multi-warehouse support (products stored across multiple locations)
- Real-time stock tracking with atomic updates
- Stock reservation system for in-progress orders
- Low-stock alerting (configurable thresholds per product)
- Full audit trail of all stock movements
- Support for product variants (size, color, etc.)
- Bulk import/export of inventory data
- Integration with order service via events
  `

func main() {
cfg := config.Load()

```
if cfg.AnthropicAPIKey == "" {
	log.Fatal("ANTHROPIC_API_KEY environment variable is required")
}

task := defaultTask
if len(os.Args) > 1 {
	task = os.Args[1]
}

outputDir := "./generated"
if len(os.Args) > 2 {
	outputDir = os.Args[2]
}

ctx := context.Background()

fmt.Println("Starting Inventory Microservice Agent Pipeline...")
fmt.Printf("Task: %s\n", task)
fmt.Printf("Output directory: %s\n\n", outputDir)

pipeline := orchestrator.NewInventoryPipeline(cfg)

result, err := pipeline.Run(ctx, task)
if err != nil {
	log.Fatalf("Pipeline failed: %v", err)
}

fmt.Printf("\n✅ Pipeline completed in %s\n", result.Duration.Round(1e9))
fmt.Printf("Saving artifacts to %s...\n", outputDir)

if err := orchestrator.SaveArtifacts(result, outputDir); err != nil {
	log.Fatalf("Failed to save artifacts: %v", err)
}

fmt.Printf("\n📁 Generated files:\n")
for _, r := range result.Results {
	fmt.Printf("  %s/ (%d artifacts)\n", r.AgentName, len(r.Artifacts))
	for _, a := range r.Artifacts {
		if a.Filename != "" {
			fmt.Printf("    └─ %s\n", a.Filename)
		}
	}
}

fmt.Printf("\n✨ Done! See %s/README.md for a full summary.\n", outputDir)
```

}