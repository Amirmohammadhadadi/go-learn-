// ============================================================================
// FILE: race_condition_guide.go
// TITLE: راهنمای کامل Race Condition و استفاده از go test -race
// HOW TO RUN: go run race_condition_guide.go
// HOW TO TEST: go test -race -v race_condition_guide_test.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Race Condition چیست و چرا خطرناک است؟
// ============================================================================
//
// Race Condition:
// زمانی که دو یا چند گوروتین به طور همزمان به یک حافظه مشترک دسترسی دارند
// و حداقل یکی از آنها در حال نوشتن است، نتیجه غیرقابل پیش‌بینی می‌شود.
//
// مشکلات ناشی از Race Condition:
// 1. دیتا کوراپشن (Data Corruption): مقادیر نادرست و غیرمنتظره
// 2. کرش برنامه (Panic): مثلاً concurrent map writes
// 3. رفتار غیرقابل پیش‌بینی: نتیجه هر بار اجرا متفاوت است
// 4. باگ‌های سخت برای دیباگ: ممکن است فقط در production ظاهر شود
//
// قانون طلایی:
// "اگر دو گوروتین بدون هماهنگی به یک متغیر دسترسی داشته باشند،
//  برنامه شما race condition دارد - حتی اگر الان کار می‌کند!"
// ============================================================================

package _11_race_condition

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// بخش 1: مثال Race Condition (بدون محافظت)
// ============================================================================

// ❌ مثال 1: Race condition ساده روی یک عدد
func demonstrateRaceCondition() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💣 RACE CONDITION EXAMPLE - UNSAFE CODE")
	fmt.Println(stringsRepeat("=", 80))

	counter := 0
	var wg sync.WaitGroup

	// 1000 گوروتین که همزمان به counter اضافه می‌کنند
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // ❌ این خط race condition دارد
		}()
	}

	wg.Wait()

	fmt.Printf("Expected counter: 1000\n")
	fmt.Printf("Actual counter:   %d\n", counter)
	fmt.Println("💡 Note: Result is unpredictable! Run multiple times to see.")
	fmt.Println("   On your machine: go run -race race_condition_guide.go")
}

// ❌ مثال 2: Race condition در map (که باعث panic می‌شود)
func demonstrateMapRaceCondition() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🗺️ MAP RACE CONDITION - Concurrent Writes")
	fmt.Println(stringsRepeat("=", 80))

	// این کد معمولاً panic می‌کند: "concurrent map writes"
	// برای جلوگیری از panic در این دمو، آن را کامنت می‌کنیم

	fmt.Println("❌ The following code would panic:")
	fmt.Println(`   m := make(map[int]int)
   var wg sync.WaitGroup
   for i := 0; i < 100; i++ {
       wg.Add(1)
       go func(i int) {
           defer wg.Done()
           m[i] = i * 2  // ❌ concurrent map write
       }(i)
   }
   wg.Wait()`)

	fmt.Println("\n💡 Go will panic: 'fatal error: concurrent map writes'")
}

// ============================================================================
// بخش 2: راه‌حل‌های جلوگیری از Race Condition
// ============================================================================

// ✅ راه‌حل 1: استفاده از Mutex
func demonstrateMutexSolution() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔒 SOLUTION 1: MUTEX")
	fmt.Println(stringsRepeat("=", 80))

	var (
		counter int
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()   // قفل کردن
			counter++   // بخش بحرانی (critical section)
			mu.Unlock() // آزاد کردن قفل
		}()
	}

	wg.Wait()
	fmt.Printf("✅ Mutex result: %d\n", counter)
}

