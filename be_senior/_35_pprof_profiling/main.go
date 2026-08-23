// ============================================================================
// FILE: pprof_profiling_guide.go
// TITLE: راهنمای کامل Profiling با pprof در Go - بهینه‌سازی عملکرد
// HOW TO RUN: go run pprof_profiling_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - Profiling چیست و چرا نیاز است؟
// ============================================================================
//
// Profiling فرآیند اندازه‌گیری عملکرد برنامه برای پیدا کردن bottleneckهاست.
// pprof ابزار استاندارد Go برای profiling است و اطلاعات زیر را ارائه می‌دهد:
//
// 1. CPU Profile   - زمان صرف شده در توابع مختلف
// 2. Memory Profile - مصرف حافظه (تخصیص‌ها)
// 3. Block Profile - زمان صرف شده در عملیات blocking (mutex, channels)
// 4. Mutex Profile - رقابت روی mutexها
// 5. Goroutine Profile - تعداد و وضعیت گوروتین‌ها
// 6. ThreadCreate Profile - تعداد threadهای ایجاد شده
// 7. Heap Profile   - تخصیص‌های heap
// 8. Alloc Profile  - کل تخصیص‌های حافظه (تاریخی)
//
// قانون طلایی:
// "قبل از بهینه‌سازی، همیشه profile بگیر. حدس نزن - اندازه بگیر.
//  از pprof برای پیدا کردن bottleneck استفاده کن، سپس بهینه‌سازی کن."
// ============================================================================

package __pprof_profiling

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// بخش 1: انواع پروفایل‌ها (Profile Types)
// ============================================================================

func demonstrateProfileTypes() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 SECTION 1: PROFILE TYPES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ PROFILE TYPES IN Go                                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  Type               | Description              | Command       │
│  -------------------|--------------------------|---------------│
│  cpu                | CPU usage by function    | go test -cpuprofile│
│  mem                | Memory allocations       | go test -memprofile│
│  heap               | Live heap objects        | pprof.Heap      │
│  allocs             | All allocations (history)| pprof.Allocs    │
│  block              | Blocking operations      | go test -blockprofile│
│  mutex              | Mutex contention         | go test -mutexprofile│
│  goroutine          | Goroutine stack traces   | pprof.Goroutine │
│  threadcreate       | Thread creation          | pprof.ThreadCreate│
│  trace              | Execution tracing        | go test -trace  │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 2: برنامه نمونه برای Profiling
// ============================================================================

// fibonacci - تابع محاسباتی سنگین (برای CPU profile)
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// heavyAllocation - تابع با تخصیص حافظه زیاد (برای Memory profile)
func heavyAllocation(size int) []byte {
	// تخصیص slice بزرگ
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	return data
}

// inefficientConcat - الحاق ناکارآمد رشته (برای CPU/Memory)
func inefficientConcat(count int) string {
	var result string
	for i := 0; i < count; i++ {
		result += "x" // هر بار رشته جدید تخصیص می‌دهد
	}
	return result
}

// efficientConcat - الحاق کارآمد
func efficientConcat(count int) string {
	var builder strings.Builder
	builder.Grow(count)
	for i := 0; i < count; i++ {
		builder.WriteString("x")
	}
	return builder.String()
}

// blockingOperation - عملیات blocking (برای Block profile)
func blockingOperation(d time.Duration) {
	time.Sleep(d)
}

// mutexContention - رقابت روی mutex (برای Mutex profile)
func mutexContention() {
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			// کار کوتاه
			time.Sleep(1 * time.Millisecond)
			mu.Unlock()
		}()
	}
	wg.Wait()
}

// ============================================================================
// بخش 3: CPU Profiling
// ============================================================================

func demonstrateCPUProfiling() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ SECTION 3: CPU PROFILING")
	fmt.Println(strings.Repeat("=", 80))

	// ایجاد فایل برای CPU profile
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()

	// شروع CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// اجرای کدی که می‌خواهیم profile بگیریم
	fmt.Println("  Running CPU-intensive task...")
	for i := 0; i < 30; i++ {
		fibonacci(30)
	}

	fmt.Println("  CPU profile saved to cpu.prof")
	fmt.Println("  Analyze with: go tool pprof cpu.prof")
}

// ============================================================================
// بخش 4: Memory Profiling
// ============================================================================

