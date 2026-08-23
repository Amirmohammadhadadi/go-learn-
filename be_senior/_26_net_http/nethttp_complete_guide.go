// ============================================================================
// FILE: nethttp_complete_guide.go
// TITLE: راهنمای کامل پکیج net/http در Go - ساخت سرور و کلاینت حرفه‌ای
// HOW TO RUN: go run nethttp_complete_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - net/http چیست؟
// ============================================================================
//
// پکیج net/http امکانات کامل برای:
// 1. ساخت HTTP سرور (Server)
// 2. درخواست‌های HTTP کلاینت (Client)
// 3. مدیریت مسیرها (Routing)
// 4. میدلورها (Middleware)
// 5. مدیریت سشن و کوکی
// 6. فایل سرور استاتیک
// 7. WebSocket (با کمک پکیج‌های دیگر)
//
// قانون طلایی:
// "net/http به تنهایی برای 90% پروژه‌ها کافی است. قبل از رفتن به فریمورک‌های
//  شخص ثالث، اول ببینید با خود Go می‌توانید کارتان را انجام دهید"
// ============================================================================

package __internal_packages

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"
)

// ============================================================================
// بخش 1: HTTP Server پایه - ساده‌ترین سرور
// ============================================================================

// 1.1 یک Handler ساده (تابع با امضای http.HandlerFunc)
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// تنظیم هدرها
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Custom-Header", "my-value")

	// تنظیم وضعیت (پیش‌فرض 200 OK)
	w.WriteHeader(http.StatusOK)

	// نوشتن پاسخ
	fmt.Fprintf(w, "Hello, %s! You requested: %s\n",
		r.URL.Path, r.Method)
}

// 1.2 Handler که JSON برمی‌گرداند
type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	msg := Message{
		Status:  "success",
		Message: "Hello from JSON API",
		Time:    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// تبدیل به JSON و ارسال
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// 1.3 Handler با پارامترهای URL
func greetHandler(w http.ResponseWriter, r *http.Request) {
	// گرفتن پارامتر از URL (مثل /greet/Ali)
	name := strings.TrimPrefix(r.URL.Path, "/greet/")
	if name == "" {
		name = "World"
	}

	fmt.Fprintf(w, "Hello, %s!\n", name)
}

// 1.4 Handler با Query Parameters
func searchHandler(w http.ResponseWriter, r *http.Request) {
	// خواندن query parameters (مثل /search?q=golang&page=2)
	query := r.URL.Query()

	q := query.Get("q")
	page := query.Get("page")
	limit := query.Get("limit")

	if q == "" {
		q = "default"
	}
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	fmt.Fprintf(w, "Search: q=%s, page=%s, limit=%s\n", q, page, limit)
}

// 1.5 Handler با دریافت JSON (POST)
type UserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// فقط POST را قبول کن
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// خواندن body
	var req UserRequest

	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// اعتبارسنجی
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// پاسخ موفقیت
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "created",
		"user":   req,
	})
}

// ============================================================================
// بخش 2: انواع Handler (http.Handler, http.HandlerFunc, http.ServeMux)
// ============================================================================

// 2.1 ساختار سفارشی که http.Handler را پیاده‌سازی می‌کند
type CustomHandler struct {
	Prefix string
}

func (h *CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "CustomHandler with prefix: %s, path: %s\n",
		h.Prefix, r.URL.Path)
}

// 2.2 Handler با closure (دسترسی به متغیرهای خارجی)
func authMiddleware(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// بررسی API key در هدر
		providedKey := r.Header.Get("X-API-Key")
		if providedKey != apiKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, "Authenticated! Path: %s\n", r.URL.Path)
	}
}

// ============================================================================
// بخش 3: میدلورها (Middleware)
// ============================================================================

// 3.1 میدلور لاگینگ
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrapper برای گرفتن وضعیت پاسخ
		ww := &responseWriterWrapper{ResponseWriter: w, status: http.StatusOK}

		// اجرای handler اصلی
		next.ServeHTTP(ww, r)

		// لاگ بعد از اجرا
		log.Printf("%s %s %d %v",
			r.Method, r.URL.Path, ww.status, time.Since(start))
	})
}

// responseWriterWrapper برای گرفتن وضعیت پاسخ
type responseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func (ww *responseWriterWrapper) WriteHeader(status int) {
	ww.status = status
	ww.ResponseWriter.WriteHeader(status)
}

