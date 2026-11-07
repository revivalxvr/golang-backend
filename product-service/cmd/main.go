package main

import (
	"crypto/tls"
	"product-service/env"
	"product-service/src/factory"
	"product-service/src/grpc/client"
	"product-service/src/grpc/interceptor"
	"product-service/src/publisher"
	"product-service/src/server"

	"product-service/database"
	"product-service/src/common/log"
	"product-service/src/common/utils"

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
	fmt.Println("Starting product-service application...")

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
		log.Logger.Fatal("Gagal memuat sertifikat server: " + err.Error())
	}
	clientCert, err := tls.LoadX509KeyPair("../certs/client.crt", "../certs/client.key")
	if err != nil {
		log.Logger.Fatal("Gagal memuat sertifikat client: " + err.Error())
	}
	serverTlsConf := utils.CreateServerTlsConf(caCert, serverCert)
	clientTlsConf := utils.CreateClientTlsConf(caCert, clientCert)

	db := database.NewPostgresql()
	defer database.ClosePostgresql(db)

	gl := publisher.NewKafkaPublisher("grpc-log")
	defer gl.Close()

	rp := publisher.NewKafkaPublisher("restful-log")
	defer rp.Close()
	i := interceptor.NewUnary(gl)

	uc := client.NewUser(i, clientTlsConf)

	f := factory.NewFactory(db, gl, uc, rp)

	restfulServer := server.NewRestful(f)
	defer restfulServer.Stop()

	go restfulServer.Start()

	grpcServer := server.NewGrpcServer(f, i, serverTlsConf)
	defer grpcServer.Stop()

	go grpcServer.Run()

	<-ctx.Done()
}
