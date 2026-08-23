
// ============================================================================
// FILE: context_management_guide.go
// TITLE: نقش Context در مدیریت گوروتین‌ها، Timeout و Cancel
// HOW TO RUN: go run context_management_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Context چیست و چرا برای مدیریت گوروتین ضروری است؟
// ============================================================================
//
// Context سه نقش اصلی در مدیریت گوروتین‌ها دارد:
//
// 1. کنسل کردن (Cancellation):
//    - ارسال سیگنال "بس کن" به همه گوروتین‌های یک درخت
//    - جلوگیری از نشت گوروتین (goroutine leak)
//    - آزادسازی منابع وقتی دیگر نیازی نیست
//
// 2. تایم‌اوت (Timeout):
//    - محدود کردن زمان اجرای عملیات
//    - محافظت در برابر عملیات‌های بینهایت
//    - بهبود تجربه کاربری با پاسخ به موقع
//
// 3. ددلاین (Deadline):
//    - تعیین زمان مطلق پایان کار (مثل "تا ساعت 14:30 تمام کن")
//    - هماهنگی با مهلت‌های خارجی (API limits, SLA)
//
// قانون طلایی: Context را به عنوان اولین پارامتر به هر تابعی که ممکن است
// مسدود شود یا زمان ببرد، پاس بدهید.
// ============================================================================

package _13_2_context_management

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// بخش 1: اصول پایه Context و مدیریت گوروتین
// ============================================================================

func demonstrateBasicContext() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 CONTEXT BASICS - Goroutine Management")
	fmt.Println(stringsRepeat("=", 80))

	// context.Background(): ریشه درخت Context
	ctx := context.Background()
	fmt.Printf("Background context: %T\n", ctx)

	// context.TODO(): وقتی نمی‌دانی چه بدهی (placeholder)
	ctxTODO := context.TODO()
	fmt.Printf("TODO context: %T\n", ctxTODO)

	// WithCancel: ایجاد Context با قابلیت کنسل کردن
	ctxCancel, cancel := context.WithCancel(ctx)
	defer cancel() // همیشه defer cancel() برای جلوگیری از نشت

	fmt.Printf("Cancel context: %T, cancel func: %T\n", ctxCancel, cancel)
}

// demonstrateContextCancellation نقش Context در کنسل کردن گوروتین‌ها
func demonstrateContextCancellation() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🛑 CONTEXT CANCELLATION - Stopping Goroutines")
	fmt.Println(stringsRepeat("=", 80))

	// ایجاد Context قابل کنسل
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	// راه‌اندازی چندین گوروتین که همه به Context گوش می‌دهند
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go workerWithContext(i, ctx, &wg)
	}

	// بگذار 500ms کار کنند
	time.Sleep(500 * time.Millisecond)

	fmt.Println("\n👉 Main: Sending cancel signal to all workers...")
	cancel() // این یک خط همه گوروتین‌ها را متوقف می‌کند

	wg.Wait()
	fmt.Println("✅ All workers stopped cleanly")
}

