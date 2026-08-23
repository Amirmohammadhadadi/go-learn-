// ============================================================================
// FILE: goroutines_management_guide.go
// TITLE: مدیریت پیشرفته گوروتین‌ها در Go - نحوه کنترل، نظارت و جلوگیری از نشت
// HOW TO RUN: go run goroutines_management_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - چرا مدیریت گوروتین مهم است؟
// ============================================================================
//
// مشکلات رایج با گوروتین‌ها:
// 1. نشت گوروتین (Goroutine Leak): گوروتین‌هایی که هیچ‌گاه تمام نمی‌شوند
// 2. گوروتین‌های بی‌کنترل (Runaway goroutines): مصرف بی‌رویه CPU
// 3. شرایط رقابتی (Race Conditions): دسترسی همزمان به داده‌های مشترک
// 4. Deadlock: گوروتین‌هایی که منتظر یکدیگر می‌مانند
//
// در این فایل: تکنیک‌های حرفه‌ای برای مدیریت، نظارت و کنترل گوروتین‌ها
// ============================================================================

package _9_goroutines_management

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// بخش 1: ردیابی تعداد گوروتین‌های فعال
// ============================================================================

// getGoroutineCount تعداد گوروتین‌های فعال را برمی‌گرداند
func getGoroutineCount() int {
	return runtime.NumGoroutine()
}

// printGoroutineCount تعداد گوروتین‌ها را به صورت فرمت شده چاپ می‌کند
func printGoroutineCount(label string) {
	fmt.Printf("📊 [%s] Active goroutines: %d\n", label, getGoroutineCount())
}

// demonstrateGoroutineTracking نمایش ردیابی تعداد گوروتین‌ها
func demonstrateGoroutineTracking() {
	fmt.Println("\n=== Tracking Goroutine Count ===")

	printGoroutineCount("start")

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
		}(i)
	}

	printGoroutineCount("after starting 5 goroutines")

	wg.Wait()
	printGoroutineCount("after all finished")
}

// ============================================================================
// بخش 2: شناسایی و رفع نشت گوروتین (Goroutine Leak)
// ============================================================================

// ❌ الگوی غلط: نشت گوروتین
func leakingGoroutine() {
	ch := make(chan int)

	go func() {
		// این گوروتین هیچ‌گاه تمام نمی‌شود چون کانال هیچ‌گاه داده نمی‌گیرد
		value := <-ch // ❌ برای همیشه منتظر می‌ماند
		fmt.Println(value)
	}()

	// فراموش کردیم به کانال داده بفرستیم یا آن را ببندیم
	// تابع برمی‌گردد اما گوروتین هنوز زنده است ← نشت حافظه
}

// ✅ الگوی درست: استفاده از context برای جلوگیری از نشت
func nonLeakingGoroutine(ctx context.Context) {
	ch := make(chan int)

	go func() {
		select {
		case value := <-ch:
			fmt.Println(value)
		case <-ctx.Done():
			fmt.Println("Goroutine cancelled, exiting cleanly")
			return
		}
	}()

	// اگر کاری نشد، context کنسل می‌شود و گوروتین خارج می‌شود
}

// ✅ الگوی درست: کانال را با ظرفیت مناسب یا بستن به موقع
func properChannelUsage() {
	ch := make(chan int, 1) // بافر دار

	go func() {
		ch <- 42
		close(ch) // بستن کانال
	}()

	value, ok := <-ch
	if ok {
		fmt.Printf("Received: %d\n", value)
	}
}

// demonstrateGoroutineLeakDetection نمایش تشخیص نشت
func demonstrateGoroutineLeakDetection() {
	fmt.Println("\n=== Goroutine Leak Detection ===")

	printGoroutineCount("before leak example")

	// این تابع نشت دارد - در کد واقعی نباید استفاده شود
	// برای نمایش، فقط یک بار اجرا می‌کنیم
	func() {
		ch := make(chan int)
		go func() {
			<-ch // منتظر می‌ماند
		}()
		// کانال را نمی‌بندیم و داده نمی‌فرستیم
		time.Sleep(10 * time.Millisecond)
	}()

	// اجازه می‌دهیم گوروتین اجرا شود
	time.Sleep(50 * time.Millisecond)

	printGoroutineCount("after leak example (one goroutine leaked)")

	fmt.Println("💡 راه حل: استفاده از context.WithTimeout یا بستن کانال")
}

