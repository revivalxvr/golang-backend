package utils

import (
	"encoding/base64"
	"order-service/env"
)

func CreateMidtransBasicAuth() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(env.Conf.Midtrans.ServerKey+":"))
}
