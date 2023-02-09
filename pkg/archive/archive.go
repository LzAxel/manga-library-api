package archive

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

const chapterImageFormats = ".jpg .jpeg"

func CreateChapterArchive(chapterPath, mangaTitle string, volume, number int) error {
	file, err := os.Create(filepath.Join(chapterPath,
		fmt.Sprintf("%v %d Том %d Глава.zip", mangaTitle, volume, number)))
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	if err := addPagesToArchive(zipWriter, chapterPath); err != nil {
		return err
	}
	return nil
}

func UnzipUploadChapterArchive(uploadPath string, file *multipart.FileHeader) (int, error) {
	fileReader, err := file.Open()
	if err != nil {
		return 0, errors.New("filed to open file")
	}
	defer fileReader.Close()

	archiveBytes, err := ioutil.ReadAll(fileReader)
	if err != nil {
		return 0, errors.New("filed to open file")
	}

	archiveReader, err := zip.NewReader(bytes.NewReader(archiveBytes), int64(len(archiveBytes)))
	if err != nil {
		return 0, errors.New("filed to get zip reader")
	}

	chapterPageCounter := 0
	for _, zipFile := range archiveReader.File {
		fmt.Println("reading file:", zipFile.Name)
		if !strings.Contains(chapterImageFormats, filepath.Ext(zipFile.Name)) {
			continue
		}
		fileBytes, err := readZipFile(zipFile)
		if err != nil {
			continue
		}
		zipFile.Name = fmt.Sprintf("%v%v", chapterPageCounter, filepath.Ext(zipFile.Name))
		fmt.Println("create file:", filepath.Join(uploadPath, zipFile.Name))
		outputFile, err := os.Create(filepath.Join(uploadPath, zipFile.Name))
		if err != nil {
			return 0, err
		}
		defer outputFile.Close()
		_, err = io.Copy(outputFile, bytes.NewReader(fileBytes))
		if err != nil {
			return 0, err
		}

		chapterPageCounter++
	}

	return chapterPageCounter, nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func addPagesToArchive(zf *zip.Writer, chapterPath string) error {
	files, err := ioutil.ReadDir(chapterPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() || !strings.Contains(chapterImageFormats, filepath.Ext(file.Name())) {
			continue
		}
		fileBytes, err := ioutil.ReadFile(filepath.Join(chapterPath, file.Name()))
		if err != nil {
			return err
		}
		fileWriter, err := zf.Create(file.Name())
		if err != nil {
			return err
		}
		fileWriter.Write(fileBytes)
	}

	return nil
}
