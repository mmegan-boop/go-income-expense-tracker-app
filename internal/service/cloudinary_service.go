package service

import (
	"bytes"
	"context"
	"fmt"
	"go-income-expense-tracker-app/internal/constant"
	"go-income-expense-tracker-app/internal/utils"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService interface {
	UploadRaw(data []byte, fileName string) (string, error)
}

type cloudinaryService struct {
	client *cloudinary.Cloudinary
}

func NewCloudinaryService() (CloudinaryService, error) {
	cloudName := utils.GetConfig(constant.CLOUDINARY_CLOUD_NAME)
	apiKey := utils.GetConfig(constant.CLOUDINARY_API_KEY)
	apiSecret := utils.GetConfig(constant.CLOUDINARY_API_SECRET)
	c, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}

	return &cloudinaryService{client: c}, nil
}

func (s *cloudinaryService) UploadRaw(data []byte, fileName string) (string, error) {
	reader := bytes.NewReader(data)

	result, err := s.client.Upload.Upload(context.Background(), reader, uploader.UploadParams{
		PublicID:     fileName,
		ResourceType: "raw",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to cloudinary: %w", err)
	}

	return result.SecureURL, nil
}
