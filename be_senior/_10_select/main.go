// ============================================================================
// FILE: select_guide.go
// TITLE: راهنمای کامل Select در Go - همه کاربردها و الگوها
// HOW TO RUN: go run select_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Select چیست و چرا ساخته شد؟
// ============================================================================
//
// select به شما اجازه می‌دهد روی چندین کانال همزمان منتظر بمانید:
// - اولین کانالی که آماده شود (داده دارد یا جا برای فرستادن) اجرا می‌شود
// - اگر چند کانال همزمان آماده شوند، یک‌کی از تصادفی انتخاب می‌شود
// - اگر هیچ کانالی آماده نباشد، select بلاک می‌شود
// - case default: باعث می‌شود select non-blocking شود
//
// کاربردهای اصلی:
// 1. منتظر ماندن برای چند کانال همزمان
// 2. تایم‌اوت برای عملیات‌های کانال
// 3. Non-blocking send/receive
// 4. تشخیص بسته شدن کانال
// 5. متریک و نظارت
// ============================================================================

package _10_select

import (
	"fmt"
	"time"
)

// ============================================================================
// بخش 1: Select پایه - منتظر ماندن برای چند کانال
// ============================================================================

func demonstrateBasicSelect() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📡 SELECT BASICS - Waiting on Multiple Channels")
	fmt.Println(stringsRepeat("=", 80))

	// ایجاد دو کانال
	ch1 := make(chan string)
	ch2 := make(chan string)

	// گوروتین اول: بعد از 100ms به ch1 می‌فرستد
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "from channel 1"
	}()

	// گوروتین دوم: بعد از 200ms به ch2 می‌فرستد
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "from channel 2"
	}()

	// select منتظر می‌ماند تا اولین کانال آماده شود
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("✅ Received: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("✅ Received: %s\n", msg2)
		}
	}
}

// demonstrateSelectRandomness نشان دادن تصادفی بودن select هنگام آماده بودن همزمان
func demonstrateSelectRandomness() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎲 SELECT RANDOMNESS - When multiple channels are ready")
	fmt.Println(stringsRepeat("=", 80))

	ch1 := make(chan string)
	ch2 := make(chan string)

	// هر دو کانال را پر می‌کنیم
	go func() {
		ch1 <- "ping"
	}()
	go func() {
		ch2 <- "pong"
	}()

	// کمی صبر می‌کنیم تا هر دو goroutine برسند
	time.Sleep(10 * time.Millisecond)

	// اجرای select چند بار برای دیدن تصادفی بودن
	count1 := 0
	count2 := 0

	for i := 0; i < 100; i++ {
		select {
		case <-ch1:
			count1++
			// دوباره پر می‌کنیم برای دفعه بعد
			go func() { ch1 <- "ping" }()
		case <-ch2:
			count2++
			go func() { ch2 <- "pong" }()
		}
	}

	fmt.Printf("📊 Channel 1 selected: %d times\n", count1)
	fmt.Printf("📊 Channel 2 selected: %d times\n", count2)
	fmt.Println("💡 Note: Distribution is roughly equal (random)")
}

// ============================================================================
// بخش 2: Select با Timeout (زمان انتظار محدود)
// ============================================================================

func demonstrateSelectTimeout() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⏰ SELECT WITH TIMEOUT")
	fmt.Println(stringsRepeat("=", 80))

	ch := make(chan string)

	// یک عملیات که ممکن است طول بکشد
	go func() {
		time.Sleep(2 * time.Second) // خیلی طولانی
		ch <- "result"
	}()

	fmt.Println("Waiting for result with 1 second timeout...")

	select {
	case res := <-ch:
		fmt.Printf("✅ Received: %s\n", res)
	case <-time.After(1 * time.Second):
		fmt.Println("❌ Timeout! Operation took too long")
	}
}

