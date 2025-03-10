package conf

// Require a top-level magic key "mediadex" in the config file
// without passing it around internally
type MagicConf struct {
	Mediadex *Conf `yaml:"mediadex"`
}

type Conf struct {
	Actions ActionConf  `yaml:"actions"`
	Backend BackendConf `yaml:"backends"`
	Paths   PathsConf   `yaml:"paths"`
}

type ActionConf struct {
	DryRun  bool `yaml:"dryrun"`  // Don't write to any backends
	Force   bool `yaml:"force"`   // Regenerate all metadata
	Vaccuum bool `yaml:"vaccuum"` // Remove entries for deleted files
}

type PathsConf struct {
	Movies []string `yaml:"movies"`
	Music  []string `yaml:"music"`
	Series []string `yaml:"series"`
}

type BackendConf struct {
	ArangoDB   *ArangoConf     `yaml:"arangodb"`
	OpenSearch *OpenSearchConf `yaml:"opensearch"`
}

type ArangoConf struct {
	Auth struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		Insecure bool   `yaml:"insecure"`
	} `yaml:"auth"`
	Settings struct {
		DbName    string `yaml:"database_name"`
		ColPrefix string `yaml:"collection_prefix"`
		MovieCol  string `yaml:"movie_collection"`
		MusicCol  string `yaml:"music_collection"`
		SeriesCol string `yaml:"series_collection"`
	} `yaml:"settings"`
}

type OpenSearchConf struct {
	Auth struct {
		Hosts    []string `yaml:"hosts"`
		User     string   `yaml:"user"`
		Pass     string   `yaml:"pass"`
		Insecure bool     `yaml:"insecure"`
	} `yaml:"auth"`
	Settings struct {
		IndexPrefix  string `yaml:"index_prefix"`
		MoviesIndex  string `yaml:"movie_index"`
		MusicIndex   string `yaml:"music_index"`
		SeriesIndex  string `yaml:"series_index"`
		ReplicaCount int    `yaml:"replicas"`
		ShardCount   int    `yaml:"shards"`
	} `yaml:"settings"`
}
