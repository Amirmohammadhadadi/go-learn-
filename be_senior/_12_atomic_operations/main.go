// ============================================================================
// FILE: atomic_operations_guide.go
// TITLE: راهنمای کامل عملیات اتمیک در Go - بدون قفل و پرسرعت
// HOW TO RUN: go run atomic_operations_guide.go
// HOW TO TEST: go test -race -v atomic_operations_guide_test.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - عملیات اتمیک چیست و چه تفاوتی با Mutex دارد؟
// ============================================================================
//
// عملیات اتمیک (Atomic Operations):
// - عملیاتی که در سطح CPU به صورت غیرقابل تقسیم (indivisible) اجرا می‌شود
// - بدون نیاز به قفل (lock-free) انجام می‌شود
// - بسیار سریع‌تر از Mutex (چندین برابر)
//
// مقایسه با Mutex:
// ┌─────────────────────────────────────────────────────────────────┐
// │ ویژگی         │ Atomic Operations │ Mutex                       │
// ├─────────────────────────────────────────────────────────────────┤
// │ سرعت          │ خیلی سریع (~ns)   │ نسبتاً سریع (~μs)            │
// │ هزینه         │ بدون قفل          │ قفل سنگین                    │
// │ کاربرد        │ انواع ساده (int)  │ هر نوع داده (struct, map)   │
// │ Blocking      │ غیر-blocking      │ blocking                     │
// │ پیچیدگی       │ ساده              │ ساده                         │
// └─────────────────────────────────────────────────────────────────┘
//
// قانون طلایی:
// "برای شمارنده‌ها، فلگ‌ها و مقادیر ساده از atomic استفاده کن.
//  برای ساختارهای پیچیده و عملیات چندمتغیره از Mutex استفاده کن."
// ============================================================================

package __atomic_operations

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// ============================================================================
// بخش 1: عملیات اتمیک پایه روی اعداد صحیح
// ============================================================================

func demonstrateBasicAtomic() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔢 BASIC ATOMIC OPERATIONS ON INTEGERS")
	fmt.Println(stringsRepeat("=", 80))

	var counter int64
	var wg sync.WaitGroup

	// افزودن اتمیک
	fmt.Println("1. Atomic Add:")
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()
	fmt.Printf("   Counter after 10 adds: %d\n", counter)

	// بارگذاری اتمیک (Load)
	fmt.Println("\n2. Atomic Load:")
	value := atomic.LoadInt64(&counter)
	fmt.Printf("   Loaded value: %d\n", value)

	// ذخیره اتمیک (Store)
	fmt.Println("\n3. Atomic Store:")
	atomic.StoreInt64(&counter, 100)
	fmt.Printf("   After store: %d\n", atomic.LoadInt64(&counter))

	// تعویض اتمیک (Swap)
	fmt.Println("\n4. Atomic Swap:")
	old := atomic.SwapInt64(&counter, 200)
	fmt.Printf("   Old: %d, New: %d\n", old, atomic.LoadInt64(&counter))

	// مقایسه و تعویض (Compare And Swap - CAS)
	fmt.Println("\n5. Atomic Compare And Swap (CAS):")
	success := atomic.CompareAndSwapInt64(&counter, 200, 300)
	fmt.Printf("   CAS success: %v, New value: %d\n", success, atomic.LoadInt64(&counter))

	// CAS ناموفق
	success = atomic.CompareAndSwapInt64(&counter, 999, 500)
	fmt.Printf("   CAS with wrong expected: %v, Value unchanged: %d\n", success, atomic.LoadInt64(&counter))
}

// ============================================================================
// بخش 2: عملیات اتمیک روی انواع مختلف
// ============================================================================

