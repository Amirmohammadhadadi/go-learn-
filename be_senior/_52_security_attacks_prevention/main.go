// ============================================================================
// FILE: security_attacks_prevention_guide.go
// TITLE: راهنمای کامل جلوگیری از حملات رایج در Go
// HOW TO RUN: go run security_attacks_prevention_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - حملات رایج و راه‌های جلوگیری
// ============================================================================
//
// حملات رایج در برنامه‌های وب:
//
// 1. SQL Injection (تزریق SQL)
//    - مهاجم کد SQL مخرب را در ورودی قرار می‌دهد
//    - راه‌حل: استفاده از prepared statements / parameterized queries
//
// 2. XSS (Cross-Site Scripting)
//    - مهاجم کد JavaScript مخرب را تزریق می‌کند
//    - راه‌حل: escaping خروجی، CSP headers
//
// 3. CSRF (Cross-Site Request Forgery)
//    - مهاجم کاربر را فریب می‌دهد تا درخواست ناخواسته ارسال کند
//    - راه‌حل: CSRF tokens, SameSite cookies
//
// 4. CORS (Cross-Origin Resource Sharing)
//    - کنترل دسترسی منابع از دامنه‌های دیگر
//    - راه‌حل: پیکربندی صحیح CORS headers
//
// 5. Rate Limiting
//    - مهاجم با درخواست‌های زیاد سرور را overwhelm می‌کند
//    - راه‌حل: محدودیت نرخ درخواست‌ها
//
// 6. Timeout و Body Limit
//    - مهاجم با درخواست‌های طولانی یا بزرگ منابع را اشغال می‌کند
//    - راه‌حل: timeout و محدودیت حجم body
//
// قانون طلایی:
// "هیچ‌وقت به ورودی کاربر اعتماد نکن. همیشه از prepared statements استفاده کن.
//  خروجی را escape کن. CSRF token را برای درخواست‌های state-changing استفاده کن.
//  از rate limiting و timeout برای جلوگیری از DoS استفاده کن."
// ============================================================================

package __security_attacks_prevention

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver برای مثال
)

// ============================================================================
// بخش 1: SQL Injection Prevention (با Prepared Statements)
// ============================================================================

// ============================================
// ❌ مثال بد: آسیب‌پذیر به SQL Injection
// ============================================
// NEVER DO THIS! این کد خطرناک است
func vulnerableQuery(db *sql.DB, userInput string) error {
	// این کد آسیب‌پذیر است! مثال: ورودی "'; DROP TABLE users; --"
	query := fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", userInput)
	_, err := db.Exec(query)
	return err
}

// ============================================
// ✅ مثال خوب: استفاده از Prepared Statements
// ============================================

// SQLSafeQuery ساختار برای کوئری امن
type SQLSafeQuery struct {
	db *sql.DB
}

