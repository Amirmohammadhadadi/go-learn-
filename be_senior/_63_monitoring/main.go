// ============================================================================
// FILE: monitoring_guide.go
// TITLE: راهنمای کامل Monitoring در Go - Prometheus, Health Checks, OpenTelemetry
// HOW TO RUN: go run monitoring_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Monitoring چیست و چرا نیاز است؟
// ============================================================================
//
// Monitoring فرآیند جمع‌آوری، ذخیره‌سازی و نمایش متریک‌های برنامه است.
//
// سه پایه اصلی Monitoring (سه‌گانه طلایی):
// 1. Metrics (متریک‌ها): داده‌های عددی در طول زمان (مثل نرخ خطا، latency)
// 2. Health Checks: بررسی سلامت سرویس و وابستگی‌های آن
// 3. Tracing: ردیابی درخواست‌ها در چندین سرویس
//
// چرا Monitoring مهم است؟
// - تشخیص مشکلات قبل از کاربران
//- درک عملکرد برنامه
// - برنامه‌ریزی برای مقیاس‌پذیری
// - بررسی Service Level Agreements (SLAs)
//
// قانون طلایی:
// "هر endpoint باید health check داشته باشد.
//  هر سرویس باید metrics (RED method) را expose کند.
//  هر درخواست مهم باید trace شود."
// ============================================================================

package __monitoring

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// ============================================================================
// بخش 1: Prometheus Metrics - تنظیمات اولیه
// ============================================================================

// 1.1 متریک‌های پیش‌فرض
var (
	// HTTP Request Metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpRequestSize = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_size_bytes",
			Help:       "HTTP request size in bytes",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "path"},
	)

	// Business Metrics
	userRegistrations = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_registrations_total",
			Help: "Total number of user registrations",
		},
	)

	activeUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of currently active users",
		},
	)

	orderValue = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "order_value_usd",
			Help:    "Order value in USD",
			Buckets: []float64{10, 25, 50, 100, 250, 500, 1000},
		},
	)

	// Database Metrics
	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type", "table"},
	)

	dbErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_errors_total",
			Help: "Total number of database errors",
		},
		[]string{"error_type"},
	)

	dbConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	// Cache Metrics
	cacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_name"},
	)

	cacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_name"},
	)

	// System Metrics
	goGoroutines = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_goroutines",
			Help: "Number of goroutines",
		},
	)

	goMemoryUsage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
	)

	goCPUCores = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "go_cpu_cores",
			Help: "Number of CPU cores",
		},
	)
)

// 1.2 Updateroutine برای متریک‌های سیستم
func updateSystemMetrics() {
	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for range ticker.C {
			goGoroutines.Set(float64(runtime.NumGoroutine()))
			goMemoryUsage.Set(float64(getMemoryUsage()))
			goCPUCores.Set(float64(runtime.NumCPU()))
		}
	}()
}

// getMemoryUsage گرفتن مصرف حافظه
func getMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// ============================================================================
// بخش 2: Prometheus Middleware برای HTTP
// ============================================================================

// PrometheusMiddleware میدلور برای جمع‌آوری متریک‌های HTTP
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// اندازه درخواست
		size := r.ContentLength
		if size < 0 {
			size = 0
		}
		httpRequestSize.WithLabelValues(r.Method, r.URL.Path).Observe(float64(size))

		// Wrapper برای گرفتن status code
		ww := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// اجرای handler بعدی
		next.ServeHTTP(ww, r)

		// ثبت متریک‌ها
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", ww.statusCode)).Inc()
	})
}

// statusWriter wrapper برای گرفتن status code
type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// ============================================================================
// بخش 3: Health Checks
// ============================================================================

// HealthStatus وضعیت سلامت
type HealthStatus struct {
	Status    string           `json:"status"`
	Timestamp time.Time        `json:"timestamp"`
	Uptime    string           `json:"uptime"`
	Version   string           `json:"version"`
	Checks    map[string]Check `json:"checks,omitempty"`
}