func demonstrateAtomicTypes() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📦 ATOMIC OPERATIONS ON DIFFERENT TYPES")
	fmt.Println(stringsRepeat("=", 80))

	// int32
	fmt.Println("\n1. int32 operations:")
	var i32 int32 = 0
	atomic.AddInt32(&i32, 5)
	fmt.Printf("   AddInt32: %d\n", i32)

	// uint32
	fmt.Println("\n2. uint32 operations:")
	var u32 uint32 = 10
	atomic.AddUint32(&u32, ^uint32(0)) // این معادل -1 است
	fmt.Printf("   AddUint32 (decrement): %d\n", u32)

	// uint64
	fmt.Println("\n3. uint64 operations:")
	var u64 uint64 = 100
	atomic.StoreUint64(&u64, 200)
	fmt.Printf("   StoreUint64: %d\n", atomic.LoadUint64(&u64))

	// uintptr
	fmt.Println("\n4. uintptr operations:")
	var ptr uintptr = uintptr(unsafe.Pointer(&i32))
	atomic.AddUintptr(&ptr, unsafe.Sizeof(i32))
	fmt.Printf("   AddUintptr: %x\n", ptr)

	// unsafe.Pointer
	fmt.Println("\n5. unsafe.Pointer operations:")
	var p unsafe.Pointer
	val1 := 42
	val2 := 100
	atomic.StorePointer(&p, unsafe.Pointer(&val1))
	fmt.Printf("   Stored pointer to: %d\n", *(*int)(atomic.LoadPointer(&p)))

	atomic.CompareAndSwapPointer(&p, unsafe.Pointer(&val1), unsafe.Pointer(&val2))
	fmt.Printf("   After CAS: %d\n", *(*int)(p))
}

// ============================================================================
// بخش 3: Atomic Bool (با استفاده از int32)
// ============================================================================

type AtomicBool struct {
	value int32
}

func NewAtomicBool(initial bool) *AtomicBool {
	val := int32(0)
	if initial {
		val = 1
	}
	return &AtomicBool{value: val}
}

func (b *AtomicBool) Set(value bool) {
	var val int32 = 0
	if value {
		val = 1
	}
	atomic.StoreInt32(&b.value, val)
}

func (b *AtomicBool) Get() bool {
	return atomic.LoadInt32(&b.value) == 1
}

func (b *AtomicBool) CompareAndSwap(old, new bool) bool {
	var oldVal, newVal int32 = 0, 0
	if old {
		oldVal = 1
	}
	if new {
		newVal = 1
	}
	return atomic.CompareAndSwapInt32(&b.value, oldVal, newVal)
}

func (b *AtomicBool) Toggle() bool {
	for {
		current := b.Get()
		next := !current
		if b.CompareAndSwap(current, next) {
			return next
		}
	}
}

func demonstrateAtomicBool() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 ATOMIC BOOL - Custom Implementation")
	fmt.Println(stringsRepeat("=", 80))

	flag := NewAtomicBool(false)
	var wg sync.WaitGroup

	// 10 گوروتین که سعی می‌کنند flag را toggle کنند
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			newVal := flag.Toggle()
			fmt.Printf("  Goroutine %d toggled to: %v\n", id, newVal)
		}(i)
	}

	wg.Wait()
	fmt.Printf("\nFinal value: %v\n", flag.Get())
}

// ============================================================================
// بخش 4: Atomic Counter (شمارنده اتمیک پیشرفته)
// ============================================================================

type AtomicCounter struct {
	value int64
}

func NewAtomicCounter() *AtomicCounter {
	return &AtomicCounter{}
}

func (c *AtomicCounter) Inc() int64 {
	return atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Dec() int64 {
	return atomic.AddInt64(&c.value, -1)
}

func (c *AtomicCounter) Add(delta int64) int64 {
	return atomic.AddInt64(&c.value, delta)
}

func (c *AtomicCounter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

func (c *AtomicCounter) Set(value int64) {
	atomic.StoreInt64(&c.value, value)
}

func (c *AtomicCounter) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&c.value, old, new)
}

func (c *AtomicCounter) Reset() {
	atomic.StoreInt64(&c.value, 0)
}

func demonstrateAtomicCounter() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📊 ATOMIC COUNTER - Production Ready")
	fmt.Println(stringsRepeat("=", 80))

	counter := NewAtomicCounter()
	var wg sync.WaitGroup

	// 1000 گوروتین، هر کدام 100 بار افزایش
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.Inc()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Final counter value: %d (expected: 100000)\n", counter.Get())

	// تست CAS
	fmt.Println("\nCAS Example:")
	counter.Set(50)
	fmt.Printf("  Current: %d\n", counter.Get())

	if counter.CompareAndSwap(50, 100) {
		fmt.Printf("  CAS successful, new value: %d\n", counter.Get())
	}

	if !counter.CompareAndSwap(50, 200) {
		fmt.Printf("  CAS failed (expected 50 but value is %d)\n", counter.Get())
	}
}

