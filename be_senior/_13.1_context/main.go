package _13_1_context

// ============================================================================
// FILE: context_guide.go
// TITLE: راهنمای کامل Context در Go - تک‌فایل اجراشونده با سکشن‌بندی
// HOW TO RUN: go run context_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Context چیست و چرا ساخته شد؟
// ============================================================================
//
// Context یک اینترفیس در Go است که سه کار اصلی انجام می‌دهد:
// 1. انتقال ددلاین (زمان پایان) و تایم‌اوت
// 2. انتقال سیگنال کنسل (cancel)
// 3. انتقال مقادیر کلید-مقدار در طول زنجیره فراخوانی
//
// هر درخواست HTTP یک Context می‌گیرد که وقتی درخواست کنسل می‌شود (مرورگر بسته شد)
// تمام گوروتین‌های آن درخواست هم کنسل می‌شوند.
//
// قانون طلایی: Context را هرگز در struct ذخیره نکنید، همیشه به عنوان اولین پارامتر تابع پاس دهید.
// ============================================================================

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ============================================================================
// بخش 1: انواع Contextهای خالی (Empty Contexts)
// ============================================================================

// demonstrateEmptyContexts نمایش دو نوع Context خالی
func demonstrateEmptyContexts() {
	// context.Background() - برای main، تست، و لایه‌های بالا
	// هرگز کنسل نمی‌شود، هیچ مقداری ندارد
	ctxBG := context.Background()
	fmt.Printf("Background: %T, Done: %v\n", ctxBG, ctxBG.Done())

	// context.TODO() - وقتی نمی‌دانی چه Context بدهی (مثل placeholder)
	// برای زمانی که قصد دارید بعداً آن را جایگزین کنید
	ctxTODO := context.TODO()
	fmt.Printf("TODO: %T, Err: %v\n", ctxTODO, ctxTODO.Err())

	// نکته: در کد نهایی، TODO نباید بماند
}

// ============================================================================
// بخش 2: Context با قابلیت Cancel (WithCancel)
// ============================================================================

// demonstrateCancel نمایش چگونگی کنسل کردن دستی یک Context
func demonstrateCancel() {
	fmt.Println("\n=== WithCancel Example ===")

	// parent context: Background
	parent := context.Background()

	// ctx: context جدید با قابلیت کنسل
	// cancel: تابعی که برای کنسل کردن صدا می‌زنیم
	ctx, cancel := context.WithCancel(parent)

	// نکته مهم: همیشه defer cancel() برای جلوگیری از نشت حافظه (goroutine leak)
	defer cancel()

	// اجرای یک گوروتین که منتظر سیگنال کنسل است
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("✅ Goroutine received cancel signal")
			fmt.Printf("   Reason: %v\n", ctx.Err())
		case <-time.After(5 * time.Second):
			fmt.Println("❌ This should not happen - timeout")
		}
	}()

	// صبر یک ثانیه سپس کنسل کردن
	time.Sleep(1 * time.Second)
	fmt.Println("👉 Calling cancel()...")
	cancel() // این خط باعث می‌شود ctx.Done() بسته شود

	// صبر می‌کنیم تا گوروتین چاپ کند
	time.Sleep(500 * time.Millisecond)
}

// ============================================================================
// بخش 3: Context با تایم‌اوت (WithTimeout)
// ============================================================================

// demonstrateTimeout نمایش تایم‌اوت خودکار
func demonstrateTimeout() {
	fmt.Println("\n=== WithTimeout Example ===")

	// Contextی که بعد از 2 ثانیه خودکار کنسل می‌شود
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // همیشه فراخوانی کن حتی اگر تایم‌اوت داشته باشی

	start := time.Now()

	// شبیه‌سازی یک عملیات طولانی
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Operation completed")
	case <-ctx.Done():
		elapsed := time.Since(start)
		fmt.Printf("✅ Operation cancelled after %v\n", elapsed)
		fmt.Printf("   Reason: %v\n", ctx.Err()) // context deadline exceeded
	}
}

// slowOperation یک تابع که کار طولانی انجام می‌دهد
// و از Context برای تشخیص کنسل شدن استفاده می‌کند
func slowOperation(ctx context.Context, duration time.Duration) error {
	select {
	case <-time.After(duration):
		return nil // عملیات با موفقیت انجام شد
	case <-ctx.Done():
		return ctx.Err() // کنسل شد یا تایم‌اوت خورد
	}
}

