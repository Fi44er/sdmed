package scrapermanager_service

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	scraper_entity "github.com/Fi44er/sdmed/internal/module/scraper/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/utils"
)

type IScraperService interface {
	Scraper(ctx context.Context, params scraper_entity.ScrapeParams, onProgress func(completed, total int, message string)) []scraper_entity.Items
}

type ITRUCodeUseCase interface {
	UpsertMany(ctx context.Context, codes []*scraper_entity.TRUCode) error
	GetByCode(ctx context.Context, code string) (*scraper_entity.TRUCode, error)
}

type IProductUseCase interface {
	CreateMany(ctx context.Context, data []*scraper_entity.Product) error
}

type ScraperManagerService struct {
	logger         *logger.Logger
	scraperService IScraperService
	productUseCase IProductUseCase
	truCodeUseCase ITRUCodeUseCase

	status      ScraperStatusDTO
	statusMu    sync.RWMutex
	mu          sync.Mutex
	cancelFunc  context.CancelFunc
	broadcaster chan ScraperStatusDTO
	subscribers map[chan ScraperStatusDTO]bool
	subMu       sync.Mutex
}

func NewScraperManagerService(
	logger *logger.Logger,
	scraper IScraperService,
	prodUC IProductUseCase,
	truUC ITRUCodeUseCase,
) *ScraperManagerService {
	return &ScraperManagerService{
		logger:         logger,
		scraperService: scraper,
		productUseCase: prodUC,
		truCodeUseCase: truUC,
		broadcaster:    make(chan ScraperStatusDTO, 10),
		subscribers:    make(map[chan ScraperStatusDTO]bool),
	}
}

func (s *ScraperManagerService) GetStatus() ScraperStatusDTO {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.status
}

func (s *ScraperManagerService) Start(ctx context.Context, params scraper_entity.ScrapeParams) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.IsRunning {
		return fmt.Errorf("парсер уже запущен")
	}

	runCtx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	go s.runWorkflow(runCtx, params)
	return nil
}

func (s *ScraperManagerService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}

func (s *ScraperManagerService) runWorkflow(ctx context.Context, params scraper_entity.ScrapeParams) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Errorf("CRITICAL PANIC: %v\nStack: %s", r, debug.Stack())
		}
		s.finalize()
	}()

	items := s.scraperService.Scraper(ctx, params, func(done, total int, msg string) {
		s.updateGlobalStatus(WeightScraping, 0, done, total, msg)
	})
	if items == nil || ctx.Err() != nil {
		return
	}

	enrichedData := s.enrichWithESNSI(ctx, items, 0.70, WeightEnriching)

	if err := s.syncWithDatabase(ctx, enrichedData, 0.85, WeightSyncDB); err != nil {
		s.logger.Errorf("Sync error: %v", err)
		return
	}

	s.updateGlobalStatus(0, 1.0, 1, 1, "Синхронизация завершена успешно")

	time.Sleep(1 * time.Second)
	s.statusMu.Lock()
	s.status.IsRunning = false
	s.statusMu.Unlock()
}

func (s *ScraperManagerService) syncWithDatabase(ctx context.Context, data []enrichedResult, start, weight float64) error {
	truGroup := make(map[string]*scraper_entity.TRUCode)
	var products []scraper_entity.Product

	for _, res := range data {
		code := res.TRU
		if code == "" {
			code = res.Item.CategoryArticle
		}

		if _, ok := truGroup[code]; !ok {
			truGroup[code] = &scraper_entity.TRUCode{
				Code:     code,
				IsCustom: false,
				Prices:   []scraper_entity.TRUCodePrice{},
			}
		}

		for _, regItem := range res.Item.Items {
			truGroup[code].Prices = append(truGroup[code].Prices, scraper_entity.TRUCodePrice{
				RegionIso: regItem.Region,
				Price:     regItem.Price,
			})
		}

		for _, p := range res.Item.Product {
			products = append(products, scraper_entity.Product{Article: p.Article, Name: p.Name})
		}
	}

	truToUpsert := make([]*scraper_entity.TRUCode, 0, len(truGroup))
	for code, entity := range truGroup {
		existing, err := s.truCodeUseCase.GetByCode(ctx, code)
		if err == nil && existing != nil {
			entity.ID = existing.ID
		}
		truToUpsert = append(truToUpsert, entity)
	}

	if len(truToUpsert) > 0 {
		truChunks := chunkSlice(truToUpsert, 100)
		for i, chunk := range truChunks {
			s.updateGlobalStatus(weight/2, start, i+1, len(truChunks),
				fmt.Sprintf("БД: Синхронизация ТРУ (%d/%d)", i+1, len(truChunks)))
			if err := s.truCodeUseCase.UpsertMany(ctx, chunk); err != nil {
				s.logger.Errorf("Failed upsert TRU: %v", err)
			}
		}
	}

	if len(products) > 0 {
		productChunks := chunkSlice(products, 500)
		productStartOffset := start + (weight / 2)

		for i, c := range productChunks {
			s.updateGlobalStatus(weight/2, productStartOffset, i+1, len(productChunks),
				fmt.Sprintf("БД: Сохранение товаров (%d/%d)", i+1, len(productChunks)))
			chunk := make([]*scraper_entity.Product, 0, len(c))
			for idx, _ := range c {
				chunk = append(chunk, &c[idx])
			}
			if err := s.productUseCase.CreateMany(ctx, chunk); err != nil {
				s.logger.Errorf("Failed to save products: %v", err)
			}
		}
	}

	return nil
}

