package main

import (
	"notification-service/env"
	"notification-service/src/factory"
	"notification-service/src/server"

	"notification-service/src/common/log"
	"notification-service/src/publisher"

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
	fmt.Println("Starting notification-service application...")

	ctx, cancel := context.WithCancel(context.Background())
	handleCloseApp(cancel)

	defer log.Logger.Sync()

	godotenv.Load()
	env.Load()

	pr := publisher.NewRabbitMQ()
	defer pr.Close()

	f := factory.NewFactory(pr)

	restfulServer := server.NewRestful(f)
	defer restfulServer.Stop()

	go restfulServer.Start()

	<-ctx.Done()
}
