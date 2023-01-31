package service

import (
	"context"
	"image"
	"manga-library/internal/domain"
	"manga-library/internal/storage"
	"manga-library/pkg/logger"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"image/jpeg"
	"image/png"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

const (
	previewUploadDir = `.\upload\preview\`
	previewUrlBase   = `preview/`
)

type PreviewService struct {
	storage storage.Preview
	logger  logger.Logger
}

func NewPreviewService(storage storage.Preview, logger logger.Logger) *PreviewService {
	return &PreviewService{storage: storage, logger: logger}
}

func (s *PreviewService) Create(ctx context.Context, file multipart.File, filename, uploaderId string) (string, error) {
	id := uuid.NewString()

	filename = id + filepath.Ext(filename)

	preview := domain.Preview{
		Id:         id,
		FileName:   filename,
		UploaderId: uploaderId,
		URL:        previewUrlBase + filename,
		CreatedAt:  time.Now(),
	}

	_, err := file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	baseImage, format, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	resizedImage := resize.Resize(250, 350, baseImage, resize.NearestNeighbor)
	outputFile, err := os.Create(filepath.Join(previewUploadDir, filename))
	if err != nil {
		return "", err
	}
	defer outputFile.Close()

	switch format {
	case "jpeg":
		jpeg.Encode(outputFile, resizedImage, &jpeg.Options{Quality: 100})
	case "png":
		png.Encode(outputFile, resizedImage)
	}

	return s.storage.Create(ctx, preview)
}

func (s *PreviewService) Delete(ctx context.Context, previewId string) error {
	return s.storage.Delete(ctx, previewId)
}
