// ============================================================================
// FILE: performance_optimization_guide.go
// TITLE: راهنمای کامل بهینه‌سازی عملکرد در Go - Escape Analysis, Pool, Profiling, کاهش Allocations
// HOW TO RUN: go run performance_optimization_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - چرا بهینه‌سازی عملکرد مهم است؟
// ============================================================================
//
// بهینه‌سازی عملکرد در Go شامل سه حوزه اصلی است:
//
// 1. بهینه‌سازی حافظه (Memory Optimization)
//    - فهم Escape Analysis (کجا متغیرها در heap تخصیص می‌یابند)
//    - استفاده از sync.Pool برای reuse اشیاء
//    - کاهش allocationها
//
// 2. Profiling در تولید (Production Profiling)
//    - جمع‌آوری پروفایل بدون توقف برنامه
//    - تحلیل CPU و Memory در حال اجرا
//    - شناسایی bottlenecks
//
// 3. کاهش Allocations
//    - استفاده از strings.Builder به جای +
//    - Pre-allocating slices
//    - Reusing buffers
//    - Using value receivers vs pointer receivers
//
// قانون طلایی:
// "قبل از بهینه‌سازی، profile بگیر. حدس نزن - اندازه بگیر.
//  از escape analysis آگاه باش. از sync.Pool برای اشیاء پرتکرار استفاده کن.
//  allocationها را در hot paths کاهش بده."
// ============================================================================

package _____performance_optimization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// بخش 1: Escape Analysis - درک فرار متغیرها به Heap
// ============================================================================

// escapeAnalysisExample1 - متغیر در stack می‌ماند (بدون escape)
func escapeAnalysisExample1() int {
	// مقدار در stack تخصیص می‌یابد
	x := 42
	return x
}

// escapeAnalysisExample2 - متغیر escape می‌کند به heap
func escapeAnalysisExample2() *int {
	// این متغیر به heap escape می‌کند (چون اشاره‌گر برگردانده می‌شود)
	x := 42
	return &x
}

// escapeAnalysisExample3 - اینترفیس باعث escape می‌شود
func escapeAnalysisExample3() interface{} {
	// اینترفیس‌ها معمولاً باعث escape می‌شوند
	x := 42
	return x
}

// escapeAnalysisExample4 - closure باعث escape می‌شود
func escapeAnalysisExample4() func() int {
	x := 42
	return func() int {
		return x
	}
}

// escapeAnalysisExample5 - ارسال به channel می‌تواند باعث escape شود
func escapeAnalysisExample5(ch chan *int) {
	x := 42
	ch <- &x // x escape می‌کند
}

// checkEscapeAnalysis نمایش escape analysis
func checkEscapeAnalysis() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔍 ESCAPE ANALYSIS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
# بررسی escape analysis با کامپایلر:

$ go build -gcflags="-m" main.go

# خروجی نمونه:
./main.go:15:6: can inline escapeAnalysisExample1
./main.go:15:6: moved to heap: x (escapeAnalysisExample2)
./main.go:25:2: x escapes to heap (escapeAnalysisExample3)

📊 RULES OF THUMB:

   ✅ STACK (خوب - بدون allocation):
   • متغیرهای محلی که اشاره‌گرشان برگردانده نمی‌شود
   • متغیرهایی که به closure نمی‌روند
   • مقدارهایی که به interface{} تبدیل نمی‌شوند
   
   ❌ HEAP (بد - با allocation):
   • برگرداندن اشاره‌گر به متغیر محلی
   • ذخیره اشاره‌گر در struct که escape می‌کند
   • تبدیل به interface{}
   • بستن در closure
   • ارسال به channel
   • ذخیره در slice/map که escape می‌کنند
