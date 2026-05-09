package parser_utils

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2/log"
)

var (
	// globalClient — единый HTTP-клиент с переиспользованием соединений
	globalClient *http.Client
	// rateLimiter — тикер для ограничения частоты запросов
	rateLimiter *time.Ticker
	// mu — мьютекс для синхронизации доступа к rate limiter
	rateMu sync.Mutex
	// initOnce — гарантирует однократную инициализацию
	initOnce sync.Once
)

const (
	// RequestInterval — минимальный интервал между запросами (1.5 секунды)
	RequestInterval = 1500 * time.Millisecond
	// MaxRetries — максимальное количество попыток
	MaxRetries = 3
	// BaseRetryDelay — базовая задержка перед повторной попыткой
	BaseRetryDelay = 5 * time.Second
	// RateLimitPause — пауза при получении 429/503
	RateLimitPause = 60 * time.Second
	// RequestTimeout — таймаут на один запрос
	RequestTimeout = 30 * time.Second
)

func Request(url string) *goquery.Document {
	log.Info("[Request] Fetching URL: ", url)
	doc := RateLimitedRequest(url)
	if doc == nil {
		log.Error("[Request] Failed to fetch URL: ", url)
	}
	return doc
}

// initGlobalClient инициализирует глобальный HTTP-клиент и rate limiter
func initGlobalClient() {
	initOnce.Do(func() {
		transport := &http.Transport{
			MaxIdleConns:        5,
			MaxIdleConnsPerHost: 2,
			MaxConnsPerHost:     2,
			IdleConnTimeout:     90 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		}

		globalClient = &http.Client{
			Transport: transport,
			Timeout:   RequestTimeout,
		}

		globalClient.Transport.(*http.Transport).MaxIdleConns = 5
		globalClient.Transport.(*http.Transport).MaxIdleConnsPerHost = 2
		globalClient.Transport.(*http.Transport).MaxConnsPerHost = 2
		globalClient.Transport.(*http.Transport).IdleConnTimeout = 90 * time.Second
		rateLimiter = time.NewTicker(RequestInterval)
		log.Info("[RateLimiter] Initialized: interval=", RequestInterval, ", maxRetries=", MaxRetries)
	})
}

// waitForRateLimit ожидает разрешения от rate limiter перед запросом
func waitForRateLimit() {
	initGlobalClient()
	rateMu.Lock()
	<-rateLimiter.C
	rateMu.Unlock()
}

// RateLimitedRequest выполняет HTTP GET запрос с rate-limiting и retry
func RateLimitedRequest(url string) *goquery.Document {
	initGlobalClient()

	var lastErr error
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		// Ожидаем rate limit
		waitForRateLimit()

		log.Info("[RateLimiter] Request attempt ", attempt, "/", MaxRetries, " -> ", url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error("[RateLimiter] Error creating request: ", err)
			return nil
		}

		resp, err := globalClient.Do(req)
		if err != nil {
			lastErr = err
			log.Warn("[RateLimiter] Request failed (attempt ", attempt, "): ", err)
			retryDelay := BaseRetryDelay * time.Duration(1<<(attempt-1)) // exponential backoff
			log.Info("[RateLimiter] Waiting ", retryDelay, " before retry...")
			time.Sleep(retryDelay)
			continue
		}

		// Обработка rate-limit ответов (429, 503)
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body.Close()
			log.Warn("[RateLimiter] Got status ", resp.StatusCode, " — pausing for ", RateLimitPause)
			time.Sleep(RateLimitPause)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			log.Warn("[RateLimiter] Got status ", resp.StatusCode, " for URL: ", url)
			retryDelay := BaseRetryDelay * time.Duration(1<<(attempt-1))
			time.Sleep(retryDelay)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			log.Warn("[RateLimiter] Error reading response body: ", err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		if err != nil {
			log.Error("[RateLimiter] Error parsing HTML: ", err)
			return nil
		}

		return doc
	}

	log.Error("[RateLimiter] All ", MaxRetries, " attempts failed for URL: ", url, " last error: ", lastErr)
	return nil
}

// RateLimitedRawRequest выполняет HTTP GET запрос и возвращает сырые байты ответа
// Используется для запросов с кастомными заголовками (например, parse_price)
func RateLimitedRawRequest(req *http.Request) (*http.Response, []byte, error) {
	initGlobalClient()

	var lastErr error
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		// Ожидаем rate limit
		waitForRateLimit()

		log.Info("[RateLimiter] Raw request attempt ", attempt, "/", MaxRetries, " -> ", req.URL.String())

		resp, err := globalClient.Do(req)
		if err != nil {
			lastErr = err
			log.Warn("[RateLimiter] Raw request failed (attempt ", attempt, "): ", err)
			retryDelay := BaseRetryDelay * time.Duration(1<<(attempt-1))
			time.Sleep(retryDelay)
			continue
		}

		// Обработка rate-limit ответов
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			resp.Body.Close()
			log.Warn("[RateLimiter] Got status ", resp.StatusCode, " — pausing for ", RateLimitPause)
			time.Sleep(RateLimitPause)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			log.Warn("[RateLimiter] Error reading raw response: ", err)
			continue
		}

		return resp, body, nil
	}

	return nil, nil, fmt.Errorf("all %d attempts failed, last error: %w", MaxRetries, lastErr)
}

// CleanupRateLimiter останавливает rate limiter (вызывать при завершении)
func CleanupRateLimiter() {
	if rateLimiter != nil {
		rateLimiter.Stop()
		log.Info("[RateLimiter] Stopped")
	}
}

// ResetRateLimiter сбрасывает инициализацию (для тестов)
func ResetRateLimiter() {
	CleanupRateLimiter()
	initOnce = sync.Once{}
	globalClient = nil
	rateLimiter = nil
}

// Helper: создаёт новый буфер из байтов (для transfer в другие функции)
func NewReaderFromBytes(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}