// ✅ راه‌حل 2: استفاده از RWMutex (بهتر برای خواندن زیاد)
func demonstrateRWMutexSolution() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📖 SOLUTION 2: RWMUTEX (Read/Write Mutex)")
	fmt.Println(stringsRepeat("=", 80))

	type SafeCounter struct {
		mu   sync.RWMutex
		data map[string]int
	}

	sc := SafeCounter{
		data: make(map[string]int),
	}

	var wg sync.WaitGroup

	// 10 نویسنده
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sc.mu.Lock() // قفل نوشتن
			sc.data[fmt.Sprintf("key%d", id)] = id
			sc.mu.Unlock()
		}(i)
	}

	// 100 خواننده (همزمان با نویسنده‌ها می‌توانند بخوانند)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sc.mu.RLock() // قفل خواندن (اجازه خواندن همزمان)
			_ = sc.data[fmt.Sprintf("key%d", id%10)]
			sc.mu.RUnlock()
		}(i)
	}

	wg.Wait()
	fmt.Printf("✅ RWMutex: %d writers, %d readers completed\n", 10, 100)
}

// ✅ راه‌حل 3: استفاده از Atomic Operations (برای انواع ساده)
func demonstrateAtomicSolution() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚡ SOLUTION 3: ATOMIC OPERATIONS")
	fmt.Println(stringsRepeat("=", 80))

	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1) // اتمیک و بدون قفل
		}()
	}

	wg.Wait()
	fmt.Printf("✅ Atomic result: %d\n", counter)
}

// ✅ راه‌حل 4: استفاده از Channel (ارتباط از طریق کانال)
func demonstrateChannelSolution() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📡 SOLUTION 4: CHANNELS (Share by communicating)")
	fmt.Println(stringsRepeat("=", 80))

	type Counter struct {
		value int
		ch    chan int
		done  chan struct{}
	}

	// گوروتین اختصاصی برای مدیریت شمارنده
	counter := &Counter{
		ch:   make(chan int),
		done: make(chan struct{}),
	}

	// مدیریت کننده شمارنده (تنها گوروتینی که به counter دسترسی دارد)
	go func() {
		for {
			select {
			case <-counter.ch:
				counter.value++
			case <-counter.done:
				return
			}
		}
	}()

	var wg sync.WaitGroup

	// 1000 گوروتین که از طریق کانال درخواست می‌دهند
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.ch <- 1 // درخواست افزایش
		}()
	}

	wg.Wait()
	close(counter.done)

	fmt.Printf("✅ Channel result: %d\n", counter.value)
}

// ============================================================================
// بخش 3: استفاده از go test -race
// ============================================================================

// این بخش برای نمایش در main است، اما تست‌ها در فایل جداگانه نوشته می‌شوند
func demonstrateGoTestRace() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🧪 GO TEST -RACE COMMAND")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
📌 HOW TO USE "go test -race":

1. Create a test file (e.g., counter_test.go):

   func TestCounter(t *testing.T) {
       counter := 0
       var wg sync.WaitGroup
       
       for i := 0; i < 100; i++ {
           wg.Add(1)
           go func() {
               defer wg.Done()
               counter++  // Race condition!
           }()
       }
       wg.Wait()
   }

2. Run with race detector:
   $ go test -race -v

3. Output if race detected:
   ==================
   WARNING: DATA RACE
   Write at 0x00c000... by goroutine 7:
   Read at 0x00c000... by goroutine 8:
   ==================

💡 TIPS:
   • Always run tests with -race in CI/CD
   • -race makes tests slower (2-10x), but worth it
   • Race detector finds bugs that are hard to reproduce
   • Fix all race conditions before production