// workerWithContext یک گوروتین که به Context گوش می‌دهد
func workerWithContext(id int, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			// سیگنال کنسل شدن دریافت شد
			fmt.Printf("  Worker %d: received cancel signal (%v)\n", id, ctx.Err())
			return
		default:
			// انجام کار معمولی
			fmt.Printf("  Worker %d: working...\n", id)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// ============================================================================
// بخش 2: Context با Timeout - محدود کردن زمان اجرا
// ============================================================================

func demonstrateContextTimeout() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⏰ CONTEXT TIMEOUT - Limiting Execution Time")
	fmt.Println(stringsRepeat("=", 80))

	// عملیاتی که 2 ثانیه طول می‌کشد
	slowOperation := func(ctx context.Context) error {
		select {
		case <-time.After(2 * time.Second):
			return nil // عملیات موفق
		case <-ctx.Done():
			return ctx.Err() // تایم‌اوت یا کنسل شد
		}
	}

	// تایم‌اوت 1 ثانیه (کمتر از زمان عملیات)
	fmt.Println("Scenario 1: Timeout 1s (operation takes 2s)")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel1()

	start := time.Now()
	err := slowOperation(ctx1)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("  ❌ Operation failed: %v (after %v)\n", err, elapsed)
	} else {
		fmt.Printf("  ✅ Operation succeeded after %v\n", elapsed)
	}

	// تایم‌اوت 3 ثانیه (بیشتر از زمان عملیات)
	fmt.Println("\nScenario 2: Timeout 3s (operation takes 2s)")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel2()

	start = time.Now()
	err = slowOperation(ctx2)
	elapsed = time.Since(start)

	if err != nil {
		fmt.Printf("  ❌ Operation failed: %v (after %v)\n", err, elapsed)
	} else {
		fmt.Printf("  ✅ Operation succeeded after %v\n", elapsed)
	}
}

// demonstrateContextTimeoutChain زنجیره تایم‌اوت در توابع تو در تو
func demonstrateContextTimeoutChain() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔗 CONTEXT TIMEOUT CHAIN - Propagating Deadlines")
	fmt.Println(stringsRepeat("=", 80))

	// لایه 1: تابع اصلی با تایم‌اوت 3 ثانیه
	mainFunction := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		fmt.Println("Main: Starting with 3s timeout")
		return businessLogic(ctx)
	}

	// لایه 2: منطق کسب و کار
	businessLogic := func(ctx context.Context) error {
		fmt.Println("  BusinessLogic: Received context with deadline")

		// می‌توانیم deadline را ببینیم
		if deadline, ok := ctx.Deadline(); ok {
			fmt.Printf("    Deadline: %v\n", deadline)
		}

		// فراخوانی لایه پایین‌تر
		return databaseQuery(ctx)
	}

	// لایه 3: کوئری دیتابیس
	databaseQuery := func(ctx context.Context) error {
		fmt.Println("    DatabaseQuery: Starting query...")

		// شبیه‌سازی کوئری سنگین
		select {
		case <-time.After(2 * time.Second):
			fmt.Println("    DatabaseQuery: Query completed")
			return nil
		case <-ctx.Done():
			fmt.Printf("    DatabaseQuery: Cancelled: %v\n", ctx.Err())
			return ctx.Err()
		}
	}

	err := mainFunction()
	if err != nil {
		fmt.Printf("\n❌ Error: %v\n", err)
	} else {
		fmt.Println("\n✅ Success")
	}
}

// ============================================================================
// بخش 3: Context با Deadline (زمان مطلق)
// ============================================================================

func demonstrateContextDeadline() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📅 CONTEXT DEADLINE - Absolute Time Limits")
	fmt.Println(stringsRepeat("=", 80))

	// تنظیم ددلاین برای 2 ثانیه بعد
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	fmt.Printf("Current time: %v\n", time.Now().Format("15:04:05"))
	fmt.Printf("Deadline set: %v\n", deadline.Format("15:04:05"))

	// حلقه‌ای که هر 500ms چک می‌کند
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\n✅ Deadline reached: %v\n", ctx.Err())
			return
		default:
			remaining := time.Until(deadline)
			fmt.Printf("  Working... %v remaining\n", remaining.Round(time.Millisecond))
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// demonstrateDeadlineVsTimeout مقایسه Deadline و Timeout
func demonstrateDeadlineVsTimeout() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚖️ DEADLINE vs TIMEOUT")
	fmt.Println(stringsRepeat("=", 80))

	// Timeout: نسبی (از الان تا X ثانیه)
	fmt.Println("Timeout: relative time from now")
	ctxTimeout, _ := context.WithTimeout(context.Background(), 5*time.Second)
	deadlineTimeout, _ := ctxTimeout.Deadline()
	fmt.Printf("  Timeout 5s from now → deadline at %v\n", deadlineTimeout.Format("15:04:05"))

	// Deadline: مطلق (ساعت مشخص)
	fmt.Println("\nDeadline: absolute time")
	fixedDeadline := time.Now().Add(5 * time.Second)
	ctxDeadline, _ := context.WithDeadline(context.Background(), fixedDeadline)
	deadlineDeadline, _ := ctxDeadline.Deadline()
	fmt.Printf("  Fixed deadline at %v\n", deadlineDeadline.Format("15:04:05"))

	fmt.Println("\n💡 When to use which:")
	fmt.Println("  • Timeout: عملیات با زمان مشخص (API call, DB query)")
	fmt.Println("  • Deadline: هماهنگی با مهلت‌های خارجی (SLA, batch job)")
}

