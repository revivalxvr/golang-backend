package main

import (
	"crypto/tls"
	"user-service/env"
	"user-service/src/factory"
	"user-service/src/server"

	"user-service/database"
	"user-service/src/common/log"
	"user-service/src/common/utils"
	"user-service/src/publisher"

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
	fmt.Println("Starting user-service application...")

	ctx, cancel := context.WithCancel(context.Background())
	handleCloseApp(cancel)

	defer log.Logger.Sync()

	godotenv.Load()
	env.Load()

	caCert, err := os.ReadFile("../certs/ca.crt")
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	serverCert, err := tls.LoadX509KeyPair("../certs/server.crt", "../certs/server.key")
	if err != nil {
		log.Logger.Fatal("Gagal memuat sertifikat: " + err.Error())
	}
	serverTlsConf := utils.CreateServerTlsConf(caCert, serverCert)

	db := database.NewPostgresql()
	defer database.ClosePostgresql(db)

	rc := database.NewRedis()
	defer database.CloseRedis(rc)

	pr := publisher.NewRabbitMQ()
	defer pr.Close()

	kf := publisher.NewKafkaPublisher("restful-log")
	defer kf.Close()

	gl := publisher.NewKafkaPublisher("grpc-log")
	defer gl.Close()

	//panggil function NewFactory yg ada di factory
	f := factory.NewFactory(db, rc, pr, kf, gl)

	//panggil function NewRestful yg ada di server
	restfulServer := server.NewRestful(f)
	defer restfulServer.Stop()

	go restfulServer.Start()

	grpcServer := server.NewGrpcServer(f, serverTlsConf)
	defer grpcServer.Stop()

	go grpcServer.Run()

	<-ctx.Done()
}
