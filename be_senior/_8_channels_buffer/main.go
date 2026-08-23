// ============================================================================
// FILE: channels_buffer_guide.go
// TITLE: راهنمای کامل کانال‌های با بافر و بدون بافر در Go
// HOW TO RUN: go run channels_buffer_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - تفاوت اساسی کانال با بافر و بدون بافر
// ============================================================================
//
// کانال بدون بافر (Unbuffered Channel):
// - ظرفیت = 0
// - فرستنده تا زمانی که گیرنده نیامده بلاک می‌شود (همگام‌سازی همزمان)
// - مانند تحویل دست به دست: فرستنده منتظر می‌ماند تا گیرنده بیاید و بگیرد
// - تضمین می‌کند که داده قبل از ادامه دادن دریافت شده است
//
// کانال با بافر (Buffered Channel):
// - ظرفیت > 0
// - فرستنده تا پر شدن بافر بلاک می‌شود
// - گیرنده تا خالی شدن بافر بلاک می‌شود
// - مانند صندوق پستی: می‌توانی چند نامه بندازی و بعداً بروی برداری
//
// قانون طلایی:
// "Unbuffered channels provide synchronous communication,
//  Buffered channels provide asynchronous communication"
// ============================================================================

package _8_channels_buffer

import (
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// بخش 1: کانال بدون بافر (Unbuffered Channel) - عمیق
// ============================================================================

// demonstrateUnbufferedBasic اصول پایه کانال بدون بافر
func demonstrateUnbufferedBasic() {
	fmt.Println("\n=== Unbuffered Channel Basics ===")

	// ایجاد کانال بدون بافر (ظرفیت پیش‌فرض 0)
	ch := make(chan int)

	// فرستنده و گیرنده باید همزمان آماده باشند
	go func() {
		fmt.Println("Sender: sending 42...")
		ch <- 42 // بلاک می‌شود تا گیرنده بیاید
		fmt.Println("Sender: sent successfully!")
	}()

	// کمی صبر می‌کنیم تا فرستنده اجرا شود
	time.Sleep(100 * time.Millisecond)

	fmt.Println("Receiver: ready to receive...")
	value := <-ch // بلاک می‌شود تا داده بیاید
	fmt.Printf("Receiver: got %d\n", value)
}

// demonstrateUnbufferedBlocking نمایش بلاک شدن در کانال بدون بافر
func demonstrateUnbufferedBlocking() {
	fmt.Println("\n=== Unbuffered Channel Blocking Behavior ===")

	ch := make(chan string)

	// سناریو 1: فرستنده بدون گیرنده (بلاک می‌شود)
	fmt.Println("Scenario 1: Sender without receiver")
	go func() {
		fmt.Println("  Sender: about to send...")
		ch <- "hello"
		fmt.Println("  Sender: unblocked!")
	}()

	time.Sleep(200 * time.Millisecond)
	fmt.Println("  Main: now creating receiver")
	msg := <-ch
	fmt.Printf("  Main: received '%s'\n", msg)

	// سناریو 2: گیرنده بدون فرستنده (بلاک می‌شود)
	fmt.Println("\nScenario 2: Receiver without sender")
	go func() {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("  Sender: sending after delay...")
		ch <- "world"
	}()

	fmt.Println("  Main: waiting to receive...")
	msg2 := <-ch
	fmt.Printf("  Main: received '%s'\n", msg2)
}

// demonstrateUnbufferedSynchronization استفاده از کانال بدون بافر برای همگام‌سازی
func demonstrateUnbufferedSynchronization() {
	fmt.Println("\n=== Unbuffered Channel for Synchronization ===")

	done := make(chan bool)

	// تابعی که می‌خواهیم منتظر تمام شدنش بمانیم
	go func() {
		fmt.Println("Working...")
		time.Sleep(1 * time.Second)
		fmt.Println("Work done!")

		done <- true // سیگنال تمام شدن
	}()

	// منتظر سیگنال می‌مانیم
	<-done
	fmt.Println("Main: received completion signal")
}

// ============================================================================
// بخش 2: کانال با بافر (Buffered Channel) - عمیق
// ============================================================================

// demonstrateBufferedBasic اصول پایه کانال با بافر
func demonstrateBufferedBasic() {
	fmt.Println("\n=== Buffered Channel Basics ===")

	// کانال با بافر ظرفیت 3
	ch := make(chan int, 3)

	// فرستادن بدون گیرنده - تا ظرفیت بافر کار می‌کند
	fmt.Println("Sending values without receiver:")
	ch <- 10
	fmt.Println("  Sent 10, buffer: 1/3")
	ch <- 20
	fmt.Println("  Sent 20, buffer: 2/3")
	ch <- 30
	fmt.Println("  Sent 30, buffer: 3/3 (full)")

	// اگر سعی کنیم چهارمی بفرستیم، بلاک می‌شود
	// ch <- 40 // ❌ این خط برنامه را قفل می‌کند

	fmt.Println("\nReceiving values:")
	fmt.Printf("  Received: %d (buffer: 2/3)\n", <-ch)
	fmt.Printf("  Received: %d (buffer: 1/3)\n", <-ch)
	fmt.Printf("  Received: %d (buffer: 0/3)\n", <-ch)
}

// demonstrateBufferedBlocking نمایش بلاک شدن در کانال با بافر
func demonstrateBufferedBlocking() {
	fmt.Println("\n=== Buffered Channel Blocking Behavior ===")

	// کانال با ظرفیت 2
	ch := make(chan int, 2)

	// پر کردن بافر
	ch <- 1
	ch <- 2
	fmt.Println("Buffer is full (2/2)")

	// تلاش برای فرستادن سومین - این بلاک می‌شود
	go func() {
		fmt.Println("Sender: trying to send 3rd value...")
		ch <- 3 // بلاک می‌شود تا جا باز شود
		fmt.Println("Sender: sent 3rd value!")
	}()

	// صبر می‌کنیم تا فرستنده برسد
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main: receiving one value...")
	<-ch
	fmt.Println("Main: received, now buffer has 1 slot free")

	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main: receiving another...")
	<-ch

	time.Sleep(100 * time.Millisecond)
	fmt.Println("Main: done")
}

// demonstrateBufferedVsUnbuffered مقایسه مستقیم
func demonstrateBufferedVsUnbuffered() {
	fmt.Println("\n=== Buffered vs Unbuffered Comparison ===")

	// بدون بافر - نیازی به sleep نیست چون همگام است
	fmt.Println("--- Unbuffered Channel ---")
	unbuf := make(chan int)
	go func() {
		unbuf <- 42
		fmt.Println("  Unbuffered: sent")
	}()
	fmt.Printf("  Unbuffered: received %d\n", <-unbuf)

	// با بافر - فرستنده بلاک نمی‌شود
	fmt.Println("\n--- Buffered Channel ---")
	buf := make(chan int, 1)
	buf <- 42
	fmt.Println("  Buffered: sent (non-blocking)")
	fmt.Printf("  Buffered: received %d\n", <-buf)
}

// ============================================================================
// بخش 3: الگوهای عملی با کانال بدون بافر
// ============================================================================

// workerSignal کارگری که با سیگنال شروع به کار می‌کند
func workerSignal(id int, start <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	<-start // منتظر سیگنال شروع
	fmt.Printf("Worker %d: started working\n", id)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Worker %d: finished\n", id)
}

func demonstrateUnbufferedAsSignal() {
	fmt.Println("\n=== Unbuffered Channel as Signal ===")

	start := make(chan struct{})
	var wg sync.WaitGroup

	// راه‌اندازی 3 کارگر که همه منتظر سیگنال هستند
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go workerSignal(i, start, &wg)
	}

	fmt.Println("All workers waiting for start signal...")
	time.Sleep(500 * time.Millisecond)

	fmt.Println("Sending start signal...")
	close(start) // بستن کانال به همه سیگنال می‌دهد

	wg.Wait()
	fmt.Println("All workers completed")
}

