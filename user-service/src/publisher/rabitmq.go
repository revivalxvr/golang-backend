package publisher

import (
	"encoding/json"
	"sync"
	"time"
	"user-service/env"
	"user-service/src/common/dto/request"
	"user-service/src/common/log"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

// jangan lupa tambahkan kafka di factory
type RabitMQ struct {
	conn       *amqp091.Connection
	chann      *amqp091.Channel
	closeChann chan *amqp091.Error
	mutex      sync.Mutex
}

func NewRabbitMQ() *RabitMQ {
	conn, err := amqp091.DialConfig(env.Conf.RabbitMQ.DSN, amqp091.Config{
		Heartbeat: 10 * time.Second,
	})
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	closeChann := conn.NotifyClose(make(chan *amqp091.Error, 1))
	chann := setupChannel(conn)

	return &RabitMQ{
		conn:       conn,
		chann:      chann,
		closeChann: closeChann,
	}
}

func (r *RabitMQ) Reconnect() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.conn != nil && !r.conn.IsClosed() {
		return
	}
	for {
		conn, err := amqp091.DialConfig(env.Conf.RabbitMQ.DSN, amqp091.Config{
			Heartbeat: 10 * time.Second,
		})
		if err != nil {
			log.Logger.Fatal(err.Error())
		}
		r.conn = conn
		r.chann = setupChannel(conn)
		r.closeChann = conn.NotifyClose(make(chan *amqp091.Error, 1))
		log.Logger.Info("RabbitMQ connected")
		break
	}
}
func (r *RabitMQ) Publish(exchange string, key string, data *request.RabbitMqMessage) error {
	if r.conn == nil || r.conn.IsClosed() {
		r.Reconnect()
	}
	jsonData, err := json.Marshal(data.Message)
	if err != nil {
		return err
	}

	msg := amqp091.Publishing{
		Headers: amqp091.Table{
			"request_id": data.RequestID,
			"user_id":    data.UserID,
		},
		AppId:       "user-service",
		ContentType: "application/json",
		Body:        jsonData,
	}

	if err := r.chann.Publish(exchange, key, true, false, msg); err != nil {
		return err
	}

	return nil
}

func (r *RabitMQ) Close() {
	if err := r.conn.Close(); err != nil {
		log.Logger.Fatal(err.Error())
	}
	if err := r.chann.Close(); err != nil {
		log.Logger.Fatal(err.Error())
	}
}

func setupChannel(conn *amqp091.Connection) *amqp091.Channel {
	chann, err := conn.Channel()
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	if err := chann.ExchangeDeclare("user", "direct", true, false, false, false, nil); err != nil {
		log.Logger.Fatal(err.Error())
	}
	return chann
}
