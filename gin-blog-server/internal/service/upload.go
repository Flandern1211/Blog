package service

import (
	"gin-blog/internal/utils/upload"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type UploadService interface {
	UploadFile(c *gin.Context, file *multipart.FileHeader) (string, error)
}

type uploadService struct{}

func NewUploadService() UploadService {
	return &uploadService{}
}

func (s *uploadService) UploadFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	oss := upload.NewOSS()
	url, _, err := oss.UploadFile(file)
	if err != nil {
		return "", err
	}
	return url, nil
}