// ============================================================================
// بخش 3: کنترل گوروتین‌ها با Context (کنسل کردن گروهی)
// ============================================================================

// workerWithContext یک گوروتین که به context گوش می‌دهد
func workerWithContext(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("🛑 Worker %d: cancelled (%v)\n", id, ctx.Err())
			return
		default:
			// انجام کار
			fmt.Printf("🔧 Worker %d: working...\n", id)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// demonstrateContextCancellation کنسل کردن گروهی گوروتین‌ها
func demonstrateContextCancellation() {
	fmt.Println("\n=== Context Cancellation (Group Cancel) ===")

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// اجرای 3 گوروتین
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go workerWithContext(ctx, i, &wg)
	}

	// بگذار 500 میلی‌ثانیه کار کنند
	time.Sleep(500 * time.Millisecond)

	// کنسل کردن همه
	fmt.Println("👉 Cancelling all workers...")
	cancel()

	wg.Wait()
	fmt.Println("✅ All workers cleaned up")
}

// ============================================================================
// بخش 4: کنترل گوروتین‌ها با تایم‌اوت سراسری
// ============================================================================

// demonstrateGlobalTimeout استفاده از تایم‌اوت برای کل گوروتین‌ها
func demonstrateGlobalTimeout() {
	fmt.Println("\n=== Global Timeout ===")

	// تایم‌اوت کلی: 1 ثانیه
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// 5 گوروتین که هر کدام 300 میلی‌ثانیه کار می‌کنند
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			select {
			case <-time.After(300 * time.Millisecond):
				fmt.Printf("Worker %d: completed work\n", id)
			case <-ctx.Done():
				fmt.Printf("Worker %d: timed out (%v)\n", id, ctx.Err())
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("✅ All workers finished or timed out")
}

// ============================================================================
// بخش 5: محدود کردن تعداد همزمان گوروتین‌ها (Semaphore Pattern)
// ============================================================================

// Semaphore یک ساختار ساده برای محدود کردن همزمانی
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore ایجاد سمافور با ظرفیت مشخص
func NewSemaphore(maxConcurrent int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, maxConcurrent),
	}
}

// Acquire گرفتن یک جایگاه (اگر پر باشد بلاک می‌شود)
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Release آزاد کردن یک جایگاه
func (s *Semaphore) Release() {
	<-s.ch
}

// demonstrateSemaphore محدود کردن همزمانی با سمافور
func demonstrateSemaphore() {
	fmt.Println("\n=== Semaphore (Limiting Concurrency) ===")

	sem := NewSemaphore(3) // حداکثر 3 گوروتین همزمان
	var wg sync.WaitGroup

	start := time.Now()

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			sem.Acquire()
			defer sem.Release()

			// شبیه‌سازی کار
			fmt.Printf("Worker %d: started\n", id)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Worker %d: finished\n", id)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("✅ All 10 workers completed in %v (max 3 concurrent)\n", elapsed)
}

// ============================================================================
// بخش 6: Worker Pool با اندازه ثابت (Fixed Worker Pool)
// ============================================================================

// FixedWorkerPool یک pool با تعداد ثابت worker
type FixedWorkerPool struct {
	jobQueue chan func()
	wg       sync.WaitGroup
}

