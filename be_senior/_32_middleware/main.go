// ============================================================================
// FILE: middleware_guide.go
// TITLE: راهنمای کامل Middleware در Go - Logging, Auth, Recovery
// HOW TO RUN: go run middleware_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Middleware چیست؟
// ============================================================================
//
// Middleware تابعی است که درخواست HTTP را قبل از رسیدن به handler اصلی پردازش می‌کند.
// کاربردهای اصلی:
// 1. Logging - ثبت اطلاعات درخواست‌ها
// 2. Authentication - بررسی احراز هویت
// 3. Recovery - بازیابی از panic
// 4. Rate Limiting - محدودیت نرخ درخواست
// 5. CORS - مدیریت درخواست‌های cross-origin
// 6. Compression - فشرده‌سازی پاسخ
// 7. Request ID - افزودن شناسه یکتا به هر درخواست
//
// قانون طلایی:
// "Middlewareها را به ترتیب درست اجرا کن (Logging اول، Recovery آخر).
//  هر middleware باید کار خاصی انجام دهد و زنجیره را ادامه دهد یا متوقف کند."
// ============================================================================

package __middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// ============================================================================
// بخش 1: Middleware در net/http خالص
// ============================================================================

// 1.1 نوع Middleware در net/http
type Middleware func(http.Handler) http.Handler

// 1.2 Chain کردن Middlewareها
func chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// 1.3 Logging Middleware (net/http)
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// wrapper برای گرفتن status code
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// اجرای handler بعدی
		next.ServeHTTP(wrapped, r)

		// لاگ بعد از اجرا
		log.Printf("[%s] %s %s %d %v",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			wrapped.status,
			time.Since(start),
		)
	})
}

// responseWriter wrapper برای گرفتن status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// 1.4 Recovery Middleware (net/http)
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// 1.5 Authentication Middleware (net/http)
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// بررسی API Key در هدر
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "secret-key-123" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// 1.6 Request ID Middleware (net/http)
type contextKey string

const RequestIDKey contextKey = "request_id"

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// خواندن از هدر یا ایجاد جدید
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}

		// افزودن به context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// 1.7 Rate Limiting Middleware (net/http)
type rateLimiter struct {
	ticker *time.Ticker
	ch     chan struct{}
}

func newRateLimiter(rate time.Duration) *rateLimiter {
	rl := &rateLimiter{
		ticker: time.NewTicker(rate),
		ch:     make(chan struct{}, 1),
	}
	go func() {
		for range rl.ticker.C {
			select {
			case rl.ch <- struct{}{}:
			default:
			}
		}
	}()
	return rl
}

func (rl *rateLimiter) Allow() bool {
	select {
	case <-rl.ch:
		return true
	default:
		return false
	}
}

func rateLimitMiddleware(limitPerSecond int) Middleware {
	rl := newRateLimiter(time.Second / time.Duration(limitPerSecond))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.Allow() {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// 1.8 CORS Middleware (net/http)
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 1.9 Timeout Middleware (net/http)
func timeoutMiddleware(timeout time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				http.Error(w, "Request timeout", http.StatusGatewayTimeout)
			}
		})
	}
}

// 1.10 Example Handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(RequestIDKey)
	w.Write([]byte(fmt.Sprintf("Hello! Request ID: %v\n", requestID)))
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("something went wrong!")
}

// 1.11 اجرای سرور با net/http
func runNetHTTPServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/panic", panicHandler)

	// اعمال middlewareها با ترتیب صحیح
	handler := chain(
		mux,
		requestIDMiddleware,     // اول: request ID
		loggingMiddleware,       // دوم: logging
		rateLimitMiddleware(10), // سوم: rate limit
		recoveryMiddleware,      // آخر: recovery (بیرونی‌ترین)
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("net/http server starting on :8080")
	log.Fatal(server.ListenAndServe())
}

// ============================================================================
// بخش 2: Middleware در Chi
// ============================================================================

/*
نصب:
go get github.com/go-chi/chi/v5
go get github.com/go-chi/chi/v5/middleware
*/

