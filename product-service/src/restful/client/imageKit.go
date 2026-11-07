package client

import (
	"context"
	"os"

	pkgIK "product-service/src/common/pkg/imagekit"

	ik "github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/packages/param"
)

type ImageKit struct {
	ik *ik.Client
}

func NewImageKit() *ImageKit {
	ikClient := pkgIK.New()

	return &ImageKit{
		ik: ikClient,
	}
}

func (c *ImageKit) UploadFile(ctx context.Context, path string, fileName string) (*ik.FileUploadResponse, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	res, err := c.ik.Files.Upload(ctx, ik.FileUploadParams{
		File:              file,
		FileName:          fileName,
		UseUniqueFileName: param.NewOpt(true),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ImageKit) DeleteFile(ctx context.Context, fileID string) error {
	err := c.ik.Files.Delete(ctx, fileID)
	return err
}
