// ============================================================================
// FILE: inmemory_queue_guide.go
// TITLE: راهنمای کامل In-Memory Message Queue با Channel و گوروتین در Go
// HOW TO RUN: go run inmemory_queue_guide.go
// ============================================================================

// ============================================================================
// بخش 0: مقدمه - In-Memory Message Queue چیست؟
// ============================================================================
//
// In-Memory Message Queue یک صف پیام ساده است که در حافظه (RAM) اجرا می‌شود.
// در Go می‌توانیم با استفاده از channelها و گوروتین‌ها به سادگی آن را پیاده‌سازی کنیم.
//
// مزایا:
// 1. سرعت بسیار بالا (بدون I/O)
// 2. سادگی پیاده‌سازی
// 3. عدم نیاز به وابستگی خارجی
//
// معایب:
// 1. از دست رفتن داده در صورت crash برنامه
// 2. محدود به حافظه
// 3. مقیاس‌پذیری محدود
//
// کاربردها:
// 1. پردازش داخلی بین گوروتین‌ها
// 2. بافر کردن موقت درخواست‌ها
// 3. Fan-out و Fan-in داخلی
//
// قانون طلایی:
// "برای ارتباط داخلی بین گوروتین‌ها از channel استفاده کن.
//  برای نیازهای ساده و موقت، in-memory queue کافی است.
//  برای نیازهای production و persistent از message broker واقعی استفاده کن."
// ============================================================================

package __inmemory_queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// بخش 1: انواع پیام و Queue پایه
// ============================================================================

// MessageType نوع پیام
type MessageType string

const (
	TypeTask   MessageType = "task"
	TypeEvent  MessageType = "event"
	TypeSignal MessageType = "signal"
)

// Message ساختار پیام
type Message struct {
	ID         string            `json:"id"`
	Type       MessageType       `json:"type"`
	Payload    interface{}       `json:"payload"`
	Headers    map[string]string `json:"headers"`
	CreatedAt  time.Time         `json:"created_at"`
	RetryCount int               `json:"retry_count"`
}

// NewMessage ایجاد پیام جدید
func NewMessage(msgType MessageType, payload interface{}) *Message {
	return &Message{
		ID:         uuid.New().String(),
		Type:       msgType,
		Payload:    payload,
		Headers:    make(map[string]string),
		CreatedAt:  time.Now(),
		RetryCount: 0,
	}
}

// SimpleQueue ساده‌ترین صف با channel
type SimpleQueue struct {
	queue chan *Message
	mu    sync.RWMutex
}

// NewSimpleQueue ایجاد صف ساده
func NewSimpleQueue(size int) *SimpleQueue {
	return &SimpleQueue{
		queue: make(chan *Message, size),
	}
}

// Enqueue افزودن پیام به صف
func (q *SimpleQueue) Enqueue(msg *Message) error {
	select {
	case q.queue <- msg:
		return nil
	default:
		return errors.New("queue is full")
	}
}

// EnqueueWithTimeout افزودن پیام با timeout
func (q *SimpleQueue) EnqueueWithTimeout(msg *Message, timeout time.Duration) error {
	select {
	case q.queue <- msg:
		return nil
	case <-time.After(timeout):
		return errors.New("enqueue timeout")
	}
}

// Dequeue خارج کردن پیام از صف (مسدودکننده)
func (q *SimpleQueue) Dequeue() *Message {
	return <-q.queue
}

// DequeueNonBlocking خارج کردن پیام بدون مسدود شدن
func (q *SimpleQueue) DequeueNonBlocking() (*Message, error) {
	select {
	case msg := <-q.queue:
		return msg, nil
	default:
		return nil, errors.New("queue is empty")
	}
}

// DequeueWithTimeout خارج کردن پیام با timeout
func (q *SimpleQueue) DequeueWithTimeout(timeout time.Duration) (*Message, error) {
	select {
	case msg := <-q.queue:
		return msg, nil
	case <-time.After(timeout):
		return nil, errors.New("dequeue timeout")
	}
}