type Check struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Latency string `json:"latency,omitempty"`
}

// HealthChecker بررسی کننده سلامت
type HealthChecker struct {
	checks    map[string]func(context.Context) error
	startTime time.Time
	mu        sync.RWMutex
}

// NewHealthChecker ایجاد health checker جدید
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks:    make(map[string]func(context.Context) error),
		startTime: time.Now(),
	}
}

// RegisterCheck ثبت یک check جدید
func (h *HealthChecker) RegisterCheck(name string, check func(context.Context) error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[name] = check
}

// RunChecks اجرای همه checks
func (h *HealthChecker) RunChecks(ctx context.Context) map[string]Check {
	h.mu.RLock()
	defer h.mu.RUnlock()

	results := make(map[string]Check)

	for name, check := range h.checks {
		start := time.Now()
		err := check(ctx)
		latency := time.Since(start)

		status := "healthy"
		msg := ""

		if err != nil {
			status = "unhealthy"
			msg = err.Error()
		}

		results[name] = Check{
			Status:  status,
			Message: msg,
			Latency: latency.String(),
		}
	}

	return results
}

// HealthHandler هندلر endpoint /health
func (h *HealthChecker) HealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	results := h.RunChecks(ctx)

	// تعیین status کلی
	overallStatus := "healthy"
	for _, check := range results {
		if check.Status == "unhealthy" {
			overallStatus = "unhealthy"
			break
		}
	}

	uptime := time.Since(h.startTime)

	status := HealthStatus{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		Version:   "1.0.0",
		Checks:    results,
	}

	w.Header().Set("Content-Type", "application/json")

	if overallStatus == "unhealthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(status)
}

// ============================================================================
// بخش 4: مثال‌های Health Checks
// ============================================================================

// DatabaseCheck بررسی اتصال دیتابیس
func DatabaseCheck(db *sql.DB) func(context.Context) error {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		return db.PingContext(ctx)
	}
}

// RedisCheck بررسی اتصال Redis
func RedisCheck(client *redis.Client) func(context.Context) error {
	return func(ctx context.Context) error {
		return client.Ping(ctx).Err()
	}
}

// DiskSpaceCheck بررسی فضای دیسک
func DiskSpaceCheck(path string, minFreeGB float64) func(context.Context) error {
	return func(ctx context.Context) error {
		var stat syscall.Statfs_t
		if err := syscall.Statfs(path, &stat); err != nil {
			return err
		}
		freeGB := float64(stat.Bavail*uint64(stat.Bsize)) / (1024 * 1024 * 1024)
		if freeGB < minFreeGB {
			return fmt.Errorf("disk space low: %.2f GB free, need %.2f GB", freeGB, minFreeGB)
		}
		return nil
	}
}

// MemoryCheck بررسی حافظه
func MemoryCheck(maxPercent float64) func(context.Context) error {
	return func(ctx context.Context) error {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		usedPercent := float64(m.Alloc) / float64(m.Sys) * 100
		if usedPercent > maxPercent {
			return fmt.Errorf("memory usage too high: %.2f%%", usedPercent)
		}
		return nil
	}
}

// ============================================================================
// بخش 5: OpenTelemetry Tracing - تنظیمات اولیه
// ============================================================================

// TracerProvider تنظیمات OpenTelemetry
func initTracerProvider() (*sdktrace.TracerProvider, error) {
	// ایجاد exporter برای stdout (در پروداکشن از Jaeger یا Zipkin استفاده کنید)
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	// ایجاد TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("myapp"),
			semconv.ServiceVersionKey.String("1.0.0"),
			attribute.String("environment", "development"),
		)),
	)

	// تنظیم به عنوان global provider
	otel.SetTracerProvider(tp)

	// تنظیم propagator برای headers (برای context propagation)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}