func demonstrateSelectTimeoutLoop() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⏰ TIMEOUT LOOP - Retry with timeout")
	fmt.Println(stringsRepeat("=", 80))

	ch := make(chan int)

	// تولید کننده آهسته
	go func() {
		for i := 1; i <= 5; i++ {
			time.Sleep(500 * time.Millisecond)
			ch <- i
		}
		close(ch)
	}()

	timeout := 300 * time.Millisecond

	for {
		select {
		case value, ok := <-ch:
			if !ok {
				fmt.Println("✅ Channel closed, done!")
				return
			}
			fmt.Printf("Received: %d\n", value)
		case <-time.After(timeout):
			fmt.Printf("⚠️  Timeout after %v, continuing...\n", timeout)
			// در حلقه می‌مانیم و دوباره تلاش می‌کنیم
		}
	}
}

// ============================================================================
// بخش 3: Select با Default (Non-Blocking Operations)
// ============================================================================

func demonstrateSelectDefault() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚀 NON-BLOCKING SELECT (default)")
	fmt.Println(stringsRepeat("=", 80))

	ch := make(chan int)

	// Non-blocking receive (اگر داده نباشد، نمی‌ماند)
	select {
	case value := <-ch:
		fmt.Printf("Received: %d\n", value)
	default:
		fmt.Println("No data available (non-blocking receive)")
	}

	// Non-blocking send (اگر جا نباشد، نمی‌ماند)
	select {
	case ch <- 42:
		fmt.Println("Sent successfully")
	default:
		fmt.Println("Channel is full or no receiver (non-blocking send)")
	}

	// کاربرد واقعی: تلاش برای دریافت، اگر نبود ادامه بده
	fmt.Println("\n📊 Real use case: collecting metrics without blocking")

	metrics := make(chan int, 10)

	// پر کردن متریک‌ها
	go func() {
		for i := 1; i <= 5; i++ {
			metrics <- i * 10
		}
		close(metrics)
	}()

	// جمع‌آوری متریک‌ها بدون بلاک شدن
	collected := make([]int, 0)
	for {
		select {
		case m, ok := <-metrics:
			if !ok {
				fmt.Printf("✅ Collected all metrics: %v\n", collected)
				return
			}
			collected = append(collected, m)
		default:
			// بدون متریک، کار دیگری بکن
			fmt.Println("  No metrics yet, doing other work...")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// ============================================================================
// بخش 4: Select با تشخیص بسته شدن کانال
// ============================================================================

func demonstrateSelectChannelClose() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔒 DETECTING CLOSED CHANNELS WITH select")
	fmt.Println(stringsRepeat("=", 80))

	jobs := make(chan int, 5)
	done := make(chan bool)

	// Worker
	go func() {
		for {
			select {
			case job, ok := <-jobs:
				if !ok {
					fmt.Println("📢 Jobs channel closed, stopping worker")
					done <- true
					return
				}
				fmt.Printf("🔧 Processing job: %d\n", job)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// ارسال چند job
	for i := 1; i <= 3; i++ {
		jobs <- i
	}
	close(jobs) // بستن کانال

	// منتظر اتمام worker
	<-done
	fmt.Println("✅ Worker stopped gracefully")
}

// ============================================================================
// بخش 5: Select در حلقه (Loop with Select)
// ============================================================================

func demonstrateSelectLoop() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔄 SELECT IN LOOP - Event Loop Pattern")
	fmt.Println(stringsRepeat("=", 80))

	type Event struct {
		Type string
		Data interface{}
	}

	events := make(chan Event, 10)
	quit := make(chan struct{})

	// تولید کننده رویداد
	go func() {
		for i := 1; i <= 5; i++ {
			events <- Event{Type: "data", Data: i}
			time.Sleep(100 * time.Millisecond)
		}
		events <- Event{Type: "flush", Data: nil}
		time.Sleep(100 * time.Millisecond)
		close(quit)
	}()

	// Event Loop
	running := true
	for running {
		select {
		case event, ok := <-events:
			if !ok {
				fmt.Println("Events channel closed")
				running = false
				continue
			}

			switch event.Type {
			case "data":
				fmt.Printf("📦 Processing data: %v\n", event.Data)
			case "flush":
				fmt.Println("💾 Flushing buffers...")
			}

		case <-quit:
			fmt.Println("🛑 Quit signal received")
			running = false

		default:
			// هیچ رویدادی نیست، کمی استراحت
			time.Sleep(10 * time.Millisecond)
		}
	}

	fmt.Println("✅ Event loop finished")
}

// ============================================================================
// بخش 6: Select با Nil Channel (مفید برای غیرفعال کردن case)
// ============================================================================

func demonstrateNilChannel() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("❌ NIL CHANNELS - Disabling select cases")
	fmt.Println(stringsRepeat("=", 80))

	ch1 := make(chan int)
	ch2 := make(chan int)

	// یک کانال را nil می‌کنیم (غیرفعال)
	var nilChannel chan int = nil

	go func() {
		ch1 <- 42
	}()

	// select با nil channel (case با nil channel هرگز انتخاب نمی‌شود)
	select {
	case v := <-ch1:
		fmt.Printf("✅ From ch1: %d\n", v)
	case v := <-nilChannel:
		fmt.Printf("This will never be selected: %d\n", v)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("Timeout")
	}

	// کاربرد واقعی: غیرفعال کردن موقت یک کانال
	fmt.Println("\n📌 Real use case: Temporarily disabling a channel")

	enableCh1 := true
	ch1Active := ch1
	ch2Active := ch2

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch1 <- 100
	}()

	for i := 0; i < 2; i++ {
		select {
		case v, ok := <-ch1Active:
			if !ok {
				ch1Active = nil // غیرفعال کردن
				continue
			}
			fmt.Printf("Got from ch1: %d\n", v)
			if enableCh1 {
				// بعد از دریافت، کانال را غیرفعال می‌کنیم
				ch1Active = nil
				fmt.Println("  ch1 disabled after receiving")
			}
		case v := <-ch2Active:
			fmt.Printf("Got from ch2: %d\n", v)
		}
	}
}

// ============================================================================
// بخش 7: Select با Ticker (Timerهای تکراری)
// ============================================================================

func demonstrateSelectTicker() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⏲️ SELECT WITH TICKER - Periodic Operations")
	fmt.Println(stringsRepeat("=", 80))

	ticker := time.NewTicker(200 * time.Millisecond)
	heartbeat := time.NewTicker(500 * time.Millisecond)
	timeout := time.After(2 * time.Second)

	counter := 0

	fmt.Println("Starting periodic tasks...")

	for {
		select {
		case <-ticker.C:
			counter++
			fmt.Printf("  🔄 Tick %d: doing work\n", counter)

		case <-heartbeat.C:
			fmt.Printf("  💓 Heartbeat: system alive, processed %d ticks\n", counter)

		case <-timeout:
			fmt.Println("⏰ Timeout reached, stopping")
			ticker.Stop()
			heartbeat.Stop()
			return
		}
	}
}

