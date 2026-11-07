package consumer

import (
	"bag-service/src/common/log"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

const MaxRetryCount = 3

func (c *RabbitMQ) OrderCreated() {
	log.Logger.Info("starting order created notification consumer")

	// 1. Konsisten menggunakan nama queue "order.created" sesuai dengan setupChannel
	queueName := "order.created"
	consumerTag := "order.created-consumer"

	orderCreatedConsumer, err := c.cart.ConsumeWithContext(
		c.ctx,
		queueName,
		consumerTag,
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Logger.Error("failed to start consuming order.created queue: " + err.Error())
		return
	}

	for {
		select {
		case msg, ok := <-orderCreatedConsumer:
			if !ok {
				log.Logger.Warn("order.created channel closed, stopping consumer loop")
				return
			}
			if c.conn == nil || c.conn.IsClosed() {
				log.Logger.Warn("RabbitMQ connection closed, reconnecting...")
				return
			}

			if err := c.service.RemoveCheckedOutItems(c.ctx, msg.Body); err != nil {
				log.Logger.Error("failed to process order.created message: " + err.Error())

				retryCount := c.getRetryCount(msg.Headers)
				if retryCount < MaxRetryCount {
					log.Logger.Warn(fmt.Sprintf("retrying midtrans tx notification (%d/%d)...", retryCount+1, MaxRetryCount))

					// Requeue = true agar pesan kembali ke antrean untuk dicoba lagi
					if err := msg.Nack(false, true); err != nil {
						log.Logger.Error("failed to nack message (requeue): " + err.Error())
					}
				} else {
					log.Logger.Error(fmt.Sprintf("max retry limit (%d) reached. Dropping message.", MaxRetryCount))

					// Requeue = false agar pesan dibuang / dikirim ke Dead Letter Exchange (DLX)
					if err := msg.Nack(false, false); err != nil {
						log.Logger.Error("failed to nack message (drop): " + err.Error())
					}
				}
				continue
			}
			// 3. Jika berhasil, kirim ACK
			if err := msg.Ack(false); err != nil {
				log.Logger.Error("failed to ack message: " + err.Error())
				continue
			}

			log.Logger.Info("successfully processed consumer")

		case <-c.ctx.Done():
			log.Logger.Info("context cancelled, exiting consumer loop")
			return
		}
	}
}

func (c *RabbitMQ) getRetryCount(headers amqp091.Table) int {
	if headers == nil {
		return 0
	}

	if xDeath, ok := headers["x-death"].([]any); ok && len(xDeath) > 0 {
		if deathMap, ok := xDeath[0].(amqp091.Table); ok {
			if count, ok := deathMap["count"].(int64); ok {
				return int(count)
			}
		}
	}
	return 0
}