// TracerHelper helper functions for tracing
type TracerHelper struct {
	tracer trace.Tracer
}

// NewTracerHelper ایجاد tracer helper جدید
func NewTracerHelper(name string) *TracerHelper {
	return &TracerHelper{
		tracer: otel.Tracer(name),
	}
}

// StartSpan شروع یک span جدید با context
func (th *TracerHelper) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return th.tracer.Start(ctx, name, opts...)
}

// AddEvent افزودن رویداد به span فعلی
func (th *TracerHelper) AddEvent(span trace.Span, name string, attrs ...attribute.KeyValue) {
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// RecordError ثبت خطا در span
func (th *TracerHelper) RecordError(span trace.Span, err error, attrs ...attribute.KeyValue) {
	span.RecordError(err, trace.WithAttributes(attrs...))
	span.SetStatus(codes.Error, err.Error())
}

// ============================================================================
// بخش 6: Tracing در HTTP Handlers (مثال)
// ============================================================================

// TraceMiddleware میدلور tracing برای HTTP
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("http-server")

		// استخراج context از هدرها
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		// ایجاد span جدید
		ctx, span := tracer.Start(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		defer span.End()

		// افزودن attributes
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.Path),
			attribute.String("http.user_agent", r.UserAgent()),
		)

		// ذخیره span در context برای استفاده در handler
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		// ثبت status code
		// (در wrapper واقعی باید status code را ذخیره کنید)
	})
}

// ============================================================================
// بخش 7: مثال کامل - Order Service با Metrics, Health, Tracing
// ============================================================================

// OrderService سرویس سفارشات با monitoring کامل
type OrderService struct {
	tracer      *TracerHelper
	healthCheck *HealthChecker
	db          *sql.DB
	redis       *redis.Client
}

// NewOrderService ایجاد سرویس سفارشات
func NewOrderService(db *sql.DB, redis *redis.Client) *OrderService {
	tracer := NewTracerHelper("order-service")
	healthCheck := NewHealthChecker()

	// ثبت health checks
	healthCheck.RegisterCheck("database", DatabaseCheck(db))
	healthCheck.RegisterCheck("redis", RedisCheck(redis))
	healthCheck.RegisterCheck("memory", MemoryCheck(90))
	healthCheck.RegisterCheck("disk", DiskSpaceCheck("/", 1.0))

	return &OrderService{
		tracer:      tracer,
		healthCheck: healthCheck,
		db:          db,
		redis:       redis,
	}
}

// CreateOrder ایجاد سفارش (با tracing کامل)
func (s *OrderService) CreateOrder(ctx context.Context, userID string, amount float64) error {
	// شروع span
	ctx, span := s.tracer.StartSpan(ctx, "CreateOrder",
		trace.WithAttributes(
			attribute.String("user_id", userID),
			attribute.Float64("amount", amount),
		))
	defer span.End()

	// مرحله 1: اعتبارسنجی
	_, validateSpan := s.tracer.StartSpan(ctx, "ValidateOrder")
	if amount <= 0 {
		s.tracer.RecordError(validateSpan, fmt.Errorf("invalid amount"))
		validateSpan.End()
		return fmt.Errorf("invalid amount")
	}
	validateSpan.AddEvent("validation_passed")
	validateSpan.End()

	// مرحله 2: بررسی موجودی
	_, stockSpan := s.tracer.StartSpan(ctx, "CheckStock")
	defer stockSpan.End()

	// شبیه‌سازی کوئری دیتابیس
	dbStart := time.Now()
	// db.Query(...)
	time.Sleep(50 * time.Millisecond)
	dbQueryDuration.WithLabelValues("select", "products").Observe(time.Since(dbStart).Seconds())

	stockSpan.AddEvent("stock_available", attribute.String("product_id", "prod-123"))

	// مرحله 3: ایجاد سفارش
	_, createSpan := s.tracer.StartSpan(ctx, "CreateOrderInDB")
	defer createSpan.End()

	// شبیه‌سازی insert
	dbStart = time.Now()
	// db.Exec(...)
	time.Sleep(100 * time.Millisecond)
	dbQueryDuration.WithLabelValues("insert", "orders").Observe(time.Since(dbStart).Seconds())

	// ثبت متریک‌های کسب و کار
	userRegistrations.Inc() // در واقع increase برای سفارش
	orderValue.Observe(amount)

	createSpan.AddEvent("order_created", attribute.String("order_id", "ORD-123"))

	return nil
}

