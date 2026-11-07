package consumer

import (
	"context"
	"encoding/json"
	"log-service/env"
	"log-service/src/client"
	"log-service/src/common/log"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type Kafka struct {
	kafka  *kgo.Client
	config *env.Consumer
	ctx    context.Context
	es     *client.ElasticSearch
}

func NewKafkaConsumer(ctx context.Context, config *env.Consumer) *Kafka {
	kafka, err := kgo.NewClient(
		kgo.SeedBrokers(config.KafkaBrokers...),
		kgo.ConsumerGroup(config.KafkaGroupID),
		kgo.ConsumeTopics(config.KafkaTopic),
		kgo.DisableAutoCommit(),
	)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	adm := kadm.NewClient(kafka)
	ensureTopicExists(adm, config.KafkaTopic)
	ensureTopicExists(adm, config.KafkaDLQTopic)

	es := client.NewElasticSearch(config.ElasticsearchUrl, config.ElasticsearchIndex)

	return &Kafka{
		kafka:  kafka,
		config: config,
		ctx:    ctx,
		es:     es,
	}
}

func ensureTopicExists(adm *kadm.Client, topic string) {
	resp, err := adm.CreateTopics(context.Background(), 3, 2, nil, topic)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	for _, r := range resp {
		if r.Err != nil {
			if r.Err == kerr.TopicAlreadyExists {
				log.Logger.Info("Topic already exists", zap.String("topic", topic))
				continue
			}
			log.Logger.Fatal(r.Err.Error())
		}

		log.Logger.Info("Topic created", zap.String("topic", topic))
	}
}

func (k *Kafka) RunConsume() {
	log.Logger.Info("Starting kafka consumer")
	for {
		select {
		case <-k.ctx.Done():
			log.Logger.Info("stopping kafka consumer")
			return
		default:
			fetches := k.kafka.PollFetches(k.ctx)
			fetches.EachError(func(s string, i int32, err error) {
				log.Logger.Error("error fetching", zap.String("topic", s), zap.Int32("partition", i), zap.Error(err))
			})
			fetches.EachRecord(func(r *kgo.Record) {
				// Safety check jika es bernilai nil
				if k.es == nil {
					log.Logger.Error("Elasticsearch client is nil")
					return
				}

				if err := k.es.Index(r.Value); err != nil {
					k.sendToDLQ(r, err)
					return
				}
			})
			if err := k.kafka.CommitUncommittedOffsets(context.Background()); err != nil {
				log.Logger.Error("error committing offsets", zap.Error(err))
			}
		}
	}
}

func (k *Kafka) sendToDLQ(r *kgo.Record, err error) {
	jsonData, _ := json.Marshal(map[string]any{
		"topic":     r.Topic,
		"value":     string(r.Value),
		"error":     err.Error(),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"partition": r.Partition,
		"offset":    r.Offset,
	})
	k.kafka.Produce(context.Background(), &kgo.Record{
		Topic: k.config.KafkaDLQTopic,
		Value: jsonData,
	}, func(r *kgo.Record, err error) {
		if err != nil {
			log.Logger.Error(err.Error())
		}
	})
	log.Logger.Warn("sent to DQL", zap.String("topic", k.config.KafkaDLQTopic))
}
func (k *Kafka) Close() {
	k.kafka.Close()
	k.es.Close()
}