// ============================================================================
// بخش 8: الگوی Quit (خروج تمیز از گوروتین‌ها)
// ============================================================================

type Worker struct {
	quit chan struct{}
	done chan struct{}
}

func NewWorker() *Worker {
	return &Worker{
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
}

func (w *Worker) Start(workChan <-chan int) {
	go func() {
		defer close(w.done)

		for {
			select {
			case job, ok := <-workChan:
				if !ok {
					fmt.Println("  Work channel closed, stopping worker")
					return
				}
				fmt.Printf("  Processing job: %d\n", job)
				time.Sleep(100 * time.Millisecond)

			case <-w.quit:
				fmt.Println("  Quit signal received, cleaning up...")
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	close(w.quit)
	<-w.done
	fmt.Println("  Worker stopped cleanly")
}

func demonstrateQuitPattern() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚪 QUIT PATTERN - Graceful Shutdown")
	fmt.Println(stringsRepeat("=", 80))

	work := make(chan int, 10)
	worker := NewWorker()

	worker.Start(work)

	// ارسال چند job
	for i := 1; i <= 3; i++ {
		work <- i
	}

	time.Sleep(200 * time.Millisecond)

	// توقف graceful
	fmt.Println("Initiating shutdown...")
	worker.Stop()

	fmt.Println("✅ Graceful shutdown complete")
}

// ============================================================================
// بخش 9: الگوی متریک و نظارت با Select
// ============================================================================

type MonitoredWorker struct {
	requests  chan int
	responses chan int
	metrics   chan string
	quit      chan struct{}
}

func NewMonitoredWorker() *MonitoredWorker {
	return &MonitoredWorker{
		requests:  make(chan int, 100),
		responses: make(chan int, 100),
		metrics:   make(chan string, 10),
		quit:      make(chan struct{}),
	}
}

func (mw *MonitoredWorker) Start() {
	go func() {
		processed := 0
		totalLatency := time.Duration(0)

		for {
			select {
			case req := <-mw.requests:
				start := time.Now()

				// شبیه‌سازی کار
				time.Sleep(50 * time.Millisecond)

				latency := time.Since(start)
				totalLatency += latency
				processed++

				mw.responses <- req * 2

				// ارسال متریک هر 5 درخواست
				if processed%5 == 0 {
					avgLatency := totalLatency / time.Duration(processed)
					mw.metrics <- fmt.Sprintf("processed=%d, avg_latency=%v",
						processed, avgLatency)
				}

			case <-mw.metrics:
				// خواندن متریک توسط monitoring system
				// (در اینجا فقط برای نمایش)

			case <-mw.quit:
				close(mw.responses)
				close(mw.metrics)
				return
			}
		}
	}()
}

func (mw *MonitoredWorker) Stop() {
	close(mw.quit)
}

func demonstrateMonitoring() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📊 MONITORING WITH SELECT")
	fmt.Println(stringsRepeat("=", 80))

	worker := NewMonitoredWorker()
	worker.Start()

	// ارسال درخواست‌ها
	go func() {
		for i := 1; i <= 12; i++ {
			worker.requests <- i
			time.Sleep(20 * time.Millisecond)
		}
	}()

	// دریافت متریک‌ها و پاسخ‌ها
	done := time.After(2 * time.Second)

	for {
		select {
		case resp, ok := <-worker.responses:
			if !ok {
				fmt.Println("Worker stopped")
				return
			}
			fmt.Printf("📥 Response: %d\n", resp)

		case metric, ok := <-worker.metrics:
			if !ok {
				return
			}
			fmt.Printf("📊 Metric: %s\n", metric)

		case <-done:
			fmt.Println("⏰ Monitoring done")
			worker.Stop()
			return
		}
	}
}

// ============================================================================
// بخش 10: الگوهای ترکیبی و پیشرفته
// ============================================================================

func demonstrateAdvancedSelect() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔧 ADVANCED SELECT PATTERNS")
	fmt.Println(stringsRepeat("=", 80))

	// الگو: Priority Select (اولویت‌بندی)
	fmt.Println("\n📌 Priority Select Pattern")

	highPriority := make(chan int, 5)
	lowPriority := make(chan int, 5)

	go func() {
		for i := 1; i <= 5; i++ {
			lowPriority <- i
			if i%2 == 0 {
				highPriority <- i * 10
			}
			time.Sleep(10 * time.Millisecond)
		}
		close(highPriority)
		close(lowPriority)
	}()

	// اولویت بالا همیشه اول پردازش می‌شود
	for {
		select {
		case h, ok := <-highPriority:
			if !ok {
				highPriority = nil
				continue
			}
			fmt.Printf("🔥 HIGH priority: %d\n", h)
		default:
			select {
			case l, ok := <-lowPriority:
				if !ok {
					fmt.Println("All done")
					return
				}
				fmt.Printf("  LOW priority: %d\n", l)
			case <-time.After(100 * time.Millisecond):
				fmt.Println("Timeout, no more messages")
				return
			}
		}
	}
}

func demonstrateSelectWithContext() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 SELECT WITH CONTEXT (Cancellation)")
	fmt.Println(stringsRepeat("=", 80))

	// در یک فایل واقعی، context را import می‌کنید
	// اینجا یک پیاده‌سازی ساده شبیه‌سازی می‌کنیم

	type Context interface {
		Done() <-chan struct{}
		Err() error
	}

	// شبیه‌سازی context.WithTimeout
	simulateContext := func(timeout time.Duration) (ctx Context, cancel func()) {
		done := make(chan struct{})
		go func() {
			time.Sleep(timeout)
			close(done)
		}()
		return &simpleContext{done: done}, func() {}
	}

	type simpleContext struct {
		done chan struct{}
	}

	func (c *simpleContext) Done() <-chan struct{} { return c.done }
	func (c *simpleContext) Err() error { return nil }

	work := make(chan int)

	go func() {
		time.Sleep(500 * time.Millisecond)
		work <- 42
	}()

	ctx, cancel := simulateContext(200 * time.Millisecond)
	defer cancel()

	select {
	case result := <-work:
		fmt.Printf("✅ Work completed: %d\n", result)
	case <-ctx.Done():
		fmt.Println("❌ Work cancelled by context")
	}
}