// HealthCheck هندلر health check
func (s *OrderService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	s.healthCheck.HealthHandler(w, r)
}

// ============================================================================
// بخش 8: Metrics Endpoint و سرور کامل
// ============================================================================

func setupMonitoringServer() {
	// به‌روزرسانی متریک‌های سیستم
	updateSystemMetrics()

	// ایجاد health checker
	healthCheck := NewHealthChecker()

	// ثبت checks
	healthCheck.RegisterCheck("memory", MemoryCheck(90))
	healthCheck.RegisterCheck("goroutine", func(ctx context.Context) error {
		if runtime.NumGoroutine() > 1000 {
			return fmt.Errorf("too many goroutines: %d", runtime.NumGoroutine())
		}
		return nil
	})

	// ایجاد router
	mux := http.NewServeMux()

	// Metrics endpoint (برای Prometheus)
	mux.Handle("/metrics", promhttp.Handler())

	// Health endpoint
	mux.HandleFunc("/health", healthCheck.HealthHandler)
	mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		// بررسی readiness (وابستگی‌ها)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	})

	// Business endpoint با middleware
	orderHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "ok"}`))
	})
	mux.Handle("/api/orders", PrometheusMiddleware(orderHandler))

	// شروع سرور برای metrics (معمولاً روی پورت جداگانه)
	metricsPort := ":9090"
	log.Printf("Metrics server starting on %s", metricsPort)
	log.Printf("  - Prometheus metrics: http://localhost%s/metrics", metricsPort)
	log.Printf("  - Health check: http://localhost%s/health", metricsPort)
	log.Printf("  - Liveness: http://localhost%s/health/live", metricsPort)
	log.Printf("  - Readiness: http://localhost%s/health/ready", metricsPort)

	go func() {
		if err := http.ListenAndServe(metricsPort, mux); err != nil {
			log.Printf("Metrics server error: %v", err)
		}
	}()

	// سرور اصلی برنامه
	appPort := ":8080"
	appMux := http.NewServeMux()
	appMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Application server"))
	})

	log.Printf("Application server starting on %s", appPort)
	if err := http.ListenAndServe(appPort, appMux); err != nil {
		log.Fatal(err)
	}
}

// ============================================================================
// بخش 9: Prometheus Configuration File
// ============================================================================

func prometheusConfigExample() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 PROMETHEUS CONFIGURATION")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# ============================================================================
# prometheus.yml
# ============================================================================

global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets: []

rule_files:
  - "alerts.yml"

scrape_configs:
  - job_name: 'myapp'
    static_configs:
      - targets: ['localhost:9090']
    metrics_path: '/metrics'
    
  - job_name: 'myapp-application'
    static_configs:
      - targets: ['app:8080']
    metrics_path: '/metrics'

# ============================================================================
# alerts.yml
# ============================================================================

groups:
  - name: myapp_alerts
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate"

      - alert: HighRequestLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High request latency"

      - alert: InstanceDown
        expr: up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Instance {{ $labels.instance }} down"
`)
}

// ============================================================================
// بخش 10: Grafana Dashboard (JSON - کوتاه شده)
// ============================================================================

