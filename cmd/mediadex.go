package main

import (
	"os"

	"internal/cli"
	"internal/mlog"
	"pkg/worker"
)

func main() {
	// get configuration
	config, err := cli.ParseConf()
	if err != nil {
		mlog.Error("Failure parsing configuration", err)
		os.Exit(1)
	}

	// create worker and start working
	cmdWorker, err := worker.MakeWorker(config)
	if err != nil {
		mlog.Error("Failure creating worker", err)
		os.Exit(2)
	}

	mlog.Info("Search paths",
		"movies", config.Paths.Features,
		"music", config.Paths.Music,
		"series", config.Paths.Episodes,
	)

	mlog.Info("Backend configuration",
		"ArangoDB", config.Backend.ArangoDB != nil,
		"OpenSearch", config.Backend.OpenSearch != nil,
	)

	mlog.Info("Starting")
	err = cmdWorker.Work()
	if err != nil {
		mlog.Error("Worker Failed", err)
		os.Exit(3)
	}
}
