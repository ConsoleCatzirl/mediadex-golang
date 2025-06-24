package backend

import (
	"context"
	"errors"
	"fmt"

	"internal/item"
	"internal/mlog"
	"pkg/conf"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/connection"
)

type ArangoClient struct {
	ctx            context.Context
	config         *conf.ArangoConf
	itemPipe       chan item.Item
	upstream       arangodb.Client
	dbName         string
	connection     connection.Connection
	database       arangodb.Database
	colNameMap     map[string]arangodb.Collection
	colNameEpisode string
	colNameFeature string
	colNameMusic   string
}

func MakeArangoClient(cfg *conf.ArangoConf, pipe chan item.Item) *ArangoClient {
	newClient := &ArangoClient{
		ctx:      context.Background(),
		config:   cfg,
		itemPipe: pipe,

		dbName: cfg.Settings.DbName,

		colNameMap:     make(map[string]arangodb.Collection),
		colNameEpisode: cfg.Settings.ColPrefix + cfg.Settings.EpisodeCol + cfg.Settings.ColSuffix,
		colNameFeature: cfg.Settings.ColPrefix + cfg.Settings.FeatureCol + cfg.Settings.ColSuffix,
		colNameMusic:   cfg.Settings.ColPrefix + cfg.Settings.MusicCol + cfg.Settings.ColSuffix,
	}
	return newClient
}

func (c *ArangoClient) Connect() error {
	err := c.connectClient()
	if err != nil {
		return err
	}

	err = c.connectDB()
	if err != nil {
		return err
	}

	err = c.connectCollection(c.colNameEpisode)
	if err != nil {
		return err
	}

	err = c.connectCollection(c.colNameFeature)
	if err != nil {
		return err
	}

	err = c.connectCollection(c.colNameMusic)
	if err != nil {
		return err
	}

	return nil
}

func (c *ArangoClient) connectClient() error {
	mlog.Trace("backend.ArangoClient.connectClient", "Connecting to ArangoDB server")

	var proto string
	if c.config.Auth.Insecure {
		proto = "http://"
	} else {
		proto = "https://"
	}
	uri := proto + c.config.Auth.Host + ":" + c.config.Auth.Port

	endpoint := connection.NewRoundRobinEndpoints([]string{uri})
	conn_conf := connection.DefaultHTTP2ConfigurationWrapper(endpoint, c.config.Auth.Insecure)
	c.connection = connection.NewHttp2Connection(conn_conf)

	creds := connection.NewBasicAuth(c.config.Auth.User, c.config.Auth.Pass)
	err := c.connection.SetAuthentication(creds)
	if err != nil {
		mlog.Error("Failed to authenticate with ArangoDB", err)
		return err
	}

	c.upstream = arangodb.NewClient(c.connection)
	return nil
}

