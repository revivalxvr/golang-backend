package env

import (
	"bag-service/src/common/log"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
)

type rabbitmq struct {
	DSN string
}
type postgres struct {
	DSN string
}
type currentApp struct {
	RestfulAddr   string
	GrpcBasicAuth string
	AllowedOrigin string
}
type productService struct {
	GrpcAddr string
	GrpcAuth string
}
type userService struct {
	GrpcAddr string
	GrpcAuth string
}
type kafka struct {
	Addr1 string
	Addr2 string
	Addr3 string
}

type jwt struct {
	PublicKey *rsa.PublicKey
}

type Config struct {
	CurrentApp     *currentApp
	ProductService *productService
	UserService    *userService
	Postgres       *postgres
	Jwt            *jwt
	Kafka          *kafka
	RabbitMQ       *rabbitmq
}

var Conf *Config

func Load() {
	if os.Getenv("MODE") == "NON-DEV" {
		LoadFromVault()
	}
	currentAppConf := new(currentApp)
	currentAppConf.RestfulAddr = os.Getenv("RESTFUL_ADDR")
	currentAppConf.AllowedOrigin = os.Getenv("ALLOWED_ORIGIN")

	postgresConf := new(postgres)
	postgresConf.DSN = os.Getenv("POSTGRES_DSN")

	base64Str, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_PUBLIC_KEY"))
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	jwtConf := new(jwt)
	jwtConf.PublicKey = loadRsaPublicKey(string(base64Str))

	productServiceConf := new(productService)
	productServiceConf.GrpcAddr = os.Getenv("PRODUCT_SERVICE_GRPC_ADDR")
	productServiceConf.GrpcAuth = os.Getenv("PRODUCT_SERVICE_GRPC_AUTH")

	userServiceConf := new(userService)
	userServiceConf.GrpcAddr = os.Getenv("USER_SERVICE_GRPC_ADDR")
	userServiceConf.GrpcAuth = os.Getenv("USER_SERVICE_GRPC_AUTH")

	kafka := new(kafka)
	kafka.Addr1 = os.Getenv("KAFKA_ADDR_1")
	kafka.Addr2 = os.Getenv("KAFKA_ADDR_2")
	kafka.Addr3 = os.Getenv("KAFKA_ADDR_3")

	rabbitmqConf := new(rabbitmq)
	rabbitmqConf.DSN = os.Getenv("RABBITMQ_DSN")

	Conf = &Config{

		CurrentApp:     currentAppConf,
		Postgres:       postgresConf,
		Jwt:            jwtConf,
		ProductService: productServiceConf,
		UserService:    userServiceConf,
		Kafka:          kafka,
		RabbitMQ:       rabbitmqConf,
	}
}

func loadRsaPublicKey(publicKeyStr string) *rsa.PublicKey {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		log.Logger.Fatal("failed to parse public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	rsaPublicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		log.Logger.Fatal(err.Error())
	}

	return rsaPublicKey
}