// ============================================================================
// بخش 11: الگوهای درست و غلط با Select
// ============================================================================

func demonstrateSelectCorrectVsWrong() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("✅❌ SELECT - CORRECT VS WRONG PATTERNS")
	fmt.Println(stringsRepeat("=", 80))

	// ✅ درست: استفاده از default برای non-blocking
	fmt.Println("\n✅ Correct: Non-blocking with default")
	ch := make(chan int)
	select {
	case <-ch:
	default:
		fmt.Println("  No data, continuing")
	}

	// ❌ غلط: select بدون case (compile error)
	// select {} // این خط کد کامپایل نمی‌شود

	// ✅ درست: select در حلقه با شرط خروج
	fmt.Println("\n✅ Correct: Select loop with exit condition")
	quit := make(chan struct{})
	ch2 := make(chan int)

	go func() {
		time.Sleep(100 * time.Millisecond)
		close(quit)
	}()

	running := true
	for running {
		select {
		case <-ch2:
			// process
		case <-quit:
			fmt.Println("  Quitting loop")
			running = false
		}
	}

	// ❌ غلط: فراموش کردن close در کانال‌های quit
	fmt.Println("\n❌ Wrong: Never closing quit channel (goroutine leak)")

	// ✅ درست: استفاده از time.After برای timeout یکبار مصرف
	fmt.Println("\n✅ Correct: One-time timeout with time.After")
	select {
	case <-ch:
	case <-time.After(1 * time.Second):
		fmt.Println("  Timeout")
	}

	// ❌ غلط: استفاده از time.After در حلقه (نشت حافظه)
	fmt.Println("\n❌ Wrong: time.After in loop (memory leak)")
	fmt.Println("  for { select { case <-time.After(1s): } } // هر بار تایمر جدید می‌سازد")
	fmt.Println("  ✅ راه حل: استفاده از time.NewTicker برای کارهای تکراری")
}