// pingPong بازی ping-pong با کانال بدون بافر
func pingPong() {
	fmt.Println("\n=== Ping-Pong with Unbuffered Channel ===")

	ch := make(chan string)

	go func() {
		for i := 0; i < 3; i++ {
			msg := <-ch
			fmt.Printf("Pong received: %s\n", msg)
			ch <- "pong"
		}
	}()

	ch <- "ping"
	for i := 0; i < 3; i++ {
		msg := <-ch
		fmt.Printf("Ping received: %s\n", msg)
		if i < 2 {
			ch <- "ping"
		}
	}
}

// ============================================================================
// بخش 4: الگوهای عملی با کانال با بافر
// ============================================================================

// demonstrateBufferedAsQueue استفاده از کانال با بافر به عنوان صف
func demonstrateBufferedAsQueue() {
	fmt.Println("\n=== Buffered Channel as Queue ===")

	queue := make(chan string, 5)

	// تولید کننده - سریع
	go func() {
		for i := 1; i <= 5; i++ {
			queue <- fmt.Sprintf("task-%d", i)
			fmt.Printf("Producer: added task-%d (buffer: %d/5)\n", i, len(queue))
		}
		close(queue)
	}()

	// مصرف کننده - آهسته
	time.Sleep(100 * time.Millisecond)
	for task := range queue {
		fmt.Printf("Consumer: processing %s\n", task)
		time.Sleep(200 * time.Millisecond)
	}
}

