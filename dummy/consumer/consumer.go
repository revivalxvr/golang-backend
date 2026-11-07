package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esutil"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaESConsumer struct {
	client      *kgo.Client
	esClient    *elasticsearch.Client
	bulkIndexer esutil.BulkIndexer
	ctx         context.Context
	stats       *Stats
}

type Stats struct {
	MessagesConsumed int
	BulkRequests     int
	FailedInserts    int
}

func NewKafkaESConsumer(ctx context.Context) *KafkaESConsumer {
	kafkaOpts := []kgo.Opt{
		kgo.SeedBrokers("localhost:9093", "localhost:9095", "localhost:9097"),
		kgo.ConsumerGroup("dummy-log-consumer"),
		kgo.ConsumeTopics("dummy-log"),
		kgo.DisableAutoCommit(),
	}

	client, err := kgo.NewClient(kafkaOpts...)
	if err != nil {
		panic(err.Error())
	}

	esCfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200", "http://localhost:9201", "http://localhost:9202"},
	}
	esClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		panic(err.Error())
	}

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        esClient,
		Index:         "dummy-log",
		NumWorkers:    4,
		FlushBytes:    5000 * 1024,
		FlushInterval: 30 * time.Second,
		OnError: func(ctx context.Context, err error) {
			log.Printf("Bulk insert error: %v", err)
		},
	})
	if err != nil {
		panic(err.Error())
	}

	return &KafkaESConsumer{
		client:      client,
		esClient:    esClient,
		bulkIndexer: bulkIndexer,
		ctx:         ctx,
		stats:       &Stats{},
	}
}

func (c *KafkaESConsumer) Run() {
	fmt.Println("Starting consumer...")

	go c.logMetrics()

	for {
		select {
		case <-c.ctx.Done():
			return

		default:
			fetches := c.client.PollFetches(context.Background())
			if fetches.IsClientClosed() {
				return
			}

			fetches.EachError(func(topic string, partition int32, err error) {
				log.Printf("Fetch error: topic=%s partition=%d error=%v", topic, partition, err)
			})

			c.processFetches(fetches)
		}
	}
}

func (c *KafkaESConsumer) processFetches(fetches kgo.Fetches) {
	fetches.EachRecord(func(r *kgo.Record) {

		err := c.bulkIndexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action: "create",
				Body:   strings.NewReader(string(r.Value)),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					c.stats.MessagesConsumed++
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					c.stats.FailedInserts++
					log.Printf("Failed to create document: %s", res.Error.Reason)
				},
			},
		)

		if err != nil {
			log.Printf("Failed to add document to bulk creater: %v", err)
			return
		}

		c.stats.BulkRequests++
	})

	if err := c.client.CommitUncommittedOffsets(context.Background()); err != nil {
		log.Printf("Failed to commit offsets: %v", err)
	}
}

func (c *KafkaESConsumer) Close() {
	log.Printf("Received signal, initiating shutdown...")

	if err := c.bulkIndexer.Close(context.Background()); err != nil {
		log.Printf("Error closing bulk indexer: %v", err)
	}

	c.client.Close()
}

func (c *KafkaESConsumer) logMetrics() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			stats := c.bulkIndexer.Stats()
			log.Printf(
				"Metrics: Consumed=%d | BulkRequests=%d | Failed=%d | ESIndexed=%d | ESFailed=%d",
				c.stats.MessagesConsumed,
				c.stats.BulkRequests,
				c.stats.FailedInserts,
				stats.NumIndexed,
				stats.NumFailed,
			)
		}
	}
}

func handleCloseApp(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
	}()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	handleCloseApp(cancel)

	c := NewKafkaESConsumer(ctx)
	defer c.Close()

	go c.Run()

	<-ctx.Done()
}
