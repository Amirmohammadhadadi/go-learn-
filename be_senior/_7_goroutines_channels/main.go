// ============================================================================
// FILE: goroutines_channels_guide.go
// TITLE: راهنمای کامل گوروتین و کانال در Go - سطح مقدماتی تا پیشرفته
// HOW TO RUN: go run goroutines_channels_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - گوروتین چیست و چرا ساخته شد؟
// ============================================================================
//
// گوروتین (Goroutine):
// - یک تابع سبک وزن که همزمان (concurrently) با توابع دیگر اجرا می‌شود
// - هزینه ساخت: ~2KB حافظه (مقایسه با thread که ~1MB است)
// - می‌توانید ده‌ها هزار گوروتین همزمان داشته باشید
// - توسط Go runtime مدیریت می‌شود، نه توسط OS
//
// کانال (Channel):
// - راهی برای ارتباط بین گوروتین‌ها
// - مانند یک لوله (pipe) که داده را از یک گوروتین به گوروتین دیگر می‌فرستد
// - به طور پیش‌فرض blocking است (فرستنده و گیرنده منتظر هم می‌مانند)
//
// قانون طلایی:
// "Do not communicate by sharing memory; instead, share memory by communicating."
// ============================================================================

package _7_goroutines_channels

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// بخش 1: ساختار گوروتین و تفاوت با تابع معمولی
// ============================================================================

// normalFunction یک تابع معمولی ساده
func normalFunction(name string) {
	fmt.Printf("Normal function: Hello %s\n", name)
}

// goroutineFunction تابعی که به عنوان گوروتین اجرا می‌شود
func goroutineFunction(name string) {
	fmt.Printf("Goroutine: Hello %s\n", name)
}

// demonstrateBasicGoroutine نمایش تفاوت اجرای معمولی و گوروتین
func demonstrateBasicGoroutine() {
	fmt.Println("\n=== Basic Goroutine vs Normal Function ===")

	// اجرای معمولی - منتظر می‌ماند تا تمام شود
	normalFunction("Ali")

	// اجرا به عنوان گوروتین - بلافاصله برمی‌گردد و موازی اجرا می‌شود
	go goroutineFunction("Reza")

	// مشکل: برنامه ممکن است قبل از اجرای گوروتین تمام شود
	// راه حل: کمی صبر می‌کنیم (در بخش بعدی راه حل درست را می‌گوییم)
	time.Sleep(10 * time.Millisecond)

	fmt.Println("Main function finished")
}

// ============================================================================
// بخش 2: منتظر ماندن برای گوروتین‌ها (sync.WaitGroup)
// ============================================================================

// worker یک تابع که کار مشخصی انجام می‌دهد
func worker(id int, wg *sync.WaitGroup) {
	// وقتی تابع تمام شد، به WaitGroup بگو که تمام کردم
	defer wg.Done()

	fmt.Printf("Worker %d: starting\n", id)
	time.Sleep(time.Duration(id) * 100 * time.Millisecond)
	fmt.Printf("Worker %d: done\n", id)
}

// demonstrateWaitGroup استفاده از WaitGroup برای منتظر ماندن چند گوروتین
func demonstrateWaitGroup() {
	fmt.Println("\n=== WaitGroup Example ===")

	// WaitGroup مانند شمارنده است که منتظر می‌ماند تا صفر شود
	var wg sync.WaitGroup

	// اجرای 3 تا گوروتین
	for i := 1; i <= 3; i++ {
		wg.Add(1)         // به شمارنده اضافه کن
		go worker(i, &wg) // اشاره‌گر فرستادیم چون داخل تابع تغییر می‌کند
	}

	// منتظر می‌مانیم تا همه گوروتین‌ها Done() را صدا بزنند
	wg.Wait()
	fmt.Println("All workers completed")
}

// ============================================================================
// بخش 3: کانال (Channel) - مبانی
// ============================================================================

// demonstrateUnbufferedChannel کانال بدون بافر (unbuffered)
func demonstrateUnbufferedChannel() {
	fmt.Println("\n=== Unbuffered Channel ===")

	// ایجاد کانال از نوع string
	// بدون بافر: فرستنده تا وقتی گیرنده نیامده بلاک می‌شود
	messages := make(chan string)

	// گوروتین فرستنده
	go func() {
		fmt.Println("Sending message...")
		messages <- "Hello from goroutine" // بلاک می‌شود تا گیرنده بیاید
		fmt.Println("Message sent successfully")
	}()

	// کمی صبر می‌کنیم تا گوروتین برسد
	time.Sleep(100 * time.Millisecond)

	// گیرنده: مقدار را از کانال می‌خواند
	fmt.Println("Receiving message...")
	msg := <-messages
	fmt.Printf("Received: %s\n", msg)
}

