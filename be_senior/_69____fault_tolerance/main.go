// ============================================================================
// FILE: fault_tolerance_guide.go
// TITLE: راهنمای کامل Fault Tolerance در Go - Retry, Backoff, Circuit Breaker
// HOW TO RUN: go run fault_tolerance_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Fault Tolerance چیست و چرا نیاز است؟
// ============================================================================
//
// Fault Tolerance توانایی سیستم برای ادامه کار در صورت بروز خطا در بخش‌های مختلف است.
//
// الگوهای اصلی Fault Tolerance:
//
// 1. Retry (تلاش مجدد)
//    - تکرار درخواست‌های ناموفق
//    - همراه با Backoff (فاصله زمانی افزایشی)
//    - مناسب برای خطاهای موقتی (نوسان شبکه، timeout)
//
// 2. Backoff (تأخیر هوشمند)
//    -线性: wait = base * attempt
//    -指数: wait = base * 2^attempt
//    -全双曲: wait = min(max, base * 2^attempt) + jitter
//    - Jitter: افزودن randomness برای جلوگیری از thundering herd
//
// 3. Circuit Breaker (قطع کننده مدار)
//    - حالت‌ها: Closed (بسته) → Open (باز) → Half-Open (نیمه‌باز)
//    - Closed: درخواست‌ها عبور می‌کنند، خطاها شمارش می‌شوند
//    - Open: درخواست‌ها بدون اجرا fail می‌شوند (مدار قطع است)
//    - Half-Open: بعد از timeout، چند درخواست تست ارسال می‌شوند
//    - مناسب برای جلوگیری از cascade failure
//
// کتابخانه‌های محبوب:
// - go-resiliency: مجموعه‌ای از الگوهای resiliency
// - gobreaker: پیاده‌سازی circuit breaker
// - retry: پیاده‌سازی retry با backoff
//
// قانون طلایی:
// "همیشه برای درخواست‌های شبکه retry با backoff پیاده‌سازی کن.
//  از circuit breaker برای محافظت از سیستم در برابر خطاهای زنجیره‌ای استفاده کن.
//  همیشه jitter اضافه کن تا از thundering herd جلوگیری شود."
// ============================================================================

package ______fault_tolerance

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// بخش 1: Retry با Backoff پایه (بدون کتابخانه)
// ============================================================================

// RetryConfig تنظیمات retry
type RetryConfig struct {
	MaxAttempts   int           // حداکثر تعداد تلاش
	InitialDelay  time.Duration // تأخیر اولیه
	MaxDelay      time.Duration // حداکثر تأخیر
	BackoffFactor float64       // ضریب رشد (مثلاً 2 برای exponential)
	Jitter        bool          // افزودن randomness
}

// DefaultRetryConfig تنظیمات پیش‌فرض
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
	}
}

// RetryableFunc تابعی که می‌توان retry کرد
type RetryableFunc func() error

// RetryWithBackoff اجرای تابع با retry و backoff
func RetryWithBackoff(fn RetryableFunc, config RetryConfig) error {
	var err error
	delay := config.InitialDelay

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		// اگر آخرین تلاش بود، خطا را برگردان
		if attempt == config.MaxAttempts-1 {
			break
		}

		// محاسبه تأخیر
		currentDelay := delay
		if config.Jitter {
			currentDelay = addJitter(currentDelay)
		}

		log.Printf("Attempt %d failed: %v, retrying in %v", attempt+1, err, currentDelay)
		time.Sleep(currentDelay)

		// محاسبه تأخیر بعدی (exponential backoff)
		delay = time.Duration(float64(delay) * config.BackoffFactor)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", config.MaxAttempts, err)
}

// addJitter افزودن randomness به تأخیر
func addJitter(delay time.Duration) time.Duration {
	jitter := time.Duration(rand.Int63n(int64(delay / 2)))
	return delay + jitter
}

// ============================================================================
// بخش 2: Retry با Context (قابل کنسل شدن)
// ============================================================================

