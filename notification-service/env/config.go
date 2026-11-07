package env

import (
	"os"
)

type rabbitmq struct {
	DSN string
}

type currentApp struct {
	RestfulAddr string
}

type midtrans struct {
	ServerKey string
}
type Config struct {
	RabbitMQ   *rabbitmq
	CurrentApp *currentApp
	Midtrans   *midtrans
}

var Conf *Config

func Load() {
	rabbitmqConf := new(rabbitmq)
	rabbitmqConf.DSN = os.Getenv("RABBITMQ_DSN")
	currentAppConf := new(currentApp)
	currentAppConf.RestfulAddr = os.Getenv("RESTFUL_ADDR")

	midtransConf := new(midtrans)
	midtransConf.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")

	Conf = &Config{
		RabbitMQ:   rabbitmqConf,
		CurrentApp: currentAppConf,
		Midtrans:   midtransConf,
	}
}
