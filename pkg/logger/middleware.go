package logger

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware() fiber.Handler {
	logger := NewLogger()
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)
		latencyStr := fmt.Sprintf("%.2fms", float64(latency.Microseconds())/1000)
		logger.Infof(
			"Request %s %s | IP: %s | User-Agent: %s | Duration: %s | Status: %d",
			c.Method(),
			c.Path(),
			c.IP(),
			c.Context().UserAgent(),
			latencyStr,
			c.Response().StatusCode(),
		)
		if err != nil {
			logger.Errorf("Error: %v", err)
		}
		return err
	}
}
