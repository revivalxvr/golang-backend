package consumer

import (
	"context"
	"time"

	"email-service/env"
	"email-service/src/broker/publisher"
	"email-service/src/broker/service"
	"email-service/src/common/dto/request"
	"email-service/src/common/log"
	"email-service/src/factory"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn                 *amqp091.Connection
	otpChann             *amqp091.Channel
	ctx                  context.Context
	service              *service.RabbitMQ
	rabbitMQLogPublisher *publisher.Kafka
}

func NewRabbitMQ(ctx context.Context, f *factory.Factory) *RabbitMQ {
	conn := connect()
	chann := SetupChannel(conn, "user", "otp", "otp")
	return &RabbitMQ{
		conn:                 conn,
		otpChann:             chann,
		ctx:                  ctx,
		service:              service.NewRabbitMQ(f),
		rabbitMQLogPublisher: f.RabbitMqPublisher,
	}
}
func (c *RabbitMQ) checkCloseChannel(closeChann chan *amqp091.Error) {
	for err := range closeChann {
		if err != nil {
			c.Reconnect()
		}
	}
}

func (c *RabbitMQ) Close() {
	if err := c.conn.Close(); err != nil {
		log.Logger.Fatal(err.Error())
	}
	if err := c.otpChann.Close(); err != nil {
		log.Logger.Fatal(err.Error())
	}
}
func (c *RabbitMQ) Reconnect() {
	if c.conn != nil && !c.conn.IsClosed() {
		return
	}
	c.conn = connect()
	c.otpChann = SetupChannel(c.conn, "user", "otp", "otp")

	go c.Otp()

	log.Logger.Info("RabbitMQ connected successfully")
}

func (c *RabbitMQ) publishLog(body any, msg amqp091.Delivery, start time.Time, tag, queue string, err error) {
	req := request.RabbitMQLogRequest{
		Timestamp:    time.Now().Format(time.RFC3339),
		Exchange:     msg.Exchange,
		RoutingKey:   msg.RoutingKey,
		Queue:        queue,
		ConsumerTag:  tag,
		Payload:      body,
		ContentType:  msg.ContentType,
		DeliveryMode: msg.DeliveryMode,
		AppID:        msg.AppId,
		Acked:        err == nil,
		Latency:      float64(time.Since(start).Nanoseconds()) / 1e6,
	}
	if reqId, ok := msg.Headers["request_id"].(string); ok {
		req.RequestId = reqId
	}
	if userId, ok := msg.Headers["user_id"].(string); ok {
		req.UserId = &userId
	}
	if err != nil {
		errStr := err.Error()
		req.Error = &errStr
	}
	if err := c.rabbitMQLogPublisher.Publish(req); err != nil {
		log.Logger.Error(err.Error())
	}
}

func SetupChannel(conn *amqp091.Connection, exchangeName, queueName, keyName string) *amqp091.Channel {
	chann, err := conn.Channel()
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	// 1. Deklarasi Exchange Utama & DLX
	dlqExchangeName := exchangeName + "-dlx"

	if err := chann.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
		log.Logger.Fatal(err.Error())
	}
	if err := chann.ExchangeDeclare(dlqExchangeName, "direct", true, false, false, false, nil); err != nil {
		log.Logger.Fatal(err.Error())
	}

	// 2. Deklarasi Queue DLQ (Tempat pesan "beristirahat" sebelum retry)
	dlqQueueName := queueName + "-dlq"
	dlqArgs := amqp091.Table{
		"x-queue-type":              "quorum",
		"x-message-ttl":             int32(5000),  // Pesan diam di DLQ selama 5 detik
		"x-dead-letter-exchange":    exchangeName, // Setelah 5 detik, kembalikan ke Exchange Utama
		"x-dead-letter-routing-key": keyName,      // Routing key kembali ke queue utama ('otp')
	}

	if _, err := chann.QueueDeclare(dlqQueueName, true, false, false, false, dlqArgs); err != nil {
		log.Logger.Fatal(err.Error())
	}
	if err := chann.QueueBind(dlqQueueName, keyName+"-dlq", dlqExchangeName, false, nil); err != nil {
		log.Logger.Fatal(err.Error())
	}

	// 3. Deklarasi Queue Utama
	mainQueueArgs := amqp091.Table{
		"x-queue-type":              "quorum",
		"x-dead-letter-exchange":    dlqExchangeName,  // Jika gagal (Nack), lempar ke DLX
		"x-dead-letter-routing-key": keyName + "-dlq", // Routing ke 'otp-dlq'
	}

	queue, err := chann.QueueDeclare(queueName, true, false, false, false, mainQueueArgs)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	if err := chann.QueueBind(queue.Name, keyName, exchangeName, false, nil); err != nil {
		log.Logger.Fatal(err.Error())
	}

	return chann
}

func connect() *amqp091.Connection {
	for {
		conn, err := amqp091.DialConfig(env.Conf.RabbitMQ.DSN, amqp091.Config{
			Heartbeat: 10 * time.Second,
		})
		if err != nil {
			log.Logger.Error(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		return conn
	}

}
