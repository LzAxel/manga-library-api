package memory

import (
	"context"
	"manga-library/internal/domain"
	"manga-library/pkg/logger"
	"sync"
)

const previewCollection = "preview"

type PreviewMemory struct {
	m      sync.Map
	logger logger.Logger
}

func NewPreviewMemory(logger logger.Logger) *PreviewMemory {
	return &PreviewMemory{logger: logger}
}

func (m *PreviewMemory) Create(ctx context.Context, preview domain.Preview) (string, error) {
	m.m.Store(preview.Id, preview)

	return preview.URL, nil
}

func (m *PreviewMemory) Delete(ctx context.Context, previewId string) error {
	m.m.Delete(previewId)

	return nil
}
