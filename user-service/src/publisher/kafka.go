package publisher

import (
	"context"
	"encoding/json"
	"net"
	"time"
	"user-service/env"
	"user-service/src/common/log"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

// jangan lupa tambahkan Kafka di factory
type Kafka struct {
	client *kgo.Client
	topic  string
}

func NewKafkaPublisher(topic string) *Kafka {
	// Inisialisasi net.Dialer khusus IPv4
	netDialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	client, err := kgo.NewClient(
		// Ubah ke IP IPv4 127.0.0.1
		kgo.SeedBrokers(env.Conf.Kafka.Addr1, env.Conf.Kafka.Addr2, env.Conf.Kafka.Addr3),
		kgo.Dialer(func(ctx context.Context, network, addr string) (net.Conn, error) {
			return netDialer.DialContext(ctx, "tcp4", addr)
		}),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
		kgo.ProducerLinger(5*time.Second),
		kgo.RequiredAcks(kgo.LeaderAck()),
		kgo.DisableIdempotentWrite(),
	)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	adm := kadm.NewClient(client)

	ensureTopicExists(adm, topic)

	return &Kafka{
		client: client,
		topic:  topic,
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
func (k *Kafka) Publish(data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	k.client.Produce(context.Background(), &kgo.Record{
		Topic: k.topic,
		Value: jsonData,
	}, func(r *kgo.Record, err error) {
		if err != nil {
			log.Logger.Error(err.Error())
		}
	})

	return nil
}

func (k *Kafka) Close() {
	//pastikan semua pesan terkirim ke broker
	if err := k.client.Flush(context.Background()); err != nil {
		log.Logger.Error(err.Error())
	}

	k.client.Close()
}