// RetryWithContext اجرای تابع با retry و پشتیبانی از context
func RetryWithContext(ctx context.Context, fn RetryableFunc, config RetryConfig) error {
	var err error
	delay := config.InitialDelay

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		// بررسی context
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = fn()
		if err == nil {
			return nil
		}

		if attempt == config.MaxAttempts-1 {
			break
		}

		currentDelay := delay
		if config.Jitter {
			currentDelay = addJitter(currentDelay)
		}

		log.Printf("Attempt %d failed: %v, retrying in %v", attempt+1, err, currentDelay)

		// منتظر ماندن با قابلیت cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(currentDelay):
		}

		delay = time.Duration(float64(delay) * config.BackoffFactor)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", config.MaxAttempts, err)
}

// ============================================================================
// بخش 3: Retry با تشخیص خطاهای مشمول retry
// ============================================================================

// RetryConfigWithPredicate تنظیمات retry با شرط
type RetryConfigWithPredicate struct {
	MaxAttempts   int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
	Jitter        bool
	ShouldRetry   func(error) bool // تابع تشخیص خطاهای قابل retry
}

// DefaultRetryConfigWithPredicate تنظیمات پیش‌فرض
func DefaultRetryConfigWithPredicate() RetryConfigWithPredicate {
	return RetryConfigWithPredicate{
		MaxAttempts:   3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      5 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
		ShouldRetry:   DefaultShouldRetry,
	}
}

// DefaultShouldRetry خطاهای قابل retry پیش‌فرض
func DefaultShouldRetry(err error) bool {
	if err == nil {
		return false
	}
	// خطاهای موقتی (network timeout, connection refused, etc.)
	return errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled)
}

// RetryWithPredicate اجرای retry با شرط خطا
func RetryWithPredicate(fn RetryableFunc, config RetryConfigWithPredicate) error {
	var err error
	delay := config.InitialDelay

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		if attempt == config.MaxAttempts-1 || !config.ShouldRetry(err) {
			break
		}

		currentDelay := delay
		if config.Jitter {
			currentDelay = addJitter(currentDelay)
		}

		log.Printf("Attempt %d failed: %v, retrying in %v", attempt+1, err, currentDelay)
		time.Sleep(currentDelay)

		delay = time.Duration(float64(delay) * config.BackoffFactor)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return fmt.Errorf("failed after %d attempts: %w", config.MaxAttempts, err)
}

// ============================================================================
// بخش 4: Circuit Breaker - پیاده‌سازی کامل
// ============================================================================

// CircuitBreakerState وضعیت‌های circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota // مدار بسته (عادی)
	StateOpen                              // مدار باز (عدم接受)
	StateHalfOpen                          // نیمه‌باز (تست)
)

// String تبدیل وضعیت به رشته
func (s CircuitBreakerState) String() string {
	switch s {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF-OPEN"
	default:
		return "UNKNOWN"
	}
}

// CircuitBreakerConfig تنظیمات circuit breaker
type CircuitBreakerConfig struct {
	FailureThreshold int           // تعداد خطاهای مجاز قبل از باز شدن مدار
	SuccessThreshold int           // تعداد موفقیت‌های مورد نیاز در حالت Half-Open
	Timeout          time.Duration // مدت زمان ماندن در حالت Open
	HalfOpenMaxCalls int           // حداکثر تعداد درخواست‌های تست در Half-Open
}

// DefaultCircuitBreakerConfig تنظیمات پیش‌فرض
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		Timeout:          30 * time.Second,
		HalfOpenMaxCalls: 1,
	}
}

// CircuitBreaker ساختار اصلی
type CircuitBreaker struct {
	config    CircuitBreakerConfig
	state     CircuitBreakerState
	mu        sync.RWMutex

	// آمار
	failureCount int
	successCount int

	// زمان آخرین failure (برای timeout)
	lastFailureTime time.Time

	// Half-Open state
	halfOpenCalls int
}

// NewCircuitBreaker ایجاد circuit breaker جدید
func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config: config,
		state:  StateClosed,
	}
}

// Execute اجرای تابع با circuit breaker
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.allowRequest() {
		return errors.New("circuit breaker is open")
	}

	err := fn()
	cb.recordResult(err == nil)
	return err
}

// allowRequest بررسی مجاز بودن درخواست
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailureTime) > cb.config.Timeout {
			cb.transitionToHalfOpen()
			return true
		}
		return false
	case StateHalfOpen:
		if cb.halfOpenCalls < cb.config.HalfOpenMaxCalls {
			cb.halfOpenCalls++
			return true
		}
		return false
	default:
		return false
	}
}

