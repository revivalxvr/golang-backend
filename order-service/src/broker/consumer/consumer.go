package consumer

import (
	"context"
	"time"

	"order-service/env"
	"order-service/src/broker/service"
	"order-service/src/common/log"
	"order-service/src/factory"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn            *amqp091.Connection
	midtransTxNotif *amqp091.Channel
	ctx             context.Context
	service         *service.RabbitMQ
}

// di panggil di main.go
func NewRabbitMQ(ctx context.Context, f *factory.Factory) *RabbitMQ {
	conn := connect()
	chann := SetupChannel(conn, "notification", "midtrans-tx", "midtrans-tx")
	return &RabbitMQ{
		conn:            conn,
		midtransTxNotif: chann,
		ctx:             ctx,
		service:         service.NewRabbitMQ(f),
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
	if err := c.midtransTxNotif.Close(); err != nil {
		log.Logger.Fatal(err.Error())
	}
}
func (c *RabbitMQ) Reconnect() {
	if c.conn != nil && !c.conn.IsClosed() {
		return
	}
	c.conn = connect()
	c.midtransTxNotif = SetupChannel(c.conn, "notification", "midtrans-tx", "midtrans-tx")

	go c.MidtransTx()

	log.Logger.Info("RabbitMQ connected successfully")
}

func SetupChannel(conn *amqp091.Connection, exchangeName, queueName, keyName string) *amqp091.Channel {
	chann, err := conn.Channel()
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	if err := chann.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
		log.Logger.Fatal(err.Error())
	}
	queue, err := chann.QueueDeclare(queueName, true, false, false, false, amqp091.Table{
		"x-queue-type": "quorum", // <-- Konfigurasi Quorum Queue
	})
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
