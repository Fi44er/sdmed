package filesystem

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/chai2010/webp"
)

type LocalFileStorage struct {
	basePath string
	logger   *logger.Logger
}

func NewLocalFileStorage(
	basePath string,
) *LocalFileStorage {
	return &LocalFileStorage{
		basePath: basePath,
	}
}

func (s *LocalFileStorage) UploadFile(name string, data []byte) error {
	outputPath := s.basePath + name
	reader := bytes.NewReader(data)

	if img, err := webp.Decode(reader); err == nil {
		return s.saveAsWebP(img, outputPath)
	}

	reader.Seek(0, io.SeekStart)
	img, _, err := image.Decode(reader)
	if err == nil {
		newPath := s.replaceExtToWebP(outputPath)
		return s.saveAsWebP(img, newPath)
	}

	if err := os.MkdirAll(s.basePath, 0755); err != nil {
		s.logger.Errorf("failed to create directory: %s", s.basePath)
		return err
	}

	if os.WriteFile(outputPath, data, 0644) != nil {
		s.logger.Errorf("failed to write file: %s", outputPath)
		return err
	}

	return nil
}

func (s *LocalFileStorage) replaceExtToWebP(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext) + ".webp"
}

func (s *LocalFileStorage) saveAsWebP(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	options := webp.Options{Quality: 80, Lossless: false}
	if err := webp.Encode(file, img, &options); err != nil {
		return err
	}

	return nil
}
