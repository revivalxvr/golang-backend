package client

import (
	"context"
	"log-service/src/common/log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esutil"
)

type ElasticSearch struct {
	client      *elasticsearch.Client
	bulkIndexer esutil.BulkIndexer
}

func NewElasticSearch(urls []string, index string) *ElasticSearch {
	client, err := elasticsearch.New(
		elasticsearch.WithAddresses(urls...),
	)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         index,
		Client:        client,
		FlushInterval: 1 * time.Second,
	})
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	return &ElasticSearch{
		client:      client,
		bulkIndexer: bulkIndexer,
	}
}
func (e *ElasticSearch) Index(data []byte) error {
	return e.bulkIndexer.Add(context.Background(), esutil.BulkIndexerItem{
		Action: "create",
		Body:   strings.NewReader(string(data)),
	})
}
func (e *ElasticSearch) Close() {
	if err := e.bulkIndexer.Close(context.Background()); err != nil {
		log.Logger.Fatal(err.Error())
	}
}
