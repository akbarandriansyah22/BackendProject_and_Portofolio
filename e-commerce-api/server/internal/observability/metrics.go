package observability

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ======================
// METRICS DEFINITIONS
// ======================

var (
	// HTTP Request Metrics
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests by method, path, and status code",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds (p50, p95, p99)",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0},
		},
		[]string{"method", "path"},
	)

	HttpRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "path"},
	)

	HttpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "path", "status"},
	)

	// Error Metrics
	HttpErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total HTTP errors by method, path, and status code",
		},
		[]string{"method", "path", "status"},
	)

	// Application Metrics
	DatabaseConnectionsActive = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	DatabaseConnectionErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_connection_errors_total",
			Help: "Total database connection errors",
		},
		[]string{"reason"},
	)

	// System Metrics (optional, can be extended)
	AppInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_info",
			Help: "Application version and build information",
		},
		[]string{"version", "environment", "build_time"},
	)
)

// ======================
// INIT METRICS
// ======================

func InitMetrics() {
	prometheus.MustRegister(
		// HTTP Metrics
		HttpRequestsTotal,
		HttpRequestDuration,
		HttpRequestSize,
		HttpResponseSize,
		HttpErrorsTotal,
		// Database Metrics
		DatabaseConnectionsActive,
		DatabaseConnectionErrors,
		// Application Info
		AppInfo,
	)

	// Set app info (gauge with static values)
	AppInfo.WithLabelValues("1.0.0", "staging", "unknown").Set(1)
}

// ======================
// /metrics ENDPOINT
// ======================

func RegisterMetricsEndpoint(app *fiber.App) {
	app.Get("/metrics", func(c *fiber.Ctx) error {
		// Set content type
		c.Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8; escaping=underscores")

		// Directly call promhttp.Handler
		// Create a custom http.ResponseWriter that writes to Fiber
		encoder := &metricsEncoder{bodyBuf: []byte{}}

		// Get the default handler and serve it
		handler := promhttp.Handler()
		handler.ServeHTTP(encoder, &http.Request{})

		// Send the response
		return c.Send(encoder.bodyBuf)
	})
}

// metricsEncoder implements http.ResponseWriter
type metricsEncoder struct {
	statusCode int
	headerMap  map[string][]string
	bodyBuf    []byte
}

func (e *metricsEncoder) Header() http.Header {
	if e.headerMap == nil {
		e.headerMap = make(map[string][]string)
	}
	return e.headerMap
}

func (e *metricsEncoder) Write(b []byte) (int, error) {
	e.bodyBuf = append(e.bodyBuf, b...)
	return len(b), nil
}

func (e *metricsEncoder) WriteHeader(statusCode int) {
	e.statusCode = statusCode
}