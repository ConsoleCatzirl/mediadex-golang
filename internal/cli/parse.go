package cli

import (
	"flag"
	"pkg/conf"
)

func ParseConf() (*conf.Conf, error) {
	// path to config file
	var confFile string
	flag.StringVar(&confFile, "config", "", "Path to config file")

	// action type boolean flags
	var dryrun, force, vaccuum bool
	flag.BoolVar(&force, "force", false, "Regenerate all metadata")
	flag.BoolVar(&dryrun, "dryrun", false, "Scan files without writing to a backend")
	flag.BoolVar(&vaccuum, "vaccuum", false, "Remove backend entries for non-existent files")

	// list of additional search paths
	var featurePaths arrayFlags
	flag.Var(&featurePaths, "movies", "Paths to search for movies (may be given multiple times)")
	var musicPaths arrayFlags
	flag.Var(&musicPaths, "music", "Paths to search for music (may be given multiple times)")
	var episodePaths arrayFlags
	flag.Var(&episodePaths, "series", "Paths to search for series (may be given multiple times)")

	// get values
	flag.Parse()

	config, err := readConfigFile(confFile)
	if err != nil {
		return nil, err
	}

	err = config.AddDefaults()
	if err != nil {
		return nil, err
	}

	// set config action toggles
	if dryrun {
		config.Actions.DryRun = true
	}
	if force {
		config.Actions.Force = true
	}
	if vaccuum {
		config.Actions.Vaccuum = true
	}

	// append search paths
	for _, path := range featurePaths {
		config.Paths.Features = append(config.Paths.Features, path)
	}

	for _, path := range musicPaths {
		config.Paths.Music = append(config.Paths.Music, path)
	}

	for _, path := range episodePaths {
		config.Paths.Episodes = append(config.Paths.Episodes, path)
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}