// NewFixedWorkerPool ایجاد pool جدید
func NewFixedWorkerPool(workerCount int, queueSize int) *FixedWorkerPool {
	pool := &FixedWorkerPool{
		jobQueue: make(chan func(), queueSize),
	}

	// راه‌اندازی workers
	for i := 0; i < workerCount; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

func (p *FixedWorkerPool) worker(id int) {
	defer p.wg.Done()
	for job := range p.jobQueue {
		fmt.Printf("Worker %d: executing job\n", id)
		job()
	}
}

// Submit ارسال job به pool
func (p *FixedWorkerPool) Submit(job func()) {
	p.jobQueue <- job
}

// Stop توقف pool و منتظر ماندن برای اتمام jobs
func (p *FixedWorkerPool) Stop() {
	close(p.jobQueue)
	p.wg.Wait()
}

func demonstrateFixedWorkerPool() {
	fmt.Println("\n=== Fixed Worker Pool ===")

	pool := NewFixedWorkerPool(3, 10)

	// ارسال 10 job
	for i := 1; i <= 10; i++ {
		jobID := i
		pool.Submit(func() {
			fmt.Printf("Job %d: processing\n", jobID)
			time.Sleep(50 * time.Millisecond)
		})
	}

	pool.Stop()
	fmt.Println("✅ All jobs completed")
}

// ============================================================================
// بخش 7: نظارت و مانیتورینگ گوروتین‌ها
// ============================================================================

// GoroutineMonitor ساختار نظارت بر گوروتین‌ها
type GoroutineMonitor struct {
	activeCount int32
	maxCount    int32
	mu          sync.Mutex
	history     []int32
}

// NewGoroutineMonitor ایجاد مانیتور جدید
func NewGoroutineMonitor() *GoroutineMonitor {
	return &GoroutineMonitor{
		history: make([]int32, 0),
	}
}

// TrackStart ثبت شروع یک گوروتین
func (m *GoroutineMonitor) TrackStart() {
	count := atomic.AddInt32(&m.activeCount, 1)
	m.mu.Lock()
	if count > m.maxCount {
		m.maxCount = count
	}
	m.history = append(m.history, count)
	m.mu.Unlock()
}

// TrackEnd ثبت پایان یک گوروتین
func (m *GoroutineMonitor) TrackEnd() {
	atomic.AddInt32(&m.activeCount, -1)
}

// GetStats گرفتن آمار
func (m *GoroutineMonitor) GetStats() (active, max int32) {
	return atomic.LoadInt32(&m.activeCount), m.maxCount
}

// demonstrateGoroutineMonitoring نمایش نظارت بر گوروتین‌ها
func demonstrateGoroutineMonitoring() {
	fmt.Println("\n=== Goroutine Monitoring ===")

	monitor := NewGoroutineMonitor()
	var wg sync.WaitGroup

	// اجرای گوروتین‌های مختلف
	for i := 0; i < 20; i++ {
		wg.Add(1)
		monitor.TrackStart()

		go func(id int) {
			defer func() {
				monitor.TrackEnd()
				wg.Done()
			}()

			duration := time.Duration(id%100+50) * time.Millisecond
			time.Sleep(duration)
		}(i)
	}

	// مانیتورینگ در حین اجرا
	go func() {
		for i := 0; i < 5; i++ {
			active, max := monitor.GetStats()
			fmt.Printf("📈 Monitor: active=%d, max=%d\n", active, max)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	wg.Wait()
	active, max := monitor.GetStats()
	fmt.Printf("✅ Final: active=%d, max concurrent=%d\n", active, max)
}

// ============================================================================
// بخش 8: مدیریت panic در گوروتین‌ها (recovery)
// ============================================================================

// SafeGoroutine اجرای امن گوروتین با recovery
func SafeGoroutine(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("🔥 Recovered from panic in goroutine: %v\n", r)
				// می‌توانی لاگ بزنی، متریک ثبت کنی، etc.
			}
		}()
		fn()
	}()
}

// demonstratePanicRecovery نمایش مدیریت panic در گوروتین
func demonstratePanicRecovery() {
	fmt.Println("\n=== Panic Recovery in Goroutines ===")

	var wg sync.WaitGroup

	// گوروتین با panic (بدون recovery - برنامه را crash می‌کند)
	// go func() {
	// 	panic("crash!") // ❌ این برنامه را می‌کشد
	// }()

	// گوروتین امن با recovery
	wg.Add(1)
	SafeGoroutine(func() {
		defer wg.Done()
		fmt.Println("Working safely...")
		panic("something went wrong!") // این recovery می‌شود
		fmt.Println("This line never runs")
	})

	wg.Wait()
	fmt.Println("✅ Program continues after panic recovery")
}

// ============================================================================
// بخش 9: الگوی ErrGroup (گروه خطا) - مدیریت خطا در گوروتین‌ها
// ============================================================================

// ErrGroup یک پیاده‌سازی ساده از errgroup
type ErrGroup struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewErrGroup ایجاد errgroup جدید
func NewErrGroup() *ErrGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return &ErrGroup{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Go اجرای تابع در گوروتین با مدیریت خطا
func (g *ErrGroup) Go(fn func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := fn(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				g.cancel() // کنسل کردن همه
			})
		}
	}()
}

// Wait منتظر ماندن و برگرداندن اولین خطا
func (g *ErrGroup) Wait() error {
	g.wg.Wait()
	g.cancel()
	return g.err
}