/*
// 2.1 Chi Server با Middlewareهای داخلی و سفارشی
func runChiServer() {
	r := chi.NewRouter()

	// Middlewareهای داخلی Chi
	r.Use(middleware.RequestID)      // افزودن Request ID
	r.Use(middleware.RealIP)         // گرفتن IP واقعی
	r.Use(middleware.Logger)         // لاگ کردن درخواست‌ها
	r.Use(middleware.Recoverer)      // بازیابی از panic
	r.Use(middleware.Timeout(60 * time.Second)) // تایم‌اوت
	r.Use(middleware.Compress(5))     // فشرده‌سازی پاسخ

	// 2.2 Custom Auth Middleware در Chi
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey != "secret-key" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"unauthorized"}`))
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// 2.3 Group-specific Middleware
	r.Route("/api", func(r chi.Router) {
		// این middleware فقط برای مسیرهای /api اعمال می‌شود
		r.Use(rateLimitMiddleware(100))

		r.Get("/users", getUsers)
		r.Get("/users/{id}", getUser)
	})

	// 2.4 Middleware برای Subrouter خاص
	admin := r.Route("/admin", func(r chi.Router) {
		// فقط ادمین‌ها
		r.Use(adminOnlyMiddleware)
		r.Get("/stats", getStats)
		r.Get("/users", listAllUsers)
	})

	_ = admin

	// 2.5 Custom Logging Middleware (پیشرفته)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			log.Printf("[CHI] %s %s - %d - %v",
				r.Method,
				r.URL.Path,
				ww.Status(),
				time.Since(start),
			)
		})
	})

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Chi!"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	log.Println("Chi server starting on :8081")
	http.ListenAndServe(":8081", r)
}

// 2.6 Admin-only Middleware
func adminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("X-User-Role")
		if role != "admin" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "admin only"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// 2.7 Cache Middleware
func cacheMiddleware(cacheDuration time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", int(cacheDuration.Seconds())))
			next.ServeHTTP(w, r)
		})
	}
}

// 2.8 Request Validation Middleware
func validateContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			contentType := r.Header.Get("Content-Type")
			if !strings.HasPrefix(contentType, "application/json") {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// Handler functions
func getUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]string{"user1", "user2"})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	json.NewEncoder(w).Encode(map[string]string{"id": id, "name": "User " + id})
}

func getStats(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]int{"users": 100, "requests": 5000})
}

func listAllUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode([]map[string]string{{"id": "1", "name": "Admin"}})
}
*/

// ============================================================================
// بخش 3: Middleware در Gin
// ============================================================================

/*
نصب:
go get github.com/gin-gonic/gin
*/

/*
// 3.1 Gin Server با Middleware
func runGinServer() {
	r := gin.Default() // شامل Logger و Recovery به صورت پیش‌فرض

	// 3.2 Custom Logging Middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next() // اجرای handler بعدی
		log.Printf("[GIN] %s %s - %d - %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
		)
	})

	// 3.3 Authentication Middleware
	r.Use(func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "secret-key" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	})

	// 3.4 Request ID Middleware
	r.Use(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	})

	// 3.5 Rate Limiting Middleware
	limiter := rateLimiter{limit: 10, interval: time.Second}
	r.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	})

	// 3.6 CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// 3.7 Group-specific Middleware
	api := r.Group("/api")
	api.Use(authMiddlewareGin)
	{
		api.GET("/users", getUsersGin)
		api.GET("/users/:id", getUserGin)
	}

	admin := r.Group("/admin")
	admin.Use(adminMiddlewareGin)
	{
		admin.GET("/stats", getStatsGin)
		admin.GET("/users", listUsersGin)
	}

	// 3.8 Conditional Middleware
	r.Use(func(c *gin.Context) {
		// فقط برای مسیرهای خاص
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Set("api_version", "v1")
		}
		c.Next()
	})

	// Routes
	r.GET("/", func(c *gin.Context) {
		requestID, _ := c.Get("request_id")
		c.String(http.StatusOK, fmt.Sprintf("Hello Gin! Request ID: %v", requestID))
	})

	r.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	log.Println("Gin server starting on :8082")
	r.Run(":8082")
}

// 3.9 Rate Limiter struct برای Gin
type rateLimiterGin struct {
	limit     int
	interval  time.Duration
	tokens    int
	lastReset time.Time
	mu        sync.Mutex
}

func newRateLimiterGin(limit int, interval time.Duration) *rateLimiterGin {
	return &rateLimiterGin{
		limit:     limit,
		interval:  interval,
		tokens:    limit,
		lastReset: time.Now(),
	}
}

func (rl *rateLimiterGin) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.Sub(rl.lastReset) >= rl.interval {
		rl.tokens = rl.limit
		rl.lastReset = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

// 3.10 Auth Middleware برای Gin
func authMiddlewareGin(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization"})
		return
	}
	// بررسی token
	c.Set("user_id", "123")
	c.Next()
}

func adminMiddlewareGin(c *gin.Context) {
	role := c.GetHeader("X-User-Role")
	if role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}
	c.Next()
}

// 3.11 Response Time Middleware
func responseTimeMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	c.Header("X-Response-Time", time.Since(start).String())
}

// 3.12 Error Handler Middleware
func errorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// بررسی خطاهای بعد از اجرای handler
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errors": c.Errors,
			})
		}
	}
}

// Handler functions
func getUsersGin(c *gin.Context) {
	c.JSON(http.StatusOK, []string{"user1", "user2"})
}

func getUserGin(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "name": "User " + id})
}

func getStatsGin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"users": 100})
}

func listUsersGin(c *gin.Context) {
	c.JSON(http.StatusOK, []gin.H{{"id": "1", "name": "Admin"}})
}
*/