// ============================================================================
// بخش 4: Context در Pipeline (انتقال در زنجیره)
// ============================================================================

func demonstrateContextInPipeline() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔗 CONTEXT IN PIPELINE - Cancellable Data Processing")
	fmt.Println(stringsRepeat("=", 80))

	// Stage 1: تولید داده
	generate := func(ctx context.Context, numbers []int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range numbers {
				select {
				case out <- n:
				case <-ctx.Done():
					fmt.Println("  Generator: cancelled")
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}()
		return out
	}

	// Stage 2: ضرب در 2
	multiply := func(ctx context.Context, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				select {
				case out <- n * 2:
				case <-ctx.Done():
					fmt.Println("  Multiplier: cancelled")
					return
				}
			}
		}()
		return out
	}

	// Stage 3: فیلتر (حذف اعداد بزرگتر از 10)
	filter := func(ctx context.Context, in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				if n > 10 {
					continue
				}
				select {
				case out <- n:
				case <-ctx.Done():
					fmt.Println("  Filter: cancelled")
					return
				}
			}
		}()
		return out
	}

	// اجرای Pipeline با قابلیت کنسل شدن
	ctx, cancel := context.WithTimeout(context.Background(), 450*time.Millisecond)
	defer cancel()

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	fmt.Println("Starting pipeline with 450ms timeout...")
	start := time.Now()

	stage1 := generate(ctx, numbers)
	stage2 := multiply(ctx, stage1)
	stage3 := filter(ctx, stage2)

	results := make([]int, 0)
	for result := range stage3 {
		results = append(results, result)
		fmt.Printf("  Result: %d\n", result)
	}

	elapsed := time.Since(start)
	fmt.Printf("\nPipeline finished after %v, got %d results\n", elapsed, len(results))
}

// ============================================================================
// بخش 5: Context با Value (عبور داده در درخت گوروتین)
// ============================================================================

type contextKey string

const (
	KeyRequestID contextKey = "request_id"
	KeyUserID    contextKey = "user_id"
	KeyTraceID   contextKey = "trace_id"
)

func demonstrateContextValue() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📦 CONTEXT WITH VALUE - Passing Data in Goroutine Tree")
	fmt.Println(stringsRepeat("=", 80))

	// ایجاد Context با مقادیر مختلف
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyRequestID, "req-123-abc")
	ctx = context.WithValue(ctx, KeyUserID, 42)
	ctx = context.WithValue(ctx, KeyTraceID, "trace-xyz-789")

	// توابع تو در تو که مقادیر را منتقل می‌کنند
	handler := func(ctx context.Context) {
		requestID := ctx.Value(KeyRequestID).(string)
		userID := ctx.Value(KeyUserID).(int)
		traceID := ctx.Value(KeyTraceID).(string)

		fmt.Printf("  Handler: request_id=%s, user_id=%d, trace_id=%s\n",
			requestID, userID, traceID)

		service := func(ctx context.Context) {
			// مقادیر همچنان در دسترس هستند
			reqID := ctx.Value(KeyRequestID).(string)
			fmt.Printf("    Service: received request_id=%s\n", reqID)

			repo := func(ctx context.Context) {
				// و اینجا هم
				userID := ctx.Value(KeyUserID).(int)
				fmt.Printf("      Repository: querying for user_id=%d\n", userID)
			}
			repo(ctx)
		}
		service(ctx)
	}

	handler(ctx)

	fmt.Println("\n⚠️  Important notes about Context values:")
	fmt.Println("  • Only use for request-scoped data (trace ID, user ID)")
	fmt.Println("  • Never use for optional parameters")
	fmt.Println("  • Values are immutable (always create new Context)")
	fmt.Println("  • Use custom types for keys to avoid collisions")
}

