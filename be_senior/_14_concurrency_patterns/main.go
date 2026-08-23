//worker pool, fan-in, fan-out, pipeline

// ============================================================================
// FILE: concurrency_patterns_guide.go
// TITLE: الگوهای رایج همزمانی در Go - Worker Pool, Fan-In, Fan-Out, Pipeline
// HOW TO RUN: go run concurrency_patterns_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - چرا به این الگوها نیاز داریم؟
// ============================================================================
//
// چهار الگوی اصلی همزمانی در Go:
//
// 1. Worker Pool:
//    - گروهی از workerها که کارها را از یک صف مشترک انجام می‌دهند
//    - محدود کردن تعداد گوروتین‌های همزمان
//    - مناسب برای پردازش دسته‌ای کارهای مستقل
//
// 2. Fan-Out:
//    - توزیع کار بین چندین گوروتین (پخش کردن)
//    - یک ورودی، چندین خروجی
//    - مناسب برای موازی‌سازی کارهای مستقل
//
// 3. Fan-In:
//    - جمع‌آوری نتایج از چندین گوروتین در یک کانال (ادغام)
//    - چندین ورودی، یک خروجی
//    - مناسب برای جمع‌آوری نتایج از منابع مختلف
//
// 4. Pipeline:
//    - زنجیره‌ای از stageها که هر stage داده را پردازش و به stage بعد می‌دهد
//    - هر stage می‌تواند همزمان اجرا شود
//    - مناسب برای پردازش جریان داده (stream processing)
// ============================================================================

package _14_concurrency_patterns

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// بخش 1: الگوی Worker Pool (استخر کارگر)
// ============================================================================

// Job یک واحد کار
type Job struct {
	ID     int
	Number int
}

// Result نتیجه پردازش یک Job
type Result struct {
	JobID   int
	Value   int
	Success bool
}

// WorkerPool ساختار اصلی Worker Pool
type WorkerPool struct {
	jobQueue    chan Job
	resultQueue chan Result
	workerCount int
	wg          sync.WaitGroup
	stopCh      chan struct{}
	stopped     atomic.Bool
}

// NewWorkerPool ایجاد Worker Pool جدید
func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
	return &WorkerPool{
		jobQueue:    make(chan Job, queueSize),
		resultQueue: make(chan Result, queueSize),
		workerCount: workerCount,
		stopCh:      make(chan struct{}),
	}
}

// worker تابعی که هر worker اجرا می‌کند
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case job, ok := <-wp.jobQueue:
			if !ok {
				return // jobQueue بسته شده
			}
			// پردازش job
			result := wp.processJob(job)
			wp.resultQueue <- result

		case <-wp.stopCh:
			return
		}
	}
}

// processJob پردازش یک job (شبیه‌سازی کار سنگین)
func (wp *WorkerPool) processJob(job Job) Result {
	// شبیه‌سازی کار محاسباتی
	time.Sleep(10 * time.Millisecond)

	return Result{
		JobID:   job.ID,
		Value:   job.Number * job.Number,
		Success: true,
	}
}

// Start راه‌اندازی Worker Pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Submit ارسال Job به Pool
func (wp *WorkerPool) Submit(job Job) bool {
	if wp.stopped.Load() {
		return false
	}

	select {
	case wp.jobQueue <- job:
		return true
	default:
		return false // Queue پر است
	}
}

// SubmitWithTimeout ارسال Job با تایم‌اوت
func (wp *WorkerPool) SubmitWithTimeout(job Job, timeout time.Duration) error {
	if wp.stopped.Load() {
		return fmt.Errorf("pool is stopped")
	}

	select {
	case wp.jobQueue <- job:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("submit timeout")
	}
}

// Results گرفتن کانال نتایج
func (wp *WorkerPool) Results() <-chan Result {
	return wp.resultQueue
}

// Stop توقف Pool و منتظر ماندن برای اتمام
func (wp *WorkerPool) Stop() {
	wp.stopped.Store(true)
	close(wp.stopCh)
	close(wp.jobQueue)
	wp.wg.Wait()
	close(wp.resultQueue)
}

