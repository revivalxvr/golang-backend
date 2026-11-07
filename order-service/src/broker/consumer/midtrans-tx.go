package consumer

import (
	"fmt"
	"order-service/src/common/log"

	"github.com/rabbitmq/amqp091-go"
)

const MaxRetryCount = 3

func (c *RabbitMQ) MidtransTx() {
	log.Logger.Info("starting midtrans tx notification consumer")

	// 1. Konsisten menggunakan nama queue "midtrans-tx" sesuai dengan setupChannel
	queueName := "midtrans-tx"
	consumerTag := "midtrans-tx-consumer"

	midtransTxConsumer, err := c.midtransTxNotif.ConsumeWithContext(
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
		log.Logger.Error("failed to start consuming midtrans-tx queue: " + err.Error())
		return
	}

	for {
		select {
		case msg, ok := <-midtransTxConsumer:
			if !ok {
				log.Logger.Warn("midtrans-tx channel closed, stopping consumer loop")
				return
			}

			if c.conn == nil || c.conn.IsClosed() {
				log.Logger.Warn("rabbitmq connection closed, stopping consumer")
				return
			}

			// 2. Olah pesan melalui service
			if err := c.service.HandlerMidtransTxNotif(c.ctx, msg.Body); err != nil {
				log.Logger.Error(fmt.Sprintf("failed to process midtrans tx notification: %v", err))

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

			log.Logger.Info("successfully processed midtrans tx notification")

		case <-c.ctx.Done():
			log.Logger.Info("context cancelled, exiting midtrans tx consumer")
			return
		}
	}
}

// Helper membaca header RabbitMQ untuk batas retry
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
