package scraper_service

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	scraper_entity "github.com/Fi44er/sdmed/internal/module/scraper/entity"
	scraper_constant "github.com/Fi44er/sdmed/internal/module/scraper/pkg/constant"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/brotli"
)

type IRateLimiterService interface {
	CleanupRateLimiter()
	RateLimitedRequest(ctx context.Context, url string) *goquery.Document
	RateLimitedRawRequest(ctx context.Context, req *http.Request) (*http.Response, []byte, error)
	Request(ctx context.Context, url string) *goquery.Document
}

type ScraperService struct {
	logger             *logger.Logger
	rateLimiterService IRateLimiterService
}

func NewScraperService(logger *logger.Logger, rateLimiterService IRateLimiterService) *ScraperService {
	return &ScraperService{
		logger:             logger,
		rateLimiterService: rateLimiterService,
	}
}

func (s *ScraperService) Scraper(ctx context.Context, onProgress func(completed, total int, message string)) []scraper_entity.Items {
	regions := scraper_constant.Regions
	mainUrl := "https://ktsr.sfr.gov.ru"

	s.logger.Info("Starting web scraper")
	// log.Info(fmt.Sprintf("[Scraper] Total articles: %d, Total regions: %d, Total price requests: ~%d",
	// 	len(articles), len(regions), len(articles)*len(regions)))

	// Гарантируем cleanup rate limiter при завершении
	defer s.rateLimiterService.CleanupRateLimiter()

	doc := s.rateLimiterService.Request(ctx, mainUrl)
	if doc == nil {
		s.logger.Error("Failed to fetch the main document")
		return nil
	}
	s.logger.Infof("Fetched main document from: %v", mainUrl)

	// 06 - /ru-RU/product/section/index/2
	sectionsMap := s.ParseSectionUrl(doc)
	s.logger.Info("Parsed section URLs")

	// Предварительное заполнение articleUrlMap

	// [06-01-01] -> {url: https://ktsr.sfr.gov.ru/ru-RU/product/product/order86n/5, name: 06-01-01 Трость опорная, регулируемая по высоте, без устройства противоскольжения}
	articleUrlMap := make(map[string]scraper_entity.Category)
	for key, value := range sectionsMap {
		url := fmt.Sprintf("%v%v", mainUrl, value)
		s.logger.Infof("Pre-parsing category article URL for article type: %v URL: %v", key, url)
		s.ParseCategoryArticleUrl(ctx, url, articleUrlMap)
	}

	// Предварительное заполнение productsMap
	productsMap := make(map[string][]scraper_entity.ParseProductsArticlesType)
	for _, info := range articleUrlMap {
		if _, ok := productsMap[info.URL]; !ok {
			s.logger.Infof("Fetching products articles for URL: %v", info.URL)
			productsMap[info.URL] = s.ParseProductsArticles(ctx, info.URL)
		}
	}

	totalTasks := len(articleUrlMap) * len(regions)
	completedTasks := 0

	items := make(map[string]scraper_entity.Items)
	results := make(chan struct {
		article string
		item    scraper_entity.Items
	}, totalTasks)

	var wg sync.WaitGroup
	var mu sync.Mutex

	maxGoroutines := 2
	sem := make(chan struct{}, maxGoroutines)

	for article, _ := range articleUrlMap {
		articleType := strings.Split(article, "-")[0]
		s.logger.Infof("Processing article: %v with type: %v", article, articleType)
		for _, region := range regions {
			wg.Add(1)
			sem <- struct{}{}

			go func(article, articleType string, region scraper_constant.Region) {
				defer wg.Done()
				defer func() { <-sem }()

				select {
				case <-ctx.Done():
					s.logger.Warnf("Context done, stopping goroutine for article: %v and region: %v", article, region.Iso3166)
					return
				default:
				}

				certificatePrice := s.ParceCertificatePriceRegion(ctx, region, article, articleType)
				if certificatePrice == nil {
					s.logger.Warnf("No certificate price found for article: %v in region: %v", article, region.Iso3166)

					mu.Lock()
					completedTasks++
					progress := completedTasks
					mu.Unlock()
					s.logger.Infof("[Progress] %d/%d (%.1f%%)", progress, totalTasks, float64(progress)/float64(totalTasks)*100)
					if onProgress != nil {
						onProgress(progress, totalTasks, fmt.Sprintf("Артикул: %s, регион: %s (нет цены)", article, region.Iso3166))
					}

					return
				}

				mu.Lock()
				if existingItems, exist := items[article]; exist {
					existingItems.Items = append(existingItems.Items, scraper_entity.Item{
						Price:  *certificatePrice,
						Region: region.Iso3166,
					})
					items[article] = existingItems
				} else {
					newItems := scraper_entity.Items{
						CategoryArticle: article,
						CategoryName:    articleUrlMap[article].Name,
						Product:         productsMap[articleUrlMap[article].URL],
						Items: []scraper_entity.Item{
							{
								Price:  *certificatePrice,
								Region: region.Iso3166,
							},
						},
					}
					items[article] = newItems
				}
				result := items[article]

				completedTasks++
				progress := completedTasks
				mu.Unlock()

				s.logger.Infof("[Progress] %d/%d (%.1f%%) — article: %s, region: %s",
					progress, totalTasks, float64(progress)/float64(totalTasks)*100, article, region.Iso3166)
				if onProgress != nil {
					onProgress(progress, totalTasks, fmt.Sprintf("Артикул: %s, регион: %s", article, region.Iso3166))
				}

				select {
				case <-ctx.Done():
					return
				case results <- struct {
					article string
					item    scraper_entity.Items
				}{article, result}:
				}
			}(article, articleType, region)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		mu.Lock()
		items[result.article] = result.item
		mu.Unlock()
	}

	itemSlice := make([]scraper_entity.Items, 0, len(items))
	for _, item := range items {
		itemSlice = append(itemSlice, item)
	}

	s.logger.Infof("Scraping completed. Total items scraped: %v", len(itemSlice))
	return itemSlice
}

func (s *ScraperService) ParseCategoryArticleUrl(ctx context.Context, url string, articleUrlMap map[string]scraper_entity.Category) {
	doc := s.rateLimiterService.Request(ctx, url)

	doc.Find("a.category-inner-item-info").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")

		var cartUrl string
		var article string
		if exists {
			s.Find("div.category-inner-item__catalog").Each(func(i int, div *goquery.Selection) {
				text := strings.TrimSpace(div.Text())
				if text != "" {
					article = text
					cartUrl = "https://ktsr.sfr.gov.ru" + href
				}
			})
			s.Find("div.category-inner-item__title").Each(func(i int, div *goquery.Selection) {
				text := strings.TrimSpace(div.Text())
				if text != "" {
					articleWithSpace := article + " "
					name := strings.Replace(text, articleWithSpace, "", 1)
					articleUrlMap[article] = scraper_entity.Category{
						URL:  cartUrl,
						Name: name,
					}
				}
			})
		}
	})
}