// Len طول صف
func (q *SimpleQueue) Len() int {
	return len(q.queue)
}

// Cap ظرفیت صف
func (q *SimpleQueue) Cap() int {
	return cap(q.queue)
}

// Close بستن صف
func (q *SimpleQueue) Close() {
	close(q.queue)
}

// ============================================================================
// بخش 2: Worker Pool با Queue
// ============================================================================

// Worker ساختار worker
type Worker struct {
	ID      int
	Queue   *SimpleQueue
	Handler func(*Message) error
	StopCh  chan struct{}
	Wg      *sync.WaitGroup
	Stats   *WorkerStats
}

type WorkerStats struct {
	ProcessedCount int64
	ErrorCount     int64
	LastProcessed  time.Time
}

// NewWorker ایجاد worker جدید
func NewWorker(id int, queue *SimpleQueue, handler func(*Message) error) *Worker {
	return &Worker{
		ID:      id,
		Queue:   queue,
		Handler: handler,
		StopCh:  make(chan struct{}),
		Wg:      &sync.WaitGroup{},
		Stats:   &WorkerStats{},
	}
}

// Start شروع worker
func (w *Worker) Start() {
	w.Wg.Add(1)
	go w.run()
}

func (w *Worker) run() {
	defer w.Wg.Done()

	for {
		select {
		case <-w.StopCh:
			log.Printf("Worker %d stopping", w.ID)
			return
		default:
			msg := w.Queue.Dequeue()
			if msg == nil {
				continue
			}

			if err := w.Handler(msg); err != nil {
				atomic.AddInt64(&w.Stats.ErrorCount, 1)
				log.Printf("Worker %d error processing message %s: %v", w.ID, msg.ID, err)
			} else {
				atomic.AddInt64(&w.Stats.ProcessedCount, 1)
				w.Stats.LastProcessed = time.Now()
			}
		}
	}
}

// Stop توقف worker
func (w *Worker) Stop() {
	close(w.StopCh)
	w.Wg.Wait()
}

// WorkerPool مجموعه workerها
type WorkerPool struct {
	workers []*Worker
	queue   *SimpleQueue
	wg      sync.WaitGroup
	stopCh  chan struct{}
	stats   *PoolStats
}

type PoolStats struct {
	TotalProcessed int64
	TotalErrors    int64
	WorkerStats    []*WorkerStats
}

// NewWorkerPool ایجاد worker pool
func NewWorkerPool(queue *SimpleQueue, workerCount int, handler func(*Message) error) *WorkerPool {
	pool := &WorkerPool{
		workers: make([]*Worker, workerCount),
		queue:   queue,
		stopCh:  make(chan struct{}),
		stats: &PoolStats{
			WorkerStats: make([]*WorkerStats, workerCount),
		},
	}

	for i := 0; i < workerCount; i++ {
		worker := NewWorker(i, queue, handler)
		pool.workers[i] = worker
		pool.stats.WorkerStats[i] = worker.Stats
	}

	return pool
}

// Start شروع worker pool
func (p *WorkerPool) Start() {
	for _, worker := range p.workers {
		worker.Start()
	}

	p.wg.Add(1)
	go p.monitor()
}

func (p *WorkerPool) monitor() {
	defer p.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			p.updateStats()
		}
	}
}

func (p *WorkerPool) updateStats() {
	var totalProcessed, totalErrors int64
	for _, worker := range p.workers {
		totalProcessed += atomic.LoadInt64(&worker.Stats.ProcessedCount)
		totalErrors += atomic.LoadInt64(&worker.Stats.ErrorCount)
	}

	atomic.StoreInt64(&p.stats.TotalProcessed, totalProcessed)
	atomic.StoreInt64(&p.stats.TotalErrors, totalErrors)
}