`)
}

// ============================================================================
// بخش 4: الگوهای پیشرفته برای جلوگیری از Race
// ============================================================================

// ✅ الگوی 1: Single Flight (اجرای یک بار)
type SingleFlight struct {
	mu     sync.Mutex
	inflig map[string]*call
}

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

func (sf *SingleFlight) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	sf.mu.Lock()
	if sf.inflig == nil {
		sf.inflig = make(map[string]*call)
	}

	if c, ok := sf.inflig[key]; ok {
		sf.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	sf.inflig[key] = c
	sf.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	sf.mu.Lock()
	delete(sf.inflig, key)
	sf.mu.Unlock()

	return c.val, c.err
}

func demonstrateSingleFlight() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚀 ADVANCED PATTERN: SINGLE FLIGHT")
	fmt.Println(stringsRepeat("=", 80))

	var sf SingleFlight
	var wg sync.WaitGroup

	expensiveOperation := func() (interface{}, error) {
		fmt.Println("  Executing expensive operation (only once!)")
		time.Sleep(100 * time.Millisecond)
		return "result", nil
	}

	// 10 گوروتین همزمان درخواست می‌دهند
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			result, _ := sf.Do("key", expensiveOperation)
			fmt.Printf("  Goroutine %d got: %v\n", id, result)
		}(i)
	}

	wg.Wait()
}

// ✅ الگوی 2: Thread-Local Storage با sync.Pool
type UserSession struct {
	ID   string
	Name string
	Data map[string]interface{}
}

func demonstrateSyncPool() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🏊 ADVANCED PATTERN: SYNC.POOL")
	fmt.Println(stringsRepeat("=", 80))

	pool := sync.Pool{
		New: func() interface{} {
			return &UserSession{
				Data: make(map[string]interface{}),
			}
		},
	}

	var wg sync.WaitGroup

	// هر گوروتین از pool یک شیء می‌گیرد (بدون race)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			session := pool.Get().(*UserSession)
			defer pool.Put(session)

			session.ID = fmt.Sprintf("user%d", id)
			session.Name = fmt.Sprintf("Name%d", id)
			session.Data["timestamp"] = time.Now()

			fmt.Printf("  Session %d: %+v\n", id, session.ID)
		}(i)
	}

	wg.Wait()
}

// ============================================================================
// بخش 5: ابزارهای تشخیص Race Condition
// ============================================================================

func demonstrateRaceDetectionTools() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔍 RACE DETECTION TOOLS")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
📦 TOOLS FOR RACE DETECTION:

1. go test -race (built-in)
   $ go test -race ./...
   
2. go run -race
   $ go run -race main.go

3. go build -race
   $ go build -race -o app main.go
   $ ./app

4. During CI/CD:
   go test -race -cover -v ./...

5. VS Code settings:
   "go.testFlags": ["-race"]

💡 LIMITATIONS OF RACE DETECTOR:
   • Only detects races that actually happen during execution
   • May miss races that require specific timing
   • Slower execution (2-10x)
   • Memory usage increases (5-10x)
   • Not available on all architectures (Windows 32-bit limited)

💡 BEST PRACTICES:
   • Run tests with -race on every PR
   • Include -race in CI pipeline
   • Test with high concurrency (use -count flag)
   • Don't ignore race detector warnings
   • Fix races as soon as they're found
`)
}

// ============================================================================
// بخش 6: مقایسه راه‌حل‌های مختلف
// ============================================================================

func compareSolutions() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📊 COMPARISON OF SOLUTIONS")
	fmt.Println(stringsRepeat("=", 80))

	benchmark := func(name string, fn func()) {
		start := time.Now()
		fn()
		elapsed := time.Since(start)
		fmt.Printf("  %-20s: %v\n", name, elapsed)
	}

	// تست Mutex
	benchmark("Mutex", func() {
		var counter int
		var mu sync.Mutex
		var wg sync.WaitGroup
		for i := 0; i < 10000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mu.Lock()
				counter++
				mu.Unlock()
			}()
		}
		wg.Wait()
	})

	// تست Atomic
	benchmark("Atomic", func() {
		var counter int64
		var wg sync.WaitGroup
		for i := 0; i < 10000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				atomic.AddInt64(&counter, 1)
			}()
		}
		wg.Wait()
	})

	// تست Channel
	benchmark("Channel", func() {
		ch := make(chan struct{}, 100)
		var wg sync.WaitGroup
		for i := 0; i < 10000; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ch <- struct{}{}
				<-ch
			}()
		}
		wg.Wait()
	})

	fmt.Println("\n💡 RECOMMENDATIONS:")
	fmt.Println("  • Simple counters → Atomic")
	fmt.Println("  • Complex shared data → Mutex")
	fmt.Println("  • Many readers, few writers → RWMutex")
	fmt.Println("  • Communication between goroutines → Channels")
	fmt.Println("  • Object pooling → sync.Pool")
}

