package service

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"manga-library/internal/domain"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	uploadMangaPath     = "./upload/"
	chapterImageFormats = ".jpg .jpeg"
)

func (s *MangaService) UploadChapter(ctx context.Context, chapterDTO domain.UploadChapterDTO) error {
	manga, err := s.storage.GetBySlug(ctx, chapterDTO.MangaSlug)
	if err != nil {
		return err
	}

	uploadPath := filepath.Join(uploadMangaPath, manga.Slug,
		fmt.Sprintf("%v/%v", chapterDTO.Volume, chapterDTO.Number))

	for _, chapter := range manga.Chapters {
		if chapter.Number == chapterDTO.Number && chapter.Volume == chapterDTO.Volume {
			return errors.New("chapter already exists")
		}
	}

	file, err := chapterDTO.File.Open()
	if err != nil {
		return errors.New("filed to open file")
	}
	defer file.Close()

	archiveBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("filed to open file")
	}
	zipReader, err := zip.NewReader(bytes.NewReader(archiveBytes), int64(len(archiveBytes)))
	if err != nil {
		return errors.New("filed to get zip reader")
	}

	if _, err := os.Stat(uploadPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(uploadPath, 0777)
		if err != nil {
			s.logger.Fatalln(err)
			return errors.New("failed to create directory")
		}
	}

	outputZip, err := os.Create(filepath.Join(uploadPath,
		fmt.Sprintf("%v %d Том %d Глава.zip", manga.Title, chapterDTO.Volume, chapterDTO.Number)))
	if err != nil {
		return err
	}
	defer outputZip.Close()

	newArchive := zip.NewWriter(outputZip)
	defer newArchive.Close()

	chapterPageCounter := 0
	for _, zipFile := range zipReader.File {
		if !strings.Contains(chapterImageFormats, filepath.Ext(zipFile.Name)) {
			continue
		}
		s.logger.Debugf("reading file: %v", zipFile.Name)
		fileBytes, err := readZipFile(zipFile)
		if err != nil {
			s.logger.Debugf("failed to read file: %v", zipFile.Name)
			continue
		}
		zipFile.Name = fmt.Sprintf("%v%v", chapterPageCounter, filepath.Ext(zipFile.Name))
		outputFile, err := os.Create(filepath.Join(uploadPath, zipFile.Name))
		if err != nil {
			return err
		}
		defer outputFile.Close()
		io.Copy(outputFile, bytes.NewReader(fileBytes))

		fileWriter, err := newArchive.Create(zipFile.Name)
		if err != nil {
			return err
		}
		fileWriter.Write(fileBytes)

		chapterPageCounter++
	}

	err = s.storage.UploadChapter(ctx, domain.Chapter{
		Volume:      chapterDTO.Volume,
		Number:      chapterDTO.Number,
		PageCount:   chapterPageCounter,
		MangaSlug:   manga.Slug,
		UploaderId:  chapterDTO.UploaderID,
		IsPublished: false,

		UploadedAt: time.Now(),
	})
	if err != nil {
		err = os.RemoveAll(uploadPath)
		if err != nil {
			s.logger.Fatalln(err)
		}
		return errors.New("failed to save chapter")
	}

	return nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
