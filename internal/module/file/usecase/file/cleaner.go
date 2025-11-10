// file_usecase/cleaner.go
package file_usecase

import (
	"context"
	"sync"
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
	running     bool
	mutex       sync.RWMutex
}

func (fc *FileCleaner) Name() string {
	return "file_cleaner"
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
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	if fc.running {
		fc.logger.Warn("File cleaner is already running")
		return
	}

	fc.stopCh = make(chan struct{})
	fc.running = true

	ticker := time.NewTicker(fc.interval)

	go func() {
		fc.logger.Infof("File cleaner started with interval: %v", fc.interval)

		fc.cleanupExpiredFiles()
		for {
			select {
			case <-ticker.C:
				fc.cleanupExpiredFiles()
			case <-fc.stopCh:
				ticker.Stop()
				fc.mutex.Lock()
				fc.running = false
				fc.mutex.Unlock()
				fc.logger.Info("File cleaner stopped")
				return
			}
		}
	}()
}

func (fc *FileCleaner) Stop(ctx context.Context) error {
	fc.mutex.Lock()
	defer fc.mutex.Unlock()

	if !fc.running {
		return nil
	}

	close(fc.stopCh)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
		return nil
	}
}

func (fc *FileCleaner) IsRunning() bool {
	fc.mutex.RLock()
	defer fc.mutex.RUnlock()
	return fc.running
}

func (fc *FileCleaner) cleanupExpiredFiles() {
	ctx := context.Background()

	expiredFiles, err := fc.repository.GetExpiredTemporaryFiles(ctx)
	if err != nil {
		fc.logger.Errorf("Failed to get expired files: %v", err)
		return
	}

	fc.logger.Debug("len expired files", len(expiredFiles))

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
