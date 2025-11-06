// file_usecase/cleaner.go
package file_usecase

import (
	"context"
	"time"

	"github.com/Fi44er/sdmed/pkg/logger"
)

type FileCleaner struct {
	repository  IFileRepository
	fileStorage IFileStorage
	logger      *logger.Logger
	interval    time.Duration
	ttl         time.Duration
	stopCh      chan struct{}
}

func NewFileCleaner(
	repository IFileRepository,
	fileStorage IFileStorage,
	logger *logger.Logger,
	interval time.Duration,
	ttl time.Duration,
) *FileCleaner {
	return &FileCleaner{
		repository:  repository,
		fileStorage: fileStorage,
		logger:      logger,
		interval:    interval,
		ttl:         ttl,
		stopCh:      make(chan struct{}),
	}
}

func (fc *FileCleaner) Start() {
	ticker := time.NewTicker(fc.interval)

	go func() {
		fc.logger.Infof("File cleaner started with interval: %v", fc.interval)

		for {
			select {
			case <-ticker.C:
				fc.cleanupExpiredFiles()
			case <-fc.stopCh:
				ticker.Stop()
				fc.logger.Info("File cleaner stopped")
				return
			}
		}
	}()
}

func (fc *FileCleaner) Stop() {
	close(fc.stopCh)
}

func (fc *FileCleaner) cleanupExpiredFiles() {
	ctx := context.Background()
	expirationTime := time.Now().Add(-fc.ttl)

	expiredFiles, err := fc.repository.GetExpiredTemporaryFiles(ctx, expirationTime)
	if err != nil {
		fc.logger.Errorf("Failed to get expired files: %v", err)
		return
	}

	if len(expiredFiles) > 0 {
		fc.logger.Infof("Cleaning up %d expired temporary files", len(expiredFiles))
	}

	for _, file := range expiredFiles {
		if err := fc.fileStorage.Delete(file.Name); err != nil {
			fc.logger.Errorf("Failed to delete file from storage %s: %v", file.Name, err)
			continue
		}

		if err := fc.repository.Delete(ctx, file.ID); err != nil {
			fc.logger.Errorf("Failed to delete file record %s: %v", file.ID, err)
			continue
		}

		fc.logger.Debugf("Cleaned up expired file: %s", file.Name)
	}
}