// ============================================================================
// بخش 4: مقایسه پیاده‌سازی Middleware
// ============================================================================

func compareMiddlewareImplementations() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("📊 MIDDLEWARE COMPARISON")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ FEATURE              │ net/http     │ Chi          │ Gin                      │
├──────────────────────┼──────────────┼──────────────┼──────────────────────────┤
│ Signature            │ func(http.Handler) http.Handler │
│                      │              │ Same         │ func(*gin.Context)       │
├──────────────────────┼──────────────┼──────────────┼──────────────────────────┤
│ Built-in Middleware  │ ❌ None      │ ✅ Many      │ ✅ Logger, Recovery      │
├──────────────────────┼──────────────┼──────────────┼──────────────────────────┤
│ Easy to write        │ ⚠️ Manual    │ ✅ Easy      │ ✅ Very Easy             │
├──────────────────────┼──────────────┼──────────────┼──────────────────────────┤
│ Chaining             │ ✅ Manual    │ ✅ Built-in  │ ✅ Built-in (Use)        │
├──────────────────────┼──────────────┼──────────────┼──────────────────────────┤
│ Group-specific       │ ⚠️ Manual    │ ✅ Easy      │ ✅ Easy                  │
├──────────────────────┼──────────────┼──────────────┼──────────────────────────┤
│ Performance          │ ✅ Fastest   │ ✅ Very Fast │ ✅ Fast                  │
├──────────────────────┼──────────────┼──────────────┼──────────────────────────┤
│ Context support      │ ✅ Native    │ ✅ Native    │ ⚠️ Wrapped               │
├──────────────────────┼──────────────┼──────────────┼──────────────────────────┤
│ Learning curve       │ Medium       │ Easy         │ Easy                     │
└──────────────────────┴──────────────┴──────────────┴──────────────────────────┘

📌 MIDDLEWARE ORDER MATTERS!

   Correct order (一般):
   1. Request ID / Tracing
   2. Logging
   3. Recovery (should be early)
   4. Rate Limiting
   5. Authentication
   6. Authorization
   7. CORS
   8. Compression
   9. Timeout
   10. Handler

   Wrong order example:
   ❌ Auth → Recovery (panic در auth不会被 recover می‌شود)
   ✅ Recovery → Auth (panic در auth recover می‌شود)