// ============================================================================
// بخش 7: Race Condition در سناریوهای واقعی
// ============================================================================

func demonstrateRealWorldRaces() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🌍 REAL-WORLD RACE CONDITION SCENARIOS")
	fmt.Println(stringsRepeat("=", 80))

	// سناریو 1: Cache بدون محافظت
	fmt.Println("\n📌 Scenario 1: Unsafe Cache")

	type UnsafeCache struct {
		data map[string]string
	}

	cache := &UnsafeCache{
		data: make(map[string]string),
	}

	var wg sync.WaitGroup

	// نویسنده
	wg.Add(1)
	go func() {
		defer wg.Done()
		cache.data["key"] = "value" // ❌ race
	}()

	// خواننده
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = cache.data["key"] // ❌ race
	}()

	wg.Wait()
	fmt.Println("  ❌ Unsafe cache has race condition!")

	// سناریو 2: Slice بدون محافظت
	fmt.Println("\n📌 Scenario 2: Unsafe Slice")

	slice := make([]int, 0, 100)

	wg.Add(2)
	go func() {
		defer wg.Done()
		slice = append(slice, 1) // ❌ race
	}()

	go func() {
		defer wg.Done()
		slice = append(slice, 2) // ❌ race
	}()

	wg.Wait()
	fmt.Println("  ❌ Concurrent append to slice is unsafe!")
}

// ============================================================================
// بخش 8: الگوهای درست و غلط
// ============================================================================

func demonstrateCorrectVsWrong() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("✅❌ RACE CONDITION - CORRECT VS WRONG")
	fmt.Println(stringsRepeat("=", 80))

	// ❌ غلط: دسترسی همزمان بدون محافظت
	fmt.Println("\n❌ Wrong: Concurrent access without protection")
	fmt.Println(`   var counter int
   go func() { counter++ }()
   go func() { counter++ }()`)

	// ✅ درست: استفاده از Mutex
	fmt.Println("\n✅ Correct: Mutex protection")
	fmt.Println(`   var mu sync.Mutex
   go func() { mu.Lock(); counter++; mu.Unlock() }()
   go func() { mu.Lock(); counter++; mu.Unlock() }()`)

	// ❌ غلط: کپی کردن Mutex
	fmt.Println("\n❌ Wrong: Copying mutex")
	fmt.Println(`   type Counter struct { mu sync.Mutex }
   c1 := Counter{}
   c2 := c1  // ❌ mutex copied!`)

	// ✅ درست: استفاده از اشاره‌گر
	fmt.Println("\n✅ Correct: Using pointer for mutex")
	fmt.Println(`   type Counter struct { mu sync.Mutex }
   c1 := &Counter{}
   c2 := c1  // ✅ pointer copied, same mutex`)

	// ❌ غلط: ذخیره اشاره‌گر به متغیر محلی در گوروتین
	fmt.Println("\n❌ Wrong: Capturing loop variable in goroutine")
	fmt.Println(`   for i := 0; i < 10; i++ {
       go func() {
           fmt.Println(i)  // ❌ race on i
       }()
   }`)

	// ✅ درست: پاس دادن به عنوان پارامتر
	fmt.Println("\n✅ Correct: Pass as parameter")
	fmt.Println(`   for i := 0; i < 10; i++ {
       go func(val int) {
           fmt.Println(val)  // ✅ safe
       }(i)
   }`)
}