// GetUserByID دریافت کاربر با ID (با prepared statement)
func (s *SQLSafeQuery) GetUserByID(ctx context.Context, id int) (*User, error) {
	// روش 1: استفاده از QueryRow با placeholder
	query := "SELECT id, name, email FROM users WHERE id = $1"
	var user User
	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsersByName جستجوی کاربران با نام (با LIKE امن)
func (s *SQLSafeQuery) GetUsersByName(ctx context.Context, name string) ([]User, error) {
	// استفاده از placeholder برای جلوگیری از injection حتی در LIKE
	query := "SELECT id, name, email FROM users WHERE name LIKE $1"
	rows, err := s.db.QueryContext(ctx, query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// CreateUser ایجاد کاربر جدید (با prepared statement)
func (s *SQLSafeQuery) CreateUser(ctx context.Context, name, email string) error {
	// روش 2: استفاده از Exec با placeholder
	query := "INSERT INTO users (name, email, created_at) VALUES ($1, $2, $3)"
	_, err := s.db.ExecContext(ctx, query, name, email, time.Now())
	return err
}

// BulkInsertUsers درج دسته‌ای کاربران (با prepared statement)
func (s *SQLSafeQuery) BulkInsertUsers(ctx context.Context, users []User) error {
	// آماده‌سازی statement
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, user := range users {
		if _, err := stmt.ExecContext(ctx, user.Name, user.Email); err != nil {
			return err
		}
	}
	return nil
}

// ============================================================================
// بخش 2: XSS Prevention (Cross-Site Scripting)
// ============================================================================

// ============================================
// ❌ مثال بد: آسیب‌پذیر به XSS
// ============================================
func vulnerableXSSHandler(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("comment")
	// این کد آسیب‌پذیر است! ورودی: <script>alert('XSS')</script>
	fmt.Fprintf(w, "Your comment: %s", userInput)
}

// ============================================
// ✅ مثال خوب: جلوگیری از XSS
// ============================================

// XSSPrevention ساختار برای جلوگیری از XSS
type XSSPrevention struct{}

// EscapeHTML escape کردن HTML (روش 1)
func (x *XSSPrevention) EscapeHTML(input string) string {
	return template.HTMLEscapeString(input)
}

// EscapeJS escape کردن JavaScript
func (x *XSSPrevention) EscapeJS(input string) string {
	return template.JSEscapeString(input)
}

// EscapeURL escape کردن URL
func (x *XSSPrevention) EscapeURL(input string) string {
	return template.URLQueryEscaper(input)
}

// SafeHTMLHandler هندلر امن با escaping
func (x *XSSPrevention) SafeHTMLHandler(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("comment")
	safeOutput := x.EscapeHTML(userInput)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Your comment: %s", safeOutput)
}

// SafeJSONHandler هندلر JSON امن
func (x *XSSPrevention) SafeJSONHandler(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("data")

	// JSON encoding به طور خودکار escape می‌کند
	response := map[string]string{"data": userInput}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CSPHeadersMiddleware میدلور برای افزودن CSP headers
func CSPHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Content-Security-Policy header
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline';")

		// X-XSS-Protection header (برای مرورگرهای قدیمی)
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// X-Content-Type-Options
		w.Header().Set("X-Content-Type-Options", "nosniff")

		next.ServeHTTP(w, r)
	})
}

// ============================================================================
// بخش 3: CSRF Prevention (Cross-Site Request Forgery)
// ============================================================================

// CSRFManager مدیریت توکن‌های CSRF
type CSRFManager struct {
	tokens map[string]csrfToken
	mu     sync.RWMutex
}

type csrfToken struct {
	Value     string
	ExpiresAt time.Time
}

// NewCSRFManager ایجاد مدیر CSRF جدید
func NewCSRFManager() *CSRFManager {
	m := &CSRFManager{
		tokens: make(map[string]csrfToken),
	}
	// شروع cleanup routine
	go m.cleanup()
	return m
}

func (m *CSRFManager) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for key, token := range m.tokens {
			if now.After(token.ExpiresAt) {
				delete(m.tokens, key)
			}
		}
		m.mu.Unlock()
	}
}

// GenerateToken تولید توکن CSRF
func (m *CSRFManager) GenerateToken(sessionID string) string {
	b := make([]byte, 32)
	rand.Read(b)
	token := base64.URLEncoding.EncodeToString(b)

	m.mu.Lock()
	defer m.mu.Unlock()
	m.tokens[sessionID] = csrfToken{
		Value:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	return token
}

// VerifyToken بررسی توکن CSRF
func (m *CSRFManager) VerifyToken(sessionID, token string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stored, exists := m.tokens[sessionID]
	if !exists {
		return false
	}
	if time.Now().After(stored.ExpiresAt) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(stored.Value), []byte(token)) == 1
}