// demonstrateWorkerPool نمایش Worker Pool در عمل
func demonstrateWorkerPool() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🧠 WORKER POOL PATTERN")
	fmt.Println(strings.Repeat("=", 80))

	// ایجاد Pool با 5 worker
	pool := NewWorkerPool(5, 10)
	pool.Start()

	// ارسال 20 Job
	go func() {
		for i := 1; i <= 20; i++ {
			job := Job{ID: i, Number: i * 2}
			pool.Submit(job)
			fmt.Printf("📤 Submitted job %d\n", i)
		}
	}()

	// جمع‌آوری نتایج
	results := make([]Result, 0)
	for result := range pool.Results() {
		results = append(results, result)
		fmt.Printf("📥 Result: Job %d -> %d\n", result.JobID, result.Value)
	}

	fmt.Printf("\n✅ Worker Pool completed: %d jobs processed by %d workers\n",
		len(results), 5)
}

// ============================================================================
// بخش 2: الگوی Fan-Out (پخش کردن کار)
// ============================================================================

// FanOut توزیع کار بین چندین گوروتین
type FanOut struct {
	inputChan <-chan int
	workers   int
}

// NewFanOut ایجاد FanOut جدید
func NewFanOut(input <-chan int, workers int) *FanOut {
	return &FanOut{
		inputChan: input,
		workers:   workers,
	}
}

// worker در FanOut
func (fo *FanOut) worker(id int, outputChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for value := range fo.inputChan {
		// پردازش موازی
		result := fo.process(value)
		outputChan <- result
	}
}

// process پردازش یک مقدار
func (fo *FanOut) process(value int) int {
	// شبیه‌سازی کار سنگین
	time.Sleep(20 * time.Millisecond)
	return value * value
}

// Start شروع FanOut
func (fo *FanOut) Start() <-chan int {
	outputChan := make(chan int, fo.workers)
	var wg sync.WaitGroup

	// راه‌اندازی workers
	for i := 0; i < fo.workers; i++ {
		wg.Add(1)
		go fo.worker(i, outputChan, &wg)
	}

	// بستن outputChan بعد از اتمام همه workers
	go func() {
		wg.Wait()
		close(outputChan)
	}()

	return outputChan
}

// demonstrateFanOut نمایش Fan-Out
func demonstrateFanOut() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📡 FAN-OUT PATTERN (Distribution)")
	fmt.Println(strings.Repeat("=", 80))

	// ایجاد داده ورودی
	input := make(chan int, 20)
	for i := 1; i <= 20; i++ {
		input <- i
	}
	close(input)

	// ایجاد FanOut با 4 worker
	fanOut := NewFanOut(input, 4)
	start := time.Now()

	// دریافت نتایج
	results := make([]int, 0)
	for result := range fanOut.Start() {
		results = append(results, result)
		fmt.Printf("🎯 Result: %d\n", result)
	}

	elapsed := time.Since(start)
	fmt.Printf("\n✅ Fan-Out completed: %d items processed by 4 workers in %v\n",
		len(results), elapsed)
}

// ============================================================================
// بخش 3: الگوی Fan-In (ادغام چند کانال)
// ============================================================================

// FanIn ادغام چند کانال در یک کانال
type FanIn struct {
	channels []<-chan int
}

// NewFanIn ایجاد FanIn جدید
func NewFanIn(channels ...<-chan int) *FanIn {
	return &FanIn{
		channels: channels,
	}
}

// Merge ادغام همه کانال‌ها
func (fi *FanIn) Merge() <-chan int {
	output := make(chan int)
	var wg sync.WaitGroup

	// برای هر کانال ورودی یک گوروتین
	for _, ch := range fi.channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for value := range c {
				output <- value
			}
		}(ch)
	}

	// بستن output بعد از اتمام همه
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// demonstrateFanIn نمایش Fan-In
func demonstrateFanIn() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔄 FAN-IN PATTERN (Merging)")
	fmt.Println(strings.Repeat("=", 80))

	// ایجاد چند کانال ورودی
	ch1 := make(chan int, 5)
	ch2 := make(chan int, 5)
	ch3 := make(chan int, 5)

	// تولید داده در کانال‌ها
	go func() {
		for i := 1; i <= 5; i++ {
			ch1 <- i
		}
		close(ch1)
	}()

	go func() {
		for i := 6; i <= 10; i++ {
			ch2 <- i
		}
		close(ch2)
	}()

	go func() {
		for i := 11; i <= 15; i++ {
			ch3 <- i
		}
		close(ch3)
	}()

	// ادغام کانال‌ها
	fanIn := NewFanIn(ch1, ch2, ch3)

	// خواندن از کانال ادغام شده
	results := make([]int, 0)
	for value := range fanIn.Merge() {
		results = append(results, value)
		fmt.Printf("🔀 Merged value: %d\n", value)
	}

	fmt.Printf("\n✅ Fan-In completed: merged %d values from 3 channels\n", len(results))
}