// ============================================================================
// بخش 5: Atomic Value (برای انواع دلخواه)
// ============================================================================

func demonstrateAtomicValue() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💎 ATOMIC VALUE - Store Any Type Safely")
	fmt.Println(stringsRepeat("=", 80))

	type Config struct {
		Host  string
		Port  int
		Debug bool
	}

	var config atomic.Value

	// تنظیم اولیه کانفیگ
	initialConfig := &Config{
		Host:  "localhost",
		Port:  8080,
		Debug: true,
	}
	config.Store(initialConfig)
	fmt.Printf("Initial config: %+v\n", config.Load().(*Config))

	// به‌روزرسانی همزمان کانفیگ توسط چند گوروتین
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(version int) {
			defer wg.Done()
			newConfig := &Config{
				Host:  fmt.Sprintf("server-%d", version),
				Port:  8080 + version,
				Debug: version%2 == 0,
			}
			config.Store(newConfig)
			time.Sleep(10 * time.Millisecond)
		}(i)
	}

	wg.Wait()
	fmt.Printf("Final config: %+v\n", config.Load().(*Config))

	// نکته مهم: atomic.Value نمی‌تواند دو نوع مختلف را ذخیره کند
	fmt.Println("\n⚠️  Note: atomic.Value requires same type for all stores")
}

// ============================================================================
// بخش 6: Atomic Flag (علامت‌دهنده اتمیک)
// ============================================================================

type AtomicFlag struct {
	value int32
}

func NewAtomicFlag() *AtomicFlag {
	return &AtomicFlag{}
}

func (f *AtomicFlag) Set() {
	atomic.StoreInt32(&f.value, 1)
}

func (f *AtomicFlag) Clear() {
	atomic.StoreInt32(&f.value, 0)
}

func (f *AtomicFlag) IsSet() bool {
	return atomic.LoadInt32(&f.value) == 1
}

func (f *AtomicFlag) TestAndSet() bool {
	// اگر 0 بود به 1 تبدیل کن و true برگردان
	return atomic.CompareAndSwapInt32(&f.value, 0, 1)
}

func (f *AtomicFlag) TestAndClear() bool {
	// اگر 1 بود به 0 تبدیل کن و true برگردان
	return atomic.CompareAndSwapInt32(&f.value, 1, 0)
}

func demonstrateAtomicFlag() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚩 ATOMIC FLAG - Test-And-Set Pattern")
	fmt.Println(stringsRepeat("=", 80))

	flag := NewAtomicFlag()
	var wg sync.WaitGroup

	// 10 گوروتین سعی می‌کنند flag را set کنند
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if flag.TestAndSet() {
				fmt.Printf("  Goroutine %d: I set the flag!\n", id)
			} else {
				fmt.Printf("  Goroutine %d: Flag was already set\n", id)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("\nFlag is set: %v\n", flag.IsSet())

	// Clear و دوباره تست
	flag.Clear()
	fmt.Printf("After clear: %v\n", flag.IsSet())

	if flag.TestAndSet() {
		fmt.Println("Successfully set the flag again")
	}
}

// ============================================================================
// بخش 7: Atomic Integer (عملیات پیشرفته)
// ============================================================================

type AtomicInt struct {
	value int64
}

func NewAtomicInt(initial int64) *AtomicInt {
	return &AtomicInt{value: initial}
}

func (a *AtomicInt) Get() int64 {
	return atomic.LoadInt64(&a.value)
}

func (a *AtomicInt) Set(value int64) {
	atomic.StoreInt64(&a.value, value)
}

func (a *AtomicInt) Add(delta int64) int64 {
	return atomic.AddInt64(&a.value, delta)
}

func (a *AtomicInt) Increment() int64 {
	return a.Add(1)
}

func (a *AtomicInt) Decrement() int64 {
	return a.Add(-1)
}

func (a *AtomicInt) Swap(new int64) int64 {
	return atomic.SwapInt64(&a.value, new)
}

func (a *AtomicInt) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&a.value, old, new)
}

func (a *AtomicInt) GetAndAdd(delta int64) int64 {
	for {
		current := a.Get()
		if a.CompareAndSwap(current, current+delta) {
			return current
		}
	}
}