// ============================================================================
// بخش 12: جمع‌بندی و جدول مرجع سریع
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 COMPLETE SELECT GUIDE IN GO")
	fmt.Println("All patterns and use cases with executable examples")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: Select پایه
	demonstrateBasicSelect()
	demonstrateSelectRandomness()

	// بخش 2: Timeout
	demonstrateSelectTimeout()
	demonstrateSelectTimeoutLoop()

	// بخش 3: Non-blocking
	demonstrateSelectDefault()

	// بخش 4: تشخیص بسته شدن
	demonstrateSelectChannelClose()

	// بخش 5: حلقه
	demonstrateSelectLoop()

	// بخش 6: nil channel
	demonstrateNilChannel()

	// بخش 7: ticker
	demonstrateSelectTicker()

	// بخش 8: quit pattern
	demonstrateQuitPattern()

	// بخش 9: monitoring
	demonstrateMonitoring()

	// بخش 10: پیشرفته
	demonstrateAdvancedSelect()
	demonstrateSelectWithContext()

	// بخش 11: درست و غلط
	demonstrateSelectCorrectVsWrong()

	// بخش 12: جمع‌بندی
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 SELECT QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ PATTERN                    │ USE CASE                         │")
	fmt.Println("├────────────────────────────┼───────────────────────────────────┤")
	fmt.Println("│ select { case <-ch1 ... }  │ منتظر ماندن برای چند کانال        │")
	fmt.Println("│ select { case <-ch: default}│ Non-blocking receive/send        │")
	fmt.Println("│ select { case <-ch: case <-time.After(d): } │ Timeout         │")
	fmt.Println("│ select { case v,ok:=<-ch: }│ تشخیص بسته شدن کانال              │")
	fmt.Println("│ select { case <-ticker.C: }│ عملیات تکراری (با Ticker)        │")
	fmt.Println("│ select { case <-quit: }    │ خروج تمیز از گوروتین‌ها           │")
	fmt.Println("│ select { default: }        │ بررسی بدون بلاک شدن               │")
	fmt.Println("│ select { case <-nilCh: }   │ غیرفعال کردن موقت یک case         │")
	fmt.Println("└────────────────────────────┴───────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. select با اولین case آماده اجرا می‌شود (تصادفی اگر چندتا باشند)")
	fmt.Println("  2. از default برای عملیات non-blocking استفاده کنید")
	fmt.Println("  3. همیشه راهی برای خروج از select loop داشته باشید (quit channel)")
	fmt.Println("  4. برای timeout یکبار مصرف از time.After استفاده کنید")
	fmt.Println("  5. برای عملیات تکراری از time.Ticker استفاده کنید")
	fmt.Println("  6. nil channel در select هرگز انتخاب نمی‌شود (برای غیرفعال کردن)")
	fmt.Println("  7. select {} بدون case باعث deadlock می‌شود")
	fmt.Println("  8. در حلقه‌های select، همیشه شرط خروج داشته باشید")

	fmt.Println("\n⚠️  COMMON PITFALLS:")
	fmt.Println("  • فراموش کردن default در non-blocking select")
	fmt.Println("  • استفاده از time.After در حلقه (نشت حافظه)")
	fmt.Println("  • بستن کانال از سمت گیرنده")
	fmt.Println("  • نداشتن راه خروج از select loop")
	fmt.Println("  • select با یک case که هیچ‌گاه آماده نمی‌شود (deadlock)")
}