func demonstrateMemoryProfiling() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💾 SECTION 4: MEMORY PROFILING")
	fmt.Println(strings.Repeat("=", 80))

	// اجرای کد با تخصیص حافظه زیاد
	fmt.Println("  Running memory-intensive task...")
	for i := 0; i < 100; i++ {
		_ = heavyAllocation(1024 * 1024) // 1MB each
	}

	// گرفتن heap profile
	f, err := os.Create("heap.prof")
	if err != nil {
		log.Fatal("could not create heap profile: ", err)
	}
	defer f.Close()

	runtime.GC() // اجبار به garbage collection
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write heap profile: ", err)
	}

	fmt.Println("  Heap profile saved to heap.prof")
	fmt.Println("  Analyze with: go tool pprof -http=:8080 heap.prof")
}

// ============================================================================
// بخش 5: Goroutine Profiling
// ============================================================================

func demonstrateGoroutineProfiling() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 SECTION 5: GOROUTINE PROFILING")
	fmt.Println(strings.Repeat("=", 80))

	// ایجاد چند گوروتین
	for i := 0; i < 10; i++ {
		go func(id int) {
			time.Sleep(5 * time.Second)
		}(i)
	}

	// گرفتن goroutine profile
	f, err := os.Create("goroutine.prof")
	if err != nil {
		log.Fatal("could not create goroutine profile: ", err)
	}
	defer f.Close()

	if err := pprof.Lookup("goroutine").WriteTo(f, 0); err != nil {
		log.Fatal("could not write goroutine profile: ", err)
	}

	fmt.Println("  Goroutine profile saved to goroutine.prof")
	fmt.Printf("  Current goroutines: %d\n", runtime.NumGoroutine())
}

// ============================================================================
// بخش 6: Block Profiling
// ============================================================================

func demonstrateBlockProfiling() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⏸️ SECTION 6: BLOCK PROFILING")
	fmt.Println(strings.Repeat("=", 80))

	// تنظیم نرخ نمونه‌برداری block profile
	runtime.SetBlockProfileRate(1)

	// اجرای عملیات blocking
	fmt.Println("  Running blocking operations...")
	for i := 0; i < 10; i++ {
		blockingOperation(10 * time.Millisecond)
	}

	// گرفتن block profile
	f, err := os.Create("block.prof")
	if err != nil {
		log.Fatal("could not create block profile: ", err)
	}
	defer f.Close()

	if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
		log.Fatal("could not write block profile: ", err)
	}

	fmt.Println("  Block profile saved to block.prof")
}

// ============================================================================
// بخش 7: Mutex Profiling
// ============================================================================

func demonstrateMutexProfiling() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔒 SECTION 7: MUTEX PROFILING")
	fmt.Println(strings.Repeat("=", 80))

	// تنظیم نرخ نمونه‌برداری mutex profile
	runtime.SetMutexProfileFraction(1)

	// اجرای عملیات با mutex contention
	fmt.Println("  Running mutex contention...")
	mutexContention()

	// گرفتن mutex profile
	f, err := os.Create("mutex.prof")
	if err != nil {
		log.Fatal("could not create mutex profile: ", err)
	}
	defer f.Close()

	if err := pprof.Lookup("mutex").WriteTo(f, 0); err != nil {
		log.Fatal("could not write mutex profile: ", err)
	}

	fmt.Println("  Mutex profile saved to mutex.prof")
}

// ============================================================================
// بخش 8: HTTP Profiling (معروف‌ترین روش)
// ============================================================================

// برای استفاده در برنامه‌های HTTP
func startHTTPProfiling() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🌐 SECTION 8: HTTP PROFILING")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ HTTP PROFILING SETUP                                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  import _ "net/http/pprof"                                    │
│                                                                │
│  func main() {                                                │
│      go func() {                                              │
│          log.Println(http.ListenAndServe("localhost:6060", nil))│
│      }()                                                      │
│      // your application code here                            │
│  }                                                            │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

📊 AVAILABLE PROFILES AT http://localhost:6060/debug/pprof/

   • /debug/pprof/           - Index page
   • /debug/pprof/cmdline    - Command line
   • /debug/pprof/profile    - CPU profile (30s)
   • /debug/pprof/heap       - Heap profile
   • /debug/pprof/goroutine  - Goroutine profile
   • /debug/pprof/block      - Block profile
   • /debug/pprof/mutex      - Mutex profile
   • /debug/pprof/threadcreate - Thread creation
   • /debug/pprof/trace      - Execution trace

