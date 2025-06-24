package cli

import (
	"internal/mlog"
	"pkg/conf"

	flag "github.com/spf13/pflag"
)

func ParseConf() (*conf.Conf, error) {
	// path to config file
	var confFile string
	flag.StringVarP(&confFile, "config", "c", "", "Path to config file")

	// verbose/debug mode (same thing)
	var verbose bool
	flag.BoolVarP(&verbose, "verbose", "v", false, "Display debugging output")

	// action type boolean flags
	var dryrun, force, vaccuum bool
	flag.BoolVarP(&force, "force", "F", false, "Regenerate all metadata")
	flag.BoolVarP(&dryrun, "dryrun", "D", false, "Scan files without writing to a backend")
	flag.BoolVarP(&vaccuum, "vaccuum", "V", false, "Remove backend entries for non-existent files")

	// list of additional search paths
	var episodePaths = make([]string, 0)
	flag.StringArrayVar(&episodePaths, "series", []string{}, "Paths to search for series (accepts multiple paths; may be given multiple times)")
	var featurePaths = make([]string, 0)
	flag.StringArrayVar(&featurePaths, "movies", []string{}, "Paths to search for movies (accepts multiple paths; may be given multiple times)")
	var musicPaths = make([]string, 0)
	flag.StringArrayVar(&musicPaths, "music", []string{}, "Paths to search for music (accepts multiple paths; may be given multiple times)")

	// get values
	flag.Parse()

	config, err := readConfigFile(confFile)
	if err != nil {
		return nil, err
	}
	config.AddDefaults()

	// set verbosity
	if verbose {
		mlog.Verbose()
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
	for _, path := range episodePaths {
		config.Paths.Episodes = append(config.Paths.Episodes, path)
	}

	for _, path := range featurePaths {
		config.Paths.Features = append(config.Paths.Features, path)
	}

	for _, path := range musicPaths {
		config.Paths.Music = append(config.Paths.Music, path)
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}