// ============================================================================
// بخش 13: توابع کمکی
// ============================================================================

// stringsRepeat یک تابع کمکی ساده (در فایل واقعی از import "strings" استفاده کنید)
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// SelectWithMultipleTimeouts select با چند تایم‌اوت مختلف
func SelectWithMultipleTimeouts(ch <-chan int, timeouts ...time.Duration) (int, error) {
	if len(timeouts) == 0 {
		return 0, fmt.Errorf("no timeouts provided")
	}

	// ساخت کانال‌های تایم‌اوت
	cases := make([]interface{}, 0, len(timeouts)+1)

	// ساخت select caseها به صورت داینامیک (در عمل باید از reflection استفاده کرد)
	// اینجا فقط یک مثال ساده می‌دهیم

	for i, tout := range timeouts {
		_ = i
		_ = tout
		// در عمل: time.After(tout)
	}

	_ = cases
	return 0, nil
}

// ForSelect یک helper برای select loop
type ForSelect struct {
	quit chan struct{}
	done chan struct{}
}

func NewForSelect() *ForSelect {
	return &ForSelect{
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
}

func (fs *ForSelect) Run(handler func(quit <-chan struct{}) bool) {
	go func() {
		defer close(fs.done)
		for {
			select {
			case <-fs.quit:
				return
			default:
				if !handler(fs.quit) {
					return
				}
			}
		}
	}()
}

func (fs *ForSelect) Stop() {
	close(fs.quit)
	<-fs.done
}