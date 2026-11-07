package env

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"product-service/src/common/log"
)

type postgres struct {
	DSN string
}
type currentApp struct {
	RestfulAddr   string
	GrpcBasicAuth string
	GrpcAddr      string
	AllowedOrigin string
}
type imagekit struct {
	PublicKey   string
	PrivateKey  string
	UrlEndpoint string
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
	CurrentApp  *currentApp
	Postgres    *postgres
	Jwt         *jwt
	ImageKit    *imagekit
	Kafka       *kafka
	UserService *userService
}

var Conf *Config

func Load() {
	if os.Getenv("MODE") == "NON-DEV" {
		LoadFromVault()
	}
	currentAppConf := new(currentApp)
	currentAppConf.RestfulAddr = os.Getenv("RESTFUL_ADDR")

	currentAppConf.GrpcBasicAuth = os.Getenv("GRPC_BASIC_AUTH")
	currentAppConf.GrpcAddr = os.Getenv("GRPC_ADDR")

	currentAppConf.AllowedOrigin = os.Getenv("ALLOWED_ORIGIN")

	postgresConf := new(postgres)
	postgresConf.DSN = os.Getenv("POSTGRES_DSN")

	jwtConf := new(jwt)

	base64Str, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_PUBLIC_KEY"))
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	jwtConf.PublicKey = loadRsaPublicKey(string(base64Str))

	imagekitConf := new(imagekit)
	imagekitConf.PublicKey = os.Getenv("IMAGEKIT_PUBLIC_KEY")
	imagekitConf.PrivateKey = os.Getenv("IMAGEKIT_PRIVATE_KEY")
	imagekitConf.UrlEndpoint = os.Getenv("IMAGEKIT_URL_ENDPOINT")

	kafka := new(kafka)
	kafka.Addr1 = os.Getenv("KAFKA_ADDR_1")
	kafka.Addr2 = os.Getenv("KAFKA_ADDR_2")
	kafka.Addr3 = os.Getenv("KAFKA_ADDR_3")

	userServiceConf := new(userService)
	userServiceConf.GrpcAddr = os.Getenv("USER_SERVICE_GRPC_ADDR")
	userServiceConf.GrpcAuth = os.Getenv("USER_SERVICE_GRPC_AUTH")

	Conf = &Config{
		CurrentApp:  currentAppConf,
		Postgres:    postgresConf,
		Jwt:         jwtConf,
		ImageKit:    imagekitConf,
		Kafka:       kafka,
		UserService: userServiceConf,
	}
	b, _ := json.MarshalIndent(Conf, "", "  ")
	fmt.Println("Loaded config: " + string(b))
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