// demonstrateBufferedChannel کانال با بافر (buffered)
func demonstrateBufferedChannel() {
	fmt.Println("\n=== Buffered Channel ===")

	// کانال با بافر 2 تایی
	// تا 2 مقدار می‌تواند بدون گیرنده بفرستد
	ch := make(chan int, 2)

	// فرستادن بدون گیرنده - تا ظرفیت بافر کار می‌کند
	ch <- 1
	ch <- 2
	fmt.Println("Sent 2 values without receiver")

	// سومی بلاک می‌شود (اگر این خط را uncomment کنید برنامه قفل می‌شود)
	// ch <- 3 // ❌ این خط برنامه را قفل می‌کند

	// خواندن مقادیر
	fmt.Printf("Received: %d\n", <-ch)
	fmt.Printf("Received: %d\n", <-ch)

	// حالا کانال خالی است
	fmt.Println("Channel is empty")
}

// ============================================================================
// بخش 4: فرستادن و گرفتن از کانال - عملی
// ============================================================================

// ping فقط می‌فرستد (send-only channel)
func ping(pings chan<- string, msg string) {
	pings <- msg
}

// pong می‌گیرد و سپس می‌فرستد (receive then send)
func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- fmt.Sprintf("pong: %s", msg)
}

// demonstrateChannelDirection نمایش کانال یک طرفه
func demonstrateChannelDirection() {
	fmt.Println("\n=== Channel Direction (Send/Receive Only) ===")

	pings := make(chan string, 1)
	pongs := make(chan string, 1)

	ping(pings, "hello")
	pong(pings, pongs)

	result := <-pongs
	fmt.Printf("Result: %s\n", result)
}

// ============================================================================
// بخش 5: حلقه بر روی کانال (Range over channel)
// ============================================================================

// producer مقادیر تولید می‌کند و کانال را می‌بندد
func producer(ch chan<- int, count int) {
	for i := 1; i <= count; i++ {
		fmt.Printf("Producing: %d\n", i)
		ch <- i
		time.Sleep(50 * time.Millisecond)
	}
	close(ch) // بستن کانال نشان می‌دهد که دیگر مقداری نمی‌آید
}

// demonstrateRangeOverChannel حلقه زدن با range روی کانال
func demonstrateRangeOverChannel() {
	fmt.Println("\n=== Range Over Channel ===")

	ch := make(chan int, 3)

	go producer(ch, 5)

	// range تا وقتی که کانال بسته شود ادامه می‌یابد
	for value := range ch {
		fmt.Printf("Consuming: %d\n", value)
	}

	fmt.Println("Channel closed, loop finished")
}

// ============================================================================
// بخش 6: select - منتظر ماندن برای چند کانال
// ============================================================================

func demonstrateSelect() {
	fmt.Println("\n=== Select Statement ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	// گوروتین اول
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "from channel 1"
	}()

	// گوروتین دوم
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "from channel 2"
	}()

	// select مانند switch برای کانال‌هاست
	// اولین کانالی که آماده باشد را انتخاب می‌کند
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("Received: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("Received: %s\n", msg2)
		}
	}
}

// demonstrateSelectWithTimeout select با timeout
func demonstrateSelectWithTimeout() {
	fmt.Println("\n=== Select with Timeout ===")

	ch := make(chan string)

	go func() {
		time.Sleep(2 * time.Second)
		ch <- "result"
	}()

	select {
	case res := <-ch:
		fmt.Printf("Received: %s\n", res)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout! Operation took too long")
	}
}

// demonstrateSelectWithDefault select با default (non-blocking)
func demonstrateSelectWithDefault() {
	fmt.Println("\n=== Select with Default (Non-blocking) ===")

	ch := make(chan string)

	select {
	case msg := <-ch:
		fmt.Printf("Received: %s\n", msg)
	default:
		fmt.Println("No message waiting, continuing...")
	}

	// بدون default، این کد بلاک می‌شد
	fmt.Println("Program continues without blocking")
}

// ============================================================================
// بخش 7: الگوی Fan-Out (توزیع کار بین چند گوروتین)
// ============================================================================

// workerPool یک تابع که کار را بین چند worker پخش می‌کند
func workerPool(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(50 * time.Millisecond) // شبیه‌سازی کار
		results <- job * 2
	}
}