func demonstrateErrGroup() {
	fmt.Println("\n=== ErrGroup (Error Handling) ===")

	group := NewErrGroup()

	// گوروتین اول - موفق
	group.Go(func() error {
		fmt.Println("Task 1: running...")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Task 1: success")
		return nil
	})

	// گوروتین دوم - خطا
	group.Go(func() error {
		fmt.Println("Task 2: running...")
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Task 2: failed!")
		return fmt.Errorf("task 2 error")
	})

	// گوروتین سوم - با context چک می‌کند
	group.Go(func() error {
		select {
		case <-time.After(500 * time.Millisecond):
			fmt.Println("Task 3: completed")
			return nil
		case <-group.ctx.Done():
			fmt.Printf("Task 3: cancelled due to error in another task\n")
			return group.ctx.Err()
		}
	})

	err := group.Wait()
	if err != nil {
		fmt.Printf("❌ ErrGroup returned: %v\n", err)
	} else {
		fmt.Println("✅ All tasks succeeded")
	}
}

// ============================================================================
// بخش 10: الگوی Backpressure (کنترل سرعت تولید)
// ============================================================================

// BackpressureWorkerPool worker pool با backpressure
type BackpressureWorkerPool struct {
	maxQueueSize int
	jobQueue     chan func()
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewBackpressurePool ایجاد pool با backpressure
func NewBackpressurePool(workers, queueSize int) *BackpressureWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	pool := &BackpressureWorkerPool{
		maxQueueSize: queueSize,
		jobQueue:     make(chan func(), queueSize),
		ctx:          ctx,
		cancel:       cancel,
	}

	for i := 0; i < workers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

func (p *BackpressureWorkerPool) worker(id int) {
	defer p.wg.Done()
	for {
		select {
		case job, ok := <-p.jobQueue:
			if !ok {
				return
			}
			job()
		case <-p.ctx.Done():
			return
		}
	}
}

// TrySubmit تلاش برای ارسال job (اگر queue پر باشد رد می‌کند)
func (p *BackpressureWorkerPool) TrySubmit(job func()) bool {
	select {
	case p.jobQueue <- job:
		return true
	default:
		return false // backpressure: queue پر است
	}
}

// SubmitWithTimeout ارسال job با timeout
func (p *BackpressureWorkerPool) SubmitWithTimeout(job func(), timeout time.Duration) error {
	select {
	case p.jobQueue <- job:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout submitting job")
	case <-p.ctx.Done():
		return p.ctx.Err()
	}
}

// Stop توقف pool
func (p *BackpressureWorkerPool) Stop() {
	p.cancel()
	close(p.jobQueue)
	p.wg.Wait()
}

func demonstrateBackpressure() {
	fmt.Println("\n=== Backpressure Pattern ===")

	pool := NewBackpressurePool(2, 3) // 2 worker, queue size 3

	// ارسال 10 job
	successCount := 0
	rejectCount := 0

	for i := 1; i <= 10; i++ {
		jobID := i
		submitted := pool.TrySubmit(func() {
			fmt.Printf("Job %d: processing\n", jobID)
			time.Sleep(100 * time.Millisecond)
		})

		if submitted {
			successCount++
		} else {
			rejectCount++
			fmt.Printf("Job %d: rejected (backpressure)\n", jobID)
		}
	}

	time.Sleep(500 * time.Millisecond) // اجازه بده jobs اجرا شوند
	pool.Stop()

	fmt.Printf("✅ Submitted: %d, Rejected: %d\n", successCount, rejectCount)
}

// ============================================================================
// بخش 11: الگوهای درست و غلط در مدیریت گوروتین
// ============================================================================

func demonstrateCorrectVsWrong() {
	fmt.Println("\n=== Correct vs Wrong Patterns in Goroutine Management ===")

	// ✅ درست: همیشه راهی برای خروج گوروتین وجود دارد
	fmt.Println("\n✅ Correct: Always have an exit path")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ch := make(chan int)
	go func() {
		select {
		case <-ch:
		case <-ctx.Done():
			return
		}
	}()

	// ❌ غلط: هیچ راه خروجی
	fmt.Println("\n❌ Wrong: No exit path (goroutine leak)")
	// ch2 := make(chan int)
	// go func() {
	// 	<-ch2  // forever blocked
	// }()

	// ✅ درست: بستن کانال توسط فرستنده
	fmt.Println("✅ Correct: Sender closes channel")
	work := make(chan int, 5)
	go func() {
		for i := 0; i < 5; i++ {
			work <- i
		}
		close(work)
	}()
	for range work {
		// consume
	}

	// ❌ غلط: بستن کانال توسط گیرنده
	fmt.Println("❌ Wrong: Receiver closes channel (can cause panic)")
	// work2 := make(chan int)
	// go func() {
	// 	value := <-work2
	// 	close(work2) // ❌ receiver closing
	// }()
}

// ============================================================================
// بخش 12: جمع‌بندی و اجرا
// ============================================================================

func main() {
	fmt.Println("========== GOROUTINE MANAGEMENT GUIDE ==========")
	fmt.Println("Complete guide for managing goroutines in production\n")

	// بخش 1: ردیابی
	demonstrateGoroutineTracking()

	// بخش 2: نشت یابی
	demonstrateGoroutineLeakDetection()

	// بخش 3: کنسل کردن با context
	demonstrateContextCancellation()

	// بخش 4: تایم‌اوت سراسری
	demonstrateGlobalTimeout()

	// بخش 5: سمافور
	demonstrateSemaphore()

	// بخش 6: worker pool ثابت
	demonstrateFixedWorkerPool()

	// بخش 7: مانیتورینگ
	demonstrateGoroutineMonitoring()

	// بخش 8: panic recovery
	demonstratePanicRecovery()

	// بخش 9: errgroup
	demonstrateErrGroup()

	// بخش 10: backpressure
	demonstrateBackpressure()

	// بخش 11: الگوهای درست و غلط
	demonstrateCorrectVsWrong()

	// جمع‌بندی نهایی
	fmt.Println("\n========== SUMMARY ==========")
	fmt.Println("| Technique              | Use Case                          |")
	fmt.Println("|------------------------|-----------------------------------|")
	fmt.Println("| Context.WithCancel     | Manual group cancellation         |")
	fmt.Println("| Context.WithTimeout    | Global deadline for all workers   |")
	fmt.Println("| Semaphore              | Limit concurrent goroutines       |")
	fmt.Println("| Fixed Worker Pool      | Reuse fixed number of workers     |")
	fmt.Println("| ErrGroup               | First error cancels all           |")
	fmt.Println("| Backpressure           | Reject work when system is busy   |")
	fmt.Println("| runtime.NumGoroutine() | Detect goroutine leaks            |")
	fmt.Println("| recover()              | Prevent crash from panic          |")

	fmt.Println("\n========== GOLDEN RULES ==========")
	fmt.Println("1. Every goroutine MUST have a way to exit")
	fmt.Println("2. Use context for cancellation, not flags")
	fmt.Println("3. Always know the maximum number of goroutines you create")
	fmt.Println("4. Monitor goroutine count in production")
	fmt.Println("5. Use worker pools instead of unbounded goroutine creation")
	fmt.Println("6. Never start a goroutine without knowing how it stops")
	fmt.Println("7. Use errgroup when you need first-error-wins semantics")
	fmt.Println("8. Always recover from panics in long-running goroutines")
}

// ============================================================================
// بخش 13: توابع کمکی برای پروژه واقعی
// ============================================================================

// RunWithTimeout اجرای یک تابع با تایم‌اوت
func RunWithTimeout(fn func() error, timeout time.Duration) error {
	done := make(chan error, 1)

	go func() {
		done <- fn()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("function timed out after %v", timeout)
	}
}

// RunWithRetry اجرای تابع با تلاش مجدد
func RunWithRetry(fn func() error, maxRetries int, backoff time.Duration) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(backoff * time.Duration(i+1))
	}
	return fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}

// LimitConcurrency محدود کردن همزمانی برای یک slice از آیتم‌ها
func LimitConcurrency[T any](items []T, maxConcurrent int, fn func(T) error) []error {
	sem := NewSemaphore(maxConcurrent)
	var wg sync.WaitGroup
	errors := make([]error, len(items))

	for i, item := range items {
		wg.Add(1)
		go func(idx int, it T) {
			defer wg.Done()
			sem.Acquire()
			defer sem.Release()
			errors[idx] = fn(it)
		}(i, item)
	}

	wg.Wait()
	return errors
}