func grafanaDashboardExample() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 GRAFANA DASHBOARD SAMPLE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`

┌─────────────────────────────────────────────────────────────────────────────┐
│ GRAFANA DASHBOARD - RECOMMENDED PANELS                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ 1. REQUEST RATE                                                            │
│    Query: rate(http_requests_total[5m])                                   │
│    Panel: Graph                                                           │
│                                                                             │
│ 2. ERROR RATE                                                              │
│    Query: rate(http_requests_total{status=~"5.."}[5m])                    │
│    Panel: Graph                                                           │
│                                                                             │
│ 3. REQUEST LATENCY (P95)                                                  │
│    Query: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))│
│    Panel: Graph                                                           │
│                                                                             │
│ 4. ACTIVE GOROUTINES                                                       │
│    Query: go_goroutines                                                   │
│    Panel: Graph                                                           │
│                                                                             │
│ 5. MEMORY USAGE                                                            │
│    Query: go_memory_usage_bytes                                           │
│    Panel: Graph                                                           │
│                                                                             │
│ 6. DATABASE QUERY LATENCY                                                  │
│    Query: rate(db_query_duration_seconds_sum[5m]) /                       │
│           rate(db_query_duration_seconds_count[5m])                       │
│    Panel: Graph                                                           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 11: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 MONITORING BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. RED METHOD (Rate, Errors, Duration)                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Rate: تعداد درخواست‌ها در ثانیه                                      │
│    • Errors: نرخ خطاها                                                     │
│    • Duration: زمان پاسخگویی (latency)                                    │
│    برای هر endpoint باید این سه متریک داشته باشید                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. METRICS NAMING CONVENTIONS                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│    • _total for counters (http_requests_total)                            │
│    • _seconds for durations (http_request_duration_seconds)               │
│    • _bytes for sizes (http_request_size_bytes)                           │
│    • Namespace: <app>_<component>_<metric>_<unit>                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. HEALTH CHECK LEVELS                                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Liveness: آیا برنامه زنده است؟ (بدون وابستگی‌ها)                    │
│    • Readiness: آیا ترافیک قبول می‌کند؟ (با وابستگی‌ها)                   │
│    • Startup: آیا برنامه شروع شده است؟                                    │
│    • Dependency: بررسی وابستگی‌های خارجی                                  │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. TRACING                                                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Always propagate context (headers)                                   │
│    • Sample only a percentage (e.g., 1%) for production                   │
│    • Add business-relevant attributes                                     │
│    • Record errors and important events                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. ALERTING                                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Alert on symptoms, not causes                                        │
│    • Use time windows (avoid flapping)                                    │
│    • Set reasonable thresholds                                            │
│    • Include runbooks in alerts                                           │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 12: Main
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 MONITORING IN GO")
	fmt.Println("Prometheus Metrics | Health Checks | OpenTelemetry Tracing")
	fmt.Println(strings.Repeat("=", 80))

	bestPractices()
	prometheusConfigExample()
	grafanaDashboardExample()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 Starting Monitoring Demo")
	fmt.Println(strings.Repeat("=", 80))

	// راه‌اندازی سرور monitoring
	go setupMonitoringServer()

	// راه‌اندازی tracing (اختیاری)
	tp, err := initTracerProvider()
	if err != nil {
		log.Printf("Failed to initialize tracer: %v", err)
	}
	defer tp.Shutdown(context.Background())

	fmt.Println("\n📊 Monitoring endpoints:")
	fmt.Println("  • Prometheus metrics: http://localhost:9090/metrics")
	fmt.Println("  • Health check: http://localhost:9090/health")
	fmt.Println("  • Liveness: http://localhost:9090/health/live")
	fmt.Println("  • Readiness: http://localhost:9090/health/ready")

	// نگه داشتن برنامه
	select {}
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// تحلی برای کامپایل
type sqlDB struct{}
type redisClient struct{}

func (c *redisClient) Ping(ctx context.Context) *redisStatusCmd { return nil }

type redisStatusCmd struct{}

func (cmd *redisStatusCmd) Err() error { return nil }

var _ = syscall.Statfs_t{}
var _ = sync.RWMutex{}