func (c *ArangoClient) connectDB() error {
	mlog.Trace("backend.ArangoClient.connectDB", "Connecting to ArangoDB database")

	exist, err := c.upstream.DatabaseExists(c.ctx, c.dbName)
	if exist {
		if err != nil {
			return err
		}
		c.database, err = c.upstream.GetDatabase(c.ctx, c.dbName, nil)
		if err != nil {
			return err
		}
	} else {
		c.database, err = c.upstream.CreateDatabase(c.ctx, c.dbName, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ArangoClient) connectCollection(name string) error {
	mlog.Trace("backend.ArangoClient.connectCollection", "Connecting to ArangoDB database", "collection", name)

	exist, err := c.database.CollectionExists(c.ctx, name)
	if exist {
		if err != nil {
			return err
		}
		clxn, err := c.database.GetCollection(c.ctx, name, nil)
		if err != nil {
			return err
		}
		c.colNameMap[name] = clxn
	} else {
		clxn, err := c.database.CreateCollection(c.ctx, name, nil)
		if err != nil {
			return err
		}
		c.colNameMap[name] = clxn
	}

	return nil
}

func (c *ArangoClient) Index() {
	for {
		it, more := <-c.itemPipe
		if !more {
			mlog.Trace("backend.ArangoClient.Index", "Arango pipe is closed")
			break
		}

		err := c.UpsertItem(it)
		if err != nil {
			mlog.Error("Error upserting item", err)
		}
	}
}

func (c *ArangoClient) LookupItem(key string) (item.Item, error) {
	episodeDoc, err := c.getDoc(c.colNameEpisode, key)
	if err != nil {
		mlog.Error(
			"Failure getting document", err,
			"type", item.EpisodeFamily,
			"backend", "ArangoDB",
		)
	} else if episodeDoc != nil {
		path := episodeDoc.FileStats.FullPath
		mlog.Trace(
			"backend.ArangoClient.LookupItem", "Found document",
			"type", item.EpisodeFamily,
			"file", path,
		)
		eItem := item.JsonToEpisode(episodeDoc)
		return eItem, nil
	}

	featureDoc, err := c.getDoc(c.colNameFeature, key)
	if err != nil {
		mlog.Error(
			"Failure getting document", err,
			"type", item.FeatureFamily,
			"backend", "ArangoDB",
		)
	} else if featureDoc != nil {
		path := featureDoc.FileStats.FullPath
		mlog.Trace(
			"backend.ArangoClient.LookupItem", "Found document",
			"type", item.FeatureFamily,
			"file", path,
		)
		fItem := item.JsonToEpisode(featureDoc)
		return fItem, nil
	}

	musicDoc, err := c.getDoc(c.colNameMusic, key)
	if err != nil {
		mlog.Error(
			"Failure getting document", err,
			"type", item.MusicFamily,
			"backend", "ArangoDB",
		)
	} else if musicDoc != nil {
		path := featureDoc.FileStats.FullPath
		mlog.Trace(
			"backend.ArangoClient.LookupItem", "Found document",
			"type", item.MusicFamily,
			"file", path,
		)
		mItem := item.JsonToEpisode(musicDoc)
		return mItem, nil
	}

	return nil, nil
}

func (c *ArangoClient) getDoc(colName string, key string) (*item.JsonDocument, error) {
	found := &item.JsonDocument{}

	clxn := c.colNameMap[colName]
	_, err := clxn.ReadDocument(c.ctx, key, &found)
	if err != nil {
		if err.Error() == "document not found" {
			mlog.Trace(
				"backend.ArangoClient.getDoc", "Item not found",
				"key", key,
			)
			return nil, nil
		} else {
			mlog.Error(
				"Error getting doc", err,
				"backend", "ArangoDB",
			)
			return nil, err
		}
	}

	if found.FileStats == nil {
		mlog.Trace(
			"backend.ArangoClient.getDoc", "Item is empty", "key", key,
		)
		return nil, nil
	}
	mlog.Trace(
		"backend.ArangoClient.getDoc", "Item found", "key", key,
	)
	return found, nil
}

func (c *ArangoClient) UpsertItem(it item.Item) error {
	var clxn string

	switch t := it.(type) {
	case *item.EpisodeItem:
		clxn = c.colNameEpisode
	case *item.FeatureItem:
		clxn = c.colNameFeature
	case *item.MusicItem:
		clxn = c.colNameMusic
	default:
		msg := fmt.Sprintf("Error: arango.UpsertItem: Unknown item type: %s", t)
		return errors.New(msg)
	}

	jDoc, err := it.JsonDoc()
	if err != nil {
		return err
	}
	aDoc := jDoc.Arango()

	err = c.upsertDoc(clxn, aDoc)
	if err != nil {
		return err
	}

	return nil
}

func (c *ArangoClient) upsertDoc(clxn string, doc *item.ArangoDocument) error {
	found, err := c.getDoc(clxn, doc.Key)
	if err != nil {
		mlog.Error("Failure", err, "backend", "ArangoDB")
		return err
	}
	if found == nil {
		_, err := c.colNameMap[clxn].CreateDocument(c.ctx, doc)
		if err != nil {
			mlog.Error(
				"Failure creating document", err, "backend", "ArangoDB",
			)
			return err
		}
		mlog.Trace("backend.ArangoClient.upsertDoc", "Document created")
	} else {
		_, err := c.colNameMap[clxn].UpdateDocument(c.ctx, doc.Key, doc)
		if err != nil {
			mlog.Error(
				"Failure updating document", err, "backend", "ArangoDB",
			)
			return err
		}
		mlog.Trace("backend.ArangoClient.upsertDoc", "Document updated")
	}

	return nil
}
