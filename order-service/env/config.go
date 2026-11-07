package env

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"order-service/src/common/log"
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
type midtrans struct {
	ApiHostUrl string
	ServerKey  string
}
type Config struct {
	RabbitMQ       *rabbitmq
	CurrentApp     *currentApp
	ProductService *productService
	UserService    *userService
	Postgres       *postgres
	Midtrans       *midtrans
	Jwt            *jwt
	Kafka          *kafka
}

var Conf *Config

func Load() {
	if os.Getenv("MODE") == "NON-DEV" {
		LoadFromVault()
	}
	rabbitmqConf := new(rabbitmq)
	rabbitmqConf.DSN = os.Getenv("RABBITMQ_DSN")

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

	midtransConf := new(midtrans)
	midtransConf.ApiHostUrl = os.Getenv("MIDTRANS_API_HOST_URL")
	midtransConf.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")

	//auth grpc nya harus sama dengan yg ada di product-service/.env GRPC_BASIC_AUTH
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

	Conf = &Config{
		RabbitMQ:       rabbitmqConf,
		CurrentApp:     currentAppConf,
		Postgres:       postgresConf,
		Jwt:            jwtConf,
		Midtrans:       midtransConf,
		ProductService: productServiceConf,
		UserService:    userServiceConf,
		Kafka:          kafka,
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