// recordResult ثبت نتیجه
func (cb *CircuitBreaker) recordResult(success bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		if success {
			cb.failureCount = 0
			cb.successCount++
		} else {
			cb.failureCount++
			cb.lastFailureTime = time.Now()
			if cb.failureCount >= cb.config.FailureThreshold {
				cb.transitionToOpen()
			}
		}
	case StateHalfOpen:
		cb.halfOpenCalls = 0
		if success {
			cb.successCount++
			if cb.successCount >= cb.config.SuccessThreshold {
				cb.transitionToClosed()
			}
		} else {
			cb.transitionToOpen()
		}
	}
}

func (cb *CircuitBreaker) transitionToOpen() {
	cb.state = StateOpen
	cb.failureCount = 0
	cb.successCount = 0
	cb.lastFailureTime = time.Now()
	log.Printf("Circuit breaker transitioned to OPEN")
}

func (cb *CircuitBreaker) transitionToHalfOpen() {
	cb.state = StateHalfOpen
	cb.halfOpenCalls = 0
	cb.successCount = 0
	log.Printf("Circuit breaker transitioned to HALF-OPEN")
}

func (cb *CircuitBreaker) transitionToClosed() {
	cb.state = StateClosed
	cb.failureCount = 0
	cb.successCount = 0
	log.Printf("Circuit breaker transitioned to CLOSED")
}

// GetState دریافت وضعیت فعلی
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// GetStats دریافت آمار
func (cb *CircuitBreaker) GetStats() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"state":          cb.state.String(),
		"failure_count":  cb.failureCount,
		"success_count":  cb.successCount,
		"half_open_calls": cb.halfOpenCalls,
	}
}

// ============================================================================
// بخش 5: Circuit Breaker با Context
// ============================================================================

// CircuitBreakerWithContext نسخه با پشتیبانی از context
type CircuitBreakerWithContext struct {
	*CircuitBreaker
}

// NewCircuitBreakerWithContext ایجاد circuit breaker جدید
func NewCircuitBreakerWithContext(config CircuitBreakerConfig) *CircuitBreakerWithContext {
	return &CircuitBreakerWithContext{
		CircuitBreaker: NewCircuitBreaker(config),
	}
}

// ExecuteWithContext اجرا با context
func (cb *CircuitBreakerWithContext) ExecuteWithContext(ctx context.Context, fn func(context.Context) error) error {
	if !cb.allowRequest() {
		return errors.New("circuit breaker is open")
	}

	// ایجاد channel برای دریافت نتیجه
	resultCh := make(chan error, 1)
	go func() {
		resultCh <- fn(ctx)
	}()

	select {
	case err := <-resultCh:
		cb.recordResult(err == nil)
		return err
	case <-ctx.Done():
		cb.recordResult(false)
		return ctx.Err()
	}
}

// ============================================================================
// بخش 6: go-resiliency (کتابخانه محبوب)
// ============================================================================

/*
// نصب:
// go get github.com/eapache/go-resiliency/retrier
// go get github.com/eapache/go-resiliency/breaker

import (
	"github.com/eapache/go-resiliency/retrier"
	"github.com/eapache/go-resiliency/breaker"
)

// 6.1 Retrier از go-resiliency
func resilientRetryExample() {
	// ایجاد retrier با backoff
	r := retrier.New(retrier.ExponentialBackoff(3, 100*time.Millisecond), nil)

	err := r.Run(func() error {
		// عملیات ممکن است fail شود
		return callExternalAPI()
	})

	if err != nil {
		log.Printf("Failed after retries: %v", err)
	}
}

// 6.2 Breaker از go-resiliency
func resilientBreakerExample() {
	// ایجاد circuit breaker
	b := breaker.New(3, 1, 5*time.Second)

	for i := 0; i < 10; i++ {
		result := b.Run(func() error {
			// عملیات خطرناک
			return callExternalAPI()
		})

		switch result {
		case nil:
			log.Println("Success")
		case breaker.ErrBreakerOpen:
			log.Println("Circuit breaker is open")
		default:
			log.Printf("Error: %v", result)
		}
	}
}
*/

// ============================================================================
// بخش 7: مثال - HTTP Client با Retry و Circuit Breaker
// ============================================================================