// GetStats دریافت آمار
func (p *WorkerPool) GetStats() PoolStats {
	return PoolStats{
		TotalProcessed: atomic.LoadInt64(&p.stats.TotalProcessed),
		TotalErrors:    atomic.LoadInt64(&p.stats.TotalErrors),
		WorkerStats:    p.stats.WorkerStats,
	}
}

// Stop توقف worker pool
func (p *WorkerPool) Stop() {
	close(p.stopCh)
	for _, worker := range p.workers {
		worker.Stop()
	}
	p.wg.Wait()
	close(p.queue.queue)
}

// ============================================================================
// بخش 3: Priority Queue (صف اولویت‌دار)
// ============================================================================

// PriorityMessage پیام با اولویت
type PriorityMessage struct {
	Message  *Message
	Priority int
	Index    int
}

// PriorityQueue صف اولویت‌دار (min-heap)
type PriorityQueue struct {
	items    []*PriorityMessage
	mu       sync.RWMutex
	notEmpty chan struct{}
	closed   bool
}

// NewPriorityQueue ایجاد صف اولویت‌دار
func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{
		items:    make([]*PriorityMessage, 0),
		notEmpty: make(chan struct{}, 1),
	}
}

// Push افزودن پیام با اولویت (عدد کمتر = اولویت بیشتر)
func (pq *PriorityQueue) Push(msg *Message, priority int) error {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if pq.closed {
		return errors.New("queue is closed")
	}

	pqItem := &PriorityMessage{
		Message:  msg,
		Priority: priority,
	}

	pq.items = append(pq.items, pqItem)
	pq.up(len(pq.items) - 1)

	// علامت دادن به Pop که پیام جدید وجود دارد
	select {
	case pq.notEmpty <- struct{}{}:
	default:
	}

	return nil
}

// Pop خارج کردن پیام با بالاترین اولویت
func (pq *PriorityQueue) Pop() (*Message, error) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if len(pq.items) == 0 {
		if pq.closed {
			return nil, errors.New("queue is closed")
		}
		// منتظر پیام جدید
		pq.mu.Unlock()
		<-pq.notEmpty
		pq.mu.Lock()

		if len(pq.items) == 0 {
			return nil, errors.New("queue is empty")
		}
	}

	min := pq.items[0]
	last := pq.items[len(pq.items)-1]
	pq.items[0] = last
	pq.items = pq.items[:len(pq.items)-1]

	if len(pq.items) > 0 {
		pq.down(0)
	}

	return min.Message, nil
}

// up حفظ خاصیت heap (بالا بردن)
func (pq *PriorityQueue) up(i int) {
	for {
		parent := (i - 1) / 2
		if parent == i || pq.items[parent].Priority <= pq.items[i].Priority {
			break
		}
		pq.items[parent], pq.items[i] = pq.items[i], pq.items[parent]
		i = parent
	}
}

// down حفظ خاصیت heap (پایین بردن)
func (pq *PriorityQueue) down(i int) {
	n := len(pq.items)
	for {
		left := 2*i + 1
		if left >= n {
			break
		}
		smallest := left
		right := left + 1
		if right < n && pq.items[right].Priority < pq.items[left].Priority {
			smallest = right
		}
		if pq.items[i].Priority <= pq.items[smallest].Priority {
			break
		}
		pq.items[i], pq.items[smallest] = pq.items[smallest], pq.items[i]
		i = smallest
	}
}

// Len طول صف
func (pq *PriorityQueue) Len() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()
	return len(pq.items)
}

// Close بستن صف
func (pq *PriorityQueue) Close() {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	pq.closed = true
	close(pq.notEmpty)
}

// ============================================================================
// بخش 4: Pub/Sub (انتشار و اشتراک)
// ============================================================================

// Subscriber مشترک پیام‌ها
type Subscriber struct {
	ID       string
	Messages chan *Message
	Topics   map[string]bool
	Active   bool
	mu       sync.RWMutex
}