func (a *AtomicInt) GetAndIncrement() int64 {
	return a.GetAndAdd(1)
}

func (a *AtomicInt) AddAndGet(delta int64) int64 {
	return a.Add(delta)
}

func demonstrateAtomicInt() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🔧 ATOMIC INT - Advanced Operations")
	fmt.Println(stringsRepeat("=", 80))

	num := NewAtomicInt(0)
	var wg sync.WaitGroup

	// Increment عملیات
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			num.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("After 100 increments: %d\n", num.Get())

	// GetAndIncrement
	old := num.GetAndIncrement()
	fmt.Printf("GetAndIncrement: old=%d, new=%d\n", old, num.Get())

	// AddAndGet
	newVal := num.AddAndGet(10)
	fmt.Printf("AddAndGet(10): %d\n", newVal)

	// Swap
	old = num.Swap(1000)
	fmt.Printf("Swap(1000): old=%d, new=%d\n", old, num.Get())
}

// ============================================================================
// بخش 8: Atomic vs Mutex - مقایسه عملکرد
// ============================================================================

func benchmarkAtomicVsMutex() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚡ ATOMIC vs MUTEX - Performance Comparison")
	fmt.Println(stringsRepeat("=", 80))

	const iterations = 1000000
	const goroutines = 100

	// تست Atomic
	atomicBench := func() time.Duration {
		var counter int64
		var wg sync.WaitGroup

		start := time.Now()
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < iterations/goroutines; j++ {
					atomic.AddInt64(&counter, 1)
				}
			}()
		}
		wg.Wait()
		return time.Since(start)
	}

	// تست Mutex
	mutexBench := func() time.Duration {
		var counter int64
		var mu sync.Mutex
		var wg sync.WaitGroup

		start := time.Now()
		for i := 0; i < goroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < iterations/goroutines; j++ {
					mu.Lock()
					counter++
					mu.Unlock()
				}
			}()
		}
		wg.Wait()
		return time.Since(start)
	}

	atomicTime := atomicBench()
	mutexTime := mutexBench()

	fmt.Printf("Atomic time: %v\n", atomicTime)
	fmt.Printf("Mutex time:  %v\n", mutexTime)
	fmt.Printf("Atomic is %.2fx faster\n", float64(mutexTime)/float64(atomicTime))

	fmt.Println("\n💡 Conclusion: Atomic is significantly faster for simple operations")
}

// ============================================================================
// بخش 9: الگوهای پیشرفته با Atomic
// ============================================================================

// الگوی 1: Sequence Generator (تولید کننده ID اتمیک)
type AtomicSequence struct {
	value int64
}

func NewAtomicSequence(initial int64) *AtomicSequence {
	return &AtomicSequence{value: initial}
}

func (s *AtomicSequence) Next() int64 {
	return atomic.AddInt64(&s.value, 1)
}

func (s *AtomicSequence) Current() int64 {
	return atomic.LoadInt64(&s.value)
}

// الگوی 2: Limit Counter (محدود کننده نرخ)
type RateLimiter struct {
	limit     int64
	current   int64
	resetTime int64
}

func NewRateLimiter(limit int64) *RateLimiter {
	return &RateLimiter{
		limit:     limit,
		current:   0,
		resetTime: time.Now().Unix(),
	}
}

func (r *RateLimiter) Allow() bool {
	now := time.Now().Unix()
	resetTime := atomic.LoadInt64(&r.resetTime)

	// هر 1 ثانیه یکبار ریست می‌شود
	if now-resetTime >= 1 {
		if atomic.CompareAndSwapInt64(&r.resetTime, resetTime, now) {
			atomic.StoreInt64(&r.current, 0)
		}
	}

	// بررسی محدودیت
	for {
		current := atomic.LoadInt64(&r.current)
		if current >= r.limit {
			return false
		}
		if atomic.CompareAndSwapInt64(&r.current, current, current+1) {
			return true
		}
	}
}

// الگوی 3: Lazy Initialization با atomic.Value
type LazyLoader struct {
	value  atomic.Value
	loader func() interface{}
}

func NewLazyLoader(loader func() interface{}) *LazyLoader {
	return &LazyLoader{
		loader: loader,
	}
}