// demonstrateBufferedBatchProcessing پردازش دسته‌ای با بافر
func demonstrateBufferedBatchProcessing() {
	fmt.Println("\n=== Batch Processing with Buffered Channel ===")

	const batchSize = 3
	ch := make(chan int, batchSize)

	// تولید کننده
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
			fmt.Printf("Produced: %d\n", i)
		}
		close(ch)
	}()

	// پردازش دسته‌ای
	batch := make([]int, 0, batchSize)
	for value := range ch {
		batch = append(batch, value)
		if len(batch) == batchSize {
			fmt.Printf("Processing batch: %v\n", batch)
			batch = batch[:0] // reset batch
			time.Sleep(100 * time.Millisecond)
		}
	}

	// پردازش باقی‌مانده
	if len(batch) > 0 {
		fmt.Printf("Processing final batch: %v\n", batch)
	}
}

// ============================================================================
// بخش 5: انتخاب بین بافر و بدون بافر - راهنمای تصمیم‌گیری
// ============================================================================

// demonstrateWhenToUseEach راهنمای انتخاب نوع کانال
func demonstrateWhenToUseEach() {
	fmt.Println("\n=== When to Use Which? ===")

	// سناریو 1: نیاز به همگام‌سازی دقیق (بدون بافر)
	fmt.Println("\nScenario 1: Exact synchronization (Unbuffered)")
	done := make(chan bool)
	go func() {
		// کار مهم
		time.Sleep(100 * time.Millisecond)
		done <- true // دقیقاً وقتی کار تمام شد سیگنال می‌دهد
	}()
	<-done // دقیقاً منتظر می‌ماند
	fmt.Println("  Synchronization guaranteed")

	// سناریو 2: بافر برای جذب نوسانات (Buffered)
	fmt.Println("\nScenario 2: Absorb bursts (Buffered)")
	bursty := make(chan int, 5)

	// نوسان ناگهانی درخواست
	for i := 0; i < 5; i++ {
		bursty <- i
	}
	fmt.Println("  All 5 requests accepted immediately")

	// پردازش آهسته
	go func() {
		for v := range bursty {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("  Processed: %d\n", v)
		}
	}()
	close(bursty)

	time.Sleep(600 * time.Millisecond)
}

// ============================================================================
// بخش 6: تله‌های رایج (Common Pitfalls)
// ============================================================================