// NewSubscriber ایجاد مشترک جدید
func NewSubscriber() *Subscriber {
	return &Subscriber{
		ID:       uuid.New().String(),
		Messages: make(chan *Message, 100),
		Topics:   make(map[string]bool),
		Active:   true,
	}
}

// TopicSubscribe مشترک در topic
func (s *Subscriber) TopicSubscribe(topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Topics[topic] = true
}

// TopicUnsubscribe لغو اشتراک
func (s *Subscriber) TopicUnsubscribe(topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Topics, topic)
}

// GetMessages دریافت کانال پیام‌ها
func (s *Subscriber) GetMessages() <-chan *Message {
	return s.Messages
}

// Close بستن مشترک
func (s *Subscriber) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Active = false
	close(s.Messages)
}

// PubSub سیستم انتشار-اشتراک
type PubSub struct {
	subscribers map[string]*Subscriber
	topics      map[string][]*Subscriber
	mu          sync.RWMutex
}

// NewPubSub ایجاد سیستم Pub/Sub
func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string]*Subscriber),
		topics:      make(map[string][]*Subscriber),
	}
}

// Subscribe ثبت مشترک
func (ps *PubSub) Subscribe(subscriber *Subscriber, topics ...string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.subscribers[subscriber.ID] = subscriber

	for _, topic := range topics {
		subscriber.TopicSubscribe(topic)
		ps.topics[topic] = append(ps.topics[topic], subscriber)
	}
}

// Unsubscribe لغو اشتراک
func (ps *PubSub) Unsubscribe(subscriberID string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subscriber, exists := ps.subscribers[subscriberID]
	if !exists {
		return
	}

	for topic := range subscriber.Topics {
		ps.removeFromTopic(topic, subscriberID)
	}

	delete(ps.subscribers, subscriberID)
	subscriber.Close()
}

func (ps *PubSub) removeFromTopic(topic, subscriberID string) {
	subscribers := ps.topics[topic]
	for i, sub := range subscribers {
		if sub.ID == subscriberID {
			ps.topics[topic] = append(subscribers[:i], subscribers[i+1:]...)
			break
		}
	}

	if len(ps.topics[topic]) == 0 {
		delete(ps.topics, topic)
	}
}

// Publish انتشار پیام به topic
func (ps *PubSub) Publish(topic string, msg *Message) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	subscribers, exists := ps.topics[topic]
	if !exists {
		return
	}

	for _, subscriber := range subscribers {
		if subscriber.Active {
			select {
			case subscriber.Messages <- msg:
			default:
				log.Printf("Subscriber %s buffer full, dropping message", subscriber.ID)
			}
		}
	}
}

// Broadcast پخش پیام به همه topics
func (ps *PubSub) Broadcast(msg *Message) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for topic := range ps.topics {
		msg.Headers["topic"] = topic
		ps.Publish(topic, msg)
	}
}

// ============================================================================
// بخش 5: Delayed Queue (صف با تأخیر)
// ============================================================================

// DelayedMessage پیام با تأخیر
type DelayedMessage struct {
	Message   *Message
	ExecuteAt time.Time
}

// DelayedQueue صف با تأخیر
type DelayedQueue struct {
	queue   chan *DelayedMessage
	stopCh  chan struct{}
	handler func(*Message) error
	wg      sync.WaitGroup
}

// NewDelayedQueue ایجاد صف با تأخیر
func NewDelayedQueue(handler func(*Message) error) *DelayedQueue {
	dq := &DelayedQueue{
		queue:   make(chan *DelayedMessage, 1000),
		stopCh:  make(chan struct{}),
		handler: handler,
	}

	dq.start()
	return dq
}