func (s *ScraperService) ParceCertificatePriceRegion(ctx context.Context, region scraper_constant.Region, article string, articleType string) *float64 {
	url := fmt.Sprintf("https://ktsr.sfr.gov.ru/ru-RU/service/compensation/product-header?region=%v&type=%v&code=%v", region.Iso3166, articleType, article)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		s.logger.Errorf("[Price] Ошибка при создании запроса: %v", err)
		return nil
	}

	scraper_constant.AddHeadersToReq(req)

	resp, body, err := s.rateLimiterService.RateLimitedRawRequest(ctx, req)
	if err != nil {
		s.logger.Errorf("[Price] Ошибка при выполнении запроса: %v", err)
		return nil
	}

	if resp.StatusCode != http.StatusOK || body == nil {
		s.logger.Warnf("[Price] Не удалось получить данные: статус %d, URL: %s", resp.StatusCode, url)
		zero := 0.0
		return &zero
	}

	br := brotli.NewReader(bytes.NewReader(body))
	doc, err := goquery.NewDocumentFromReader(br)
	if err != nil {
		s.logger.Errorf("[Price] Ошибка при создании goquery документа: %v", err)
		return nil
	}

	price := s.ParcePrice(doc)
	return &price
}

func (s *ScraperService) ParcePrice(doc *goquery.Document) float64 {
	selector := "div.catalog-products__info-price.catalog-products__info-price_space span"

	priceRaw := doc.Find(selector).First().Text()
	priceRaw = strings.TrimSpace(priceRaw)

	if priceRaw == "" {
		return 0
	}

	cleaned := s.removeUnwantedCharacters(priceRaw)
	cleaned = strings.ReplaceAll(cleaned, ",", ".")

	value, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		s.logger.Errorf("[Price] Ошибка при конвертации цены '%s' в число: %v", priceRaw, err)
		return 0
	}

	return value
}

