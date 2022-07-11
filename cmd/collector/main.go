package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
	"log"
)

var (
	svcName   = "collector"
	envPrefix = "orb_collector"
	httpPort  = "8204"
	// TO CHANGE TO GET THE VERSION CORRECTLY
	version = " 0.17.0"
)

func OrbComponents() (component.Factories, error) {
	return component.Factories{}, nil
}

func main() {
	mainContext := context.Background()
	factories, err := OrbComponents()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
	}

	info := component.BuildInfo{
		Command:     "collector",
		Description: "Orb Collector",
		Version:     version,
	}

	if params, err := service.CollectorSettings{BuildInfo: info, Factories: factories}); err != nil {
		log.Fatal(err)
	}

	cmd := service.NewCommand(params)
	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("collector server run finished with error: %w", err)
	}

	return nil

}