// Add افزودن پیام با تأخیر
func (dq *DelayedQueue) Add(msg *Message, delay time.Duration) error {
	delayedMsg := &DelayedMessage{
		Message:   msg,
		ExecuteAt: time.Now().Add(delay),
	}

	select {
	case dq.queue <- delayedMsg:
		return nil
	default:
		return errors.New("delayed queue is full")
	}
}

func (dq *DelayedQueue) start() {
	dq.wg.Add(1)
	go dq.process()
}

func (dq *DelayedQueue) process() {
	defer dq.wg.Done()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	messages := make([]*DelayedMessage, 0)

	for {
		select {
		case <-dq.stopCh:
			return
		case <-ticker.C:
			now := time.Now()
			var ready []*DelayedMessage
			var remaining []*DelayedMessage

			// دسته‌بندی پیام‌های آماده
			for len(dq.queue) > 0 {
				select {
				case msg := <-dq.queue:
					if msg.ExecuteAt.Before(now) || msg.ExecuteAt.Equal(now) {
						ready = append(ready, msg)
					} else {
						remaining = append(remaining, msg)
					}
				default:
					break
				}
			}

			// بازگرداندن پیام‌های باقیمانده
			for _, msg := range remaining {
				dq.queue <- msg
			}

			// پردازش پیام‌های آماده
			for _, msg := range ready {
				if err := dq.handler(msg.Message); err != nil {
					log.Printf("Error processing delayed message: %v", err)
				}
			}
		}
	}
}

// Stop توقف صف
func (dq *DelayedQueue) Stop() {
	close(dq.stopCh)
	dq.wg.Wait()
	close(dq.queue)
}

// ============================================================================
// بخش 6: Multi-Queue (چندین صف با یک دیسپچر)
// ============================================================================

// MultiQueue چندین صف با توزیع round-robin
type MultiQueue struct {
	queues      []*SimpleQueue
	current     int
	mu          sync.Mutex
	workerCount int
}

// NewMultiQueue ایجاد چند صف
func NewMultiQueue(queueCount, queueSize int) *MultiQueue {
	mq := &MultiQueue{
		queues:      make([]*SimpleQueue, queueCount),
		workerCount: queueCount,
	}

	for i := 0; i < queueCount; i++ {
		mq.queues[i] = NewSimpleQueue(queueSize)
	}

	return mq
}

// Enqueue افزودن پیام به صف بعدی (round-robin)
func (mq *MultiQueue) Enqueue(msg *Message) error {
	mq.mu.Lock()
	idx := mq.current
	mq.current = (mq.current + 1) % mq.workerCount
	mq.mu.Unlock()

	return mq.queues[idx].Enqueue(msg)
}

// GetQueue دریافت صف با index
func (mq *MultiQueue) GetQueue(index int) *SimpleQueue {
	if index < 0 || index >= len(mq.queues) {
		return nil
	}
	return mq.queues[index]
}

// Len مجموع طول همه صف‌ها
func (mq *MultiQueue) Len() int {
	total := 0
	for _, q := range mq.queues {
		total += q.Len()
	}
	return total
}

// ============================================================================
// بخش 7: Batch Queue (صف دسته‌ای)
// ============================================================================