func (s *ScraperService) removeUnwantedCharacters(str string) string {
	re := regexp.MustCompile("[^0-9.,]+")
	return re.ReplaceAllString(str, "")
}

func (s *ScraperService) ParseProductsArticles(ctx context.Context, url string) []scraper_entity.ParseProductsArticlesType {
	s.logger.Infof("[Products] Fetching products articles from URL: %v", url)
	doc := s.rateLimiterService.Request(ctx, url)
	if doc == nil {
		s.logger.Errorf("[Products] Failed to fetch the initial document from URL: %v", url)
		return nil
	}

	paginationUrls := s.parsePaginationUrl(doc)
	if len(paginationUrls) == 0 {
		s.logger.Warn("[Products] No pagination URLs found, using the initial URL for parsing")
		paginationUrls = []string{url}
	} else {
		s.logger.Infof("[Products] Found %v pagination URLs", len(paginationUrls))
	}

	var productsArticles []scraper_entity.ParseProductsArticlesType

	for i, paginationUrl := range paginationUrls {
		s.logger.Infof("[Products] Parsing pagination page %v/%v: %v", i+1, len(paginationUrls), paginationUrl)

		pageProducts := s.parseProductsArticlesFromPage(ctx, paginationUrl)
		productsArticles = append(productsArticles, pageProducts...)

		s.logger.Infof("[Products] Got %v products from page %v", len(pageProducts), i+1)
	}

	s.logger.Infof("[Products] Total products articles parsed: %v", len(productsArticles))
	return productsArticles
}

func (s *ScraperService) parseProductsArticlesFromPage(ctx context.Context, url string) []scraper_entity.ParseProductsArticlesType {
	s.logger.Infof("[Products] Fetching document from URL: %v", url)
	doc := s.rateLimiterService.Request(ctx, url)
	if doc == nil {
		s.logger.Errorf("[Products] Failed to fetch document from URL: %s", url)
		return nil
	}
	s.logger.Infof("[Products] Fetched document success from URL: %v", url)

	var results []scraper_entity.ParseProductsArticlesType

	doc.Find("a.product-item-info").Each(func(i int, sel *goquery.Selection) {
		article := strings.TrimSpace(sel.Find("div.product-item__article").Text())
		name := strings.TrimSpace(sel.Find("div.product-item__title").Text())

		if article != "" && name != "" {
			results = append(results, scraper_entity.ParseProductsArticlesType{
				Article: article,
				Name:    name,
			})
			s.logger.Infof("[Products] Extracted product - Article: %v, Name: %v", article, name)
		} else {
			s.logger.Warnf("[Products] Missing article or name for product at index: %v", i)
		}
	})

	return results
}

func (s *ScraperService) parsePaginationUrl(doc *goquery.Document) []string {
	var paginationUrls []string
	doc.Find("li.numeric a").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if exists {
			paginationUrls = append(paginationUrls, fmt.Sprintf("https://ktsr.sfr.gov.ru%s", href))
			s.logger.Infof("[Products] Found pagination URL: %v", href)
		} else {
			s.logger.Warnf("[Products] No href attribute found for pagination link at index: %v", i)
		}
	})

	return paginationUrls
}

func (s *ScraperService) ParseSectionUrl(doc *goquery.Document) map[string]string {
	section := make(map[string]string)

	doc.Find("a.catalog-home-item").Each(func(i int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if exists {
			sel.Find("div.catalog-home-item__section").Each(func(i int, div *goquery.Selection) {
				text := strings.TrimSpace(div.Text())
				if text != "" {
					section[s.removeUnwantedCharacters(text)] = href
				}
			})
		}
	})

	return section
}
