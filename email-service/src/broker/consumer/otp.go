package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	"email-service/src/common/dto/request"
	"email-service/src/common/log"

	"github.com/rabbitmq/amqp091-go"
)

const (
	MaxRetryCount = 3
	ServiceName   = "otp-consumer"
	QueueName     = "otp"
)

func (c *RabbitMQ) Otp() {
	log.Logger.Info("starting otp consumer")

	closeChann := c.conn.NotifyClose(make(chan *amqp091.Error, 1))
	go c.checkCloseChannel(closeChann)

	otpConsumer, err := c.otpChann.ConsumeWithContext(
		c.ctx,
		QueueName,
		ServiceName,
		false, // autoAck = false
		false, false, false, nil,
	)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	for {
		select {
		case msg, ok := <-otpConsumer:
			if !ok {
				return
			}
			c.handleOtpMessage(msg)

		case <-c.ctx.Done():
			return
		}
	}
}

// handleOtpMessage memproses satu pesan OTP dengan alur retry & logging yang bersih
func (c *RabbitMQ) handleOtpMessage(msg amqp091.Delivery) {
	startTime := time.Now()
	var logBody request.SendOtp

	// 1. Validasi Format JSON
	if err := json.Unmarshal(msg.Body, &logBody); err != nil {
		log.Logger.Error("Gagal unmarshal JSON: " + err.Error())
		c.publishLog(logBody, msg, startTime, ServiceName, QueueName, err)
		_ = msg.Nack(false, false) // JSON rusak langsung dibuang (tanpa retry)
		return
	}

	maskedBody := logBody
	maskedBody.Otp = "******" // Masking OTP untuk privacy log

	// 2. Cek Jumlah Retry
	retryCount := c.getRetryCount(msg.Headers)
	if retryCount >= MaxRetryCount {
		log.Logger.Error(fmt.Sprintf("Batas maksimum retry (%d) tercapai untuk email %s. Pesan dihentikan.", MaxRetryCount, logBody.Email))

		errMaxRetry := fmt.Errorf("exceeded max retry limit (%d)", MaxRetryCount)
		c.publishLog(maskedBody, msg, startTime, ServiceName, QueueName, errMaxRetry)

		// Buang dari antrean agar tidak terjadi infinite retry loop
		_ = msg.Ack(false)
		return
	}

	// 3. Proses Kirim OTP via Service
	if err := c.service.SendOtp(msg.Body); err != nil {
		log.Logger.Warn(fmt.Sprintf("Gagal kirim OTP ke %s (Percobaan %d/%d): %v. Mengirim ke DLQ untuk ditunda 5 detik...",
			logBody.Email, retryCount+1, MaxRetryCount, err))

		c.publishLog(maskedBody, msg, startTime, ServiceName, QueueName, err)

		// Nack (requeue=false) -> Memindahkan pesan ke 'otp-dlq' untuk tunda 5s sebelum dikembalikan
		_ = msg.Nack(false, false)
		return
	}

	// Jika Berhasil (Sukses)
	if err := msg.Ack(false); err != nil {
		log.Logger.Error("Gagal Ack message: " + err.Error())
		c.publishLog(maskedBody, msg, startTime, ServiceName, QueueName, err)
		return
	}

	log.Logger.Info(fmt.Sprintf("OTP sent successfully to %s", logBody.Email))

	// Log sukses ke Kafka (err = nil)
	c.publishLog(maskedBody, msg, startTime, ServiceName, QueueName, nil)
}

// Helper membaca header RabbitMQ untuk mengukur jumlah percobaan dari x-death header
func (c *RabbitMQ) getRetryCount(headers amqp091.Table) int {
	if headers == nil {
		return 0
	}

	xDeath, exists := headers["x-death"]
	if !exists {
		return 0
	}

	switch v := xDeath.(type) {
	case []interface{}:
		if len(v) > 0 {
			if deathMap, ok := v[0].(amqp091.Table); ok {
				if count, ok := deathMap["count"].(int64); ok {
					return int(count)
				}
			}
		}
	case amqp091.Table:
		if count, ok := v["count"].(int64); ok {
			return int(count)
		}
	}

	return 0
}