func demonstrateCommonPitfalls() {
	fmt.Println("\n=== Common Pitfalls ===")

	// تله 1: Deadlock با کانال بدون بافر در همان گوروتین
	fmt.Println("\nPitfall 1: Deadlock in same goroutine")
	fmt.Println("  ❌ ch := make(chan int); ch <- 42 // deadlock!")
	fmt.Println("  ✅ Use buffered channel or separate goroutine")

	// تله 2: فراموش کردن بستن کانال
	fmt.Println("\nPitfall 2: Forgetting to close channel")
	fmt.Println("  ❌ range on channel that never closes → deadlock")
	fmt.Println("  ✅ close(ch) after done sending")

	// تله 3: بستن کانال از سمت گیرنده
	fmt.Println("\nPitfall 3: Closing from receiver side")
	fmt.Println("  ❌ Receiver closes channel (can cause panic)")
	fmt.Println("  ✅ Only sender should close")

	// تله 4: فرستادن به کانال بسته شده
	fmt.Println("\nPitfall 4: Send on closed channel")
	fmt.Println("  ❌ close(ch); ch <- 42 // panic!")
	fmt.Println("  ✅ Check with sync/atomic or use context")

	// تله 5: اندازه بافر اشتباه
	fmt.Println("\nPitfall 5: Wrong buffer size")
	fmt.Println("  ❌ make(chan int, 0) // same as unbuffered")
	fmt.Println("  ❌ make(chan int, 1000000) // memory waste")
	fmt.Println("  ✅ Choose buffer size based on expected burst")
}

// ============================================================================
// بخش 7: اندازه بافر - چگونه انتخاب کنیم؟
// ============================================================================

// demonstrateBufferSizeSelection انتخاب اندازه بافر مناسب
func demonstrateBufferSizeSelection() {
	fmt.Println("\n=== Buffer Size Selection Guide ===")

	// مثال: سیستم با نوسان ترافیک
	type Metric struct {
		burstSize       int
		processingTime  time.Duration
		recommendedSize int
	}

	metrics := []Metric{
		{burstSize: 5, processingTime: 100 * time.Millisecond, recommendedSize: 5},
		{burstSize: 100, processingTime: 10 * time.Millisecond, recommendedSize: 50},
		{burstSize: 1000, processingTime: 1 * time.Millisecond, recommendedSize: 500},
	}

	fmt.Println("Rule of thumb: buffer = expected_burst * safety_factor(0.5-1.0)")
	fmt.Println("\nExamples:")
	for _, m := range metrics {
		fmt.Printf("  Burst %d, processing %v → recommended buffer: %d\n",
			m.burstSize, m.processingTime, m.recommendedSize)
	}

	fmt.Println("\n⚠️  Warning: Too large buffer = memory waste + hidden backpressure")
	fmt.Println("✅ Start small, measure, adjust")
}

// ============================================================================
// بخش 8: تبدیل کانال بدون بافر به با بافر و برعکس
// ============================================================================

// unbufferedToBuffered تبدیل با الگوی تبدیل
func unbufferedToBuffered() {
	fmt.Println("\n=== Converting Patterns ===")

	// اگر کد با کانال بدون بافر نوشته شده، می‌توان با افزودن بافر
	// عملکرد را بدون تغییر منطق بهبود داد

	// نسخه بدون بافر (همگام)
	unbuf := make(chan int)
	go func() {
		for i := 0; i < 3; i++ {
			unbuf <- i
			fmt.Printf("Unbuffered: sent %d\n", i)
		}
		close(unbuf)
	}()

	for v := range unbuf {
		fmt.Printf("Unbuffered: received %d\n", v)
	}

	// نسخه با بافر (آسنکرون) - همان منطق اما کارآمدتر
	fmt.Println()
	buf := make(chan int, 3)
	go func() {
		for i := 0; i < 3; i++ {
			buf <- i
			fmt.Printf("Buffered: sent %d (non-blocking)\n", i)
		}
		close(buf)
	}()

	for v := range buf {
		fmt.Printf("Buffered: received %d\n", v)
	}
}

// ============================================================================
// بخش 9: تست و دیباگ کانال‌ها
// ============================================================================

// debugChannelLength نمایش وضعیت کانال
func debugChannelLength(ch chan int, name string) {
	fmt.Printf("📊 [%s] len=%d, cap=%d\n", name, len(ch), cap(ch))
}

func demonstrateChannelDebugging() {
	fmt.Println("\n=== Channel Debugging ===")

	ch := make(chan int, 5)

	debugChannelLength(ch, "empty channel")

	ch <- 1
	ch <- 2
	debugChannelLength(ch, "after 2 sends")

	<-ch
	debugChannelLength(ch, "after 1 receive")

	// تشخیص بلاک شدن با select و default
	select {
	case ch <- 3:
		fmt.Println("Send successful")
	default:
		fmt.Println("Channel is full, would block")
	}

	select {
	case v := <-ch:
		fmt.Printf("Receive successful: %d\n", v)
	default:
		fmt.Println("Channel is empty, would block")
	}
}