func (l *LazyLoader) Get() interface{} {
	if v := l.value.Load(); v != nil {
		return v
	}

	// بارگذاری با قفل اتمیک (Double-checked locking)
	v := l.loader()
	if l.value.CompareAndSwap(nil, v) {
		return v
	}
	return l.value.Load()
}

func demonstrateAdvancedPatterns() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎨 ADVANCED ATOMIC PATTERNS")
	fmt.Println(stringsRepeat("=", 80))

	// 1. Sequence Generator
	fmt.Println("\n1. Atomic Sequence Generator:")
	seq := NewAtomicSequence(0)
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				fmt.Printf("  G%d got ID: %d\n", id, seq.Next())
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("  Final sequence value: %d\n", seq.Current())

	// 2. Rate Limiter
	fmt.Println("\n2. Rate Limiter (10 requests/second):")
	limiter := NewRateLimiter(10)

	allowed := 0
	denied := 0
	for i := 0; i < 20; i++ {
		if limiter.Allow() {
			allowed++
		} else {
			denied++
		}
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Printf("  Allowed: %d, Denied: %d\n", allowed, denied)

	// 3. Lazy Initialization
	fmt.Println("\n3. Lazy Initialization:")
	loader := NewLazyLoader(func() interface{} {
		fmt.Println("  Loading expensive resource...")
		time.Sleep(100 * time.Millisecond)
		return "expensive_resource"
	})

	// اولین دسترسی - بارگذاری می‌کند
	fmt.Printf("  First get: %v\n", loader.Get())

	// دسترسی‌های بعدی - از کش استفاده می‌کند
	fmt.Printf("  Second get: %v\n", loader.Get())
	fmt.Printf("  Third get: %v\n", loader.Get())
}

// ============================================================================
// بخش 10: نکات مهم و اشتباهات رایج
// ============================================================================

func demonstrateAtomicPitfalls() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("⚠️ ATOMIC PITFALLS AND BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	// اشتباه 1: فراموش کردن بارگذاری اتمیک
	fmt.Println("\n❌ Pitfall 1: Not using atomic Load")
	fmt.Println(`   var counter int64
   go func() { atomic.AddInt64(&counter, 1) }()
   value := counter  // ❌这可能不是最新的值！`)

	// اشتباه 2: کپی کردن atomic.Value
	fmt.Println("\n❌ Pitfall 2: Copying atomic.Value")
	fmt.Println(`   var v1 atomic.Value
   v1.Store(42)
   v2 := v1  // ❌ Copying atomic.Value is not safe`)

	// اشتباه 3: استفاده از نوع متفاوت در atomic.Value
	fmt.Println("\n❌ Pitfall 3: Different types in atomic.Value")
	fmt.Println(`   var v atomic.Value
   v.Store(42)      // int
   v.Store("hello") // ❌ panic: different type`)

	// اشتباه 4: انتظار ترتیب (ordering)
	fmt.Println("\n❌ Pitfall 4: Assuming memory ordering")
	fmt.Println(`   var flag int32
   var data string
   
   // Goroutine 1
   data = "hello"
   atomic.StoreInt32(&flag, 1)  // No guarantee of ordering!
   
   // Goroutine 2
   if atomic.LoadInt32(&flag) == 1 {
       fmt.Println(data)  // May not see "hello"!
   }`)

	fmt.Println("\n✅ Correct: Use atomic.Value for complex types")
	fmt.Println("✅ Correct: Use mutex when ordering matters")
	fmt.Println("✅ Correct: Always use atomic operations for shared variables")
}