// ============================================================================
// بخش 9: جمع‌بندی و نکات نهایی
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 RACE CONDITION IN GO - Complete Guide")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: مثال Race Condition
	demonstrateRaceCondition()
	demonstrateMapRaceCondition()

	// بخش 2: راه‌حل‌ها
	demonstrateMutexSolution()
	demonstrateRWMutexSolution()
	demonstrateAtomicSolution()
	demonstrateChannelSolution()

	// بخش 3: go test -race
	demonstrateGoTestRace()

	// بخش 4: الگوهای پیشرفته
	demonstrateSingleFlight()
	demonstrateSyncPool()

	// بخش 5: ابزارها
	demonstrateRaceDetectionTools()

	// بخش 6: مقایسه
	compareSolutions()

	// بخش 7: سناریوهای واقعی
	demonstrateRealWorldRaces()

	// بخش 8: درست و غلط
	demonstrateCorrectVsWrong()

	// بخش 9: جمع‌بندی
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 RACE CONDITION QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ PATTERN              │ WHEN TO USE                             │")
	fmt.Println("├──────────────────────┼─────────────────────────────────────────┤")
	fmt.Println("│ sync.Mutex           │ General purpose mutual exclusion        │")
	fmt.Println("│ sync.RWMutex         │ Many readers, few writers               │")
	fmt.Println("│ atomic.*             │ Simple counters, flags, integers        │")
	fmt.Println("│ channels             │ Communication between goroutines        │")
	fmt.Println("│ sync.Once            │ One-time initialization                 │")
	fmt.Println("│ sync.Pool            │ Reusable objects, reduce allocations    │")
	fmt.Println("│ Single Flight        │ Prevent duplicate work                  │")
	fmt.Println("└──────────────────────┴─────────────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Always run tests with -race flag")
	fmt.Println("  2. Never ignore race detector warnings")
	fmt.Println("  3. Mutexes cannot be copied - use pointers")
	fmt.Println("  4. Use atomic for simple counters")
	fmt.Println("  5. Use channels for communication patterns")
	fmt.Println("  6. Don't share memory; communicate instead")
	fmt.Println("  7. Race conditions can be timing-dependent - test many times")
	fmt.Println("  8. Fix races as soon as they're found")
	fmt.Println("  9. Include -race in CI/CD pipeline")
	fmt.Println("  10. When in doubt, add synchronization")

	fmt.Println("\n⚠️  COMMON PITFALLS:")
	fmt.Println("  • Forgetting to lock/unlock mutex")
	fmt.Println("  • Copying mutex (use pointer instead)")
	fmt.Println("  • Concurrent map writes (use sync.Map or mutex)")
	fmt.Println("  • Closing channel multiple times")
	fmt.Println("  • Sending on closed channel")
	fmt.Println("  • Capturing loop variables in goroutines")
}

// ============================================================================
// بخش 10: توابع کمکی
// ============================================================================

func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// ============================================================================
// FILE: race_condition_guide_test.go
// این فایل تست برای نمایش go test -race است
// ============================================================================

/*
package main

import (
	"sync"
	"testing"
)

// ❌ تستی که race condition دارد
func TestCounterWithRace(t *testing.T) {
	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // ❌ RACE CONDITION
		}()
	}

	wg.Wait()

	if counter != 100 {
		t.Errorf("Expected 100, got %d", counter)
	}
}

// ✅ تستی که race condition ندارد (با Mutex)
func TestCounterSafe(t *testing.T) {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()

	if counter != 100 {
		t.Errorf("Expected 100, got %d", counter)
	}
}

// ✅ تستی با atomic
func TestCounterAtomic(t *testing.T) {
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()

	if counter != 100 {
		t.Errorf("Expected 100, got %d", counter)
	}
}

// اجرای تست با race detector:
// $ go test -race -v
//
// خروجی برای تست اول (با race):
// ==================
// WARNING: DATA RACE
// Write at 0x00c0000a4008 by goroutine 8:
//   main.TestCounterWithRace.func1()
//       /path/race_condition_guide_test.go:20 +0x3c
// ==================
*/