📝 COMMANDS:

   # CPU profile for 30 seconds
   $ go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

   # Heap profile
   $ go tool pprof http://localhost:6060/debug/pprof/heap

   # Goroutine profile
   $ go tool pprof http://localhost:6060/debug/pprof/goroutine

   # Interactive web UI
   $ go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile
`)
}

// ============================================================================
// بخش 9: تحلیل پروفایل با go tool pprof
// ============================================================================

func demonstratePprofAnalysis() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📈 SECTION 9: pprof ANALYSIS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ BASIC pprof COMMANDS                                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  # Interactive mode                                            │
│  $ go tool pprof cpu.prof                                     │
│                                                                │
│  # Top functions (most CPU time)              │
│  (pprof) top                                                  │
│  (pprof) top10                                                │
│  (pprof) top -cum                                             │
│                                                                │
│  # List function source code                                   │
│  (pprof) list fibonacci                                        │
│                                                                │
│  # Call graph (text)                                          │
│  (pprof) web                                                  │
│  (pprof) pdf                                                  │
│  (pprof) gif                                                  │
│                                                                │
│  # Web UI (recommended)                     │
│  $ go tool pprof -http=:8080 cpu.prof                         │
│                                                                │
│  # Compare two profiles                                       │
│  $ go tool pprof -base baseline.prof current.prof             │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ pprof COMMANDS (Interactive)                                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  top [n] [--cum]              - Show top n functions          │
│  list <func>                  - Show source code              │
│  weblist <func>               - Show source with disassembly  │
│  callgrind                    - Generate callgrind format     │
│  peek <func>                  - Show callers/callees          │
│  traces                       - Show all stack traces         │
│  disasm <addr>                - Disassemble                   │
│  tags                         - List tags                     │
│  help                         - Show help                     │
│  quit                         - Exit                          │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ WEB UI FEATURES                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  • Top - List of hottest functions                             │
│  • Graph - Interactive call graph                              │
│  • Flame Graph - Flame chart visualization                     │
│  • Source - Annotated source code                              │
│  • Disassembly - Assembly with annotations                     │
│  • Peek - Callers/callees                                      │
│                                                                │
│  View options:                                                 │
│  • Flat - Total time per function                              │
│  • Cum - Cumulative time (including called functions)          │
│  • diff - Compare two profiles                                 │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 10: Benchmark Profiling
// ============================================================================

func demonstrateBenchmarkProfiling() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 SECTION 10: BENCHMARK PROFILING")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ PROFILING WITH BENCHMARKS                                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  # CPU profiling with benchmark                                │
│  $ go test -bench=. -cpuprofile=cpu.prof                      │
│                                                                │
│  # Memory profiling with benchmark                             │
│  $ go test -bench=. -memprofile=mem.prof                      │
│                                                                │
│  # Block profiling                                            │
│  $ go test -bench=. -blockprofile=block.prof                  │
│                                                                │
│  # Mutex profiling                                            │
│  $ go test -bench=. -mutexprofile=mutex.prof                  │
│                                                                │
│  # Example benchmark function:                                │
│  func BenchmarkInefficientConcat(b *testing.B) {              │
│      for i := 0; i < b.N; i++ {                               │
│          inefficientConcat(100)                               │
│      }                                                        │
│  }                                                            │
│                                                                │
│  func BenchmarkEfficientConcat(b *testing.B) {                │
│      for i := 0; i < b.N; i++ {                               │
│          efficientConcat(100)                                 │
│      }                                                        │
│  }                                                            │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 11: مثال کامل Profiling Workflow
// ============================================================================

func demonstrateFullWorkflow() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 SECTION 11: COMPLETE PROFILING WORKFLOW")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ STEP-BY-STEP PROFILING WORKFLOW                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  STEP 1: Build with profiling support                         │
│  $ go build -gcflags="-N -l" -o app .                         │
│                                                                │
│  STEP 2: Add HTTP pprof import in main.go                     │
│  import _ "net/http/pprof"                                    │
│                                                                │
│  STEP 3: Start HTTP server for profiling                      │
│  go func() {                                                  │
│      log.Println(http.ListenAndServe(":6060", nil))           │
│  }()                                                          │
│                                                                │
│  STEP 4: Run your application                                 │
│  $ ./app                                                      │
│                                                                │
│  STEP 5: Generate traffic (for web apps)                      │
│  $ hey -n 10000 -c 100 http://localhost:8080/                │
│                                                                │
│  STEP 6: Capture CPU profile                                  │
│  $ go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30│
│                                                                │
│  STEP 7: Analyze profile                                      │
│  $ go tool pprof -http=:8080 cpu.prof                         │
│                                                                │
│  STEP 8: Identify bottlenecks                                 │
│  - Look for functions with high flat/cum values               │
│  - Check call graph for hot paths                             │
│  - Examine flame graph                                        │
│                                                                │
│  STEP 9: Optimize code                                        │
│  - Fix identified bottlenecks                                 │
│  - Reduce allocations                                         │
│  - Improve algorithms                                         │
│                                                                │
│  STEP 10: Re-profile and verify improvement                   │
│  - Compare profiles: go tool pprof -base old.prof new.prof    │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 12: Optimization Techniques
// ============================================================================

func demonstrateOptimizations() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ SECTION 12: OPTIMIZATION TECHNIQUES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ COMMON OPTIMIZATIONS BASED ON PROFILES                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  1. CPU Bound Problems:                                       │
│     • Optimize algorithms (O(n²) → O(n log n))                │
│     • Use memoization/caching                                 │
│     • Parallelize with goroutines                             │
│     • Precompute values                                       │
│                                                                │
│  2. Memory Bound Problems:                                    │
│     • Reduce allocations (use sync.Pool)                      │
│     • Pre-allocate slices: make([]T, 0, capacity)             │
│     • Use strings.Builder instead of +                       │
│     • Reuse buffers                                           │
│                                                                │
│  3. Blocking Problems:                                        │
│     • Reduce lock contention                                  │
│     • Use RWMutex for read-heavy workloads                    │
│     • Use channels for communication                          │
│     • Implement backpressure                                  │
│                                                                │
│  4. GC Pressure:                                             │
│     • Reduce pointer usage                                    │
│     • Use value types instead of pointers                     │
│     • Reuse objects with sync.Pool                            │
│     • Tune GOGC (environment variable)                        │
│                                                                │
│  5. I/O Bound Problems:                                      │
│     • Use buffered I/O                                        │
│     • Parallel I/O operations                                 │
│     • Implement caching                                       │
│     • Use async writes                                        │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ BEFORE vs AFTER EXAMPLE                                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  BEFORE (inefficient):                                        │
│  func buildString(items []string) string {                   │
│      var result string                                        │
│      for _, item := range items {                            │
│          result += item + ","  // Many allocations!          │
│      }                                                        │
│      return result                                            │
│  }                                                            │
│                                                                │
│  AFTER (optimized):                                          │
│  func buildString(items []string) string {                   │
│      var builder strings.Builder                             │
│      builder.Grow(estimateSize(items))  // Pre-allocate      │
│      for i, item := range items {                            │
│          if i > 0 {                                           │
│              builder.WriteString(",")                         │
│          }                                                    │
│          builder.WriteString(item)                            │
│      }                                                        │
│      return builder.String()                                  │
│  }                                                            │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 13: اشتباهات رایج
// ============================================================================

func demonstrateCommonMistakes() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("❌ SECTION 13: COMMON MISTAKES")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ COMMON PROFILING MISTAKES                                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  ❌ Mistake 1: Profiling optimized binaries                    │
│     Inlined functions and optimizations hide information      │
│     ✅ Use -gcflags="-N -l" when building for profiling       │
│                                                                │
│  ❌ Mistake 2: Profiling for too short                        │
│     1 second may not be enough for accurate results           │
│     ✅ Profile for at least 30 seconds                        │
│                                                                │
│  ❌ Mistake 3: Ignoring cumulative time (cum)                 │
│     Flat time may be low but cum time high (callers)          │
│     ✅ Use top -cum to see cumulative time                    │
│                                                                │
│  ❌ Mistake 4: Profiling on different hardware                │
│     Results vary between machines                             │
│     ✅ Profile on production-like hardware                    │
│                                                                │
│  ❌ Mistake 5: Optimizing without measuring                   │
│     "Premature optimization is the root of all evil"          │
│     ✅ Always profile before and after optimization           │
│                                                                │
│  ❌ Mistake 6: Profiling only CPU, not memory                 │
│     Memory leaks and high allocations cause problems          │
│     ✅ Always check both CPU and memory profiles              │
│                                                                │
│  ❌ Mistake 7: Not using web UI                               │
│     Text output is harder to interpret                        │
│     ✅ Use go tool pprof -http=:8080 profile.prof             │
│                                                                │
│  ❌ Mistake 8: Profiling in isolation                         │
│     Real-world usage patterns differ from tests               │
│     ✅ Profile under realistic load                           │
│                                                                │
└─────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 14: جمع‌بندی و جدول مرجع
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLETE pprof PROFILING GUIDE")
	fmt.Println("CPU | Memory | Goroutine | Block | Mutex")
	fmt.Println(strings.Repeat("=", 80))

	// نمایش انواع پروفایل
	demonstrateProfileTypes()

	// نمایش CPU profiling
	demonstrateCPUProfiling()

	// نمایش Memory profiling
	demonstrateMemoryProfiling()

	// نمایش Goroutine profiling
	demonstrateGoroutineProfiling()

	// نمایش Block profiling
	demonstrateBlockProfiling()

	// نمایش Mutex profiling
	demonstrateMutexProfiling()

	// نمایش HTTP profiling
	startHTTPProfiling()

	// نمایش تحلیل pprof
	demonstratePprofAnalysis()

	// نمایش benchmark profiling
	demonstrateBenchmarkProfiling()

	// نمایش workflow کامل
	demonstrateFullWorkflow()

	// نمایش تکنیک‌های بهینه‌سازی
	demonstrateOptimizations()

	// نمایش اشتباهات رایج
	demonstrateCommonMistakes()

	// جمع‌بندی نهایی
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 QUICK REFERENCE")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────┐
│ COMMAND CHEAT SHEET                                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                │
│  # CPU Profile                                                │
│  $ go test -bench=. -cpuprofile=cpu.prof                      │
│  $ go tool pprof -http=:8080 cpu.prof                         │
│                                                                │
│  # Memory Profile                                             │
│  $ go test -bench=. -memprofile=mem.prof                      │
│  $ go tool pprof -http=:8080 mem.prof                         │
│                                                                │
│  # HTTP Endpoints                                             │
│  $ go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30│
│  $ go tool pprof http://localhost:6060/debug/pprof/heap       │
│  $ go tool pprof http://localhost:6060/debug/pprof/goroutine  │
│                                                                │
│  # Compare Profiles                                           │
│  $ go tool pprof -base baseline.prof current.prof             │
│                                                                │
│  # Build with debug info                                      │
│  $ go build -gcflags="-N -l" -o app .                         │
│                                                                │
└─────────────────────────────────────────────────────────────────┘

💡 GOLDEN RULES:

   1. Always profile before optimizing - never guess!
   2. Use -gcflags="-N -l" to disable optimizations for profiling
   3. Profile for at least 30 seconds for accurate results
   4. Use web UI: go tool pprof -http=:8080 profile.prof
   5. Compare profiles before and after optimization
   6. Profile on production-like hardware and load
   7. Check both CPU and memory profiles
   8. Look at cumulative time (cum) not just flat time
   9. Use flame graphs for easier visualization
   10. Automate profiling in CI/CD pipelines

🎯 WHEN TO PROFILE:

   • Before major optimizations
   • When investigating performance issues
   • After adding new features
   • Before and after dependency updates
   • When preparing for scaling
   • During load testing
   • When investigating memory leaks
   • For root cause analysis of production issues

📁 PROFILE FILES GENERATED:

   cpu.prof      - CPU profile
   mem.prof      - Memory profile
   heap.prof     - Heap profile
   block.prof    - Block profile
   mutex.prof    - Mutex profile
   goroutine.prof - Goroutine profile
   trace.out     - Execution trace
`)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎯 pprof PROFILING - COMPLETE GUIDE")
	fmt.Println("Ready to optimize your Go programs!")
	fmt.Println(strings.Repeat("=", 80))
}

// اضافه کردن sync import برای mutexContention
var _ = sync.Mutex{}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

/*
خلاصه دستورات مهم pprof
نوع پروفایل	دستور گرفتن پروفایل	دستور تحلیل
CPU	go test -cpuprofile=cpu.prof	go tool pprof -http=:8080 cpu.prof
Memory	go test -memprofile=mem.prof	go tool pprof -http=:8080 mem.prof
Block	go test -blockprofile=block.prof	go tool pprof -http=:8080 block.prof
Mutex	go test -mutexprofile=mutex.prof	go tool pprof -http=:8080 mutex.prof
HTTP CPU	curl localhost:6060/debug/pprof/profile?seconds=30 > cpu.prof	go tool pprof -http=:8080 cpu.prof
HTTP Heap	curl localhost:6060/debug/pprof/heap > heap.prof	go tool pprof -http=:8080 heap.prof

*/