// ============================================================================
// بخش 6: جلوگیری از نشت گوروتین (Goroutine Leak Prevention)
// ============================================================================

func demonstrateGoroutineLeakPrevention() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🛡️ GOROUTINE LEAK PREVENTION with Context")
	fmt.Println(stringsRepeat("=", 80))

	// ❌ مثال نشت گوروتین
	fmt.Println("❌ Bad: Goroutine leak (no cancellation)")
	leakyFunc := func() {
		ch := make(chan int)
		go func() {
			// این گوروتین هیچ‌گاه خارج نمی‌شود
			value := <-ch // منتظر می‌ماند برای همیشه
			fmt.Println(value)
		}()
		// تابع برمی‌گردد ولی گوروتین هنوز زنده است
	}
	leakyFunc()
	time.Sleep(10 * time.Millisecond)
	fmt.Printf("  Goroutines after leaky call: %d\n", runtime.NumGoroutine())

	// ✅ درست: با Context
	fmt.Println("\n✅ Good: Context prevents leak")
	properFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		ch := make(chan int)
		go func() {
			select {
			case value := <-ch:
				fmt.Println(value)
			case <-ctx.Done():
				fmt.Println("  Goroutine: cancelled, exiting cleanly")
				return
			}
		}()

		// تابع برمی‌گردد، اما گوروتین بعد از timeout خارج می‌شود
	}
	properFunc()

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("  Goroutines after proper call: %d\n", runtime.NumGoroutine())
}

// ============================================================================
// بخش 7: الگوهای ترکیبی (Cancellation + Timeout + Value)
// ============================================================================

func demonstrateCombinedPatterns() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎭 COMBINED PATTERNS - Full Power of Context")
	fmt.Println(stringsRepeat("=", 80))

	// API Client شبیه‌سازی شده
	type APIClient struct {
		name string
	}

	func (c *APIClient) Call(ctx context.Context, endpoint string) (string, error) {
		// خواندن trace ID از Context
		traceID, _ := ctx.Value(KeyTraceID).(string)

		select {
		case <-time.After(100 * time.Millisecond):
			return fmt.Sprintf("%s response from %s (trace: %s)", c.name, endpoint, traceID), nil
		case <-ctx.Done():
			return "", fmt.Errorf("%s: %w", c.name, ctx.Err())
		}
	}

	// Orchestrator که چند API را همزمان فراخوانی می‌کند
	orchestrator := func(ctx context.Context, endpoints []string) ([]string, error) {
		results := make([]string, len(endpoints))
		var wg sync.WaitGroup

		// ایجاد Context با تایم‌اوت برای کل عملیات
		ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer cancel()

		clients := []APIClient{{name: "API-1"}, {name: "API-2"}, {name: "API-3"}}

		for i, endpoint := range endpoints {
			wg.Add(1)
			go func(idx int, ep string) {
				defer wg.Done()

				// هر client با Context خودش (ارث‌بری از parent)
				client := clients[idx%len(clients)]
				result, err := client.Call(ctx, ep)
				if err != nil {
					fmt.Printf("  Error calling %s: %v\n", ep, err)
					return
				}
				results[idx] = result
			}(i, endpoint)
		}

		wg.Wait()

		// بررسی اینکه آیا به دلیل timeout کنسل شده
		if ctx.Err() != nil {
			return results, fmt.Errorf("operation cancelled: %w", ctx.Err())
		}

		return results, nil
	}

	// اجرا
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyTraceID, "trace-orchestrator-001")
	ctx = context.WithValue(ctx, KeyRequestID, "req-orch-001")

	endpoints := []string{"/users", "/products", "/orders", "/payments", "/reviews"}

	fmt.Println("Calling multiple APIs with 500ms timeout...")
	results, err := orchestrator(ctx, endpoints)

	if err != nil {
		fmt.Printf("❌ Orchestration failed: %v\n", err)
	}

	fmt.Println("\nResults:")
	for i, r := range results {
		if r != "" {
			fmt.Printf("  %d: %s\n", i, r)
		}
	}
}