// ============================================================================
// بخش 4: الگوی Pipeline (خط لوله)
// ============================================================================

// PipelineStage یک مرحله از Pipeline
type PipelineStage func(<-chan int) <-chan int

// Pipeline ساختار Pipeline
type Pipeline struct {
	stages []PipelineStage
}

// NewPipeline ایجاد Pipeline جدید
func NewPipeline(stages ...PipelineStage) *Pipeline {
	return &Pipeline{
		stages: stages,
	}
}

// Execute اجرای Pipeline
func (p *Pipeline) Execute(input <-chan int) <-chan int {
	output := input

	for _, stage := range p.stages {
		output = stage(output)
	}

	return output
}

// stage1: تولید اعداد (Generator Stage)
func generatorStage(nums ...int) PipelineStage {
	return func(input <-chan int) <-chan int {
		// این stage نیازی به input ندارد (اولین stage است)
		output := make(chan int)

		go func() {
			for _, n := range nums {
				output <- n
			}
			close(output)
		}()

		return output
	}
}

// stage2: ضرب در 2
func multiplyStage(factor int) PipelineStage {
	return func(input <-chan int) <-chan int {
		output := make(chan int)

		go func() {
			for value := range input {
				output <- value * factor
			}
			close(output)
		}()

		return output
	}
}

// stage3: جمع با مقدار
func addStage(amount int) PipelineStage {
	return func(input <-chan int) <-chan int {
		output := make(chan int)

		go func() {
			for value := range input {
				output <- value + amount
			}
			close(output)
		}()

		return output
	}
}

// stage4: فیلتر کردن (حذف اعداد زوج)
func filterEvenStage() PipelineStage {
	return func(input <-chan int) <-chan int {
		output := make(chan int)

		go func() {
			for value := range input {
				if value%2 != 0 { // فقط فردها
					output <- value
				}
			}
			close(output)
		}()

		return output
	}
}

// stage5: نمایش نتایج
func printStage() PipelineStage {
	return func(input <-chan int) <-chan int {
		output := make(chan int)

		go func() {
			for value := range input {
				fmt.Printf("  📍 Pipeline output: %d\n", value)
				output <- value
			}
			close(output)
		}()

		return output
	}
}

// demonstratePipeline نمایش Pipeline
func demonstratePipeline() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔗 PIPELINE PATTERN")
	fmt.Println(strings.Repeat("=", 80))

	// ایجاد Pipeline با چند stage
	pipeline := NewPipeline(
		generatorStage(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
		multiplyStage(2),  // ضرب در 2: 2,4,6,8,10,12,14,16,18,20
		addStage(1),       // جمع با 1: 3,5,7,9,11,13,15,17,19,21
		filterEvenStage(), // حذف زوج‌ها: 3,5,7,9,11,13,15,17,19,21 (همگی فرد)
		printStage(),      // چاپ نتایج
	)

	start := time.Now()

	// اجرای Pipeline (ورودی nil چون generator خودش تولید می‌کند)
	var nilChan chan int
	for range pipeline.Execute(nilChan) {
		// مصرف نتایج
	}

	elapsed := time.Since(start)
	fmt.Printf("\n✅ Pipeline completed in %v\n", elapsed)
}

// ============================================================================
// بخش 5: Pipeline پیشرفته با Worker Pool در هر Stage
// ============================================================================

// ParallelStage مرحله موازی با Worker Pool داخلی
type ParallelStage struct {
	workerCount int
	processFunc func(int) int
}

