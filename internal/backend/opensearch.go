package backend

import (
	"context"
	"internal/item"
	"pkg/conf"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type OpenSearchClient struct {
	ctx       context.Context
	config    *conf.OpenSearchConf
	itemPipe  chan item.Item
	upstream  *opensearchapi.Client
	moviesIdx string
	musicIdx  string
	seriesIdx string
}

type OpenSearchIndexSettings struct {
	settings struct {
		index struct {
			number_of_shards   int
			number_of_replicas int
		}
	}
}

func MakeOpenSearchClient(cfg *conf.OpenSearchConf, pipe chan item.Item) *OpenSearchClient {
	newClient := &OpenSearchClient{
		ctx:      context.Background(),
		config:   cfg,
		itemPipe: pipe,

		moviesIdx: cfg.Settings.IndexPrefix + cfg.Settings.MoviesIndex,
		musicIdx:  cfg.Settings.IndexPrefix + cfg.Settings.MusicIndex,
		seriesIdx: cfg.Settings.IndexPrefix + cfg.Settings.SeriesIndex,
	}

	return newClient
}

func (c *OpenSearchClient) Connect() error {
	// connect to cluster
	err := c.connectCluster()
	if err != nil {
		return err
	}

	// ensure indices exist
	settings := c.indexSettings()

	err = c.assertIndex(c.moviesIdx, settings)
	if err != nil {
		return err
	}

	err = c.assertIndex(c.musicIdx, settings)
	if err != nil {
		return err
	}

	err = c.assertIndex(c.seriesIdx, settings)
	if err != nil {
		return err
	}

	return nil
}

func (c *OpenSearchClient) indexSettings() *OpenSearchIndexSettings {
	osis := &OpenSearchIndexSettings{}

	osis.settings.index.number_of_replicas = c.config.Settings.ReplicaCount
	if c.config.Settings.ShardCount > 0 {
		osis.settings.index.number_of_shards = c.config.Settings.ShardCount
	}

	return osis
}

func (c *OpenSearchClient) connectCluster() error {
	return nil
}

func (c *OpenSearchClient) assertIndex(name string, settings *OpenSearchIndexSettings) error {
	return nil
}

func (c *OpenSearchClient) Index() error {
	return nil
}

func (c *OpenSearchClient) LookupItem(id string) (item.Item, error) {
	return nil, nil
}