// ============================================================================
// بخش 8: Context در HTTP Server (مدیریت درخواست‌ها)
// ============================================================================

func demonstrateContextInHTTPServer() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🌐 CONTEXT IN HTTP SERVER - Request Management")
	fmt.Println(stringsRepeat("=", 80))

	// شبیه‌سازی یک HTTP handler
	type Request struct {
		ID string
	}

	type Response struct {
		Status int
		Body   string
	}

	// Handler که از Context استفاده می‌کند
	handler := func(ctx context.Context, req Request) Response {
		// خواندن مقادیر از Context
		requestID, _ := ctx.Value(KeyRequestID).(string)

		fmt.Printf("  Handler: processing request %s (req_id=%s)\n", req.ID, requestID)

		// شبیه‌سازی کار سنگین
		select {
		case <-time.After(200 * time.Millisecond):
			return Response{Status: 200, Body: "OK"}
		case <-ctx.Done():
			fmt.Printf("  Handler: request %s cancelled: %v\n", req.ID, ctx.Err())
			return Response{Status: 499, Body: "Client Closed Request"}
		}
	}

	// Middleware که Context را آماده می‌کند
	middleware := func(next func(ctx context.Context, req Request) Response) func(ctx context.Context, req Request) Response {
		return func(ctx context.Context, req Request) Response {
			// افزودن مقادیر به Context
			ctx = context.WithValue(ctx, KeyRequestID, "mid-"+req.ID)
			ctx = context.WithValue(ctx, KeyTraceID, fmt.Sprintf("trace-%d", time.Now().UnixNano()))

			// افزودن تایم‌اوت 150ms
			ctx, cancel := context.WithTimeout(ctx, 150*time.Millisecond)
			defer cancel()

			return next(ctx, req)
		}
	}

	// اجرا
	handlerWithMiddleware := middleware(handler)

	requests := []Request{
		{ID: "req-1"},
		{ID: "req-2"},
		{ID: "req-3"},
	}

	for _, req := range requests {
		fmt.Printf("\nIncoming request: %s\n", req.ID)
		ctx := context.Background()
		resp := handlerWithMiddleware(ctx, req)
		fmt.Printf("  Response: status=%d, body=%s\n", resp.Status, resp.Body)
	}
}

// ============================================================================
// بخش 9: Context در Database Operations
// ============================================================================