// NewParallelStage ایجاد مرحله موازی
func NewParallelStage(workerCount int, processFunc func(int) int) *ParallelStage {
	return &ParallelStage{
		workerCount: workerCount,
		processFunc: processFunc,
	}
}

// Run اجرای مرحله موازی
func (ps *ParallelStage) Run(input <-chan int) <-chan int {
	output := make(chan int, ps.workerCount)
	var wg sync.WaitGroup

	// راه‌اندازی worker pool داخلی
	for i := 0; i < ps.workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for value := range input {
				output <- ps.processFunc(value)
			}
		}()
	}

	// بستن output بعد از اتمام
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// demonstrateAdvancedPipeline Pipeline پیشرفته با stageهای موازی
func demonstrateAdvancedPipeline() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⚡ ADVANCED PIPELINE (Parallel Stages)")
	fmt.Println(strings.Repeat("=", 80))

	// مرحله 1: تولید اعداد
	input := make(chan int, 100)
	go func() {
		for i := 1; i <= 100; i++ {
			input <- i
		}
		close(input)
	}()

	// مرحله 2: پردازش سنگین با 4 worker
	heavyProcess := NewParallelStage(4, func(x int) int {
		time.Sleep(5 * time.Millisecond) // شبیه‌سازی کار سنگین
		return x * x
	})

	// مرحله 3: فیلتر با 2 worker
	filterProcess := NewParallelStage(2, func(x int) int {
		return x // اینجا فقط عبور می‌دهد
	})

	start := time.Now()

	// اجرای pipeline
	stage1 := heavyProcess.Run(input)
	stage2 := filterProcess.Run(stage1)

	// جمع‌آوری نتایج
	count := 0
	for range stage2 {
		count++
	}

	elapsed := time.Since(start)
	fmt.Printf("✅ Advanced pipeline processed %d items in %v with parallel stages\n",
		count, elapsed)
}

// ============================================================================
// بخش 6: مقایسه و انتخاب الگوی مناسب
// ============================================================================

func comparePatterns() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 PATTERN COMPARISON & SELECTION GUIDE")
	fmt.Println(strings.Repeat("=", 80))

	comparison := []struct {
		Pattern    string
		BestFor    string
		Example    string
		Complexity string
	}{
		{
			Pattern:    "Worker Pool",
			BestFor:    "پردازش دسته‌ای کارهای مستقل با محدودیت همزمانی",
			Example:    "پردازش تصاویر، ارسال ایمیل، درخواست‌های HTTP",
			Complexity: "ساده",
		},
		{
			Pattern:    "Fan-Out",
			BestFor:    "موازی‌سازی سریع یک جریان داده بین چند پردازشگر",
			Example:    "پردازش همزمان چند فایل، جستجوی موازی",
			Complexity: "متوسط",
		},
		{
			Pattern:    "Fan-In",
			BestFor:    "جمع‌آوری نتایج از چند منبع موازی",
			Example:    "جمع‌آوری نتایج جستجو از چند API, merge sort",
			Complexity: "متوسط",
		},
		{
			Pattern:    "Pipeline",
			BestFor:    "پردازش جریان داده با چند مرحله متوالی",
			Example:    "ETL pipeline, stream processing, تبدیل داده",
			Complexity: "پیشرفته",
		},
	}

	fmt.Println("\n| Pattern     | Best For                          | Example                    |")
	fmt.Println("|-------------|-----------------------------------|----------------------------|")
	for _, p := range comparison {
		fmt.Printf("| %-11s | %-33s | %-26s |\n",
			p.Pattern, p.BestFor[:33], p.Example[:26])
	}

	fmt.Println("\n📌 انتخاب بر اساس نیاز:")
	fmt.Println("  • کارهای مستقل زیاد → Worker Pool")
	fmt.Println("  • یک ورودی، چند پردازش → Fan-Out")
	fmt.Println("  • چند ورودی، یک خروجی → Fan-In")
	fmt.Println("  • چند مرحله متوالی → Pipeline")
	fmt.Println("  • ترکیبی → Fan-Out + Fan-In + Worker Pool در Pipeline")
}

// ============================================================================
// بخش 7: الگوی ترکیبی (Fan-Out + Fan-In + Worker Pool)
// ============================================================================

// ComplexWorkflow ترکیب همه الگوها
type ComplexWorkflow struct {
	workerCount int
}