// demonstrateTimeoutInFunction استفاده از Context در توابع
func demonstrateTimeoutInFunction() {
	fmt.Println("\n=== WithTimeout in Function ===")

	// تایم‌اوت 1 ثانیه
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// تلاش برای عملیات 2 ثانیه‌ای
	err := slowOperation(ctx, 2*time.Second)
	if err != nil {
		fmt.Printf("✅ Function returned error: %v\n", err)
	} else {
		fmt.Println("Operation succeeded")
	}

	// تلاش با عملیات 500 میلی‌ثانیه‌ای (موفق)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()

	err = slowOperation(ctx2, 500*time.Millisecond)
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
	} else {
		fmt.Println("✅ Fast operation succeeded within timeout")
	}
}

// ============================================================================
// بخش 4: Context با ددلاین (WithDeadline)
// ============================================================================

// demonstrateDeadline نمایش ددلاین مطلق (زمان مشخص)
func demonstrateDeadline() {
	fmt.Println("\n=== WithDeadline Example ===")

	// ددلاین: 3 ثانیه بعد از الان
	deadline := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	start := time.Now()

	// حلقه‌ای که هر ثانیه چک می‌کند
	for {
		select {
		case <-ctx.Done():
			elapsed := time.Since(start)
			fmt.Printf("✅ Deadline reached after %v\n", elapsed)
			fmt.Printf("   Deadline was: %v\n", deadline.Format("15:04:05"))
			fmt.Printf("   Error: %v\n", ctx.Err())
			return
		default:
			fmt.Printf("   Working... (%v elapsed)\n", time.Since(start).Round(time.Second))
			time.Sleep(1 * time.Second)
		}
	}
}

// ============================================================================
// بخش 5: Context با مقدار (WithValue)
// ============================================================================

// type contextKey تعریف نوع اختصاصی برای کلیدهای Context
// نکته مهم: برای جلوگیری از collision، همیشه از نوع سفارشی استفاده کن
type contextKey string

// ثابت‌های کلیدها - به جای string مستقیم
const (
	KeyUserID    contextKey = "user_id"
	KeyRequestID contextKey = "request_id"
	KeyTraceID   contextKey = "trace_id"
)

// demonstrateValue نمایش عبور مقدار با Context
func demonstrateValue() {
	fmt.Println("\n=== WithValue Example ===")

	// ساختن زنجیره Context با مقادیر مختلف
	ctx := context.Background()

	// افزودن user_id
	ctx = context.WithValue(ctx, KeyUserID, 12345)

	// افزودن request_id (توجه: Context immutable است، مقدار جدید برمی‌گرداند)
	ctx = context.WithValue(ctx, KeyRequestID, "req-abc-123")

	// افزودن trace_id
	ctx = context.WithValue(ctx, KeyTraceID, "trace-xyz-789")

	// حالا توابع پایین‌دست می‌توانند این مقادیر را بخوانند
	processRequest(ctx)
}

// processRequest یک تابع که مقادیر را از Context می‌خواند
func processRequest(ctx context.Context) {
	// خواندن مقدار با type assertion
	userID, ok := ctx.Value(KeyUserID).(int)
	if !ok {
		fmt.Println("❌ User ID not found or wrong type")
	} else {
		fmt.Printf("✅ User ID: %d\n", userID)
	}

	requestID, ok := ctx.Value(KeyRequestID).(string)
	if !ok {
		fmt.Println("❌ Request ID not found")
	} else {
		fmt.Printf("✅ Request ID: %s\n", requestID)
	}

	traceID, ok := ctx.Value(KeyTraceID).(string)
	if !ok {
		fmt.Println("❌ Trace ID not found")
	} else {
		fmt.Printf("✅ Trace ID: %s\n", traceID)
	}

	// ❌ اشتباه: استفاده از string لخت به عنوان کلید
	// این کار ممکن است با کتابخانه‌های دیگر تداخل کند
	_ = ctx.Value("user_id") // این کار را نکن

	// ✅ درست: همیشه از نوع تعریف‌شده استفاده کن
	_ = ctx.Value(KeyUserID) // این درست است
}

// ============================================================================
// بخش 6: ترکیب قابلیت‌های Context (Cancel + Value + Timeout)
// ============================================================================

// apiCall شبیه‌سازی فراخوانی API با Context
func apiCall(ctx context.Context, id int) (string, error) {
	// خواندن مقادیر از Context
	requestID, _ := ctx.Value(KeyRequestID).(string)

	// چک کردن کنسل شدن قبل از شروع کار سنگین
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("api call cancelled before start: %w", ctx.Err())
	default:
	}

	// شبیه‌سازی کار سنگین
	select {
	case <-time.After(500 * time.Millisecond):
		return fmt.Sprintf("Result for id %d (request: %s)", id, requestID), nil
	case <-ctx.Done():
		return "", fmt.Errorf("api call interrupted: %w", ctx.Err())
	}
}

