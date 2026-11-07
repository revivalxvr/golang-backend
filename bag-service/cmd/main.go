package main

import (
	"bag-service/env"
	"bag-service/src/broker/consumer"
	"bag-service/src/broker/publisher"
	"bag-service/src/factory"
	"bag-service/src/grpc/client"
	"bag-service/src/grpc/interceptor"
	"bag-service/src/server"
	"crypto/tls"

	"bag-service/database"
	"bag-service/src/common/log"
	"bag-service/src/common/utils"

	"context"
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
	fmt.Println("Starting bag-service application...")

	ctx, cancel := context.WithCancel(context.Background())
	handleCloseApp(cancel)

	defer log.Logger.Sync()

	godotenv.Load()
	env.Load()

	caCert, err := os.ReadFile("../certs/ca.crt")
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	clientCert, err := tls.LoadX509KeyPair("../certs/client.crt", "../certs/client.key")
	if err != nil {
		log.Logger.Fatal("Gagal memuat sertifikat client: " + err.Error())
	}

	clientTlsConf := utils.CreateClientTlsConf(caCert, clientCert)

	db := database.NewPostgresql()
	defer database.ClosePostgresql(db)

	i := interceptor.NewUnary()
	cartClient := client.NewCart(i, clientTlsConf)
	defer cartClient.Close()

	uc := client.NewUser(i, clientTlsConf)
	defer uc.Close()

	kf := publisher.NewKafkaPublisher("restful-log")
	defer kf.Close()

	f := factory.NewFactory(db, uc, cartClient, kf)

	restfulServer := server.NewRestful(f)
	defer restfulServer.Stop()

	rabbitMQConsumer := consumer.NewRabbitMQ(ctx, f)
	defer rabbitMQConsumer.Close()
	go rabbitMQConsumer.OrderCreated()

	go restfulServer.Start()

	<-ctx.Done()
}