func demonstrateFanOut() {
	fmt.Println("\n=== Fan-Out Pattern (Worker Pool) ===")

	const numJobs = 10
	const numWorkers = 3

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	// راه‌اندازی workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go workerPool(w, jobs, results, &wg)
	}

	// ارسال jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// منتظر تمام شدن workers
	go func() {
		wg.Wait()
		close(results)
	}()

	// جمع‌آوری نتایج
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}
}

// ============================================================================
// بخش 8: الگوی Fan-In (ادغام چند کانال در یک کانال)
// ============================================================================

// producer2 تولید کننده ساده
func producer2(id int, ch chan<- int, count int) {
	for i := 1; i <= count; i++ {
		ch <- id*100 + i
		time.Sleep(20 * time.Millisecond)
	}
}

// fanIn ادغام چند کانال ورودی در یک کانال خروجی
func fanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	// برای هر کانال ورودی، یک گوروتین می‌زنیم که داده‌ها را به out بفرستد
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for value := range c {
				out <- value
			}
		}(ch)
	}

	// وقتی همه تمام شدند، out را می‌بندیم
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func demonstrateFanIn() {
	fmt.Println("\n=== Fan-In Pattern ===")

	ch1 := make(chan int, 5)
	ch2 := make(chan int, 5)

	go producer2(1, ch1, 3)
	go producer2(2, ch2, 3)

	// ادغام
	merged := fanIn(ch1, ch2)

	// خواندن از کانال ادغام شده
	for value := range merged {
		fmt.Printf("Merged value: %d\n", value)
	}
}

// ============================================================================
// بخش 9: الگوی Pipeline (زنجیره کانال‌ها)
// ============================================================================

// generate اعداد تولید می‌کند
func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// square اعداد را به توان دو می‌رساند
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// printNumbers اعداد را چاپ می‌کند
func printNumbers(in <-chan int) {
	for n := range in {
		fmt.Printf("Pipeline result: %d\n", n)
	}
}

func demonstratePipeline() {
	fmt.Println("\n=== Pipeline Pattern ===")

	// زنجیره: generate -> square -> print
	numbers := generate(1, 2, 3, 4, 5)
	squared := square(numbers)
	printNumbers(squared)
}

// ============================================================================
// بخش 10: الگوی Worker Pool با کانال‌ها (نسخه کامل)
// ============================================================================

type Job struct {
	ID     int
	Number int
}

type Result struct {
	JobID int
	Value int
}

func workerPoolAdvanced(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// شبیه‌سازی کار سنگین
		time.Sleep(10 * time.Millisecond)
		result := Result{
			JobID: job.ID,
			Value: job.Number * job.Number,
		}
		results <- result
	}
}

func demonstrateAdvancedWorkerPool() {
	fmt.Println("\n=== Advanced Worker Pool ===")

	const numJobs = 10
	const numWorkers = 4

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)
	var wg sync.WaitGroup

	// راه‌اندازی workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go workerPoolAdvanced(w, jobs, results, &wg)
	}

	// ارسال jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Number: j * 5}
	}
	close(jobs)

	// منتظر تمام شدن workers در یک گوروتین جدا
	go func() {
		wg.Wait()
		close(results)
	}()

	// جمع‌آوری نتایج
	for result := range results {
		fmt.Printf("Job %d result: %d\n", result.JobID, result.Value)
	}
}

// ============================================================================
// بخش 11: الگوهای درست و غلط
// ============================================================================

// ✅ الگوی درست: استفاده از WaitGroup برای منتظر ماندن
func correctWaitGroup() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// کار
	}()
	wg.Wait()
}

// ❌ الگوی غلط: استفاده از time.Sleep برای منتظر ماندن
func wrongWait() {
	go func() {
		// کار
	}()
	time.Sleep(1 * time.Second) // ❌ هیچ تضمینی نیست، و زمان را هدر می‌دهد
}

// ✅ الگوی درست: بستن کانال توسط فرستنده
func correctClose() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch) // ✅ فرستنده کانال را می‌بندد
	}()

	for v := range ch {
		fmt.Println(v)
	}
}

// ❌ الگوی غلط: بستن کانال توسط گیرنده یا بستن دوباره
func wrongClose() {
	ch := make(chan int)
	close(ch)
	// close(ch) // ❌ panic: close of closed channel
	// ch <- 1   // ❌ panic: send on closed channel
}

// ✅ الگوی درست: بررسی بسته بودن کانال
func correctCheckClosed() {
	ch := make(chan int, 1)
	ch <- 42
	close(ch)

	value, ok := <-ch
	fmt.Printf("Value: %d, OK: %v\n", value, ok)

	value, ok = <-ch
	fmt.Printf("Value: %d, OK: %v (channel closed)\n", value, ok)
}