// ============================================================================
// بخش 10: مثال‌های ترکیبی و پیشرفته
// ============================================================================

// producerConsumerWithBuffer تولیدکننده-مصرف‌کننده با بافر
func producerConsumerWithBuffer() {
	fmt.Println("\n=== Producer-Consumer with Buffer ===")

	const (
		bufferSize   = 5
		numProducers = 3
		numConsumers = 2
		numTasks     = 20
	)

	tasks := make(chan int, bufferSize)
	var wg sync.WaitGroup

	// تولیدکننده‌ها
	for p := 1; p <= numProducers; p++ {
		wg.Add(1)
		go func(producerID int) {
			defer wg.Done()
			for i := 1; i <= numTasks/numProducers; i++ {
				task := producerID*100 + i
				tasks <- task
				fmt.Printf("Producer %d: produced %d (buffer: %d/%d)\n",
					producerID, task, len(tasks), bufferSize)
				time.Sleep(20 * time.Millisecond)
			}
		}(p)
	}

	// مصرف‌کننده‌ها
	for c := 1; c <= numConsumers; c++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			for task := range tasks {
				fmt.Printf("Consumer %d: consuming %d\n", consumerID, task)
				time.Sleep(50 * time.Millisecond)
			}
		}(c)
	}

	// منتظر تمام شدن تولیدکننده‌ها
	go func() {
		wg.Wait()
		close(tasks)
	}()

	wg.Wait()
	fmt.Println("All producers and consumers finished")
}

// demonstrateBufferFullHandling مدیریت وضعیت پر بودن بافر
func demonstrateBufferFullHandling() {
	fmt.Println("\n=== Handling Buffer Full Situations ===")

	ch := make(chan int, 2)

	// پر کردن بافر
	ch <- 1
	ch <- 2

	// روش 1: استفاده از select با default (non-blocking send)
	select {
	case ch <- 3:
		fmt.Println("Sent 3 successfully")
	default:
		fmt.Println("Buffer full, rejected 3 (backpressure)")
	}

	// روش 2: استفاده از select با timeout
	select {
	case ch <- 4:
		fmt.Println("Sent 4 successfully")
	case <-time.After(100 * time.Millisecond):
		fmt.Println("Timeout: could not send 4")
	}

	// روش 3: دریافت یک مقدار برای جا باز کردن
	fmt.Println("Making room by receiving one value...")
	<-ch
	fmt.Println("Now buffer has space")

	ch <- 5
	fmt.Println("Sent 5 successfully")
}

// ============================================================================
// بخش 11: جمع‌بندی و مقایسه نهایی
// ============================================================================

