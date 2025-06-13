package backend

import (
	"context"
	"errors"
	"fmt"
	"log"

	"internal/item"
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
		colNameEpisode: cfg.Settings.ColPrefix + cfg.Settings.EpisodeCol,
		colNameFeature: cfg.Settings.ColPrefix + cfg.Settings.FeatureCol,
		colNameMusic:   cfg.Settings.ColPrefix + cfg.Settings.MusicCol,
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
	log.Printf("Trace: connecting to ArangoDB server")

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
		log.Printf("Error: failed to authenticate with ArangoDB")
		return err
	}

	c.upstream = arangodb.NewClient(c.connection)
	return nil
}

func (c *ArangoClient) connectDB() error {
	log.Printf("Trace: connecting to ArangoDB database")

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
	log.Printf("Trace: connecting to ArangoDB collection %s", name)

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
			log.Printf("Trace: arango pipe is empty and closed")
			break
		}

		err := c.UpsertItem(it)
		if err != nil {
			log.Printf("Error upserting item: %v", err)
		}
	}
}

func (c *ArangoClient) LookupItem(key string) (item.Item, error) {
	episodeDoc, err := c.getDoc(c.colNameEpisode, key)
	if err != nil {
		log.Printf("Error getting episode document: %v", err)
	} else if episodeDoc != nil {
		log.Printf("Trace: found episode document")
		eItem := item.JsonToEpisode(episodeDoc.JsonDocument)
		return eItem, nil
	}

	featureDoc, err := c.getDoc(c.colNameEpisode, key)
	if err != nil {
		log.Printf("Error getting feature document: %v", err)
	} else if featureDoc != nil {
		log.Printf("Trace: found feature document")
		fItem := item.JsonToEpisode(featureDoc.JsonDocument)
		return fItem, nil
	}

	musicDoc, err := c.getDoc(c.colNameEpisode, key)
	if err != nil {
		log.Printf("Error getting music document: %v", err)
	} else if musicDoc != nil {
		log.Printf("Trace: found music document")
		mItem := item.JsonToEpisode(musicDoc.JsonDocument)
		return mItem, nil
	}

	return nil, nil
}

func (c *ArangoClient) getDoc(colName string, key string) (*item.ArangoDocument, error) {
	found := &item.JsonDocument{}

	clxn := c.colNameMap[colName]
	_, err := clxn.ReadDocument(c.ctx, key, &found)
	if err != nil {
		if err.Error() == "document not found" {
			return nil, nil
		} else {
			return nil, err
		}
	}

	aDoc := found.Arango()
	return aDoc, nil
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
		msg := fmt.Sprintf("Error: Unknown item type: %s", t)
		return errors.New(msg)
	}

	jDoc, err := it.JsonDoc()
	if err != nil {
		return err
	}
	aDoc := jDoc.Arango()

	c.upsertDoc(clxn, aDoc)

	return nil
}

func (c *ArangoClient) upsertDoc(clxn string, doc *item.ArangoDocument) error {
	found, _ := c.getDoc(clxn, doc.Key)
	if found == nil {
		_, err := c.colNameMap[clxn].CreateDocument(c.ctx, doc)
		if err != nil {
			return err
		}
		log.Printf("Trace: document created")
	} else {
		_, err := c.colNameMap[clxn].UpdateDocument(c.ctx, doc.Key, doc)
		if err != nil {
			return err
		}
		log.Printf("Trace: document updated")
	}

	return nil
}
