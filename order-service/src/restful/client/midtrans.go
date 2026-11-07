package client

import (
	"context"
	"encoding/json"
	"fmt"
	"order-service/env"
	"order-service/src/common/dto/request"
	"order-service/src/common/dto/response"
	"order-service/src/common/errors"
	"order-service/src/common/pkg/cbreaker"
	"order-service/src/common/utils"

	"github.com/gofiber/fiber/v3/client"
	"github.com/sony/gobreaker/v2"
)

type Midtrans interface {
	CreateTransaction(ctx context.Context, order *request.SnapTransactionRequest) (*response.MidtransTx, error)
}

type midtransImpl struct {
	cb *gobreaker.CircuitBreaker[any]
}

func NewMidtrans() Midtrans {
	cb := cbreaker.NewRestfulBreaker("midtrans service")
	return &midtransImpl{cb: cb}
}

func (c *midtransImpl) CreateTransaction(ctx context.Context, order *request.SnapTransactionRequest) (*response.MidtransTx, error) {
	res, err := c.cb.Execute(func() (any, error) {
		// PERUBAHAN: Gunakan endpoint Snap Transactions Midtrans
		// Sandbox: https://app.sandbox.midtrans.com/snap/v1/transactions
		// Production: https://app.midtrans.com/snap/v1/transactions
		apiHost := env.Conf.Midtrans.ApiHostUrl // Pastikan env ini mengarah ke base URL midtrans api (misal: https://app.sandbox.midtrans.com)
		url := apiHost + "/snap/v1/transactions"

		jsondata, err := json.Marshal(order)
		if err != nil {
			return nil, err
		}

		auth := utils.CreateMidtransBasicAuth()
		cc := client.New()

		resp, err := cc.Post(url, client.Config{
			Header: map[string]string{
				"Authorization": auth,
				"Content-Type":  "application/json",
				"Accept":        "application/json",
			},
			Body: jsondata,
		})
		if err != nil {
			return nil, err
		}

		code := resp.StatusCode()
		body := resp.Body()

		fmt.Println("MIDTRANS RAW RESPONSE:", string(body))

		if code != 200 && code != 201 {
			return nil, &errors.Response{HttpCode: code, Message: string(body)}
		}

		txRes := new(response.MidtransTx)
		if err := json.Unmarshal(body, txRes); err != nil {
			return nil, err
		}

		return txRes, nil
	})

	if err != nil {
		return nil, err
	}

	txRes, ok := res.(*response.MidtransTx)
	if !ok {
		return nil, &errors.Response{HttpCode: 500, Message: "failed to cast response to midtrans transaction response"}
	}

	return txRes, nil
}
