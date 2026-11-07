package env

import (
	"context"
	"fmt"
	"os"
	"product-service/src/common/log"

	"github.com/hashicorp/vault/api"
)

func LoadFromVault() {
	client, err := api.NewClient(&api.Config{
		Address: os.Getenv("VAULT_ADDR"),
	})
	if err != nil {
		log.Logger.Fatal("Failed to create Vault client: " + err.Error())
	}
	client.SetToken(os.Getenv("VAULT_TOKEN"))

	secret, err := client.KVv2("golang-microservices").Get(context.Background(), "product-service")
	if err != nil {
		log.Logger.Fatal("Failed to get secret from Vault: " + err.Error())
	}
	for key, val := range secret.Data {
		strVal := fmt.Sprintf("%v", val)
		os.Setenv(key, strVal)
	}
}
