package backend

import (
	"context"
	"log"

	"internal/item"
	"pkg/conf"

	"github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/connection"
)

type ArangoClient struct {
	ctx           context.Context
	config        *conf.ArangoConf
	itemPipe      chan item.Item
	upstream      arangodb.Client
	dbName        string
	connection    connection.Connection
	database      arangodb.Database
	colNameMap    map[string]arangodb.Collection
	colNameMovies string
	colNameMusic  string
	colNameSeries string
}

func MakeArangoClient(cfg *conf.ArangoConf, pipe chan item.Item) *ArangoClient {
	newClient := &ArangoClient{
		ctx:      context.Background(),
		config:   cfg,
		itemPipe: pipe,

		dbName: cfg.Settings.DbName,

		colNameMap:    make(map[string]arangodb.Collection),
		colNameMovies: cfg.Settings.ColPrefix + cfg.Settings.MovieCol,
		colNameMusic:  cfg.Settings.ColPrefix + cfg.Settings.MusicCol,
		colNameSeries: cfg.Settings.ColPrefix + cfg.Settings.SeriesCol,
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

	err = c.connectCollection(c.colNameMovies)
	if err != nil {
		return err
	}

	err = c.connectCollection(c.colNameMusic)
	if err != nil {
		return err
	}

	err = c.connectCollection(c.colNameSeries)
	if err != nil {
		return err
	}

	return nil
}

func (c *ArangoClient) connectClient() error {
	log.Printf("Trace: connecting to ArangoDB server")
	return nil
}

func (c *ArangoClient) connectDB() error {
	log.Printf("Trace: connecting to ArangoDB database")
	return nil
}

func (c *ArangoClient) connectCollection(name string) error {
	log.Printf("Trace: connecting to ArangoDB collection %s", name)
	return nil
}

func (c *ArangoClient) Index() error {
	return nil
}

func (c *ArangoClient) LookupItem(key string) (item.Item, error) {
	return nil, nil
}