func demonstrateContextInDatabase() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🗄️ CONTEXT IN DATABASE - Query with Timeout")
	fmt.Println(stringsRepeat("=", 80))

	// شبیه‌سازی یک کوئری دیتابیس
	type DBQuery struct {
		name     string
		duration time.Duration
	}

	executeQuery := func(ctx context.Context, query DBQuery) (string, error) {
		select {
		case <-time.After(query.duration):
			return fmt.Sprintf("Result from %s", query.name), nil
		case <-ctx.Done():
			return "", fmt.Errorf("query %s cancelled: %w", query.name, ctx.Err())
		}
	}

	// کوئری با تایم‌اوت
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	queries := []DBQuery{
		{name: "users", duration: 100 * time.Millisecond},
		{name: "orders", duration: 200 * time.Millisecond},
		{name: "products", duration: 400 * time.Millisecond}, // این یکی تایم‌اوت می‌خورد
		{name: "reviews", duration: 150 * time.Millisecond},
	}

	var wg sync.WaitGroup
	results := make(chan string, len(queries))

	for _, q := range queries {
		wg.Add(1)
		go func(query DBQuery) {
			defer wg.Done()
			result, err := executeQuery(ctx, query)
			if err != nil {
				fmt.Printf("  ❌ %s\n", err)
				return
			}
			results <- result
		}(q)
	}

	// منتظر اتمام همه یا timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		close(results)
		fmt.Println("\n✅ All queries completed (or failed)")
	case <-ctx.Done():
		fmt.Println("\n⚠️  Timeout reached, some queries cancelled")
	}

	// جمع‌آوری نتایج
	for result := range results {
		fmt.Printf("  📥 %s\n", result)
	}
}

// ============================================================================
// بخش 10: الگوهای درست و غلط با Context
// ============================================================================

func demonstrateCorrectVsWrong() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("✅❌ CONTEXT - CORRECT VS WRONG PATTERNS")
	fmt.Println(stringsRepeat("=", 80))

	// ✅ درست: Context به عنوان اولین پارامتر
	fmt.Println("\n✅ Correct: Context as first parameter")
	func correct(ctx context.Context, id int) error {
		return nil
	}

	// ❌ غلط: Context در جای دیگر
	fmt.Println("\n❌ Wrong: Context not as first parameter")
	// func wrong(id int, ctx context.Context) error // این را نکن

	// ✅ درست: استفاده از defer cancel()
	fmt.Println("\n✅ Correct: Always defer cancel()")
	func withDefer() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() // این خط را هرگز فراموش نکن
		_ = ctx
	}

	// ❌ غلط: فراموش کردن cancel (نشت حافظه)
	fmt.Println("\n❌ Wrong: Forgetting to call cancel (memory leak)")
	// func noDefer() {
	//     ctx, cancel := context.WithCancel(context.Background())
	//     _ = ctx
	//     // cancel() فراموش شده
	// }

	// ✅ درست: بررسی ctx.Done() در حلقه‌ها
	fmt.Println("\n✅ Correct: Check ctx.Done() in loops")
	func checkDone(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// انجام کار
			}
		}
	}

	// ❌ غلط: نادیده گرفتن ctx.Done()
	fmt.Println("\n❌ Wrong: Ignoring ctx.Done()")
	// func ignoreDone(ctx context.Context) {
	//     for {
	//         // کار طولانی بدون چک کردن ctx.Done()
	//     }
	// }

	// ✅ درست: استفاده از نوع سفارشی برای کلیدها
	fmt.Println("\n✅ Correct: Custom type for context keys")
	type myKey string
	const keyUser myKey = "user"

	// ❌ غلط: استفاده از string لخت به عنوان کلید
	fmt.Println("\n❌ Wrong: String literal as key (collision risk)")
	// ctx = context.WithValue(ctx, "user", value) // این را نکن
}

// ============================================================================
// بخش 11: توابع کمکی و ابزارها
// ============================================================================

// ContextWithDefaultTimeout ایجاد Context با تایم‌اوت پیش‌فرض
func ContextWithDefaultTimeout(parent context.Context, timeout, defaultTimeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return context.WithTimeout(parent, timeout)
}

// IsCancelled بررسی اینکه آیا Context کنسل شده است
func IsCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// GetValueOrDefault خواندن مقدار از Context با مقدار پیش‌فرض
func GetValueOrDefault[T any](ctx context.Context, key interface{}, defaultValue T) T {
	if val := ctx.Value(key); val != nil {
		if typed, ok := val.(T); ok {
			return typed
		}
	}
	return defaultValue
}

// MergeContexts ادغام دو Context (برای موارد خاص)
func MergeContexts(ctx1, ctx2 context.Context) context.Context {
	merged, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-ctx1.Done():
			cancel()
		case <-ctx2.Done():
			cancel()
		}
	}()

	return merged
}