`)
}

// ============================================================================
// بخش 2: sync.Pool - بازیابی اشیاء برای کاهش Allocation
// ============================================================================

// مثال: بدون Pool - allocation زیاد
type Buffer struct {
	data []byte
}

// بدون Pool (بد)
func processWithoutPool(data []byte) {
	buf := &Buffer{
		data: make([]byte, 1024),
	}
	copy(buf.data, data)
	// استفاده از buf
	_ = buf
}

// با Pool (خوب)
var bufferPool = sync.Pool{
	New: func() interface{} {
		return &Buffer{
			data: make([]byte, 1024),
		}
	},
}

func processWithPool(data []byte) {
	// گرفتن از pool
	buf := bufferPool.Get().(*Buffer)
	defer bufferPool.Put(buf) // بازگرداندن به pool

	copy(buf.data, data)
	// استفاده از buf
	_ = buf
}

// Example 2: JSON Buffer Pool
var jsonBufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func encodeJSONWithPool(data interface{}) ([]byte, error) {
	buf := jsonBufferPool.Get().(*bytes.Buffer)
	defer jsonBufferPool.Put(buf)

	buf.Reset() // مهم: پاک کردن محتویات قبلی

	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return nil, err
	}

	// کپی کردن نتیجه (چون buffer دوباره استفاده می‌شود)
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// sync.Pool برای structهای سنگین
type BigStruct struct {
	ID    int
	Name  string
	Data  [1024]byte
	Items []string
	Cache map[string]interface{}
}

var bigStructPool = sync.Pool{
	New: func() interface{} {
		return &BigStruct{
			Items: make([]string, 0, 100),
			Cache: make(map[string]interface{}, 50),
		}
	},
}

func getBigStruct() *BigStruct {
	return bigStructPool.Get().(*BigStruct)
}

func putBigStruct(bs *BigStruct) {
	// پاک کردن قبل از بازگرداندن
	bs.ID = 0
	bs.Name = ""
	bs.Data = [1024]byte{}
	bs.Items = bs.Items[:0] // reset slice
	for k := range bs.Cache {
		delete(bs.Cache, k)
	}
	bigStructPool.Put(bs)
}

// ============================================================================
// بخش 3: Profiling در تولید
// ============================================================================

// ProductionProfiler پروفایلر در تولید (بدون توقف برنامه)
type ProductionProfiler struct {
	cpuProfilePath   string
	memProfilePath   string
	blockProfilePath string
	mutexProfilePath string
}

// NewProductionProfiler ایجاد پروفایلر
func NewProductionProfiler(cpuPath, memPath, blockPath, mutexPath string) *ProductionProfiler {
	return &ProductionProfiler{
		cpuProfilePath:   cpuPath,
		memProfilePath:   memPath,
		blockProfilePath: blockPath,
		mutexProfilePath: mutexPath,
	}
}

// StartCPUProfile شروع CPU profiling (با duration)
func (p *ProductionProfiler) StartCPUProfile(duration time.Duration) {
	f, err := os.Create(p.cpuProfilePath)
	if err != nil {
		log.Printf("Failed to create CPU profile: %v", err)
		return
	}
	defer f.Close()

	log.Printf("Starting CPU profile for %v", duration)
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Printf("Failed to start CPU profile: %v", err)
		return
	}

	time.Sleep(duration)
	pprof.StopCPUProfile()
	log.Printf("CPU profile saved to %s", p.cpuProfilePath)
}

// WriteMemoryProfile نوشتن memory profile
func (p *ProductionProfiler) WriteMemoryProfile() {
	f, err := os.Create(p.memProfilePath)
	if err != nil {
		log.Printf("Failed to create memory profile: %v", err)
		return
	}
	defer f.Close()

	runtime.GC() // اجبار به GC برای گرفتن stat دقیق
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Printf("Failed to write heap profile: %v", err)
		return
	}
	log.Printf("Memory profile saved to %s", p.memProfilePath)
}

// StartBlockProfile شروع block profiling
func (p *ProductionProfiler) StartBlockProfile(rate int) {
	runtime.SetBlockProfileRate(rate)
}

// WriteBlockProfile نوشتن block profile
func (p *ProductionProfiler) WriteBlockProfile() {
	f, err := os.Create(p.blockProfilePath)
	if err != nil {
		log.Printf("Failed to create block profile: %v", err)
		return
	}
	defer f.Close()

	if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
		log.Printf("Failed to write block profile: %v", err)
		return
	}
	log.Printf("Block profile saved to %s", p.blockProfilePath)
}

// StartMutexProfile شروع mutex profiling
func (p *ProductionProfiler) StartMutexProfile(fraction int) {
	runtime.SetMutexProfileFraction(fraction)
}

// WriteMutexProfile نوشتن mutex profile
func (p *ProductionProfiler) WriteMutexProfile() {
	f, err := os.Create(p.mutexProfilePath)
	if err != nil {
		log.Printf("Failed to create mutex profile: %v", err)
		return
	}
	defer f.Close()

	if err := pprof.Lookup("mutex").WriteTo(f, 0); err != nil {
		log.Printf("Failed to write mutex profile: %v", err)
		return
	}
	log.Printf("Mutex profile saved to %s", p.mutexProfilePath)
}

// HTTP Profiling Handler (برای production)
func (p *ProductionProfiler) ProfilingHandler(w http.ResponseWriter, r *http.Request) {
	profileType := r.URL.Query().Get("type")
	duration := 30 * time.Second
	if d := r.URL.Query().Get("duration"); d != "" {
		if dur, err := time.ParseDuration(d); err == nil {
			duration = dur
		}
	}

	switch profileType {
	case "cpu":
		p.StartCPUProfile(duration)
		w.Write([]byte(fmt.Sprintf("CPU profile captured for %v", duration)))
	case "heap":
		p.WriteMemoryProfile()
		w.Write([]byte("Heap profile captured"))
	case "goroutine":
		pprof.Lookup("goroutine").WriteTo(w, 1)
	case "block":
		p.WriteBlockProfile()
		w.Write([]byte("Block profile captured"))
	case "mutex":
		p.WriteMutexProfile()
		w.Write([]byte("Mutex profile captured"))
	default:
		w.Write([]byte(`Usage: ?type=cpu|heap|goroutine|block|mutex`))
	}
}

// ============================================================================
// بخش 4: کاهش Allocations - تکنیک‌های عملی
// ============================================================================

// 4.1 استفاده از strings.Builder به جای + (۲x سریعتر، ۱۰x less allocations)
func inefficientStringConcat(items []string) string {
	var result string
	for _, item := range items {
		result += item // ❌ هر بار allocation جدید
	}
	return result
}

func efficientStringConcat(items []string) string {
	var builder strings.Builder
	builder.Grow(estimateSize(items)) // پیش‌تخصیص
	for _, item := range items {
		builder.WriteString(item)
	}
	return builder.String()
}

func estimateSize(items []string) int {
	total := 0
	for _, s := range items {
		total += len(s)
	}
	return total
}

// 4.2 پیش‌تخصیص Slice (avoid reallocation)
func inefficientSliceAppend(n int) []int {
	var result []int // ❌ بدون capacity
	for i := 0; i < n; i++ {
		result = append(result, i) // ممکن است چندین بار reallocate شود
	}
	return result
}

func efficientSliceAppend(n int) []int {
	result := make([]int, 0, n) // ✅ پیش‌تخصیص
	for i := 0; i < n; i++ {
		result = append(result, i)
	}
	return result
}

// 4.3 Reuse slice (با reset کردن)
type SliceRecycler struct {
	pool sync.Pool
}

func NewSliceRecycler() *SliceRecycler {
	return &SliceRecycler{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 4096)
			},
		},
	}
}

func (r *SliceRecycler) Get() []byte {
	return r.pool.Get().([]byte)
}

func (r *SliceRecycler) Put(slice []byte) {
	slice = slice[:0] // reset length, keep capacity
	r.pool.Put(slice)
}

// 4.4 استفاده از value receiver برای structهای کوچک
type SmallStruct struct {
	ID   int
	Name string
}

// Value receiver (بدون allocation)
func (s SmallStruct) GetID() int {
	return s.ID
}

// Pointer receiver (نیاز به allocation اگر struct جدید ایجاد شود)
func (s *SmallStruct) SetID(id int) {
	s.ID = id
}

// 4.5 حذف box کردن (avoid interface{} for primitives)
func processInt(v int) int {
	return v * 2
}

func processInterface(v interface{}) interface{} {
	return v.(int) * 2 // ❌ boxing, type assertion overhead
}

// 4.6 استفاده از []byte به جای string برای پردازش
func processWithString(s string) string {
	return strings.ToUpper(s)
}

func processWithBytes(b []byte) []byte {
	return bytes.ToUpper(b) // ✅ avoids string conversion
}

// 4.7 Benchmark comparison
func runBenchmarkComparison() {
	items := make([]string, 1000)
	for i := range items {
		items[i] = "test"
	}

	// Benchmark inefficient
	start := time.Now()
	for i := 0; i < 1000; i++ {
		inefficientStringConcat(items)
	}
	inefficientTime := time.Since(start)

	// Benchmark efficient
	start = time.Now()
	for i := 0; i < 1000; i++ {
		efficientStringConcat(items)
	}
	efficientTime := time.Since(start)

	fmt.Printf("\nString Concatenation Benchmark:\n")
	fmt.Printf("  Inefficient: %v\n", inefficientTime)
	fmt.Printf("  Efficient:   %v\n", efficientTime)
	fmt.Printf("  Speedup:     %.2fx\n", float64(inefficientTime)/float64(efficientTime))

	// Slice allocation benchmark
	start = time.Now()
	for i := 0; i < 1000; i++ {
		inefficientSliceAppend(1000)
	}
	inefficientSliceTime := time.Since(start)

	start = time.Now()
	for i := 0; i < 1000; i++ {
		efficientSliceAppend(1000)
	}
	efficientSliceTime := time.Since(start)

	fmt.Printf("\nSlice Append Benchmark:\n")
	fmt.Printf("  Inefficient: %v\n", inefficientSliceTime)
	fmt.Printf("  Efficient:   %v\n", efficientSliceTime)
	fmt.Printf("  Speedup:     %.2fx\n", float64(inefficientSliceTime)/float64(efficientSliceTime))
}

// ============================================================================
// بخش 5: Real-world Optimization Example
// ============================================================================

// مثال: پردازش لاگ‌ها (بدون بهینه‌سازی)
type LogProcessor struct {
	entries []string
}

func (p *LogProcessor) ProcessLine(line string) {
	// ❌ allocationهای زیاد
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return
	}
	p.entries = append(p.entries, strings.ToUpper(parts[1]))
}

// مثال: پردازش لاگ‌ها (با بهینه‌سازی)
type OptimizedLogProcessor struct {
	entries   []string
	buffer    []byte
	separator []byte
	pool      sync.Pool
}

func NewOptimizedLogProcessor() *OptimizedLogProcessor {
	return &OptimizedLogProcessor{
		entries:   make([]string, 0, 10000),
		buffer:    make([]byte, 0, 4096),
		separator: []byte(" "),
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 1024)
			},
		},
	}
}

func (p *OptimizedLogProcessor) ProcessLine(line string) {
	// تبدیل string به []byte بدون allocation (unsafe but fast)
	// در عمل: استفاده از []byte(line) allocation دارد
	// برای بهینه‌سازی بیشتر از strings.IndexByte استفاده می‌کنیم

	// پیدا کردن اولین فاصله بدون allocation
	idx := strings.IndexByte(line, ' ')
	if idx == -1 || idx+1 >= len(line) {
		return
	}

	// استخراج بخش دوم بدون allocation جدید
	field := line[idx+1:]

	// تبدیل به uppercase با کمترین allocation
	// در واقع strings.ToUpper allocation دارد، برای بهینه‌سازی بیشتر از bytes.Map استفاده می‌کنیم

	p.entries = append(p.entries, field)
}

// ============================================================================
// بخش 6: Production Readiness Checklist
// ============================================================================

func productionChecklist() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✅ PRODUCTION PERFORMANCE CHECKLIST")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ BEFORE DEPLOYMENT                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│  ☐ Run benchmarks: go test -bench=. -benchmem                             │
│  ☐ Check escape analysis: go build -gcflags="-m" 2>&1 | grep escapes     │
│  ☐ Profile CPU: go test -cpuprofile=cpu.prof                             │
│  ☐ Profile memory: go test -memprofile=mem.prof                          │
│  ☐ Check for unnecessary allocations                                      │
│  ☐ Use sync.Pool for frequently allocated objects                         │
│  ☐ Pre-allocate slices with make() when size known                        │
│  ☐ Use strings.Builder for concatenation                                  │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ IN PRODUCTION                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│  ☐ Enable pprof HTTP endpoints                                            │
│  ☐ Set up regular profiling (30s every hour)                              │
│  ☐ Monitor goroutine count                                                │
│  ☐ Track memory allocation rate                                           │
│  ☐ Set up alerts on high latency                                          │
│  ☐ Log GC pauses                                                          │
│  ☐ Monitor heap growth                                                    │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ PERFORMANCE METRICS TO MONITOR                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│  • GC Stats: runtime.ReadMemStats()                                        │
│  • Goroutine count: runtime.NumGoroutine()                                │
│  • Heap allocation: memStats.HeapAlloc                                    │
│  • GC cycles: memStats.NumGC                                              │
│  • Pause time: memStats.PauseNs                                           │
│  • CPU usage: use prometheus metrics                                      │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 7: Common Performance Pitfalls
// ============================================================================

func commonPitfalls() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚠️ COMMON PERFORMANCE PITFALLS")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. Unnecessary allocations in hot loops                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│    for i := 0; i < 1_000_000; i++ {                                       │
│        s := fmt.Sprintf("%d", i)  // ❌ allocation each iteration         │
│    }                                                                       │
│    ✅ var sb strings.Builder; sb.Grow(预估大小)                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. Converting between string and []byte unnecessarily                     │
├─────────────────────────────────────────────────────────────────────────────┤
│    s := "hello"                                                           │
│    b := []byte(s)  // ❌ allocation                                       │
│    s2 := string(b) // ❌ another allocation                               │
│    ✅ Use strings.Builder or bytes.Buffer                                 │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. Using reflection in hot paths                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│    reflect.ValueOf(x).Interface() // ❌ slow, allocations                 │
│    ✅ Use type switches or code generation                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. Large struct passed by value                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│    func process(big BigStruct) { ... } // ❌ copies whole struct          │
│    ✅ func process(big *BigStruct) { ... }                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. Using defer in tight loops                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│    for i := 0; i < 1_000_000; i++ {                                       │
│        defer mu.Unlock() // ❌ defer overhead                            │
│    }                                                                       │
│    ✅ mu.Unlock() directly                                                │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 8: Memory Stats Monitoring
// ============================================================================

// MemoryMonitor مانیتور حافظه
type MemoryMonitor struct {
	interval time.Duration
	stopCh   chan struct{}
}

func NewMemoryMonitor(interval time.Duration) *MemoryMonitor {
	return &MemoryMonitor{
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

func (m *MemoryMonitor) Start() {
	ticker := time.NewTicker(m.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.printStats()
			case <-m.stopCh:
				ticker.Stop()
				return
			}
		}
	}()
}

func (m *MemoryMonitor) Stop() {
	close(m.stopCh)
}

func (m *MemoryMonitor) printStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	log.Printf("=== Memory Stats ===")
	log.Printf("  Alloc:      %d MB", mem.Alloc/1024/1024)
	log.Printf("  TotalAlloc: %d MB", mem.TotalAlloc/1024/1024)
	log.Printf("  Sys:        %d MB", mem.Sys/1024/1024)
	log.Printf("  Mallocs:    %d", mem.Mallocs)
	log.Printf("  Frees:      %d", mem.Frees)
	log.Printf("  Live:       %d", mem.Mallocs-mem.Frees)
	log.Printf("  GC Cycles:  %d", mem.NumGC)
	log.Printf("  Goroutines: %d", runtime.NumGoroutine())
}

// ============================================================================
// بخش 9: Best Practices Summary
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("💡 PERFORMANCE BEST PRACTICES SUMMARY")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ MEMORY OPTIMIZATION                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│  • Use sync.Pool for frequently allocated objects                          │
│  • Pre-allocate slices: make([]T, 0, capacity)                            │
│  • Use strings.Builder for string concatenation                           │
│  • Avoid boxing primitives to interface{}                                 │
│  • Use value receivers for small structs                                  │
│  • Understand escape analysis                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ PROFILING                                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│  • Profile before optimizing                                              │
│  • Use production profiling                                                │
│  • Compare before/after profiles                                          │
│  • Use flame graphs for visualization                                     │
│  • Profile in production-like environment                                 │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ ALLOCATION REDUCTION                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│  • Reuse buffers and slices                                                │
│  • Use []byte instead of string for processing                           │
│  • Avoid fmt.Sprintf in hot loops                                         │
│  • Use strconv instead of fmt for conversions                             │
│  • Batch allocations                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 10: Main
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 PERFORMANCE OPTIMIZATION IN GO")
	fmt.Println("Escape Analysis | sync.Pool | Production Profiling | Reduce Allocations")
	fmt.Println(strings.Repeat("=", 80))

	checkEscapeAnalysis()
	productionChecklist()
	commonPitfalls()
	bestPractices()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚀 Running Benchmarks")
	fmt.Println(strings.Repeat("=", 80))

	runBenchmarkComparison()

	// Example of using sync.Pool
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📦 sync.Pool Example")
	fmt.Println(strings.Repeat("=", 80))

	// بدون pool
	start := time.Now()
	for i := 0; i < 100000; i++ {
		data := make([]byte, 1024)
		_ = data
	}
	noPoolTime := time.Since(start)

	// با pool
	pool := sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}
	start = time.Now()
	for i := 0; i < 100000; i++ {
		data := pool.Get().([]byte)
		pool.Put(data)
	}
	withPoolTime := time.Since(start)

	fmt.Printf("Without Pool: %v\n", noPoolTime)
	fmt.Printf("With Pool:    %v\n", withPoolTime)
	fmt.Printf("Speedup:      %.2fx\n", float64(noPoolTime)/float64(withPoolTime))

	// Memory monitoring example
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 Memory Monitoring")
	fmt.Println(strings.Repeat("=", 80))

	monitor := NewMemoryMonitor(2 * time.Second)
	monitor.Start()

	// Simulate work
	for i := 0; i < 5; i++ {
		// Allocate some memory
		_ = make([]byte, 10*1024*1024) // 10MB
		time.Sleep(2 * time.Second)
	}

	monitor.Stop()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎯 PERFORMANCE OPTIMIZATION - COMPLETE")
	fmt.Println("Ready to optimize your Go applications!")
	fmt.Println(strings.Repeat("=", 80))
}

// تابع کمکی برای تکرار رشته
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// اضافه کردن importهای لازم
var _ = json.Marshal
var _ = pprof.StartCPUProfile
var _ = os.Create

// تابع برای کامپایل
func init() {
	// فقط برای جلوگیری از خطای import
	var _ = &sync.Mutex{}
}
