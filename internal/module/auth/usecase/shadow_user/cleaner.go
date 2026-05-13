package shadow_user

import (
	"context"
	"sync"
	"time"

	"github.com/Fi44er/sdmed/pkg/logger"
)

type ShadowUserCleaner struct {
	service  IShadowUserService
	logger   *logger.Logger
	interval time.Duration
	stopCh   chan struct{}
	running  bool
	mutex    sync.RWMutex
}

func (sc *ShadowUserCleaner) Name() string {
	return "shadow_user_cleaner"
}

func NewShadowUserCleaner(
	service IShadowUserService,
	logger *logger.Logger,
	interval time.Duration,
) *ShadowUserCleaner {
	return &ShadowUserCleaner{
		service:  service,
		logger:   logger,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

func (sc *ShadowUserCleaner) Start() {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if sc.running {
		sc.logger.Warn("Shadow user cleaner is already running")
		return
	}

	sc.stopCh = make(chan struct{})
	sc.running = true

	ticker := time.NewTicker(sc.interval)

	go func() {
		sc.logger.Infof("Shadow user cleaner started with interval: %v", sc.interval)

		sc.cleanup()
		for {
			select {
			case <-ticker.C:
				sc.cleanup()
			case <-sc.stopCh:
				ticker.Stop()
				sc.mutex.Lock()
				sc.running = false
				sc.mutex.Unlock()
				sc.logger.Info("Shadow user cleaner stopped")
				return
			}
		}
	}()
}

func (sc *ShadowUserCleaner) Stop(ctx context.Context) error {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if !sc.running {
		return nil
	}

	close(sc.stopCh)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
		return nil
	}
}

func (sc *ShadowUserCleaner) IsRunning() bool {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()
	return sc.running
}

func (sc *ShadowUserCleaner) cleanup() {
	ctx := context.Background()

	if err := sc.service.CleanupExpiredShadows(ctx); err != nil {
		sc.logger.Errorf("Shadow user cleaner: failed to cleanup expired shadows: %v", err)
	}
}
