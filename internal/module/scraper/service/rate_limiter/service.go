package ratelimiter_service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	scraper_entity "github.com/Fi44er/sdmed/internal/module/scraper/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2/log"
)

type IRateLimiterService interface {
	CleanupRateLimiter()
	RateLimitedRequest(ctx context.Context, url string) *goquery.Document
	RateLimitedRawRequest(ctx context.Context, req *http.Request) (*http.Response, []byte, error)
	Request(ctx context.Context, url string) *goquery.Document
}

type IScraperSettingsRepository interface {
	Get(ctx context.Context) (*scraper_entity.ScraperSettings, error)
}

type RateLimiterService struct {
	logger     *logger.Logger
	repository IScraperSettingsRepository

	rateLimiter  *time.Ticker
	rateMu       sync.Mutex
	globalClient *http.Client
}

func NewRateLimiterService(logger *logger.Logger, repository IScraperSettingsRepository) IRateLimiterService {
	cfg, err := repository.Get(context.Background())
	if err != nil {
		logger.Errorf("[RateLimiter] Failed to get scraper settings: %v", err)
		return nil
	}

	transport := &http.Transport{
		MaxIdleConns:        5,
		MaxIdleConnsPerHost: cfg.MaxGoroutines,
		MaxConnsPerHost:     cfg.MaxGoroutines,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}

	return &RateLimiterService{
		logger:     logger,
		repository: repository,
		globalClient: &http.Client{
			Transport: transport,
			Timeout:   cfg.TimeoutS,
		},
		rateLimiter: time.NewTicker(cfg.IntervalMs),
	}
}

func NewReaderFromBytes(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}

func (s *RateLimiterService) Request(ctx context.Context, url string) *goquery.Document {
	log.Info("[Request] Fetching URL: ", url)
	doc := s.RateLimitedRequest(ctx, url)
	if doc == nil {
		log.Error("[Request] Failed to fetch URL: ", url)
	}
	return doc
}

func (s *RateLimiterService) RateLimitedRequest(ctx context.Context, url string) *goquery.Document {
	cfg, err := s.repository.Get(ctx)
	if err != nil {
		s.logger.Errorf("[RateLimiter] Failed to get scraper settings: %v", err)
		return nil
	}

	s.globalClient.Timeout = cfg.TimeoutS
	s.globalClient.Transport.(*http.Transport).MaxIdleConnsPerHost = cfg.MaxGoroutines
	s.globalClient.Transport.(*http.Transport).MaxConnsPerHost = cfg.MaxGoroutines
	s.rateLimiter = time.NewTicker(cfg.IntervalMs)

	var lastErr error
	for attempt := 1; attempt <= cfg.MaxRetries; attempt++ {
		// Ожидаем rate limit
		s.waitForRateLimit(ctx)

		s.logger.Infof("[RateLimiter] Request attempt %v/%v -> %s", attempt, cfg.MaxRetries, url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error("[RateLimiter] Error creating request: ", err)
			return nil
		}

		resp, err := s.globalClient.Do(req)
		if err != nil {
			lastErr = err
			s.logger.Warnf("[RateLimiter] Request failed (attempt %v): %v", attempt, err)
			retryDelay := cfg.RetryDelayS * time.Duration(1<<(attempt-1)) // exponential backoff
			s.logger.Infof("[RateLimiter] Waiting %v before retry...", retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		// Обработка rate-limit ответов (429, 503)
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body.Close()
			s.logger.Warnf("[RateLimiter] Got status %v — pausing for %v", resp.StatusCode, cfg.PauseS)
			time.Sleep(cfg.PauseS)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			s.logger.Warnf("[RateLimiter] Got status %v for URL: %v", resp.StatusCode, url)
			retryDelay := cfg.RetryDelayS * time.Duration(1<<(attempt-1))
			time.Sleep(retryDelay)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			s.logger.Warnf("[RateLimiter] Error reading response body: %v", err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		if err != nil {
			s.logger.Errorf("[RateLimiter] Error parsing HTML: %v", err)
			return nil
		}

		return doc
	}

	s.logger.Errorf("[RateLimiter] All %v attempts failed for URL: %v last error: %v", cfg.MaxRetries, url, lastErr)
	return nil
}

func (s *RateLimiterService) RateLimitedRawRequest(ctx context.Context, req *http.Request) (*http.Response, []byte, error) {
	cfg, err := s.repository.Get(ctx)
	if err != nil {
		s.logger.Errorf("[RateLimiter] Failed to get scraper settings: %v", err)
		return nil, nil, err
	}

	s.globalClient.Timeout = cfg.TimeoutS
	s.globalClient.Transport.(*http.Transport).MaxIdleConnsPerHost = cfg.MaxGoroutines
	s.globalClient.Transport.(*http.Transport).MaxConnsPerHost = cfg.MaxGoroutines
	s.rateLimiter = time.NewTicker(cfg.IntervalMs)

	var lastErr error
	for attempt := 1; attempt <= cfg.MaxRetries; attempt++ {
		// Ожидаем rate limit
		s.waitForRateLimit(ctx)

		s.logger.Infof("[RateLimiter] Raw request attempt %d/%d -> %s", attempt, cfg.MaxRetries, req.URL.String())

		resp, err := s.globalClient.Do(req)
		if err != nil {
			lastErr = err
			s.logger.Warnf("[RateLimiter] Raw request failed (attempt %d/%d): %v", attempt, cfg.MaxRetries, err)
			retryDelay := cfg.RetryDelayS * time.Duration(1<<(attempt-1))
			time.Sleep(retryDelay)
			continue
		}

		// Обработка rate-limit ответов
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body.Close()
			s.logger.Warnf("[RateLimiter] Got status %d — pausing for %s", resp.StatusCode, cfg.PauseS)
			time.Sleep(cfg.PauseS)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			s.logger.Warnf("[RateLimiter] Error reading raw response: %v", err)
			continue
		}

		return resp, body, nil
	}

	return nil, nil, fmt.Errorf("all %d attempts failed, last error: %w", cfg.MaxRetries, lastErr)
}

func (s *RateLimiterService) CleanupRateLimiter() {
	if s.rateLimiter != nil {
		s.rateLimiter.Stop()
		log.Info("[RateLimiter] Stopped")
	}
}

func (s *RateLimiterService) waitForRateLimit(ctx context.Context) {
	s.updCfg(ctx)
	s.rateMu.Lock()
	<-s.rateLimiter.C
	s.rateMu.Unlock()
}

func (s *RateLimiterService) updCfg(ctx context.Context) {
	cfg, err := s.repository.Get(ctx)
	if err != nil {
		s.logger.Errorf("[RateLimiter] Failed to get scraper settings: %v", err)
	}

	s.globalClient.Timeout = cfg.TimeoutS
	s.globalClient.Transport.(*http.Transport).MaxIdleConnsPerHost = cfg.MaxGoroutines
	s.globalClient.Transport.(*http.Transport).MaxConnsPerHost = cfg.MaxGoroutines
	s.rateLimiter = time.NewTicker(cfg.IntervalMs)
}