// ResilientHTTPClient کلاینت HTTP مقاوم در برابر خطا
type ResilientHTTPClient struct {
	client         *http.Client
	circuitBreaker *CircuitBreaker
	retryConfig    RetryConfigWithPredicate
}

// NewResilientHTTPClient ایجاد کلاینت مقاوم
func NewResilientHTTPClient() *ResilientHTTPClient {
	return &ResilientHTTPClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		circuitBreaker: NewCircuitBreaker(DefaultCircuitBreakerConfig()),
		retryConfig:    DefaultRetryConfigWithPredicate(),
	}
}

// Do انجام درخواست HTTP با resilience
func (c *ResilientHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	// اجرا با circuit breaker و retry
	err = c.circuitBreaker.Execute(func() error {
		return RetryWithPredicate(func() error {
			resp, err = c.client.Do(req)
			if err != nil {
				return err
			}

			// بررسی status code
			if resp.StatusCode >= 500 {
				resp.Body.Close()
				return fmt.Errorf("server error: %d", resp.StatusCode)
			}

			return nil
		}, c.retryConfig)
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Get انجام درخواست GET
func (c *ResilientHTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// ============================================================================
// بخش 8: Bulkhead Pattern (ایزوله کردن منابع)
// ============================================================================

// Bulkhead محدودکننده همزمانی (برای ایزوله کردن)
type Bulkhead struct {
	semaphore chan struct{}
	timeout   time.Duration
}

// NewBulkhead ایجاد bulkhead جدید
func NewBulkhead(maxConcurrent int, timeout time.Duration) *Bulkhead {
	return &Bulkhead{
		semaphore: make(chan struct{}, maxConcurrent),
		timeout:   timeout,
	}
}

// Execute اجرای تابع با محدودیت همزمانی
func (b *Bulkhead) Execute(fn func() error) error {
	select {
	case b.semaphore <- struct{}{}:
		defer func() { <-b.semaphore }()
		return fn()
	case <-time.After(b.timeout):
		return errors.New("bulkhead timeout: too many concurrent requests")
	}
}

// ============================================================================
// بخش 9: Combined Patterns (ترکیب الگوها)
// ============================================================================

// ResilientService سرویس مقاوم با ترکیب الگوها
type ResilientService struct {
	circuitBreaker *CircuitBreaker
	retryConfig    RetryConfig
	bulkhead       *Bulkhead
}

// NewResilientService ایجاد سرویس مقاوم
func NewResilientService() *ResilientService {
	return &ResilientService{
		circuitBreaker: NewCircuitBreaker(DefaultCircuitBreakerConfig()),
		retryConfig:    DefaultRetryConfig(),
		bulkhead:       NewBulkhead(10, 5*time.Second),
	}
}

// CallExternalAPI فراخوانی API خارجی با تمام محافظت‌ها
func (s *ResilientService) CallExternalAPI(ctx context.Context, url string) error {
	// لایه 1: Bulkhead (محدودیت همزمانی)
	return s.bulkhead.Execute(func() error {
		// لایه 2: Circuit Breaker
		return s.circuitBreaker.Execute(func() error {
			// لایه 3: Retry with Backoff و Context
			return RetryWithContext(ctx, func() error {
				// درخواست واقعی
				req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
				if err != nil {
					return err
				}

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				if resp.StatusCode >= 500 {
					return fmt.Errorf("server error: %d", resp.StatusCode)
				}

				return nil
			}, s.retryConfig)
		})
	})
}

// ============================================================================
// بخش 10: مثال‌های عملی
// ============================================================================

func demonstrateRetry() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔄 RETRY DEMONSTRATION")
	fmt.Println(stringsRepeat("=", 80))

	// شبیه‌سازی یک عملیات که گاهی خطا می‌دهد
	attempt := 0
	operation := func() error {
		attempt++
		if attempt < 3 {
			return fmt.Errorf("temporary error (attempt %d)", attempt)
		}
		return nil
	}

	config := RetryConfig{
		MaxAttempts:   5,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      1 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
	}

	err := RetryWithBackoff(operation, config)
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	} else {
		fmt.Printf("Success after %d attempts\n", attempt)
	}
}

func demonstrateCircuitBreaker() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔌 CIRCUIT BREAKER DEMONSTRATION")
	fmt.Println(stringsRepeat("=", 80))

	config := CircuitBreakerConfig{
		FailureThreshold: 3,
		SuccessThreshold: 2,
		Timeout:          5 * time.Second,
		HalfOpenMaxCalls: 1,
	}

	cb := NewCircuitBreaker(config)

	// شبیه‌سازی یک سرویس معیوب
	failureCount := 0
	operation := func() error {
		failureCount++
		if failureCount <= 5 {
			return errors.New("service error")
		}
		return nil
	}

	for i := 0; i < 10; i++ {
		err := cb.Execute(operation)
		fmt.Printf("Request %d: state=%s, err=%v\n",
			i+1, cb.GetState(), err)
		time.Sleep(100 * time.Millisecond)
	}
}

func demonstrateResilientService() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🏭 RESILIENT SERVICE DEMONSTRATION")
	fmt.Println(stringsRepeat("=", 80))

	service := NewResilientService()

	// آدرس API که ممکن است خطا بدهد
	url := "http://httpbin.org/status/500"

	ctx := context.Background()
	err := service.CallExternalAPI(ctx, url)
	if err != nil {
		fmt.Printf("Final error: %v\n", err)
	} else {
		fmt.Println("Success!")
	}
}