// CSRFMiddleware میدلور CSRF
func (m *CSRFManager) CSRFMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// فقط برای درخواست‌های POST, PUT, DELETE, PATCH
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next(w, r)
			return
		}

		// دریافت توکن از هدر یا فرم
		token := r.Header.Get("X-CSRF-Token")
		if token == "" {
			token = r.FormValue("csrf_token")
		}

		// دریافت session ID (از کوکی)
		sessionID := r.Context().Value("session_id").(string)

		if !m.VerifyToken(sessionID, token) {
			http.Error(w, "CSRF token invalid", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// CSRFFormField تولید فیلد فرم CSRF
func (m *CSRFManager) CSRFFormField(sessionID string) template.HTML {
	token := m.GenerateToken(sessionID)
	return template.HTML(fmt.Sprintf(`<input type="hidden" name="csrf_token" value="%s">`, token))
}

// ============================================================================
// بخش 4: CORS Configuration
// ============================================================================

// CORSConfig تنظیمات CORS
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig تنظیمات پیش‌فرض CORS
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           86400,
	}
}

// CORSMiddleware میدلور CORS
func CORSMiddleware(config *CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// بررسی origin مجاز
			allowed := false
			for _, o := range config.AllowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}

			if allowed && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			if config.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
			w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ============================================================================
// بخش 5: Rate Limiting
// ============================================================================

// RateLimiter محدودکننده نرخ درخواست
type RateLimiter struct {
	requests map[string][]time.Time
	limit    int
	window   time.Duration
	mu       sync.RWMutex
}

// NewRateLimiter ایجاد محدودکننده نرخ جدید
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
	// شروع cleanup
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, timestamps := range rl.requests {
			var valid []time.Time
			for _, t := range timestamps {
				if now.Sub(t) <= rl.window {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = valid
			}
		}
		rl.mu.Unlock()
	}
}

// Allow بررسی اجازه درخواست
func (rl *RateLimiter) Allow(key string) (bool, int, int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// پاک کردن درخواست‌های قدیمی
	var valid []time.Time
	for _, t := range rl.requests[key] {
		if now.Sub(t) <= rl.window {
			valid = append(valid, t)
		}
	}
	rl.requests[key] = valid

	if len(valid) >= rl.limit {
		return false, rl.limit, len(valid)
	}

	rl.requests[key] = append(rl.requests[key], now)
	return true, rl.limit, len(rl.requests[key])
}

// RateLimitMiddleware میدلور rate limiting
func (rl *RateLimiter) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// استفاده از IP کلاینت به عنوان کلید
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		key := ip

		// همچنین می‌توان از User ID برای کاربران احراز هویت شده استفاده کرد
		if userID := r.Header.Get("X-User-ID"); userID != "" {
			key = "user:" + userID
		}

		allowed, limit, current := rl.Allow(key)

		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(limit-current))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rl.window).Unix(), 10))

		if !allowed {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}

// ============================================================================
// بخش 6: Timeout و Body Limit
// ============================================================================

// TimeoutMiddleware میدلور timeout
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r.WithContext(ctx))
				done <- struct{}{}
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

// BodyLimitMiddleware میدلور محدودیت حجم body
func BodyLimitMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}

// TimeoutAndLimitHandler مثال هندلر با timeout و body limit
func TimeoutAndLimitHandler(w http.ResponseWriter, r *http.Request) {
	// خواندن body با محدودیت
	body, err := io.ReadAll(r.Body)
	if err != nil {
		if errors.Is(err, http.ErrBodyTooLarge) {
			http.Error(w, "Body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	// شبیه‌سازی کار طولانی (با context)
	select {
	case <-time.After(100 * time.Millisecond):
		w.Write([]byte("Success"))
	case <-r.Context().Done():
		http.Error(w, "Request cancelled", http.StatusGatewayTimeout)
	}
}

// ============================================================================
// بخش 7: Security Headers Middleware یکپارچه
// ============================================================================

// SecurityHeadersMiddleware میدلور امنیتی یکپارچه
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// HSTS (HTTP Strict Transport Security)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// X-Frame-Options (جلوگیری از clickjacking)
		w.Header().Set("X-Frame-Options", "DENY")

		// X-Content-Type-Options
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// X-XSS-Protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Referrer-Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content-Security-Policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// Permissions-Policy
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		next.ServeHTTP(w, r)
	})
}