// 3.2 میدلور recovery (جلوگیری از panic)
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// 3.3 میدلور CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 3.4 میدلور rate limiting ساده
func rateLimitMiddleware(next http.Handler) http.Handler {
	limiter := time.NewTicker(100 * time.Millisecond) // 10 requests per second
	queue := make(chan struct{}, 1)

	go func() {
		for range limiter.C {
			select {
			case queue <- struct{}{}:
			default:
			}
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-queue:
			next.ServeHTTP(w, r)
		case <-time.After(1 * time.Second):
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		}
	})
}

// 3.5 میدلور timeout
func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r)
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

// ============================================================================
// بخش 4: HTTP Client - ارسال درخواست‌ها
// ============================================================================

// 4.1 GET درخواست ساده
func simpleGetRequest() {
	resp, err := http.Get("https://api.github.com")
	if err != nil {
		log.Printf("GET error: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read error: %v", err)
		return
	}

	log.Printf("Status: %s, Body length: %d", resp.Status, len(body))
}

// 4.2 GET با هدرها و timeout
func advancedGetRequest(url string) error {
	// ایجاد کلاینت با timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	// ایجاد درخواست
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	// تنظیم هدرها
	req.Header.Set("User-Agent", "MyApp/1.0")
	req.Header.Set("Accept", "application/json")

	// اجرا
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// بررسی وضعیت
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %s", resp.Status)
	}

	// خواندن پاسخ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Response: %s\n", string(body[:min(200, len(body))]))
	return nil
}

// 4.3 POST درخواست با JSON
func postJSONRequest(url string, data interface{}) error {
	// تبدیل داده به JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// ایجاد درخواست
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// اجرا
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// خواندن پاسخ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("POST response: %s\n", string(body))
	return nil
}

// 4.4 POST با form data
func postFormRequest(url string, formData map[string]string) error {
	// ساخت form data
	data := url.Values{}
	for k, v := range formData {
		data.Set(k, v)
	}

	// ارسال
	resp, err := http.PostForm(url, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Form POST response: %s\n", string(body))
	return nil
}

// ============================================================================
// بخش 5: سرور فایل استاتیک (Static Files)
// ============================================================================

// 5.1 سرور فایل ساده
func fileServerExample() {
	// مسیر فایل‌های استاتیک (css, js, images)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// فایل خاص (مثل index.html)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/index.html")
	})
}

// 5.2 ServeFile برای فایل خاص
func serveFileHandler(w http.ResponseWriter, r *http.Request) {
	// با قفل (محتوا را کش نمی‌کند)
	http.ServeFile(w, r, "./files/document.pdf")
}

// ============================================================================
// بخش 6: کوکی‌ها (Cookies) و سشن
// ============================================================================

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	// ایجاد کوکی
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "abc123xyz",
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   86400,
		Secure:   true, // فقط HTTPS
		HttpOnly: true, // قابل دسترسی توسط JavaScript نیست
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
	fmt.Fprintln(w, "Cookie set!")
}

func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	// خواندن کوکی
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintln(w, "No cookie found")
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Cookie value: %s\n", cookie.Value)
}

func deleteCookieHandler(w http.ResponseWriter, r *http.Request) {
	// حذف کوکی (با تنظیم expiration در گذشته)
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}

	http.SetCookie(w, cookie)
	fmt.Fprintln(w, "Cookie deleted!")
}

// ============================================================================
// بخش 7: تست HTTP (httptest)
// ============================================================================

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "test response")
}

func demonstrateTesting() {
	// ایجاد یک تست سرور
	handler := http.HandlerFunc(testHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// ارسال درخواست تست
	resp, err := http.Get(server.URL)
	if err != nil {
		log.Printf("Test error: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Test response: %s\n", body)
}

// ============================================================================
// بخش 8: HTTP/2 و TLS (HTTPS)
// ============================================================================

func startHTTPServer() {
	// سرور با پشتیبانی HTTP/2
	server := &http.Server{
		Addr:         ":8443",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      loggingMiddleware(http.DefaultServeMux),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	// راه‌اندازی HTTPS (نیاز به فایل‌های certificate)
	// go func() {
	//     if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
	//         log.Fatal(err)
	//     }
	// }()
}

// ============================================================================
// بخش 9: راه‌اندازی کامل سرور
// ============================================================================

func startFullServer() {
	// ایجاد router اصلی
	mux := http.NewServeMux()

	// ثبت مسیرها
	mux.HandleFunc("/", helloHandler)
	mux.HandleFunc("/json", jsonHandler)
	mux.HandleFunc("/greet/", greetHandler)
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/users", createUserHandler)
	mux.HandleFunc("/cookie/set", setCookieHandler)
	mux.HandleFunc("/cookie/get", getCookieHandler)
	mux.HandleFunc("/cookie/delete", deleteCookieHandler)

	// استفاده از میدلورها
	handler := recoveryMiddleware(
		loggingMiddleware(
			corsMiddleware(
				timeoutMiddleware(30 * time.Second)(mux),
			),
		),
	)

	// سرور با تنظیمات کامل
	server := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// راه‌اندازی graceful shutdown
	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// منتظر سیگنال برای shutdown (در مثال واقعی از signal استفاده می‌شود)
	// <-sigChan
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// server.Shutdown(ctx)
}

// ============================================================================
// بخش 10: مثال‌های کاربردی
// ============================================================================

// 10.1 دانلود فایل
func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "file parameter required", http.StatusBadRequest)
		return
	}

	// تنظیم هدر برای دانلود
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/octet-stream")

	// ارسال فایل (در مثال واقعی فایل واقعی می‌فرستد)
	http.ServeFile(w, r, "./files/"+filename)
}