// ============================================================================
// بخش 12: جمع‌بندی نهایی
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 CONTEXT IN GO - Complete Guide for Goroutine Management")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: مبانی
	demonstrateBasicContext()
	demonstrateContextCancellation()

	// بخش 2: Timeout
	demonstrateContextTimeout()
	demonstrateContextTimeoutChain()

	// بخش 3: Deadline
	demonstrateContextDeadline()
	demonstrateDeadlineVsTimeout()

	// بخش 4: Pipeline
	demonstrateContextInPipeline()

	// بخش 5: Value
	demonstrateContextValue()

	// بخش 6: جلوگیری از نشت
	demonstrateGoroutineLeakPrevention()

	// بخش 7: ترکیبی
	demonstrateCombinedPatterns()

	// بخش 8: HTTP Server
	demonstrateContextInHTTPServer()

	// بخش 9: Database
	demonstrateContextInDatabase()

	// بخش 10: درست و غلط
	demonstrateCorrectVsWrong()

	// بخش 11: جمع‌بندی
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 CONTEXT QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FUNCTION              │ PURPOSE                                      │")
	fmt.Println("├───────────────────────┼──────────────────────────────────────────────┤")
	fmt.Println("│ context.Background()  │ Root context for main, tests, init          │")
	fmt.Println("│ context.TODO()        │ Placeholder when unsure                     │")
	fmt.Println("│ WithCancel(parent)    │ Manual cancellation                         │")
	fmt.Println("│ WithTimeout(parent,d) │ Auto-cancel after duration                  │")
	fmt.Println("│ WithDeadline(parent,t)│ Auto-cancel at absolute time                │")
	fmt.Println("│ WithValue(parent,k,v) │ Request-scoped values (trace ID, user ID)   │")
	fmt.Println("│ ctx.Done()            │ Channel that closes on cancel               │")
	fmt.Println("│ ctx.Err()             │ Reason for cancellation                     │")
	fmt.Println("│ ctx.Deadline()        │ Returns deadline (if any)                   │")
	fmt.Println("└───────────────────────┴──────────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES FOR CONTEXT:")
	fmt.Println("  1. Context is always the FIRST parameter of a function")
	fmt.Println("  2. NEVER store Context in structs - pass it explicitly")
	fmt.Println("  3. ALWAYS call cancel() - use defer cancel()")
	fmt.Println("  4. Use custom types for context keys (not string literals)")
	fmt.Println("  5. Don't pass nil Context - use TODO() instead")
	fmt.Println("  6. Check ctx.Done() in long-running operations")
	fmt.Println("  7. Context is for cancellation and request-scoped data only")
	fmt.Println("  8. Don't use Context for optional parameters")
	fmt.Println("  9. Propagate Context through your call chain")
	fmt.Println("  10. When in doubt, add Context (you can always add it later)")

	fmt.Println("\n⚠️  COMMON PITFALLS:")
	fmt.Println("  • Forgetting to call cancel() → goroutine leak")
	fmt.Println("  • Storing Context in structs → coupling and bugs")
	fmt.Println("  • Using string literal as key → collisions")
	fmt.Println("  • Ignoring ctx.Done() in loops → unresponsive goroutines")
	fmt.Println("  • Passing nil Context → use TODO() instead")
	fmt.Println("  • Using Context for business logic data → misuse")
}

// ============================================================================
// بخش 13: توابع کمکی (runtime)
// ============================================================================

// runtime.NumGoroutine شبیه‌سازی (در فایل واقعی از import "runtime" استفاده کنید)
func numGoroutine() int {
	return 0 // Placeholder
}

var runtime = struct {
	NumGoroutine func() int
}{
	NumGoroutine: numGoroutine,
}

// stringsRepeat (در فایل واقعی از import "strings" استفاده کنید)
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}