package conf

func (c *Conf) AddDefaults() {

	// action defaults are all false (zero-value)

	// thread defaults are minimum values
	if c.Threads.Workers == 0 {
		c.Threads.Workers = 1
	}

	if c.Threads.FileBuffer == 0 {
		c.Threads.FileBuffer = 1
	}

	if c.Threads.BackendBuffer == 0 {
		c.Threads.BackendBuffer = 1
	}

	// arango client
	if c.Backend.ArangoDB != nil {

		// no defaults for user or pass

		// add host and port
		if c.Backend.ArangoDB.Auth.Host == "" {
			c.Backend.ArangoDB.Auth.Host = "localhost"
		}
		if c.Backend.ArangoDB.Auth.Port == "" {
			c.Backend.ArangoDB.Auth.Port = "8529"
		}

		// add db and collection names
		if c.Backend.ArangoDB.Settings.DbName == "" {
			c.Backend.ArangoDB.Settings.DbName = "mediadex"
		}
		if c.Backend.ArangoDB.Settings.FeatureCol == "" {
			c.Backend.ArangoDB.Settings.FeatureCol = "movies"
		}
		if c.Backend.ArangoDB.Settings.MusicCol == "" {
			c.Backend.ArangoDB.Settings.MusicCol = "music"
		}
		if c.Backend.ArangoDB.Settings.EpisodeCol == "" {
			c.Backend.ArangoDB.Settings.EpisodeCol = "series"
		}
	}

	// opensearch client
	if c.Backend.OpenSearch != nil {

		// no defaults for user or pass

		// add host and port
		if len(c.Backend.OpenSearch.Auth.Hosts) == 0 {
			c.Backend.OpenSearch.Auth.Hosts = []string{"https://localhost:9200"}
		}

		// add index names
		if c.Backend.OpenSearch.Settings.FeatureIndex == "" {
			c.Backend.OpenSearch.Settings.FeatureIndex = "movies"
		}
		if c.Backend.OpenSearch.Settings.MusicIndex == "" {
			c.Backend.OpenSearch.Settings.MusicIndex = "music"
		}
		if c.Backend.OpenSearch.Settings.EpisodeIndex == "" {
			c.Backend.OpenSearch.Settings.EpisodeIndex = "series"
		}
	}
}