func main() {
	fmt.Println("========== BUFFERED & UNBUFFERED CHANNELS GUIDE ==========")
	fmt.Println("Complete guide with examples\n")

	// بخش 1: کانال بدون بافر
	demonstrateUnbufferedBasic()
	demonstrateUnbufferedBlocking()
	demonstrateUnbufferedSynchronization()

	// بخش 2: کانال با بافر
	demonstrateBufferedBasic()
	demonstrateBufferedBlocking()
	demonstrateBufferedVsUnbuffered()

	// بخش 3: الگوهای عملی بدون بافر
	demonstrateUnbufferedAsSignal()
	pingPong()

	// بخش 4: الگوهای عملی با بافر
	demonstrateBufferedAsQueue()
	demonstrateBufferedBatchProcessing()

	// بخش 5: راهنمای انتخاب
	demonstrateWhenToUseEach()

	// بخش 6: تله‌های رایج
	demonstrateCommonPitfalls()

	// بخش 7: انتخاب اندازه بافر
	demonstrateBufferSizeSelection()

	// بخش 8: تبدیل نوع کانال
	unbufferedToBuffered()

	// بخش 9: دیباگ
	demonstrateChannelDebugging()

	// بخش 10: مثال‌های پیشرفته
	producerConsumerWithBuffer()
	demonstrateBufferFullHandling()

	// بخش 11: جمع‌بندی نهایی
	fmt.Println("\n========== FINAL COMPARISON ==========")
	fmt.Println("| Feature              | Unbuffered        | Buffered           |")
	fmt.Println("|----------------------|-------------------|--------------------|")
	fmt.Println("| Capacity             | 0                 | >0                 |")
	fmt.Println("| Send blocks until    | Receiver ready    | Buffer has space   |")
	fmt.Println("| Receive blocks until | Sender ready      | Buffer has data    |")
	fmt.Println("| Synchronization      | Yes (exact)       | No (decoupled)     |")
	fmt.Println("| Use case             | Handshake, signal | Queue, burst       |")
	fmt.Println("| Latency              | Lower             | Higher (buffer)    |")
	fmt.Println("| Throughput           | Lower             | Higher             |")
	fmt.Println("| Memory               | No buffer         | Buffer allocated   |")

	fmt.Println("\n========== QUICK REFERENCE ==========")
	fmt.Println("| Need                           | Choose            |")
	fmt.Println("|--------------------------------|-------------------|")
	fmt.Println("| Synchronization / handshake    | Unbuffered        |")
	fmt.Println("| Signal between goroutines      | Unbuffered        |")
	fmt.Println("| Rate limiting / token bucket   | Unbuffered        |")
	fmt.Println("| Queue / task distribution      | Buffered          |")
	fmt.Println("| Absorb traffic bursts          | Buffered          |")
	fmt.Println("| Batch processing               | Buffered          |")
	fmt.Println("| Decouple producer from consumer| Buffered          |")

	fmt.Println("\n========== GOLDEN RULES ==========")
	fmt.Println("1. Unbuffered = synchronous communication")
	fmt.Println("2. Buffered = asynchronous communication")
	fmt.Println("3. Buffer size = expected burst × safety factor")
	fmt.Println("4. Never close channel from receiver side")
	fmt.Println("5. Always close channel from sender side")
	fmt.Println("6. Use select + default for non-blocking operations")
	fmt.Println("7. Monitor len(ch) for debugging, not for logic")
	fmt.Println("8. When in doubt, start with unbuffered (simpler reasoning)")
}

// ============================================================================
// بخش 12: توابع کمکی برای پروژه واقعی
// ============================================================================

// SafeSend ارسال امن به کانال با بررسی بسته بودن
func SafeSend(ch chan<- int, value int, timeout time.Duration) error {
	select {
	case ch <- value:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("send timeout after %v", timeout)
	}
}

// SafeReceive دریافت امن از کانال با timeout
func SafeReceive(ch <-chan int, timeout time.Duration) (int, error) {
	select {
	case value, ok := <-ch:
		if !ok {
			return 0, fmt.Errorf("channel closed")
		}
		return value, nil
	case <-time.After(timeout):
		return 0, fmt.Errorf("receive timeout after %v", timeout)
	}
}

// DrainChannel خالی کردن کامل یک کانال
func DrainChannel(ch <-chan int) []int {
	result := make([]int, 0)
	for {
		select {
		case value, ok := <-ch:
			if !ok {
				return result
			}
			result = append(result, value)
		default:
			return result
		}
	}
}

// IsChannelFull بررسی پر بودن کانال
func IsChannelFull(ch chan int) bool {
	return len(ch) == cap(ch)
}

// IsChannelEmpty بررسی خالی بودن کانال
func IsChannelEmpty(ch chan int) bool {
	return len(ch) == 0
}

// NewAutoScalingBuffer کانال با بافر خودتطبیق (برای موارد خاص)
func NewAutoScalingBuffer(initialSize, maxSize int) chan int {
	ch := make(chan int, initialSize)

	// مانیتورینگ و scaling خودکار (اختیاری)
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			usage := float64(len(ch)) / float64(cap(ch))
			if usage > 0.8 && cap(ch) < maxSize {
				// این فقط نمایشی است - در عمل نمی‌توان کانال را resize کرد
				fmt.Printf("⚠️  Channel usage %.0f%%, consider increasing buffer\n", usage*100)
			}
		}
	}()

	return ch
}
