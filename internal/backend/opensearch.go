package backend

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"internal/item"
	"internal/mlog"
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

		episodeIdx: cfg.Settings.IndexPrefix + cfg.Settings.EpisodeIndex + cfg.Settings.IndexSuffix,
		featureIdx: cfg.Settings.IndexPrefix + cfg.Settings.FeatureIndex + cfg.Settings.IndexSuffix,
		musicIdx:   cfg.Settings.IndexPrefix + cfg.Settings.MusicIndex + cfg.Settings.IndexSuffix,
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
	mlog.Trace("backend.OpenSearchClient.connectCluster", "Connecting to OpenSearch cluster")

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
		mlog.Error("Error checking for existing index", err, "index", name)
		return err
	}
	defer exists.Body.Close()

	if exists.StatusCode == http.StatusOK {
		mlog.Trace(
			"backend.OpenSearchClient.assureIndex", "Index exists",
			"index", name,
		)
	} else {
		// create index
		jsonBytes, err := json.Marshal(settings)
		if err != nil {
			mlog.Error("Failure marshalling index settings", err)
			return err
		}
		createReq := opensearchapi.IndicesCreateReq{
			Index: name,
			Body:  bytes.NewReader(jsonBytes),
		}
		create, err := c.upstream.Client.Do(c.ctx, createReq, nil)
		if err != nil {
			mlog.Error("Failure creating index", err, "index", name)
			return err
		}
		defer create.Body.Close()

		if create.IsError() {
			err = errors.New(create.String())
			mlog.Error("Failure creating index", err, "index", name)
			return err
		}
		mlog.Trace(
			"backend.OpenSearchClient.assureIndex", "Index created",
			"index", name,
		)
	}
	return nil
}

func (c *OpenSearchClient) Index() {
	for {
		it, more := <-c.itemPipe
		if !more {
			mlog.Trace("backend.OpenSearchClient.Index", "OpenSearch pipe is empty and closed")
			break
		}

		err := c.UpsertItem(it)
		if err != nil {
			mlog.Error("Failure upserting item", err)
		}
	}
}

func (c *OpenSearchClient) LookupItem(id string) (item.Item, error) {
	episodeDoc, err := c.getDoc(c.episodeIdx, id)
	if err != nil {
		return nil, err
	}
	if episodeDoc != nil {
		mlog.Trace(
			"backend.OpenSearchClient.LookupItem", "Found document",
			"type", item.EpisodeFamily, "ID", id,
		)
		eItem := item.JsonToEpisode(episodeDoc)
		return eItem, nil
	}

	featureDoc, err := c.getDoc(c.featureIdx, id)
	if err != nil {
		return nil, err
	}
	if featureDoc != nil {
		mlog.Trace(
			"backend.OpenSearchClient.LookupItem", "Found document",
			"type", item.FeatureFamily, "ID", id,
		)
		fItem := item.JsonToFeature(featureDoc)
		return fItem, nil
	}

	musicDoc, err := c.getDoc(c.musicIdx, id)
	if err != nil {
		return nil, err
	}
	if musicDoc != nil {
		mlog.Trace(
			"backend.OpenSearchClient.LookupItem", "Found document",
			"type", item.MusicFamily, "ID", id,
		)
		mItem := item.JsonToMusic(musicDoc)
		return mItem, nil
	}

	return nil, nil
}

func (c *OpenSearchClient) getDoc(index string, docId string) (*item.JsonDocument, error) {
	var foundBytes []byte
	foundDoc := &item.JsonDocument{}

	mlog.Trace("backend.OpenSearchClient.getDoc", "Getting document",
		"ID", docId, "index", index,
	)

	req := opensearchapi.DocumentGetReq{
		Index:      index,
		DocumentID: docId,
	}

	resp, err := c.upstream.Document.Get(c.ctx, req)
	if !resp.Found {
		return nil, nil
	} else {
		if resp.Source != nil {
			mlog.Trace("backend.OpenSearchClient.getDoc", "Using 'source'",
				"ID", docId, "index", index,
			)
			foundBytes = resp.Source
		} else if resp.Fields != nil {
			mlog.Trace("backend.OpenSearchClient.getDoc", "Using 'fields'",
				"ID", docId, "index", index,
			)
			foundBytes = resp.Fields
		} else {
			new_err := errors.New("document 'source' and 'fields' both empty")
			mlog.Warn(
				"No document data", err,
				"ID", docId, "index", index,
			)

			respAsJson, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				mlog.Error(
					"Failure marshalling document data", err,
					"ID", docId, "index", index,
				)
			}
			mlog.Trace(
				"backend.OpenSearchClient.getDoc", "Found empty document",
				"ID", docId, "index", index,
				"document", respAsJson,
			)

			return nil, new_err
		}
	}

	err = json.Unmarshal(foundBytes, foundDoc)
	if err != nil {
		return nil, err
	}

	return foundDoc, nil
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

	err = c.upsertDoc(index, osDoc)
	if err != nil {
		return err
	}
	return nil
}

func (c *OpenSearchClient) upsertDoc(idx string, doc *item.OpenSearchDocument) error {
	jsonBytes, err := json.Marshal(doc.JsonDocument)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBytes)

	mlog.Trace(
		"backend.OpenSearchClient.upsertDoc", "Upserting document",
		"index", idx, "ID", doc.ID,
	)
	req := opensearchapi.IndexReq{
		Index:      idx,
		DocumentID: doc.ID,
		Body:       bodyReader,
	}
	resp, err := c.upstream.Index(c.ctx, req)
	if err != nil {
		mlog.Error("Failure upserting document", err,
			"index", idx, "ID", doc.ID,
		)
		return err
	}

	mlog.Trace(
		"backend.OpenSearchClient.upsertDoc", "Document upserted",
		"index", idx, "ID", doc.ID, "result", resp.Result,
	)

	return nil
}
