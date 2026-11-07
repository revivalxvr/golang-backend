package main

import (
	"context"
	"email-service/env"
	"email-service/src/broker/consumer"
	"email-service/src/broker/publisher"
	"email-service/src/common/log"
	"email-service/src/factory"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func handleCloseApp(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
	}()
}

func main() {
	fmt.Println("Starting email-service application...")
	ctx, cancel := context.WithCancel(context.Background())
	handleCloseApp(cancel)

	defer log.Logger.Sync()

	godotenv.Load()
	env.Load()

	pr := publisher.NewKafkaPublisher("rabbitmq-log")
	defer pr.Close()

	f := factory.NewFactory(pr)
	rabbitMqConsumer := consumer.NewRabbitMQ(ctx, f)

	defer rabbitMqConsumer.Close()

	go rabbitMqConsumer.Otp()

	<-ctx.Done()
}