// BatchQueue صف برای پردازش دسته‌ای
type BatchQueue struct {
	queue     chan *Message
	batchSize int
	batchWait time.Duration
	handler   func([]*Message) error
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

// NewBatchQueue ایجاد صف دسته‌ای
func NewBatchQueue(batchSize int, batchWait time.Duration, handler func([]*Message) error) *BatchQueue {
	bq := &BatchQueue{
		queue:     make(chan *Message, 10000),
		batchSize: batchSize,
		batchWait: batchWait,
		handler:   handler,
		stopCh:    make(chan struct{}),
	}

	bq.start()
	return bq
}

// Add افزودن پیام
func (bq *BatchQueue) Add(msg *Message) error {
	select {
	case bq.queue <- msg:
		return nil
	default:
		return errors.New("batch queue is full")
	}
}

func (bq *BatchQueue) start() {
	bq.wg.Add(1)
	go bq.process()
}

func (bq *BatchQueue) process() {
	defer bq.wg.Done()

	batch := make([]*Message, 0, bq.batchSize)
	ticker := time.NewTicker(bq.batchWait)
	defer ticker.Stop()

	for {
		select {
		case <-bq.stopCh:
			if len(batch) > 0 {
				bq.handler(batch)
			}
			return

		case msg := <-bq.queue:
			batch = append(batch, msg)

			if len(batch) >= bq.batchSize {
				if err := bq.handler(batch); err != nil {
					log.Printf("Error processing batch: %v", err)
				}
				batch = make([]*Message, 0, bq.batchSize)
			}

		case <-ticker.C:
			if len(batch) > 0 {
				if err := bq.handler(batch); err != nil {
					log.Printf("Error processing batch: %v", err)
				}
				batch = make([]*Message, 0, bq.batchSize)
			}
		}
	}
}

// Stop توقف صف
func (bq *BatchQueue) Stop() {
	close(bq.stopCh)
	bq.wg.Wait()
	close(bq.queue)
}

// ============================================================================
// بخش 8: Example Application (مثال کاربردی)
// ============================================================================

// Task نوع task
type Task struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}

// OrderTaskHandler handler tasks
func OrderTaskHandler(msg *Message) error {
	var task Task
	if err := json.Unmarshal(msg.Payload.([]byte), &task); err != nil {
		return err
	}

	log.Printf("Processing task %s: %s (duration: %dms)", task.ID, task.Name, task.Duration)
	time.Sleep(time.Duration(task.Duration) * time.Millisecond)

	return nil
}