`)
}

// ============================================================================
// بخش 5: Middlewareهای پرکاربرد
// ============================================================================

func commonMiddlewares() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🔧 COMMON MIDDLEWARE PATTERNS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ 1. Request ID Middleware                                       │
├─────────────────────────────────────────────────────────────────┤
│ func requestID(next http.Handler) http.Handler {               │
│     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {│
│         id := r.Header.Get("X-Request-ID")                     │
│         if id == "" { id = uuid.New().String() }               │
│         ctx := context.WithValue(r.Context(), "reqID", id)     │
│         w.Header().Set("X-Request-ID", id)                     │
│         next.ServeHTTP(w, r.WithContext(ctx))                  │
│     })                                                         │
│ }                                                              │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ 2. Security Headers Middleware                                 │
├─────────────────────────────────────────────────────────────────┤
│ func securityHeaders(next http.Handler) http.Handler {         │
│     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {│
│         w.Header().Set("X-Frame-Options", "DENY")              │
│         w.Header().Set("X-Content-Type-Options", "nosniff")    │
│         w.Header().Set("X-XSS-Protection", "1; mode=block")    │
│         w.Header().Set("Strict-Transport-Security", "max-age=31536000")│
│         next.ServeHTTP(w, r)                                   │
│     })                                                         │
│ }                                                              │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ 3. Metrics Middleware (Prometheus)                             │
├─────────────────────────────────────────────────────────────────┤
│ var (                                                          │
│     requestsTotal = prometheus.NewCounterVec(...)              │
│     requestDuration = prometheus.NewHistogramVec(...)          │
│ )                                                              │
│                                                                │
│ func metrics(next http.Handler) http.Handler {                 │
│     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {│
│         start := time.Now()                                    │
│         next.ServeHTTP(w, r)                                   │
│         duration := time.Since(start)                          │
│         requestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()│
│         requestDuration.WithLabelValues(r.Method).Observe(duration.Seconds())│
│     })                                                         │
│ }                                                              │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ 4. Circuit Breaker Middleware                                  │
├─────────────────────────────────────────────────────────────────┤
│ func circuitBreaker(next http.Handler) http.Handler {          │
│     var failures int                                           │
│     var mu sync.Mutex                                          │
│     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {│
│         mu.Lock()                                              │
│         if failures >= 5 {                                     │
│             mu.Unlock()                                        │
│             http.Error(w, "Service unavailable", http.StatusServiceUnavailable)│
│             return                                             │
│         }                                                      │
│         mu.Unlock()                                            │
│                                                                │
│         // Wrap response writer to track status                │
│         next.ServeHTTP(w, r)                                   │
│     })                                                         │
│ }                                                              │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 6: تست Middleware
// ============================================================================

func testMiddleware() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🧪 TESTING MIDDLEWARE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
// Example test for logging middleware
func TestLoggingMiddleware(t *testing.T) {
    // Create a handler that returns 200
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    // Apply middleware
    middleware := loggingMiddleware(handler)

    // Create test request
    req := httptest.NewRequest("GET", "/test", nil)
    rr := httptest.NewRecorder()

    // Execute
    middleware.ServeHTTP(rr, req)

    // Assert
    assert.Equal(t, http.StatusOK, rr.Code)
}

// Example test for auth middleware
func TestAuthMiddleware(t *testing.T) {
    tests := []struct {
        name       string
        apiKey     string
        wantStatus int
    }{
        {"valid key", "secret-key-123", http.StatusOK},
        {"invalid key", "wrong", http.StatusUnauthorized},
        {"no key", "", http.StatusUnauthorized},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            handler := authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
            }))

            req := httptest.NewRequest("GET", "/", nil)
            if tt.apiKey != "" {
                req.Header.Set("X-API-Key", tt.apiKey)
            }
            rr := httptest.NewRecorder()

            handler.ServeHTTP(rr, req)
            assert.Equal(t, tt.wantStatus, rr.Code)
        })
    }
}
`)
}

// ============================================================================
// بخش 7: جمع‌بندی
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE MIDDLEWARE GUIDE")
	fmt.Println("net/http | Chi | Gin")
	fmt.Println(strings.Repeat("=", 80))

	compareMiddlewareImplementations()
	commonMiddlewares()
	testMiddleware()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
💡 NET/HTTP MIDDLEWARE TEMPLATE:

   type Middleware func(http.Handler) http.Handler

   func MyMiddleware(next http.Handler) http.Handler {
       return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
           // Before handler
           next.ServeHTTP(w, r)
           // After handler
       })
   }

   // Usage
   handler := MyMiddleware(http.HandlerFunc(myHandler))

💡 CHI MIDDLEWARE TEMPLATE:

   r := chi.NewRouter()
   r.Use(MyMiddleware)
   r.Use(chiMiddleware.Logger())
   r.Use(chiMiddleware.Recoverer)

   // Group-specific
   r.Route("/api", func(r chi.Router) {
       r.Use(authMiddleware)
       r.Get("/users", handler)
   })

💡 GIN MIDDLEWARE TEMPLATE:

   r := gin.Default()
   r.Use(MyMiddleware)

   func MyMiddleware(c *gin.Context) {
       // Before
       c.Next() // Call next handler
       // After
   }

   // Group-specific
   api := r.Group("/api")
   api.Use(authMiddleware)

💡 BEST PRACTICES:

   1. Keep middleware focused (single responsibility)
   2. Order matters - put recovery first, logging early
   3. Use context to pass data between middleware and handlers
   4. Always call next() or abort properly
   5. Don't block in middleware (use goroutines if needed)
   6. Handle errors gracefully
   7. Add timeouts to prevent hanging
   8. Test middleware in isolation
   9. Document what each middleware does
   10. Use built-in middleware when available

💡 COMMON MIDDLEWARE LIBRARIES:

   • github.com/go-chi/chi/v5/middleware - Chi middleware
   • github.com/gorilla/handlers - Gorilla handlers
   • github.com/rs/cors - CORS
   • github.com/didip/tollbooth - Rate limiting
   • github.com/unrolled/secure - Security headers
   • github.com/urfave/negroni - Middleware framework
   • github.com/justinas/alice - Middleware chaining
`)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎯 MIDDLEWARE GUIDE - COMPLETE")
	fmt.Println("Ready to build secure and observable Go services!")
	fmt.Println(strings.Repeat("=", 80))
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
