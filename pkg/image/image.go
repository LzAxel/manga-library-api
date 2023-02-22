package image

import (
	"errors"
	"image"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"image/jpeg"

	"github.com/nfnt/resize"
)

const (
	imageFormats = ".jpg .jpeg"
)

func ValidateExtansion(filename string) error {
	if !strings.Contains(imageFormats, filepath.Ext(filename)) {
		return errors.New("invalid file extension")
	}
	return nil
}

func UploadImage(file multipart.File, filename, dst string) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}

	filename = "preview" + ".jpeg"

	baseImage, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	resizedImage := resize.Resize(250, 350, baseImage, resize.NearestNeighbor)
	outputFile, err := os.Create(filepath.Join(dst, filename))
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if err := jpeg.Encode(outputFile, resizedImage, &jpeg.Options{Quality: 100}); err != nil {
		return err
	}

	return nil
}
