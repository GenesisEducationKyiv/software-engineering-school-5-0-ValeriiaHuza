package metrics

import (
	"fmt"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		metrics.GetOrCreateCounter(fmt.Sprintf(
			`http_requests_total{endpoint="%s", method="%s", status="%d"}`,
			path, c.Request.Method, status,
		)).Inc()

		metrics.GetOrCreateCounter(`http_requests_total_all`).Inc()

		metrics.GetOrCreateHistogram(fmt.Sprintf(
			`http_request_duration_seconds{endpoint="%s", method="%s"}`,
			path, c.Request.Method,
		)).Update(duration)
	}
}