// ============================================================================
// بخش 11: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 FAULT TOLERANCE BEST PRACTICES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ RETRY BEST PRACTICES                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│  • Set max attempts (3-5 for most cases)                                  │
│  • Use exponential backoff with jitter                                    │
│  • Don't retry on permanent errors (4xx, auth failures)                   │
│  • Always set timeout for each attempt                                    │
│  • Make retries idempotent                                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ CIRCUIT BREAKER BEST PRACTICES                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│  • Monitor state changes (log, metrics)                                   │
│  • Configure thresholds based on your SLOs                                │
│  • Use half-open state with small test calls                              │
│  • Combine with retry (retry before circuit opens)                        │
│  • Set appropriate timeouts                                               │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ RECOMMENDED VALUES                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│  • Max retries: 3-5                                                       │
│  • Initial delay: 100-500ms                                               │
│  • Max delay: 5-30s                                                       │
│  • Backoff factor: 2                                                      │
│  • Failure threshold: 5-10                                                │
│  • Circuit timeout: 30-60s                                                │
│  • Half-open test calls: 1-3                                              │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 12: Common Pitfalls
// ============================================================================

func commonPitfalls() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚠️ COMMON PITFALLS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ ❌ Retrying non-idempotent operations                                     │
│    • POST request without idempotency key                                 │
│    • ✅ Make operations idempotent or retry only safe methods              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ ❌ No jitter in backoff                                                   │
│    • All retries happen at same time (thundering herd)                    │
│    • ✅ Add random jitter to spread retries                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ ❌ Too short timeout                                                       │
│    • Operations timeout before retry can succeed                          │
│    • ✅ Set timeouts longer than expected operation time                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ ❌ Circuit breaker threshold too low                                      │
│    • Opens on normal transient failures                                   │
│    • ✅ Calibrate based on normal error rate                              │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 13: Main
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 FAULT TOLERANCE IN GO")
	fmt.Println("Retry | Backoff | Circuit Breaker | go-resiliency")
	fmt.Println(strings.Repeat("=", 80))

	bestPractices()
	commonPitfalls()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 Running Demonstrations")
	fmt.Println(strings.Repeat("=", 80))

	demonstrateRetry()
	demonstrateCircuitBreaker()
	demonstrateResilientService()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# Simple Retry
err := RetryWithBackoff(func() error {
    return callAPI()
}, DefaultRetryConfig())

# Circuit Breaker
cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())
err := cb.Execute(func() error {
    return unstableOperation()
})

# With Context
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
err := RetryWithContext(ctx, operation, config)

# Combined (Bulkhead + Circuit Breaker + Retry)
service := NewResilientService()
err := service.CallExternalAPI(ctx, url)

# Install go-resiliency
go get github.com/eapache/go-resiliency/retrier
go get github.com/eapache/go-resiliency/breaker
`)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎯 FAULT TOLERANCE - COMPLETE")
	fmt.Println("Build resilient Go applications!")
	fmt.Println(strings.Repeat("=", 80))
}

// تابع کمکی
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// اضافه کردن importهای لازم
var _ = sync.Mutex{}
var _ = atomic.Int64{}