// ============================================================================
// بخش 11: جمع‌بندی و جدول مرجع
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("⚡ ATOMIC OPERATIONS IN GO - Complete Guide")
	fmt.Println(stringsRepeat("=", 80))

	// بخش 1: عملیات پایه
	demonstrateBasicAtomic()

	// بخش 2: انواع مختلف
	demonstrateAtomicTypes()

	// بخش 3: Atomic Bool
	demonstrateAtomicBool()

	// بخش 4: Atomic Counter
	demonstrateAtomicCounter()

	// بخش 5: Atomic Value
	demonstrateAtomicValue()

	// بخش 6: Atomic Flag
	demonstrateAtomicFlag()

	// بخش 7: Atomic Int پیشرفته
	demonstrateAtomicInt()

	// بخش 8: مقایسه عملکرد
	benchmarkAtomicVsMutex()

	// بخش 9: الگوهای پیشرفته
	demonstrateAdvancedPatterns()

	// بخش 10: اشتباهات رایج
	demonstrateAtomicPitfalls()

	// بخش 11: جمع‌بندی
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("📚 ATOMIC OPERATIONS QUICK REFERENCE")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println("\n┌─────────────────────────────────────────────────────────────────┐")
	fmt.Println("│ FUNCTION                    │ PURPOSE                          │")
	fmt.Println("├─────────────────────────────┼──────────────────────────────────┤")
	fmt.Println("│ AddInt32/64/Uint32/64       │ Atomic addition                  │")
	fmt.Println("│ LoadInt32/64/Uint32/64      │ Atomic read                      │")
	fmt.Println("│ StoreInt32/64/Uint32/64     │ Atomic write                     │")
	fmt.Println("│ SwapInt32/64/Uint32/64      │ Atomic exchange                  │")
	fmt.Println("│ CompareAndSwapInt32/64      │ Conditional atomic update        │")
	fmt.Println("│ AddUintptr                  │ Atomic pointer arithmetic        │")
	fmt.Println("│ Load/StorePointer           │ Atomic pointer operations        │")
	fmt.Println("│ CompareAndSwapPointer       │ Atomic pointer CAS               │")
	fmt.Println("│ Value                       │ Atomic storage for any type      │")
	fmt.Println("└─────────────────────────────┴──────────────────────────────────┘")

	fmt.Println("\n💡 GOLDEN RULES:")
	fmt.Println("  1. Atomic for simple counters, flags, and integers")
	fmt.Println("  2. Mutex for complex data structures and multi-variable updates")
	fmt.Println("  3. Always use atomic.Load to read atomic variables")
	fmt.Println("  4. Never copy atomic.Value")
	fmt.Println("  5. atomic.Value requires consistent types")
	fmt.Println("  6. CAS (CompareAndSwap) is the foundation of lock-free algorithms")
	fmt.Println("  7. Atomic operations are faster but limited")
	fmt.Println("  8. Use -race flag to detect atomic misuse")
	fmt.Println("  9. For memory ordering, use mutex or channels")
	fmt.Println("  10. When in doubt, use mutex (simpler and safer)")

	fmt.Println("\n🎯 WHEN TO USE ATOMIC vs MUTEX:")
	fmt.Println("  • Simple counter → atomic")
	fmt.Println("  • Flag/status → atomic")
	fmt.Println("  • ID generator → atomic")
	fmt.Println("  • Struct with multiple fields → mutex")
	fmt.Println("  • Map operations → mutex or sync.Map")
	fmt.Println("  • Complex transactions → mutex")
	fmt.Println("  • High-performance scenarios → atomic (after benchmarking)")
}

// ============================================================================
// بخش 12: توابع کمکی
// ============================================================================

func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// ============================================================================
// FILE: atomic_operations_guide_test.go (فایل تست)
// ============================================================================

/*
package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

// Benchmark Atomic Increment
func BenchmarkAtomicIncrement(b *testing.B) {
	var counter int64
	for i := 0; i < b.N; i++ {
		atomic.AddInt64(&counter, 1)
	}
}

// Benchmark Mutex Increment
func BenchmarkMutexIncrement(b *testing.B) {
	var counter int64
	var mu sync.Mutex
	for i := 0; i < b.N; i++ {
		mu.Lock()
		counter++
		mu.Unlock()
	}
}

// Benchmark Atomic Load/Store
func BenchmarkAtomicLoadStore(b *testing.B) {
	var value int64
	atomic.StoreInt64(&value, 42)
	for i := 0; i < b.N; i++ {
		_ = atomic.LoadInt64(&value)
		atomic.StoreInt64(&value, int64(i))
	}
}

// Benchmark Atomic CAS
func BenchmarkAtomicCAS(b *testing.B) {
	var value int64
	for i := 0; i < b.N; i++ {
		atomic.CompareAndSwapInt64(&value, int64(i), int64(i+1))
	}
}

// اجرای benchmark:
// $ go test -bench=. -benchmem
*/

/*
# اجرای معمولی
go run atomic_operations_guide.go

# اجرا با race detector
go run -race atomic_operations_guide.go

# اجرای benchmarkها
go test -bench=. -benchmem

# اجرای تست‌ها با race detector
go test -race -v
*/
