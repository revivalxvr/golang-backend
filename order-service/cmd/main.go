package main

import (
	"crypto/tls"
	"order-service/env"
	"order-service/src/broker/consumer"
	"order-service/src/broker/publisher"
	"order-service/src/factory"
	"order-service/src/grpc/client"
	"order-service/src/grpc/interceptor"
	"order-service/src/server"

	"order-service/database"
	"order-service/src/common/log"
	"order-service/src/common/utils"

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
	fmt.Println("Starting order-service application...")

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

	//daftarkan grpc client
	i := interceptor.NewUnary()

	productClient := client.NewProduct(i, clientTlsConf)
	defer productClient.Close()

	uc := client.NewUser(i, clientTlsConf)
	defer uc.Close()

	cartClient := client.NewCart(i, clientTlsConf)
	defer cartClient.Close()
	//end of grpc client

	kf := publisher.NewKafkaPublisher("restful-log")
	defer kf.Close()

	pr := publisher.NewRabbitMQ()
	defer pr.Close()

	f := factory.NewFactory(db, productClient, kf, uc, pr, cartClient)

	restfulServer := server.NewRestful(f)
	defer restfulServer.Stop()

	go restfulServer.Start()

	rabbitMQConsumer := consumer.NewRabbitMQ(ctx, f)
	defer rabbitMQConsumer.Close()
	go rabbitMQConsumer.MidtransTx()

	<-ctx.Done()
}
