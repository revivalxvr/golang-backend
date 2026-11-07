package main

import (
	"context"
	"log-service/env"
	"log-service/src/common/log"
	"log-service/src/consumer"
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
	ctx, cancel := context.WithCancel(context.Background())
	handleCloseApp(cancel)

	defer log.Logger.Sync()

	godotenv.Load()
	env.Load()

	rlc := consumer.NewKafkaConsumer(ctx, env.Conf.RestfulLogConsumer)
	defer rlc.Close()

	go rlc.RunConsume()

	rmqc := consumer.NewKafkaConsumer(ctx, env.Conf.RabbitMqLogConsumer)
	defer rmqc.Close()

	grpc := consumer.NewKafkaConsumer(ctx, env.Conf.GrpcLogConsumer)
	defer grpc.Close()

	go grpc.RunConsume()

	go rmqc.RunConsume()
	<-ctx.Done()
}
