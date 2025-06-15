package main

import (
	"log"

	"internal/cli"
	"pkg/worker"
)

func main() {
	// get configuration
	config, err := cli.ParseConf()
	if err != nil {
		log.Fatalf("Failed to parse configuration: %v", err)
	}

	// create worker and start working
	cmdWorker, err := worker.MakeWorker(config)
	if err != nil {
		log.Fatalf("Failed to create worker: %s", err)
	}

	log.Printf("Movie paths: %v", config.Paths.Features)
	log.Printf("Music paths: %v", config.Paths.Music)
	log.Printf("Series paths: %v", config.Paths.Episodes)

	if config.Backend.ArangoDB != nil {
		log.Printf("Found configuration for ArangoDB")
	}
	if config.Backend.OpenSearch != nil {
		log.Printf("Found configuration for OpenSearch")
	}

	log.Println("Starting")
	err = cmdWorker.Work()
	if err != nil {
		log.Fatalf("Failed: %s", err)
	}
}