func runQueueExample() {
	log.Println("=== In-Memory Message Queue Example ===\n")

	// 1. Simple Queue
	log.Println("1. Testing Simple Queue...")
	simpleQueue := NewSimpleQueue(10)

	go func() {
		for i := 0; i < 5; i++ {
			msg := NewMessage(TypeTask, []byte(fmt.Sprintf("Task %d", i)))
			if err := simpleQueue.Enqueue(msg); err != nil {
				log.Printf("Enqueue error: %v", err)
			}
		}
	}()

	for i := 0; i < 5; i++ {
		msg := simpleQueue.Dequeue()
		log.Printf("Dequeued: %s", msg.ID)
	}

	// 2. Worker Pool
	log.Println("\n2. Testing Worker Pool...")
	workerQueue := NewSimpleQueue(100)
	pool := NewWorkerPool(workerQueue, 3, OrderTaskHandler)
	pool.Start()

	// ارسال taskها
	for i := 0; i < 10; i++ {
		task := Task{
			ID:       uuid.New().String(),
			Name:     fmt.Sprintf("Task %d", i),
			Duration: 100 + i*10,
		}
		data, _ := json.Marshal(task)
		msg := NewMessage(TypeTask, data)
		workerQueue.Enqueue(msg)
	}

	time.Sleep(2 * time.Second)
	stats := pool.GetStats()
	log.Printf("Pool stats - Processed: %d, Errors: %d", stats.TotalProcessed, stats.TotalErrors)

	pool.Stop()

	// 3. Priority Queue
	log.Println("\n3. Testing Priority Queue...")
	priorityQueue := NewPriorityQueue()

	go func() {
		priorities := []int{3, 1, 4, 1, 5, 2, 3}
		for i, pri := range priorities {
			msg := NewMessage(TypeTask, fmt.Sprintf("Task %d", i))
			priorityQueue.Push(msg, pri)
			log.Printf("Pushed task %d with priority %d", i, pri)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	for priorityQueue.Len() > 0 {
		msg, _ := priorityQueue.Pop()
		log.Printf("Popped message %s", msg.ID)
	}

	// 4. Pub/Sub
	log.Println("\n4. Testing Pub/Sub...")
	pubsub := NewPubSub()

	// ایجاد مشترکین
	sub1 := NewSubscriber()
	sub2 := NewSubscriber()
	sub3 := NewSubscriber()

	pubsub.Subscribe(sub1, "orders", "payments")
	pubsub.Subscribe(sub2, "orders")
	pubsub.Subscribe(sub3, "notifications")

	// مصرف پیام‌ها
	go func() {
		for msg := range sub1.GetMessages() {
			log.Printf("Subscriber 1 received: %s (type: %s)", msg.ID, msg.Type)
		}
	}()
	go func() {
		for msg := range sub2.GetMessages() {
			log.Printf("Subscriber 2 received: %s (type: %s)", msg.ID, msg.Type)
		}
	}()
	go func() {
		for msg := range sub3.GetMessages() {
			log.Printf("Subscriber 3 received: %s (type: %s)", msg.ID, msg.Type)
		}
	}()

	// انتشار پیام‌ها
	pubsub.Publish("orders", NewMessage(TypeEvent, "New order created"))
	pubsub.Publish("payments", NewMessage(TypeEvent, "Payment received"))
	pubsub.Publish("notifications", NewMessage(TypeSignal, "Notification sent"))

	time.Sleep(100 * time.Millisecond)

	// 5. Batch Queue
	log.Println("\n5. Testing Batch Queue...")

	batchHandler := func(messages []*Message) error {
		log.Printf("Processing batch of %d messages", len(messages))
		return nil
	}

	batchQueue := NewBatchQueue(3, 500*time.Millisecond, batchHandler)

	for i := 0; i < 7; i++ {
		msg := NewMessage(TypeTask, fmt.Sprintf("Batch item %d", i))
		batchQueue.Add(msg)
	}

	time.Sleep(2 * time.Second)
	batchQueue.Stop()
}

// ============================================================================
// بخش 9: Best Practices
// ============================================================================

func bestPractices() {
	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("💡 IN-MEMORY QUEUE BEST PRACTICES")
	fmt.Println(stringsRepeat("=", 80))

	fmt.Println(`
┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. BUFFER SIZE                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Choose appropriate buffer size for your use case                     │
│    • Too small: blocking producers                                         │
│    • Too large: memory waste                                              │
│    • Monitor queue length in production                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. GRACEFUL SHUTDOWN                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Always drain queues before shutdown                                  │
│    • Use sync.WaitGroup to wait for workers                               │
│    • Close channels to signal completion                                  │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. ERROR HANDLING                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Implement retry with backoff                                         │
│    • Use dead letter queues for failed messages                           │
│    • Log errors with context                                              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 4. MONITORING                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│    • Track queue length                                                   │
│    • Monitor worker utilization                                           │
│    • Alert on queue buildup                                               │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 5. WHEN TO USE (vs REAL BROKER)                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│    ✅ In-memory: internal communication, prototyping, low-volume         │
│    ❌ RabbitMQ/Kafka: cross-service, persistence, high-volume, replay     │
└─────────────────────────────────────────────────────────────────────────────┘
`)
}

// ============================================================================
// بخش 10: Main
// ============================================================================

func main() {
	fmt.Println(stringsRepeat("=", 80))
	fmt.Println("🎯 IN-MEMORY MESSAGE QUEUE WITH CHANNELS & GOROUTINES")
	fmt.Println("Simple, Priority, Pub/Sub, Batch, Delayed Queues")
	fmt.Println(stringsRepeat("=", 80))

	bestPractices()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🚀 Running Queue Examples")
	fmt.Println(stringsRepeat("=", 80))

	runQueueExample()

	fmt.Println("\n" + stringsRepeat("=", 80))
	fmt.Println("🎯 IN-MEMORY QUEUE - COMPLETE")
	fmt.Println("Build efficient internal message queues in Go!")
	fmt.Println(stringsRepeat("=", 80))
}

// تابع کمکی
func stringsRepeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