// 10.2 آپلود فایل
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// حداکثر 10 MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "Uploaded file: %s, size: %d bytes\n",
		handler.Filename, handler.Size)

	// در مثال واقعی فایل را ذخیره می‌کنیم
	// dst, _ := os.Create("./uploads/" + handler.Filename)
	// io.Copy(dst, file)
}

// 10.3 Proxy معکوس ساده
func reverseProxyHandler(targetURL string) http.HandlerFunc {
	url, _ := url.Parse(targetURL)

	return func(w http.ResponseWriter, r *http.Request) {
		// ساخت درخواست جدید
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.RequestURI = ""

		// ارسال به سرور مقصد
		client := &http.Client{}
		resp, err := client.Do(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// کپی هدرها
		for k, v := range resp.Header {
			w.Header()[k] = v
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

// ============================================================================
// بخش 11: اشتباهات رایج
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("❌ COMMON MISTAKES WITH net/http")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n❌ Mistake 1: Not closing response body")
	fmt.Println("   resp, _ := http.Get(url)")
	fmt.Println("   // missing defer resp.Body.Close() → memory leak")
	fmt.Println("   ✅ defer resp.Body.Close()")

	fmt.Println("\n❌ Mistake 2: Ignoring context cancellation")
	fmt.Println("   ✅ Use r.Context() for cancellation support")

	fmt.Println("\n❌ Mistake 3: Not setting timeouts")
	fmt.Println("   server := &http.Server{Addr: \":8080\"}  // no timeouts!")
	fmt.Println("   ✅ Set ReadTimeout, WriteTimeout, IdleTimeout")

	fmt.Println("\n❌ Mistake 4: Reading large bodies without limit")
	fmt.Println("   io.ReadAll(r.Body)  // could consume all memory")
	fmt.Println("   ✅ Use http.MaxBytesReader() to limit")

	fmt.Println("\n❌ Mistake 5: Not handling panics in handlers")
	fmt.Println("   ✅ Use recovery middleware")

	fmt.Println("\n❌ Mistake 6: Logging sensitive data")
	fmt.Println("   log.Println(r)  // logs headers, possibly API keys")
	fmt.Println("   ✅ Log only what's necessary")
}

// ============================================================================
// بخش 12: جمع‌بندی و اجرا
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE net/http GUIDE IN GO")
	fmt.Println(stringsRepeat("=", 80))

	// نمایش مثال‌های کلاینت (اختیاری)
	fmt.Println("\n📡 Client examples (dry run):")
	simpleGetRequest()

	// راه‌اندازی سرور (در گوروتین جدا)
	go startFullServer()

	// صبر برای نمایش
	time.Sleep(2 * time.Second)

	fmt.Println("\n💡 Server is running on http://localhost:8080")
	fmt.Println("   Test endpoints:")
	fmt.Println("   - GET  /             -> Hello message")
	fmt.Println("   - GET  /json         -> JSON response")
	fmt.Println("   - GET  /greet/Ali    -> Greeting")
	fmt.Println("   - GET  /search?q=go  -> Search with params")
	fmt.Println("   - POST /users        -> Create user (JSON)")
	fmt.Println("   - GET  /cookie/set   -> Set cookie")
	fmt.Println("   - GET  /cookie/get   -> Get cookie")

	// نگه داشتن برنامه
	fmt.Println("\n⚠️  Press Ctrl+C to stop the server")
	select {}
}

// ============================================================================
// بخش 13: توابع کمکی
// ============================================================================

func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