// demonstrateCombined ترکیب همه ویژگی‌های Context
func demonstrateCombined() {
	fmt.Println("\n=== Combined Context Example ===")

	// Context اصلی با Background
	baseCtx := context.Background()

	// افزودن مقادیر
	ctxWithValues := context.WithValue(baseCtx, KeyRequestID, "req-combined-001")
	ctxWithValues = context.WithValue(ctxWithValues, KeyUserID, 999)

	// افزودن تایم‌اوت 1 ثانیه
	ctx, cancel := context.WithTimeout(ctxWithValues, 1*time.Second)
	defer cancel()

	// فراخوانی API
	result, err := apiCall(ctx, 42)
	if err != nil {
		fmt.Printf("❌ API call failed: %v\n", err)
	} else {
		fmt.Printf("✅ API call succeeded: %s\n", result)
	}

	// تست با تایم‌اوت کمتر از زمان عملیات
	ctxShort, cancelShort := context.WithTimeout(ctxWithValues, 200*time.Millisecond)
	defer cancelShort()

	result, err = apiCall(ctxShort, 100)
	if err != nil {
		fmt.Printf("✅ Short timeout correctly cancelled: %v\n", err)
	} else {
		fmt.Printf("Result: %s\n", result)
	}
}

// ============================================================================
// بخش 7: استفاده از Context در HTTP Server (عملیاتی و واقعی)
// ============================================================================

// handleRequest یک هندلر HTTP که از Context استفاده می‌کند
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Context درخواست را بگیر (زمانی که مرورگر بسته شود، این Context کنسل می‌شود)
	ctx := r.Context()

	// افزودن مقدار به Context
	ctx = context.WithValue(ctx, KeyRequestID, "http-001")

	// انجام عملیات با Context
	result, err := doDatabaseQuery(ctx)
	if err != nil {
		// بررسی نوع خطا
		if err == context.DeadlineExceeded {
			http.Error(w, "Request timeout", http.StatusGatewayTimeout)
		} else if err == context.Canceled {
			http.Error(w, "Request cancelled", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "Result: %s\n", result)
}

// doDatabaseQuery شبیه‌سازی کوئری دیتابیس با Context
func doDatabaseQuery(ctx context.Context) (string, error) {
	// استفاده از select برای چک کردن Context
	select {
	case <-time.After(500 * time.Millisecond):
		return "query result", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// startHTTPServer راه‌اندازی سرور HTTP با هندلرهای آگاه از Context
func startHTTPServer() {
	http.HandleFunc("/", handleRequest)

	// سرور با timeoutهای مناسب
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// اجرا در گوروتین (برای non-blocking)
	go func() {
		fmt.Println("🌐 HTTP Server starting on :8080")
		fmt.Println("   Test with: curl http://localhost:8080/")
		fmt.Println("   Or press Ctrl+C to stop")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// صبر 2 ثانیه برای نمایش پیام
	time.Sleep(2 * time.Second)
}

// ============================================================================
// بخش 8: الگوهای درست و غلط با Context
// ============================================================================

// ✅ الگوی درست: Context به عنوان اولین پارامتر
func correctPattern(ctx context.Context, id int) error {
	// کار درست
	return nil
}

// ❌ الگوی غلط: Context در جای دیگر (دوم، سوم، یا داخل struct)
func wrongPattern(id int, ctx context.Context) error {
	// این کار را نکن
	return nil
}

// ❌ الگوی غلط: ذخیره Context در struct
type BadStruct struct {
	ctx context.Context // ❌ هرگز این کار را نکن
}

// ✅ الگوی درست: Context را از پارامتر دریافت کن
type GoodStruct struct {
	// هیچ فیلد ctx ای اینجا نگذار
}

func (g *GoodStruct) DoWork(ctx context.Context) error {
	// Context را به توابع پایین‌دست بفرست
	return nil
}

// ❌ الگوی غلط: نادیده گرفتن Done channel
func ignoreDone(ctx context.Context) {
	// کار طولانی بدون چک کردن ctx.Done()
	time.Sleep(10 * time.Second) // ❌ این ممکن است forever بماند
}

// ✅ الگوی درست: چک کردن ctx.Done() در حلقه‌های طولانی
func checkDoneCorrect(ctx context.Context) error {
	for i := 0; i < 1000; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// انجام کار کوچک
			time.Sleep(1 * time.Millisecond)
		}
	}
	return nil
}

// ✅ الگوی درست: همیشه defer cancel()
func alwaysDeferCancel() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // ✅ این خط حیاتی است

	// استفاده از ctx
	_ = ctx
}

// ============================================================================
// بخش 9: جلوگیری از نشت گوروتین (Goroutine Leak)
// ============================================================================

// leakExample این تابع نشت حافظه دارد - ❌ الگوی غلط
func leakExample() {
	ctx, cancel := context.WithCancel(context.Background())

	// گوروتینی که هیچ‌گاه بسته نمی‌شود (چون cancel فراخوانی نمی‌شود)
	go func() {
		<-ctx.Done()
		fmt.Println("This will never print")
	}()
	cancel()
	// فراموش کردن فراخوانی cancel()
	// _ = cancel // اگر این را نزنی، گوروتین forever می‌ماند
}

// ✅ الگوی درست: همیشه cancel را فراخوانی کن
func noLeakExample() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // ✅ این باعث می‌شود گوروتین بسته شود

	go func() {
		<-ctx.Done()
		fmt.Println("✅ Goroutine cleaned up properly")
	}()

	time.Sleep(100 * time.Millisecond)
}

