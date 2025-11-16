package filesystem

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"

	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/internal/module/file/pkg/constant"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/chai2010/webp"
)

type LocalFileStorage struct {
	logger *logger.Logger
	config *config.Config
}

func NewLocalFileStorage(
	logger *logger.Logger,
	config *config.Config,
) *LocalFileStorage {
	return &LocalFileStorage{
		logger: logger,
		config: config,
	}
}

func (s *LocalFileStorage) Upload(name *string, data []byte) error {
	outputPath := s.config.FileDir + *name
	reader := bytes.NewReader(data)

	if err := os.MkdirAll(s.config.FileDir, 0755); err != nil {
		s.logger.Errorf("failed to create directory: %s", s.config.FileDir)
		return err
	}

	if img, err := webp.Decode(reader); err == nil {
		return s.saveAsWebP(img, outputPath)
	}

	reader.Seek(0, io.SeekStart)
	img, _, err := image.Decode(reader)
	if err == nil {
		newPath := s.replaceExtToWebP(outputPath)
		*name = *name + ".webp"
		return s.saveAsWebP(img, newPath)
	}

	kind, _ := filetype.Match(data)
	if os.WriteFile(outputPath+"."+kind.Extension, data, 0644) != nil {
		s.logger.Errorf("failed to write file: %s", outputPath)
		return err
	}

	return nil
}

func (s *LocalFileStorage) Delete(name string) error {
	filePath := s.config.FileDir + name
	if err := os.Remove(filePath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete file %s: %w", name, err)
		}
	}
	return nil
}

func (s *LocalFileStorage) Get(name string) ([]byte, error) {
	filePath := s.config.FileDir + name
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, constant.ErrFileNotFound
	}
	return os.ReadFile(filePath)
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