func (s *ScraperManagerService) enrichWithESNSI(_ context.Context, items []scraper_entity.Items, start, weight float64) []enrichedResult {
	total := len(items)
	var completed int32
	results := make([]enrichedResult, total)
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for i, item := range items {
		wg.Add(1)
		go func(idx int, itm scraper_entity.Items) {
			defer wg.Done()
			sem <- struct{}{}

			tru := s.fetchTRUCodeFromAPI(itm.CategoryName)
			<-sem

			done := atomic.AddInt32(&completed, 1)
			s.updateGlobalStatus(weight, start, int(done), total, "Обогащение данными ТРУ...")
			results[idx] = enrichedResult{Item: itm, TRU: tru}
		}(i, item)
	}
	wg.Wait()
	return results
}

func (s *ScraperManagerService) fetchTRUCodeFromAPI(categoryName string) string {
	options := utils.RequestOptions{
		Method: "GET",
		URL:    "https://esnsi.gosuslugi.ru/rest/ext/v1/classifiers/10616/data",
		Query:  map[string]string{"query": categoryName},
	}
	res, err := utils.MakeRequest(options)
	if err != nil {
		return ""
	}
	var apiRes ApiResponse
	if err := json.Unmarshal(res, &apiRes); err != nil || len(apiRes.Body) == 0 {
		return ""
	}
	return apiRes.Body[0].Elements[3].Value
}

const (
	WeightScraping  = 0.70 // 70%
	WeightEnriching = 0.15 // 15%
	WeightSyncDB    = 0.15 // 15%
)

func (s *ScraperManagerService) Subscribe() chan ScraperStatusDTO {
	s.subMu.Lock()
	defer s.subMu.Unlock()
	ch := make(chan ScraperStatusDTO, 10)
	s.subscribers[ch] = true
	return ch
}

func (s *ScraperManagerService) Unsubscribe(ch chan ScraperStatusDTO) {
	s.subMu.Lock()
	defer s.subMu.Unlock()
	delete(s.subscribers, ch)
	close(ch)
}

func (s *ScraperManagerService) updateGlobalStatus(phase float64, phaseStart float64, done, total int, msg string) {
	s.statusMu.Lock()
	phaseProgress := 0.0
	if total > 0 {
		phaseProgress = float64(done) / float64(total)
	}

	globalPercent := (phaseStart + (phase * phaseProgress)) * 100

	status := ScraperStatusDTO{
		IsRunning:      true,
		CompletedTasks: done,
		TotalTasks:     total,
		Percent:        globalPercent,
		Message:        msg,
	}
	s.status = status
	s.statusMu.Unlock()

	s.subMu.Lock()
	for ch := range s.subscribers {
		select {
		case ch <- status:
		default:
		}
	}
	s.subMu.Unlock()
}

func (s *ScraperManagerService) finalize() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cancelFunc = nil
}

func chunkSlice[T any](slice []T, size int) [][]T {
	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

type enrichedResult struct {
	Item scraper_entity.Items
	TRU  string
}

type ScraperStatusDTO struct {
	IsRunning      bool    `json:"is_running"`
	CompletedTasks int     `json:"completed_tasks"`
	TotalTasks     int     `json:"total_tasks"`
	Percent        float64 `json:"percent"`
	Message        string  `json:"message"`
}

type CreateProductDTO struct {
	Article string
	Name    string
}

type ApiResponse struct {
	Body []struct {
		Elements []struct {
			Value string `json:"value"`
		} `json:"elements"`
	} `json:"body"`
}
