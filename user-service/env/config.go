package env

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"user-service/src/common/log"
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
	GrpcAddr      string
	AllowedOrigin string
}
type redis struct {
	AddrNode1 string
	AddrNode2 string
	AddrNode3 string
	AddrNode4 string
	AddrNode5 string
	AddrNode6 string
	Password  string
}
type kafka struct {
	Addr1 string
	Addr2 string
	Addr3 string
}
type jwt struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

type Config struct {
	RabbitMQ   *rabbitmq
	CurrentApp *currentApp
	Postgres   *postgres
	Redis      *redis
	Jwt        *jwt
	Kafka      *kafka
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
	currentAppConf.GrpcAddr = os.Getenv("GRPC_ADDR")
	currentAppConf.GrpcBasicAuth = os.Getenv("GRPC_BASIC_AUTH")
	currentAppConf.AllowedOrigin = os.Getenv("ALLOWED_ORIGIN") // Membaca langsung dari .env / env variable

	postgresConf := new(postgres)
	postgresConf.DSN = os.Getenv("POSTGRES_DSN")

	redis := new(redis)
	redis.AddrNode1 = os.Getenv("REDIS_ADDR_NODE_1")
	redis.AddrNode2 = os.Getenv("REDIS_ADDR_NODE_2")
	redis.AddrNode3 = os.Getenv("REDIS_ADDR_NODE_3")
	redis.AddrNode4 = os.Getenv("REDIS_ADDR_NODE_4")
	redis.AddrNode5 = os.Getenv("REDIS_ADDR_NODE_5")
	redis.AddrNode6 = os.Getenv("REDIS_ADDR_NODE_6")
	redis.Password = os.Getenv("REDIS_PASSWORD")

	kafka := new(kafka)
	kafka.Addr1 = os.Getenv("KAFKA_ADDR_1")
	kafka.Addr2 = os.Getenv("KAFKA_ADDR_2")
	kafka.Addr3 = os.Getenv("KAFKA_ADDR_3")

	jwtConf := new(jwt)
	base64Str, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_PRIVATE_KEY"))
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	jwtConf.PrivateKey = loadRsaPrivateKey(string(base64Str))

	base64Str, err = base64.StdEncoding.DecodeString(os.Getenv("JWT_PUBLIC_KEY"))
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	jwtConf.PublicKey = loadRsaPublicKey(string(base64Str))

	Conf = &Config{
		RabbitMQ:   rabbitmqConf,
		CurrentApp: currentAppConf,
		Postgres:   postgresConf,
		Redis:      redis,
		Kafka:      kafka,
		Jwt:        jwtConf,
	}
	b, _ := json.MarshalIndent(Conf, "", "  ")
	fmt.Println("Loaded config: " + string(b))
}

func loadRsaPrivateKey(privateKeyStr string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		log.Logger.Fatal("Failed to decode PEM block")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	return privateKey
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
