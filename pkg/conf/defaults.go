package conf

func (c *Conf) AddDefaults() error {

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
		if c.Backend.ArangoDB.Settings.MovieCol == "" {
			c.Backend.ArangoDB.Settings.MovieCol = "movies"
		}
		if c.Backend.ArangoDB.Settings.MusicCol == "" {
			c.Backend.ArangoDB.Settings.MusicCol = "music"
		}
		if c.Backend.ArangoDB.Settings.SeriesCol == "" {
			c.Backend.ArangoDB.Settings.SeriesCol = "series"
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
		if c.Backend.OpenSearch.Settings.MoviesIndex == "" {
			c.Backend.OpenSearch.Settings.MoviesIndex = "movies"
		}
		if c.Backend.OpenSearch.Settings.MusicIndex == "" {
			c.Backend.OpenSearch.Settings.MusicIndex = "music"
		}
		if c.Backend.OpenSearch.Settings.SeriesIndex == "" {
			c.Backend.OpenSearch.Settings.SeriesIndex = "series"
		}
	}

	return nil
}
