package service

import (
	"context"
	"errors"
	"fmt"
	"manga-library/internal/domain"
	"manga-library/pkg/archive"
	"os"
	"path/filepath"
	"time"
)

const (
	uploadMangaPath     = "./upload/"
	chapterImageFormats = ".jpg .jpeg"
)

func (s *MangaService) UploadChapter(ctx context.Context, chapterDTO domain.UploadChapterDTO, roles domain.Roles) error {
	manga, err := s.storage.GetBySlug(ctx, chapterDTO.MangaSlug)
	if err != nil {
		return err
	}
	if manga.UploaderId != chapterDTO.UploaderID && !roles.IsAdmin && !roles.IsEditor {
		return domain.ErrNotEditor
	}

	uploadPath := filepath.Join(uploadMangaPath, manga.Slug,
		fmt.Sprintf("%v/%v", chapterDTO.Volume, chapterDTO.Number))

	for _, chapter := range manga.Chapters {
		if chapter.Number == chapterDTO.Number && chapter.Volume == chapterDTO.Volume {
			return errors.New("chapter already exists")
		}
	}

	if _, err := os.Stat(uploadPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(uploadPath, 0777)
		if err != nil {
			s.logger.Fatalln(err)
			return errors.New("failed to create directory")
		}
	}

	pageCount, err := archive.UnzipUploadChapterArchive(uploadPath, chapterDTO.File)

	if err := archive.CreateChapterArchive(uploadPath, manga.Title, chapterDTO.Volume, chapterDTO.Number); err != nil {
		return errors.New("failed to create chapter archive")
	}

	err = s.storage.UploadChapter(ctx, domain.Chapter{
		Volume:      chapterDTO.Volume,
		Number:      chapterDTO.Number,
		PageCount:   pageCount,
		MangaSlug:   manga.Slug,
		UploaderId:  chapterDTO.UploaderID,
		IsPublished: false,
		UploadedAt:  time.Now(),
	})
	if err != nil {
		deleteChapterFiles(uploadPath)
		return err
	}

	return nil
}

func deleteChapterFiles(chapterPath string) error {
	err := os.RemoveAll(chapterPath)
	if err != nil {
		return err
	}
	return nil
}