// ============================================================================
// بخش 10: جمع‌بندی و نکات نهایی
// ============================================================================

func main() {
	fmt.Println("========== CONTEXT GUIDE IN GO ==========")
	fmt.Println("This file demonstrates all Context patterns")
	fmt.Println("Run each section to see live examples\n")

	// بخش 1: Empty Contexts
	demonstrateEmptyContexts()

	// بخش 2: WithCancel
	demonstrateCancel()

	// بخش 3: WithTimeout
	demonstrateTimeout()
	demonstrateTimeoutInFunction()

	// بخش 4: WithDeadline
	demonstrateDeadline()

	// بخش 5: WithValue
	demonstrateValue()

	// بخش 6: ترکیب قابلیت‌ها
	demonstrateCombined()

	// بخش 7: HTTP Server (اجرا در پس‌زمینه)
	startHTTPServer()

	// بخش 8: الگوهای درست و غلط (در کد بالا نشان داده شده)
	fmt.Println("\n=== Correct vs Wrong Patterns ===")
	fmt.Println("✅ Context as first parameter")
	fmt.Println("✅ Always defer cancel()")
	fmt.Println("✅ Check ctx.Done() in long operations")
	fmt.Println("❌ Never store Context in struct")
	fmt.Println("❌ Never use string as key (use custom type)")
	fmt.Println("❌ Never ignore ctx.Done()")

	// بخش 9: جلوگیری از نشت
	noLeakExample()

	// جمع‌بندی نهایی
	fmt.Println("\n========== SUMMARY ==========")
	fmt.Println("| Function              | Use Case                          |")
	fmt.Println("|-----------------------|-----------------------------------|")
	fmt.Println("| Background()          | Main, tests, top-level            |")
	fmt.Println("| TODO()                | Placeholder, future replacement   |")
	fmt.Println("| WithCancel()          | Manual cancellation               |")
	fmt.Println("| WithTimeout()         | Auto-cancel after duration        |")
	fmt.Println("| WithDeadline()        | Auto-cancel at absolute time      |")
	fmt.Println("| WithValue()           | Request-scoped values (trace ID)  |")
	fmt.Println("| Done()                | Channel that closes on cancel     |")
	fmt.Println("| Err()                 | Reason for cancellation           |")

	fmt.Println("\n========== GOLDEN RULES ==========")
	fmt.Println("1. Context is always the FIRST parameter")
	fmt.Println("2. NEVER store Context in structs")
	fmt.Println("3. ALWAYS call cancel() (use defer)")
	fmt.Println("4. Use custom type for keys (not string)")
	fmt.Println("5. Don't pass nil Context - use TODO() instead")
	fmt.Println("6. Check ctx.Done() in long-running operations")
	fmt.Println("7. Context is for cancellation and values, not for optional params")
}

// ============================================================================
// بخش 11: توابع کمکی برای استفاده در پروژه واقعی
// ============================================================================

// WithTimeoutOrDefault یک helper که اگر timeout صفر بود از مقدار پیش‌فرض استفاده می‌کند
func WithTimeoutOrDefault(parent context.Context, timeout time.Duration, defaultTimeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return context.WithTimeout(parent, timeout)
}

// GetRequestID خواندن Request ID از Context با fallback
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(KeyRequestID).(string); ok {
		return id
	}
	return "unknown"
}

// IsCancelled چک کردن اینکه آیا Context کنسل شده است
func IsCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