// ❌ الگوی غلط: نشت گوروتین (goroutine leak)
func wrongGoroutineLeak() {
	ch := make(chan int)
	go func() {
		// این گوروتین تا ابد منتظر می‌ماند و هرگز بسته نمی‌شود
		value := <-ch // ❌ هیچ‌کس به این کانال نمی‌فرستد
		fmt.Println(value)
	}()
	// فراموش کردیم ch <- 1 را بفرستیم
}

// ============================================================================
// بخش 12: جمع‌بندی و اجرای همه مثال‌ها
// ============================================================================

func main() {
	fmt.Println("========== GOROUTINES & CHANNELS GUIDE ==========")
	fmt.Println("Complete guide with executable examples\n")

	// بخش 1: مبانی
	demonstrateBasicGoroutine()

	// بخش 2: WaitGroup
	demonstrateWaitGroup()

	// بخش 3: کانال‌ها
	demonstrateUnbufferedChannel()
	demonstrateBufferedChannel()

	// بخش 4: جهت کانال
	demonstrateChannelDirection()

	// بخش 5: range روی کانال
	demonstrateRangeOverChannel()

	// بخش 6: select
	demonstrateSelect()
	demonstrateSelectWithTimeout()
	demonstrateSelectWithDefault()

	// بخش 7: Fan-Out
	demonstrateFanOut()

	// بخش 8: Fan-In
	demonstrateFanIn()

	// بخش 9: Pipeline
	demonstratePipeline()

	// بخش 10: Worker Pool پیشرفته
	demonstrateAdvancedWorkerPool()

	// بخش 11: الگوهای درست و غلط (نمایش)
	fmt.Println("\n=== Correct vs Wrong Patterns ===")
	fmt.Println("✅ Use WaitGroup, not time.Sleep")
	fmt.Println("✅ Sender closes channel, not receiver")
	fmt.Println("✅ Check ok when reading from closed channel")
	fmt.Println("❌ Never close channel twice")
	fmt.Println("❌ Never send on closed channel")
	fmt.Println("❌ Avoid goroutine leaks")

	// بخش 12: جمع‌بندی
	fmt.Println("\n========== SUMMARY ==========")
	fmt.Println("| Concept              | Description                          |")
	fmt.Println("|----------------------|--------------------------------------|")
	fmt.Println("| Goroutine            | Lightweight concurrent function      |")
	fmt.Println("| go func()            | Start a goroutine                    |")
	fmt.Println("| sync.WaitGroup       | Wait for goroutines to finish        |")
	fmt.Println("| make(chan T)         | Unbuffered channel                   |")
	fmt.Println("| make(chan T, N)      | Buffered channel with capacity N     |")
	fmt.Println("| ch <- value          | Send to channel                      |")
	fmt.Println("| value := <-ch        | Receive from channel                 |")
	fmt.Println("| close(ch)            | Close channel (sender does this)     |")
	fmt.Println("| select               | Wait on multiple channels            |")
	fmt.Println("| <-time.After(d)      | Timeout in select                    |")
	fmt.Println("| default:             | Non-blocking select                  |")
	fmt.Println("| fan-out              | Distribute work to multiple workers  |")
	fmt.Println("| fan-in               | Merge multiple channels into one     |")
	fmt.Println("| pipeline             | Chain channels together              |")

	fmt.Println("\n========== GOLDEN RULES ==========")
	fmt.Println("1. Don't communicate by sharing memory; share memory by communicating")
	fmt.Println("2. Use WaitGroup, never time.Sleep() to wait for goroutines")
	fmt.Println("3. The sender, not receiver, should close the channel")
	fmt.Println("4. Buffered channels are not a performance optimization")
	fmt.Println("5. A nil channel blocks forever (useful in select)")
	fmt.Println("6. Always assume channel operations will block")
	fmt.Println("7. Use select for timeout, non-blocking, and multiple channels")
}

// ============================================================================
// بخش 13: توابع کمکی برای پروژه واقعی
// ============================================================================

// SafeGo اجرای امن گوروتین با recovery
func SafeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic in goroutine: %v\n", r)
			}
		}()
		fn()
	}()
}

// MergeChannels ادغام چند کانال با نوع دلخواه (generic-like)
func MergeChannels[T any](channels ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan T) {
			defer wg.Done()
			for value := range c {
				out <- value
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// TimeoutChannel کانالی که بعد از timeout بسته می‌شود
func TimeoutChannel(duration time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		time.Sleep(duration)
		close(ch)
	}()
	return ch
}
