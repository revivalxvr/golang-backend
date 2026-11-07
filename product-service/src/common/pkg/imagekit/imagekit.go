package imagekit

import (
	"product-service/env"

	ik "github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

func New() *ik.Client {
	// Hanya gunakan Option WithPrivateKey
	// Hapus WithBaseURL agar SDK menggunakan default API Upload endpoint (https://upload.imagekit.io)
	client := ik.NewClient(
		option.WithPrivateKey(env.Conf.ImageKit.PrivateKey),
	)

	return &client
}
