package utils

import (
	"fmt"
	"math/big"
	"math/rand"
)

func GenerateOTP() (string, error) {
	max := big.NewInt(int64(999999))
	otp := rand.Int63n(max.Int64())

	return fmt.Sprintf("%06d", otp), nil
}