// NewComplexWorkflow ایجاد workflow ترکیبی
func NewComplexWorkflow(workers int) *ComplexWorkflow {
	return &ComplexWorkflow{
		workerCount: workers,
	}
}

// Process اجرای workflow ترکیبی
func (cw *ComplexWorkflow) Process(items []int) []int {
	fmt.Println("\n🔄 Complex Workflow: Fan-Out → Worker Pool → Fan-In")

	// مرحله 1: Fan-Out (توزیع بین کانال‌ها)
	inputCh := make(chan int, len(items))
	for _, item := range items {
		inputCh <- item
	}
	close(inputCh)

	// مرحله 2: Worker Pool (پردازش موازی)
	jobCh := make(chan int, cw.workerCount)
	resultCh := make(chan int, cw.workerCount)

	var wg sync.WaitGroup
	for i := 0; i < cw.workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobCh {
				// پردازش سنگین
				time.Sleep(10 * time.Millisecond)
				result := job * job
				resultCh <- result
				fmt.Printf("  Worker %d processed %d -> %d\n", id, job, result)
			}
		}(i)
	}

	// توزیع jobها بین workerها
	go func() {
		for item := range inputCh {
			jobCh <- item
		}
		close(jobCh)
		wg.Wait()
		close(resultCh)
	}()

	// مرحله 3: جمع‌آوری نتایج (Fan-In)
	results := make([]int, 0)
	for result := range resultCh {
		results = append(results, result)
	}

	return results
}

func demonstrateComplexWorkflow() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎯 COMPLEX WORKFLOW (All Patterns Combined)")
	fmt.Println(strings.Repeat("=", 80))

	workflow := NewComplexWorkflow(4)
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	start := time.Now()
	results := workflow.Process(items)
	elapsed := time.Since(start)

	fmt.Printf("\n✅ Results: %v\n", results)
	fmt.Printf("✅ Processed %d items with %d workers in %v\n",
		len(results), 4, elapsed)
}

// ============================================================================
// بخش 8: الگوهای درست و غلط
// ============================================================================

func demonstrateCorrectVsWrong() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✅❌ CORRECT VS WRONG PATTERNS")
	fmt.Println(strings.Repeat("=", 80))

	// ✅ درست: Worker Pool با اندازه محدود
	fmt.Println("\n✅ Correct: Fixed-size worker pool")
	correctPool := make(chan struct{}, 5)
	for i := 0; i < 10; i++ {
		correctPool <- struct{}{}
		go func() {
			defer func() { <-correctPool }()
			// کار
		}()
	}

	// ❌ غلط: بدون محدودیت (نشت گوروتین)
	fmt.Println("\n❌ Wrong: Unlimited goroutine creation")
	fmt.Println("  for i := 0; i < 100000; i++ { go work() } // Bad!")

	// ✅ درست: بستن کانال در Fan-In بعد از اتمام همه
	fmt.Println("\n✅ Correct: Close merged channel after all sources done")

	// ❌ غلط: بستن کانال در Fan-Out قبل از اتمام همه
	fmt.Println("\n❌ Wrong: Close output channel before all workers done")

	// ✅ درست: استفاده از select برای non-blocking در Worker Pool
	fmt.Println("\n✅ Correct: Non-blocking submit with select")

	// ❌ غلط: نادیده گرفتن context در Pipeline
	fmt.Println("\n❌ Wrong: Pipeline stages without cancellation support")
}