// ============================================================================
// بخش 8: Full Secure Server Setup
// ============================================================================

func setupSecureServer() {
	// ایجاد rate limiter
	rateLimiter := NewRateLimiter(100, time.Minute) // 100 requests per minute

	// ایجاد CSRF manager
	csrfManager := NewCSRFManager()

	// تنظیمات CORS
	corsConfig := &CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3000", "https://example.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           86400,
	}

	// ایجاد handler اصلی
	mux := http.NewServeMux()

	// اعمال middlewareها با ترتیب صحیح
	handler := http.Handler(mux)
	handler = SecurityHeadersMiddleware(handler)
	handler = CORSMiddleware(corsConfig)(handler)
	handler = TimeoutMiddleware(30 * time.Second)(handler)
	handler = BodyLimitMiddleware(10 << 20)(handler)   // 10 MB
	handler = rateLimiter.RateLimitMiddleware(handler) // تبدیل به http.HandlerFunc مناسب

	// Routes
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Secure data"))
	})

	mux.HandleFunc("/api/submit", csrfManager.CSRFMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// پردازش درخواست
		w.Write([]byte("Submitted"))
	}))

	// سرور با تنظیمات امنیتی
	server := &http.Server{
		Addr:              ":8443",
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	log.Println("Secure server starting on :8443")
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}

// ============================================================================
// بخش 9: Example Handlers با تمام محافظت‌ها
// ============================================================================

type SecureHandler struct {
	db          *sql.DB
	rateLimiter *RateLimiter
	csrfManager *CSRFManager
	xssProtect  *XSSPrevention
}

func NewSecureHandler(db *sql.DB) *SecureHandler {
	return &SecureHandler{
		db:          db,
		rateLimiter: NewRateLimiter(100, time.Minute),
		csrfManager: NewCSRFManager(),
		xssProtect:  &XSSPrevention{},
	}
}

// SecureQueryHandler هندلر کوئری امن
func (h *SecureHandler) SecureQueryHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Rate limiting
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if allowed, _, _ := h.rateLimiter.Allow(ip); !allowed {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// 2. CSRF check (برای POST)
	if r.Method == http.MethodPost {
		token := r.Header.Get("X-CSRF-Token")
		sessionID := r.Context().Value("session_id").(string)
		if !h.csrfManager.VerifyToken(sessionID, token) {
			http.Error(w, "Invalid CSRF token", http.StatusForbidden)
			return
		}
	}

	// 3. Get and validate input
	userInput := r.URL.Query().Get("name")
	if userInput == "" {
		http.Error(w, "Name required", http.StatusBadRequest)
		return
	}

	// 4. استفاده از prepared statement (امن در برابر SQL injection)
	query := "SELECT id, name, email FROM users WHERE name = $1"
	rows, err := h.db.QueryContext(r.Context(), query, userInput)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// 5. Escape خروجی (جلوگیری از XSS)
	var results []string
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			continue
		}
		// Escape قبل از نمایش
		results = append(results, h.xssProtect.EscapeHTML(user.Name))
	}

	// 6. پاسخ امن
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// ============================================================================
// بخش 10: Best Practices Summary
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 SECURITY BEST PRACTICES SUMMARY")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. SQL INJECTION                                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✅ ALWAYS use prepared statements/parameterized queries                │
│    ✅ Use QueryRowContext with placeholders ($1, $2, etc.)                │
│    ✅ Validate and sanitize input before use                              │
│    ❌ NEVER concatenate user input into SQL strings                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. XSS (Cross-Site Scripting)                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✅ Escape all output: template.HTMLEscapeString                        │
│    ✅ Use Content-Security-Policy headers                                 │
│    ✅ Set X-XSS-Protection: 1; mode=block                                 │
│    ❌ NEVER output user input directly without escaping                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. CSRF (Cross-Site Request Forgery)                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✅ Generate unique CSRF tokens per session                             │
│    ✅ Validate tokens for state-changing requests (POST, PUT, DELETE)     │
│    ✅ Use SameSite=Strict cookies                                         │
│    ❌ NEVER accept state-changing requests without CSRF validation        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. CORS (Cross-Origin Resource Sharing)                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✅ Configure specific allowed origins (not '*')                        │
│    ✅ Set allowed methods and headers                                      │
│    ✅ Use preflight (OPTIONS) requests                                     │
│    ❌ NEVER allow all origins without validation                          │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. RATE LIMITING                                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✅ Implement per-IP or per-user rate limits                            │
│    ✅ Return X-RateLimit headers                                          │
│    ✅ Use 429 status code when limit exceeded                             │
│    ❌ NEVER leave endpoints without rate limiting                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 6. TIMEOUT & BODY LIMIT                                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✅ Set read/write timeouts on server                                   │
│    ✅ Limit request body size                                             │
│    ✅ Use context for cancellation                                        │
│    ❌ NEVER rely on default unlimited timeouts                            │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 11: Common Mistakes
// ============================================================================

func commonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚠️ COMMON SECURITY MISTAKES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 1: String concatenation in SQL                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ fmt.Sprintf("SELECT * FROM users WHERE name = '%s'", name)          │
│    ✅ db.Query("SELECT * FROM users WHERE name = $1", name)               │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 2: Not escaping output                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ fmt.Fprintf(w, userInput)                                           │
│    ✅ fmt.Fprintf(w, template.HTMLEscapeString(userInput))                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 3: No CSRF protection                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ POST endpoints without token validation                             │
│    ✅ Always validate CSRF token for state-changing requests              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 4: CORS with '*'                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ Access-Control-Allow-Origin: *                                      │
│    ✅ Access-Control-Allow-Origin: https://trusted-domain.com             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 5: No rate limiting                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ No limits on API endpoints                                          │
│    ✅ Implement rate limiting for all public endpoints                     │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MISTAKE 6: No timeouts                                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│    ❌ Default server with no timeouts                                     │
│    ✅ Set ReadTimeout, WriteTimeout, IdleTimeout                          │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 12: Types
// ============================================================================

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// برای کامپایل (در واقع باید import "database/sql" داشته باشید)
type sqlDB struct{}

func (db *sqlDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (db *sqlDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}

func (db *sqlDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (db *sqlDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return nil, nil
}

type sqlRows struct{}
type sqlRow struct{}
type sqlResult struct{}

// ============================================================================
// بخش 13: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 SECURITY ATTACKS PREVENTION GUIDE")
	fmt.Println("SQL Injection | XSS | CSRF | CORS | Rate Limiting | Timeout")
	fmt.Println(stringsRepeat("=", 80))

	bestPractices()
	commonMistakes()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 QUICK CHECKLIST")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ DEPLOYMENT CHECKLIST                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ☐ SQL Injection: All queries use prepared statements                      │
│  ☐ XSS: Output is escaped, CSP headers set                                 │
│  ☐ CSRF: Tokens validated for POST/PUT/DELETE                              │
│  ☐ CORS: Configured with specific origins                                  │
│  ☐ Rate Limiting: Implemented for APIs                                     │
│  ☐ Timeout: Server and request timeouts set                                │
│  ☐ Body Limit: MaxBytesReader configured                                   │
│  ☐ Security Headers: HSTS, X-Frame-Options, etc.                           │
│  ☐ HTTPS: TLS enabled in production                                        │
│  ☐ Input Validation: All user input validated                              │
│  ☐ Error Handling: No stack traces exposed                                 │
│  ☐ Dependencies: Regularly updated                                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
`)

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 SECURITY GUIDE - COMPLETE")
	fmt.Println("Build secure Go applications!")
	fmt.Println(stringsRepeat("=", 80))
}

// تابع کمکی
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
