package backend

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"internal/item"
	"log"
	"net/http"
	"pkg/conf"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type OpenSearchClient struct {
	ctx        context.Context
	config     *conf.OpenSearchConf
	itemPipe   chan item.Item
	upstream   *opensearchapi.Client
	episodeIdx string
	featureIdx string
	musicIdx   string
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

		episodeIdx: cfg.Settings.IndexPrefix + cfg.Settings.EpisodeIndex,
		featureIdx: cfg.Settings.IndexPrefix + cfg.Settings.FeatureIndex,
		musicIdx:   cfg.Settings.IndexPrefix + cfg.Settings.MusicIndex,
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

	err = c.assureIndex(c.episodeIdx, settings)
	if err != nil {
		return err
	}

	err = c.assureIndex(c.featureIdx, settings)
	if err != nil {
		return err
	}

	err = c.assureIndex(c.musicIdx, settings)
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
	log.Println("Trace: connecting to OpenSearch cluster")

	// create upstream config
	upstreamConfig := opensearch.Config{
		Addresses: c.config.Auth.Hosts,
		Username:  c.config.Auth.User,
		Password:  c.config.Auth.Pass,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.config.Auth.Insecure,
			},
		},
	}

	// connect to cluster
	var err error
	c.upstream, err = opensearchapi.NewClient(
		opensearchapi.Config{Client: upstreamConfig},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *OpenSearchClient) assureIndex(name string, settings *OpenSearchIndexSettings) error {
	// check if index exists
	existReq := opensearchapi.IndicesExistsReq{
		Indices: []string{name},
	}
	exists, err := c.upstream.Client.Do(c.ctx, existReq, nil)
	if err != nil {
		log.Printf("Error checking for existing index: %s", err)
		return err
	}
	defer exists.Body.Close()

	if exists.StatusCode == http.StatusOK {
		log.Printf("Trace: index %s exists", name)
	} else {
		// create index
		jsonBytes, err := json.Marshal(settings)
		if err != nil {
			log.Printf("Error marshalling index settings")
		}
		createReq := opensearchapi.IndicesCreateReq{
			Index: name,
			Body:  bytes.NewReader(jsonBytes),
		}
		create, err := c.upstream.Client.Do(c.ctx, createReq, nil)
		if err != nil {
			log.Printf("Error creating index: %s", err)
			return err
		}
		defer create.Body.Close()

		if create.IsError() {
			log.Printf("Error creating index: %s", create.String())
		}
	}
	return nil
}

func (c *OpenSearchClient) Index() {
	for {
		it, more := <-c.itemPipe
		if !more {
			log.Printf("Trace: opensearch pipe is empty and closed")
			break
		}

		err := c.UpsertItem(it)
		if err != nil {
			log.Printf("Error upserting item: %v", err)
		}
	}
}

func (c *OpenSearchClient) LookupItem(id string) (item.Item, error) {
	episodeDoc, err := c.getDoc(c.episodeIdx, id)
	if err != nil {
		return nil, err
	}
	if episodeDoc != nil {
		eItem := item.JsonToEpisode(episodeDoc.JsonDocument)
		return eItem, nil
	}

	featureDoc, err := c.getDoc(c.featureIdx, id)
	if err != nil {
		return nil, err
	}
	if featureDoc != nil {
		eItem := item.JsonToFeature(featureDoc.JsonDocument)
		return eItem, nil
	}

	musicDoc, err := c.getDoc(c.musicIdx, id)
	if err != nil {
		return nil, err
	}
	if musicDoc != nil {
		eItem := item.JsonToMusic(musicDoc.JsonDocument)
		return eItem, nil
	}

	return nil, nil
}

func (c *OpenSearchClient) getDoc(index string, docId string) (*item.OpenSearchDocument, error) {
	found := &item.OpenSearchDocument{}

	req := opensearchapi.DocumentGetReq{
		Index:      index,
		DocumentID: docId,
	}
	resp, err := c.upstream.Client.Do(c.ctx, req, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		var readBytes []byte

		_, err := resp.Body.Read(readBytes)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(readBytes, &found)
		if err != nil {
			return nil, err
		}
	} else if resp.StatusCode != 404 {
		msg := fmt.Sprintf("Error getting OpenSearch document: Status Code: %v", resp.StatusCode)
		log.Printf(msg)
		return nil, errors.New(msg)
	}

	return found, nil
}

func (c *OpenSearchClient) UpsertItem(it item.Item) error {
	var index string
	switch t := it.(type) {
	case *item.EpisodeItem:
		index = c.episodeIdx
	case *item.FeatureItem:
		index = c.featureIdx
	case *item.MusicItem:
		index = c.musicIdx
	default:
		msg := fmt.Sprintf("Unknown item type: %s", t)
		return errors.New(msg)
	}

	jDoc, err := it.JsonDoc()
	if err != nil {
		return err
	}
	osDoc := jDoc.OpenSearch()

	c.upsertDoc(index, osDoc)
	return nil
}

func (c *OpenSearchClient) upsertDoc(idx string, doc *item.OpenSearchDocument) error {
	jsonBytes, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBytes)

	req := opensearchapi.IndexReq{
		Index:      idx,
		DocumentID: doc.ID,
		Body:       bodyReader,
	}
	resp, err := c.upstream.Client.Do(c.ctx, req, nil)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
