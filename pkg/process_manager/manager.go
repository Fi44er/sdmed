package process_manager

import (
	"context"
	"sync"
	"time"
)

type IBackgroundProcess interface {
	Start()
	Stop(ctx context.Context) error
	Name() string
}

type iLogger interface {
	Infof(format string, args ...any)
	Info(msg string)
	Errorf(format string, args ...any)
}

type IProcessManager interface {
	Register(process IBackgroundProcess)
	StartAll()
	StopAll(ctx context.Context) error
}

type ProcessManager struct {
	processes []IBackgroundProcess
	logger    iLogger
	wg        sync.WaitGroup
}

func NewProcessManager(logger iLogger) IProcessManager {
	return &ProcessManager{
		processes: make([]IBackgroundProcess, 0),
		logger:    logger,
	}
}

func (pm *ProcessManager) Register(process IBackgroundProcess) {
	pm.processes = append(pm.processes, process)
}

func (pm *ProcessManager) StartAll() {
	pm.logger.Info("Starting all background processes...")

	for _, process := range pm.processes {
		pm.wg.Add(1)
		go func(p IBackgroundProcess) {
			defer pm.wg.Done()
			pm.logger.Infof("Starting process: %s", p.Name())
			p.Start()
		}(process)
	}

	pm.logger.Infof("Started %d background processes", len(pm.processes))
}

func (pm *ProcessManager) StopAll(ctx context.Context) error {
	pm.logger.Info("Stopping all background processes...")

	stopCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	errorCh := make(chan error, len(pm.processes))
	var wg sync.WaitGroup

	for _, process := range pm.processes {
		wg.Add(1)
		go func(p IBackgroundProcess) {
			defer wg.Done()
			pm.logger.Infof("Stopping process: %s", p.Name())
			if err := p.Stop(stopCtx); err != nil {
				pm.logger.Errorf("Failed to stop process %s: %v", p.Name(), err)
				errorCh <- err
			} else {
				pm.logger.Infof("Successfully stopped process: %s", p.Name())
			}
		}(process)
	}

	wg.Wait()
	close(errorCh)

	pm.wg.Wait()

	if len(errorCh) > 0 {
		return <-errorCh
	}

	pm.logger.Info("All background processes stopped")
	return nil
}
