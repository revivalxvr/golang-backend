package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

type DummyLog struct {
	RequestID    string  `json:"request_id"`
	Timestamp    string  `json:"@timestamp"`
	Method       string  `json:"method"`
	URL          string  `json:"url"`
	Body         any     `json:"body"`
	StatusCode   int     `json:"status_code"`
	ClientIP     string  `json:"client_ip"`
	UserAgent    string  `json:"user_agent"`
	Latency      float64 `json:"latency"`
	ErrorMessage string  `json:"error_message,omitempty"`
}

var methods = []string{"GET", "POST", "PUT", "DELETE"}
var urls = []string{"/api/v1/auth/login", "/api/v1/auth/register", "/api/v1/product", "/api/v1/order"}
var agents = []string{"PostmanRuntime/7.43.0", "Mozilla/5.0", "Go-http-client/1.1"}
var ips = []string{"127.0.0.1", "192.168.1.10", "10.0.0.5"}

func randomLog() DummyLog {
	return DummyLog{
		RequestID:  uuid.New().String(),
		Timestamp:  time.Now().Format(time.RFC3339),
		Method:     methods[rand.Intn(len(methods))],
		URL:        urls[rand.Intn(len(urls))],
		StatusCode: []int{200, 400, 401, 404, 500}[rand.Intn(5)],
		ClientIP:   ips[rand.Intn(len(ips))],
		UserAgent:  agents[rand.Intn(len(agents))],
		Latency:    float64(rand.Intn(1000)), // ms
	}
}

func worker(client *kgo.Client) {
	for {
		logData := randomLog()
		value, _ := json.Marshal(logData)

		record := &kgo.Record{Value: value}
		client.Produce(context.Background(), record, nil)

		// jeda random pakai distribusi exponential (mean ~10s)
		delay := rand.ExpFloat64() * 10000 // ms
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}

func SendRealLogTraffic() {
	// 1. Inisialisasi net.Dialer khusus IPv4
	netDialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	client, err := kgo.NewClient(
		// Ubah ke IP IPv4 127.0.0.1
		kgo.SeedBrokers("127.0.0.1:9093", "127.0.0.1:9095", "127.0.0.1:9097"),

		// Gunakan kgo.Dialer dengan memanggil netDialer.DialContext pakai "tcp4"
		kgo.Dialer(func(ctx context.Context, network, addr string) (net.Conn, error) {
			return netDialer.DialContext(ctx, "tcp4", addr)
		}),

		kgo.DefaultProduceTopic("dummy-log"),
		kgo.ProducerBatchCompression(kgo.SnappyCompression()),
		kgo.ProducerBatchMaxBytes(1000*1000),
		kgo.ProducerLinger(5*time.Second),
	)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	admin := kadm.NewClient(client)
	if err := ensureTopicExists(admin, "dummy-log"); err != nil {
		panic(err)
	}

	numWorkers := 10
	fmt.Printf("Starting %d simulated users...\n", numWorkers)

	for i := 0; i < numWorkers; i++ {
		go worker(client)
	}

	time.Sleep(2 * time.Minute)
}

func ensureTopicExists(admin *kadm.Client, topic string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := admin.CreateTopics(ctx, 3, 2, nil, topic)
	if err != nil {
		return err
	}

	for _, r := range resp {
		if r.Err != nil {
			if r.Err == kerr.TopicAlreadyExists {
				fmt.Printf("topic %s already exists\n", r.Topic)
				continue
			}
			return err
		}
		fmt.Printf("topic %s created\n", r.Topic)
	}

	return nil
}