// ============================================================================
// بخش 9: جمع‌بندی و نکات نهایی
// ============================================================================

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("🎯 CONCURRENCY PATTERNS IN GO")
	fmt.Println("Worker Pool | Fan-Out | Fan-In | Pipeline")
	fmt.Println(strings.Repeat("=", 80))

	// بخش 1: Worker Pool
	demonstrateWorkerPool()

	// بخش 2: Fan-Out
	demonstrateFanOut()

	// بخش 3: Fan-In
	demonstrateFanIn()

	// بخش 4: Pipeline
	demonstratePipeline()

	// بخش 5: Pipeline پیشرفته
	demonstrateAdvancedPipeline()

	// بخش 6: مقایسه الگوها
	comparePatterns()

	// بخش 7: الگوی ترکیبی
	demonstrateComplexWorkflow()

	// بخش 8: الگوهای درست و غلط
	demonstrateCorrectVsWrong()

	// جمع‌بندی نهایی
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📚 FINAL SUMMARY")
	fmt.Println(strings.Repeat("=", 80))

	fmt.Println("\n🎯 When to use which pattern:")
	fmt.Println("  ┌─────────────────────────────────────────────────────────────┐")
	fmt.Println("  │ Worker Pool → محدود کردن همزمانی، پردازش دسته‌ای            │")
	fmt.Println("  │ Fan-Out      → توزیع کار بین چند پردازشگر موازی              │")
	fmt.Println("  │ Fan-In       → جمع‌آوری نتایج از چند منبع                    │")
	fmt.Println("  │ Pipeline     → پردازش جریان داده با چند مرحله متوالی         │")
	fmt.Println("  │ Combined     → سناریوهای پیچیده واقعی (ETL, Web Crawler)    │")
	fmt.Println("  └─────────────────────────────────────────────────────────────┘")

	fmt.Println("\n💡 Golden Rules:")
	fmt.Println("  1. Worker Pool: همیشه اندازه pool را محدود کنید")
	fmt.Println("  2. Fan-Out: هر worker باید کانال خروجی خود را داشته باشد")
	fmt.Println("  3. Fan-In: از WaitGroup برای هماهنگی بستن کانال استفاده کنید")
	fmt.Println("  4. Pipeline: هر stage باید قابلیت کنسل شدن داشته باشد")
	fmt.Println("  5. همه: همیشه راهی برای خروج گوروتین‌ها وجود داشته باشد")

	fmt.Println("\n🚀 Performance Tips:")
	fmt.Println("  • Worker count = runtime.NumCPU() برای کارهای محاسباتی")
	fmt.Println("  • Worker count = محدودیت منبع خارجی (DB, API) برای I/O")
	fmt.Println("  • Buffer size = نرخ تولید / نرخ مصرف برای جذب نوسانات")
	fmt.Println("  • همیشه benchmark کنید، حدس نزنید")
}

// ============================================================================
// بخش 10: توابع کمکی
// ============================================================================

// strings.Repeat برای استفاده در این فایل (چون نمی‌خواهیم import کنیم)
// در فایل واقعی از import "strings" استفاده کنید
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// جایگزین ساده برای strings.Repeat
var strings = struct {
	Repeat func(string, int) string
}{
	Repeat: stringsRepeat,
}

// BenchmarkWorkerPool benchmark ساده برای Worker Pool
func BenchmarkWorkerPool(items int, workers int) time.Duration {
	pool := NewWorkerPool(workers, items)
	pool.Start()

	for i := 0; i < items; i++ {
		pool.Submit(Job{ID: i, Number: i})
	}

	start := time.Now()
	count := 0
	for range pool.Results() {
		count++
	}
	elapsed := time.Since(start)
	pool.Stop()

	fmt.Printf("  Benchmark: %d items with %d workers → %v\n", items, workers, elapsed)
	return elapsed
}

// AdaptiveWorkerPool Worker Pool با تعداد worker تطبیقی
type AdaptiveWorkerPool struct {
	minWorkers int
	maxWorkers int
	queueSize  int
	current    atomic.Int32
}

// NewAdaptiveWorkerPool ایجاد Worker Pool تطبیقی
func NewAdaptiveWorkerPool(min, max, queueSize int) *AdaptiveWorkerPool {
	return &AdaptiveWorkerPool{
		minWorkers: min,
		maxWorkers: max,
		queueSize:  queueSize,
	}
}

// Start راه‌اندازی با مانیتورینگ خودکار
func (awp *AdaptiveWorkerPool) Start() {
	go awp.monitor()
}

func (awp *AdaptiveWorkerPool) monitor() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		current := awp.current.Load()
		queueLoad := 0 // محاسبه بار صف

		if queueLoad > 80 && current < int32(awp.maxWorkers) {
			// افزایش worker
			awp.current.Add(1)
		} else if queueLoad < 20 && current > int32(awp.minWorkers) {
			// کاهش worker
			awp.current.Add(-1)
		}
	}
}
